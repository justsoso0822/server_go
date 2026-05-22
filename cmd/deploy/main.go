// 蓝绿部署工具
//
// 部署流程：
//  1. 确保 Traefik 网关运行
//  2. 构建镜像（仅 local + IMAGE_SOURCE=local）
//  3. 启动目标颜色容器（blue/green）
//  4. 等待新容器健康检查通过
//  5. 通知旧容器触发 traffic-shift（/health/lb 返回 503，Traefik 自动摘流）
//  6. 轮询网关确认流量已切至新容器
//  7. 通知旧容器拒绝新请求
//  8. 等待旧容器排空存量请求后移除
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	defaultAppName               = "server-go"
	defaultImageName             = "server-go"
	defaultLocalVersion          = "1.0.0"
	defaultGatewayHostPort       = "7001"
	defaultGatewayInternalPort   = "7001"
	defaultAppPort               = "7001"
	defaultDashboardPort         = "18080"
	defaultHealthTimeoutSeconds  = 60
	defaultCutoverTimeoutSeconds = 30
	defaultCutoverConfirmations  = 9
	defaultDrainTimeoutSeconds   = 15
	defaultKeepImages            = 10
	defaultImageSource           = "remote"
	defaultTraefikComposeFile    = "manifest/docker/compose/traefik.yml"
	defaultTraefikDashboardFile  = "manifest/docker/compose/traefik.dashboard.yml"
	defaultDockerfile            = "manifest/docker/Dockerfile"
	defaultLocalDBComposeFile    = "manifest/docker/compose/local.yml"
	defaultComposeDir            = "manifest/docker/compose"
	defaultLockTimeoutMinutes    = 30
	defaultGoProxy               = "https://proxy.golang.org,direct"
	defaultGoSumDB               = "sum.golang.org"
	defaultGoPrivate             = ""
	defaultGetOutputTimeout      = 60 * time.Second
)

var registryByEnv = map[string]string{
	"local":      "ccr.ccs.tencentyun.com/justsoso-local",
	"test":       "ccr.ccs.tencentyun.com/justsoso-test",
	"production": "ccr.ccs.tencentyun.com/justsoso-production",
}

// activeLockFile 记录当前持有的锁文件路径，用于异常退出时自动清理。
var (
	activeLockFile   string
	activeLockFileMu sync.Mutex
)

// projectRootDir 缓存自动探测到的项目根目录。
var projectRootDir string

// httpClient 用于网关健康探测，避免默认 client 无超时阻塞部署流程。
var httpClient = &http.Client{Timeout: 5 * time.Second}

// 当前正在运行的子进程引用，信号到达时优先转发给它，再退出。
var (
	currentChildMu sync.Mutex
	currentChild   *exec.Cmd
)

type deployConfig struct {
	Env                     string
	EnvFile                 string
	AppName                 string
	ImageName               string
	ImageSource             string
	Registry                string
	Version                 string
	GatewayHostPort         string
	GatewayInternalPort     string
	AppPort                 string
	DashboardPort           string
	DashboardEnabled        bool
	HealthTimeout           time.Duration
	CutoverTimeout          time.Duration
	CutoverConfirmations    int
	DrainTimeout            time.Duration
	KeepImages              int
	TraefikComposeFile      string
	TraefikDashboardFile    string
	Dockerfile              string
	ComposeDir              string
	LocalDBComposeFile      string
	ForceGatewayReplacement bool
}

type cleanupImageLine struct {
	// ref 是完整镜像引用，例如 registry/image:v1.2.3。
	ref string
	// tag 是镜像标签，例如 v1.2.3；latest 会在清理时跳过。
	tag string
	// createdAt 是 docker images 返回的原始创建时间文本。
	createdAt string
	// createdTS 是解析后的创建时间，用于按时间排序清理旧镜像。
	createdTS time.Time
}

// ============================================================================
// Entry Point & CLI
// ============================================================================

// main 顶层 recover：打印错误、释放锁、清理临时文件后退出。
func main() {
	defer func() {
		if r := recover(); r != nil {
			if ue, ok := r.(usageError); ok {
				fmt.Println(ue.message)
				printUsage()
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
			cleanupDeployTmpDir()
			releaseDeployLock()
			os.Exit(1)
		}
	}()
	run()
}

// run 解析子命令并分发到对应处理函数。
func run() {
	// 忽略 compose orphan 容器警告
	os.Setenv("COMPOSE_IGNORE_ORPHANS", "true")

	if len(os.Args) < 2 {
		panic(usageError{message: "Error: command required"})
	}

	initProjectRoot()

	cmd := os.Args[1]
	switch cmd {
	case "build":
		build()
	case "push":
		push()
	case "deploy":
		deploy()
	case "status":
		status()
	case "start-local-db":
		ensureNoExtraArgs(cmd)
		startLocalDB()
	case "stop-local-db":
		ensureNoExtraArgs(cmd)
		stopLocalDB()
	default:
		panic(usageError{message: fmt.Sprintf("Unknown command: %s", cmd)})
	}
}

// printUsage 输出命令行帮助信息。
func printUsage() {
	fmt.Println(`Usage: go run cmd/deploy/main.go <command> <env> [options]

Commands:
  build <env> [version=xxx]         构建镜像
  push <env> [version=xxx]          推送镜像
  deploy <env> [version=xxx] [-f]   部署（蓝绿切换，构建失败会保留旧实例）
  status [env]                      查看容器状态
  start-local-db                    启动本地数据库
  stop-local-db                     停止本地数据库

Environment:
  local       本地环境
  test        测试环境
  production  生产环境

Options:
  version=xxx  指定版本标签。test/production 的 build/push/deploy 必须显式指定
  -f           网关端口冲突时强制替换网关

Common config from .env.<env>:
  IMAGE_REGISTRY                   镜像仓库，未配置时按环境使用默认仓库
  IMAGE_NAME                       镜像名称，默认 server-go
  IMAGE_SOURCE                     local 或 remote
  HOST_GATEWAY_PORT                网关宿主机端口
  GATEWAY_INTERNAL_PORT            Traefik 容器内入口端口
  APP_PORT                         应用容器内 HTTP 端口
  DEPLOY_HEALTH_TIMEOUT_SECONDS    新实例健康检查超时
  DEPLOY_CUTOVER_TIMEOUT_SECONDS   切流确认超时
  DEPLOY_CUTOVER_CONFIRMATIONS     连续命中新颜色次数
  DEPLOY_DRAIN_TIMEOUT_SECONDS     旧实例排水超时
  DEPLOY_KEEP_IMAGES               本地保留镜像数量

Examples:
  go run main.go                                         # 本地开发，自动使用 config.yaml
  go run cmd/deploy/main.go start-local-db               # 启动本地数据库
  go run cmd/deploy/main.go build local                  # 构建本地镜像
  go run cmd/deploy/main.go build test version=v1.2.3    # 构建测试镜像并指定版本
  go run cmd/deploy/main.go push local                   # 推送本地镜像，默认 version=1.0.0
  go run cmd/deploy/main.go push test version=v1.2.3     # 推送测试镜像并指定版本
  go run cmd/deploy/main.go deploy test version=v1.2.3   # 部署到测试环境
  go run cmd/deploy/main.go deploy local -f              # 强制替换本地网关后部署
  go run cmd/deploy/main.go status local                 # 查看本地容器状态`)
}

// parseArgs 解析命令行参数，返回环境名和选项键值对（version=xxx, -f 等）。
func parseArgs() (string, map[string]string) {
	env := ""
	if len(os.Args) > 2 {
		env = os.Args[2]
	}

	options := make(map[string]string)
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "-f" {
			options["force"] = "true"
			continue
		}
		parts := strings.SplitN(os.Args[i], "=", 2)
		if len(parts) == 2 {
			options[parts[0]] = parts[1]
		}
	}
	return env, options
}

// ensureNoExtraArgs 校验 start-local-db / stop-local-db 这类不接受 env 参数的命令，
// 避免用户错误输入被静默忽略导致误操作。
func ensureNoExtraArgs(cmd string) {
	if len(os.Args) > 2 {
		panic(usageError{message: fmt.Sprintf("'%s' does not accept any arguments, got: %s", cmd, strings.Join(os.Args[2:], " "))})
	}
}

// ============================================================================
// Commands
// ============================================================================

// build 构建 Docker 镜像，同时打 version 和 latest 两个 tag。
func build() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	setupSignalHandler()
	acquireDeployLock(cfg.Env)
	defer releaseDeployLock()

	fmt.Printf("Building for environment: %s with version: %s\n", env, cfg.Version)
	buildImage(cfg, cfg.Version)
	fmt.Printf("Build completed: %s\n", formatImageName(cfg, cfg.Version))
}

// push 将 version 和 latest 两个 tag 的镜像推送到远程仓库。
func push() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	setupSignalHandler()
	acquireDeployLock(cfg.Env)
	defer releaseDeployLock()

	image := formatImageName(cfg, cfg.Version)
	imageLatest := formatImageName(cfg, "latest")

	fmt.Printf("Pushing image: %s\n", image)
	mustRun("docker", "push", image)
	mustRun("docker", "push", imageLatest)
	fmt.Printf("Push completed: %s and latest\n", image)
}

// deploy 执行蓝绿部署：启动新颜色 -> 健康检查 -> 网关确认 -> 切流 -> 排水 -> 移除旧颜色。
// 首次部署（无活跃容器）时跳过旧容器切流和排水步骤，但仍确认网关已路由到新容器。
// 通过文件锁防止同一环境的并发部署。
func deploy() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	setupSignalHandler()
	acquireDeployLock(cfg.Env)
	defer releaseDeployLock()

	fmt.Println("[release] [1/8] ensure traefik gateway")
	ensureGateway(cfg)

	currentColor, targetColor := detectDeploymentColors(cfg)
	if currentColor == "" {
		fmt.Printf("No running blue/green service found, deploying to %s\n", targetColor)
	} else {
		fmt.Printf("Current active: %s, deploying to: %s\n", currentColor, targetColor)
	}

	if cfg.Env == "local" && cfg.ImageSource == "local" {
		fmt.Println("[release] [2/8] local image source detected, building image")
		buildImage(cfg, cfg.Version)
		fmt.Printf("Build completed: %s\n", formatImageName(cfg, cfg.Version))
	}
	fmt.Printf("[release] [3/8] start %s (version=%s)\n", targetColor, cfg.Version)
	startColor(cfg, targetColor)

	fmt.Printf("[release] [4/8] wait for %s to be healthy (timeout=%s)\n", targetColor, cfg.HealthTimeout)
	if err := waitForHealthy(cfg, targetColor); err != nil {
		fmt.Printf("ERROR: %v, rolling back new %s deployment...\n", err, targetColor)
		_ = downColor(cfg, targetColor)
		panic(fmt.Errorf("rollback completed, deployment failed: %w", err))
	}

	if currentColor != "" {
		if err := cutover(cfg, currentColor, targetColor); err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("[release] [5/8] confirm gateway routes to %s (%d consecutive, timeout=%s)\n", targetColor, cfg.CutoverConfirmations, cfg.CutoverTimeout)
		if err := confirmCutover(cfg, targetColor); err != nil {
			fmt.Printf("ERROR: %v, rolling back new %s deployment...\n", err, targetColor)
			_ = downColor(cfg, targetColor)
			panic(fmt.Errorf("rollback completed, deployment failed: %w", err))
		}
	}

	fmt.Printf("\n[release] SUCCESS: %s now served by %s (version=%s)\n", cfg.Env, targetColor, cfg.Version)
	fmt.Printf("Gateway: http://localhost:%s\n", cfg.GatewayHostPort)
	if cfg.DashboardEnabled {
		fmt.Printf("Traefik Dashboard: http://localhost:%s/dashboard/\n", cfg.DashboardPort)
	}

	cleanupOldImages(cfg)
}

// status 显示容器健康状态、镜像版本、当前活跃颜色等运行信息。
func status() {
	env, _ := parseArgs()
	if len(os.Args) > 3 {
		panic(usageError{message: fmt.Sprintf("'status' accepts at most one environment argument, got: %s", strings.Join(os.Args[2:], " "))})
	}
	envFile := ""
	if env != "" {
		envFile = projectPath(fmt.Sprintf(".env.%s", env))
	}
	appName := getAppName(envFile)

	fmt.Printf("=== Deploy Status ===\n")
	if env != "" {
		fmt.Printf("Environment: %s\n\n", env)
	} else {
		fmt.Println("Environment: auto-detect")
	}

	fmt.Println("[Containers]")
	mustRun("docker", "ps", "--filter", fmt.Sprintf("name=%s", appName),
		"--format", "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}")

	containerNames := getContainerNames(appName)
	if len(containerNames) > 0 {
		fmt.Println("\n[Versions]")
		for _, name := range containerNames {
			version := getContainerEnv(name, "APP_VERSION")
			color := getContainerEnv(name, "APP_COLOR")
			if version == "" {
				version = "unknown"
			}
			if color == "" {
				color = "-"
			}
			fmt.Printf("  %-30s  color=%-5s  version=%s\n", name, color, version)
		}
	}

	fmt.Println("\n[Gateway Route]")
	gatewayPort := defaultGatewayHostPort
	if env != "" {
		gatewayPort = getEnvVar(projectPath(fmt.Sprintf(".env.%s", env)), "HOST_GATEWAY_PORT", defaultGatewayHostPort)
	}
	activeColor, err := gatewayActiveColor(gatewayPort)
	if err != nil {
		fmt.Printf("  Gateway not reachable or color unknown: %v\n", err)
	} else {
		fmt.Printf("  Active color: %s (via http://localhost:%s/health)\n", activeColor, gatewayPort)
	}

	fmt.Println("\n[Networks]")
	mustRun("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", appName))

	fmt.Println("\n[Volumes]")
	mustRun("docker", "volume", "ls", "--filter", fmt.Sprintf("name=%s", appName))
}

// startLocalDB 启动本地开发用的 MySQL 和 Redis 容器。
func startLocalDB() {
	cfg := loadDeployConfig("local", map[string]string{})
	fmt.Println("Starting local database services...")
	mustRun("docker", "compose", "-f", cfg.LocalDBComposeFile, "--env-file", cfg.EnvFile, "up", "-d")

	fmt.Println("Local database services started:")
	fmt.Println("  MySQL: 127.0.0.1:330")
	fmt.Println("  Redis: 127.0.0.1:637")
	fmt.Println("")
	fmt.Println("You can now run the application with:")
	fmt.Println("  go run main.go")
}

// stopLocalDB 停止本地开发数据库容器。
func stopLocalDB() {
	cfg := loadDeployConfig("local", map[string]string{})
	fmt.Println("Stopping local database services...")
	mustRun("docker", "compose", "-f", cfg.LocalDBComposeFile, "--env-file", cfg.EnvFile, "down")
	fmt.Println("Local database services stopped")
}

// ============================================================================
// Deploy Steps
// ============================================================================

// ensureGateway 确保 Traefik 网关容器运行且对外服务端口与配置一致。
func ensureGateway(cfg deployConfig) {
	gatewayRunning, err := containerExists(gatewayContainerName(cfg.AppName))
	if err != nil {
		panic(fmt.Errorf("failed to inspect gateway: %w", err))
	}
	currentGatewayHostPort := getGatewayHostPort(cfg)
	args := traefikComposeArgs(cfg)

	switch {
	case !gatewayRunning:
		mustRun("docker", append(args, "up", "-d")...)
		waitForGatewayHealthy(cfg)
	case currentGatewayHostPort == cfg.GatewayHostPort:
		fmt.Printf("[release] gateway already aligned on host port %s\n", cfg.GatewayHostPort)
	case cfg.ForceGatewayReplacement:
		fmt.Printf("[release] gateway config mismatch: current gateway=%s, desired gateway=%s, force replacing gateway\n", currentGatewayHostPort, cfg.GatewayHostPort)
		mustRun("docker", append(args, "up", "-d", "--force-recreate")...)
		waitForGatewayHealthy(cfg)
	default:
		panic(fmt.Errorf("ERROR: gateway config mismatch: current gateway=%s, desired gateway=%s\nRefusing to replace gateway automatically. Re-run with -f to force replace the gateway.", currentGatewayHostPort, cfg.GatewayHostPort))
	}
}

// waitForGatewayHealthy 轮询 docker inspect 等待 Traefik 容器 healthcheck 通过。
func waitForGatewayHealthy(cfg deployConfig) {
	name := gatewayContainerName(cfg.AppName)
	deadline := time.Now().Add(cfg.HealthTimeout)
	for time.Now().Before(deadline) {
		status, err := getOutput("docker", "inspect", "--format", "{{.State.Health.Status}}", name)
		if err == nil && status == "healthy" {
			fmt.Printf("[release] gateway %s healthy\n", name)
			return
		}
		if err == nil && status == "" {
			running, _ := containerExists(name)
			if running {
				fmt.Printf("[release] gateway %s running (no healthcheck configured)\n", name)
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	panic(fmt.Errorf("gateway %s did not become healthy within %s", name, cfg.HealthTimeout))
}

// buildImage 执行 docker build，注入 APP_PORT 和 Go 模块代理参数，同时打 version 和 latest tag。
func buildImage(cfg deployConfig, version string) {
	buildArgs := []string{
		"build",
		"--build-arg", fmt.Sprintf("APP_PORT=%s", cfg.AppPort),
		"--build-arg", fmt.Sprintf("GOPROXY=%s", getEnvWithDefault("GOPROXY", defaultGoProxy)),
		"--build-arg", fmt.Sprintf("GOSUMDB=%s", getEnvWithDefault("GOSUMDB", defaultGoSumDB)),
		"--build-arg", fmt.Sprintf("GOPRIVATE=%s", getEnvWithDefault("GOPRIVATE", defaultGoPrivate)),
		"-t", formatImageName(cfg, version),
		"-t", formatImageName(cfg, "latest"),
		"-f", cfg.Dockerfile,
		".",
	}
	mustRun("docker", buildArgs...)
}

// startColor 生成包含运行时变量的临时 env 文件，启动指定颜色的 compose 服务。
func startColor(cfg deployConfig, color string) {
	releaseEnvFile := writeReleaseEnvFile(cfg)
	defer os.Remove(releaseEnvFile)
	composeArgs := []string{"compose", "-f", composeFile(cfg, color), "--env-file", releaseEnvFile}
	if cfg.ImageSource == "remote" {
		mustRun("docker", append(composeArgs, "pull")...)
	}
	mustRun("docker", append(composeArgs, "up", "-d")...)
}

// writeReleaseEnvFile 基于原始 .env 文件生成临时发布 env 文件，
// 追加 APP_IMAGE、APP_VERSION 等运行时变量供 compose 使用。
func writeReleaseEnvFile(cfg deployConfig) string {
	content, err := os.ReadFile(cfg.EnvFile)
	if err != nil {
		panic(fmt.Errorf("read env file %s: %w", cfg.EnvFile, err))
	}

	tmpDir := projectPath(".deploy.tmp")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		panic(fmt.Errorf("create release env dir %s: %w", tmpDir, err))
	}

	file, err := os.CreateTemp(tmpDir, fmt.Sprintf("%s-release-*.env", cfg.AppName))
	if err != nil {
		panic(fmt.Errorf("create release env file: %w", err))
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		os.Remove(file.Name())
		panic(fmt.Errorf("write env content: %w", err))
	}

	extra := fmt.Sprintf("\nAPP_IMAGE=%s\nAPP_VERSION=%s\nAPP_PORT=%s\nGATEWAY_INTERNAL_PORT=%s\n",
		formatImageName(cfg, cfg.Version), cfg.Version, cfg.AppPort, cfg.GatewayInternalPort)
	if _, err := file.WriteString(extra); err != nil {
		os.Remove(file.Name())
		panic(fmt.Errorf("append release env values: %w", err))
	}
	return file.Name()
}

// detectDeploymentColors 检测当前活跃颜色和目标部署颜色。
func detectDeploymentColors(cfg deployConfig) (string, string) {
	blueRunning, err := containerExists(appContainerName(cfg.AppName, "blue"))
	if err != nil {
		panic(fmt.Errorf("failed to inspect blue container: %w", err))
	}
	greenRunning, err := containerExists(appContainerName(cfg.AppName, "green"))
	if err != nil {
		panic(fmt.Errorf("failed to inspect green container: %w", err))
	}

	switch {
	case blueRunning && greenRunning:
		active, err := gatewayActiveColor(cfg.GatewayHostPort)
		if err != nil {
			panic(fmt.Errorf("both blue and green are running, but active color cannot be determined from gateway: %w", err))
		}
		return active, oppositeColor(active)
	case blueRunning:
		return "blue", "green"
	case greenRunning:
		return "green", "blue"
	default:
		return "", "blue"
	}
}

// waitForHealthy 轮询 docker ps 等待目标容器健康检查通过，超时返回 error。
func waitForHealthy(cfg deployConfig, color string) error {
	deadline := time.Now().Add(cfg.HealthTimeout)
	for time.Now().Before(deadline) {
		output, err := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=^%s$", appContainerName(cfg.AppName, color)), "--filter", "health=healthy", "--format", "{{.Names}}")
		if err != nil {
			return err
		}
		if hasLine(output, appContainerName(cfg.AppName, color)) {
			fmt.Printf("[release] %s healthy\n", color)
			return nil
		}
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	return fmt.Errorf("%s service failed to become healthy", color)
}

// cutover 执行流量切换：通知旧容器摘流 -> 确认网关路由到新容器 -> 排水 -> 移除旧容器。
// 任何步骤失败都会保留旧容器，避免服务中断。
func cutover(cfg deployConfig, currentColor, targetColor string) error {
	oldContainerName := appContainerName(cfg.AppName, currentColor)

	fmt.Printf("[release] [5/8] http control -> %s: trigger traffic-shift, /health/lb now returns 503\n", currentColor)
	if err := postControl(oldContainerName, cfg.AppPort, "traffic-shift"); err != nil {
		return fmt.Errorf("failed to call traffic-shift on %s: %w. Keeping old container running", currentColor, err)
	}

	fmt.Printf("[release] [6/8] confirm gateway routes to %s (%d consecutive, timeout=%s)\n", targetColor, cfg.CutoverConfirmations, cfg.CutoverTimeout)
	if err := confirmCutover(cfg, targetColor); err != nil {
		resumeErr := postControl(oldContainerName, cfg.AppPort, "resume-traffic")
		if resumeErr != nil {
			return fmt.Errorf("cutover confirmation failed: %v. Failed to resume old container traffic: %w. Manual intervention required", err, resumeErr)
		}
		if downErr := downColor(cfg, targetColor); downErr != nil {
			return fmt.Errorf("cutover confirmation failed: %v. Resumed old container traffic, but failed to remove new %s container: %w. Manual cleanup required", err, targetColor, downErr)
		}
		return fmt.Errorf("cutover confirmation failed: %w. Resumed old container traffic and removed new %s container", err, targetColor)
	}

	fmt.Printf("[release] [7/8] http control -> %s: reject any remaining new requests\n", currentColor)
	if err := postControl(oldContainerName, cfg.AppPort, "reject-new-requests"); err != nil {
		return fmt.Errorf("failed to reject new requests on %s after cutover was confirmed: %w. New container remains active; old container may still be running and needs manual cleanup", currentColor, err)
	}

	fmt.Printf("[release] waiting %s in-flight requests (timeout=%s)\n", currentColor, cfg.DrainTimeout)
	if err := waitForDrain(oldContainerName, cfg.AppPort, cfg.DrainTimeout); err != nil {
		return fmt.Errorf("drain failed after cutover was confirmed: %w. New container remains active; old container was already removed from load balancing but was not removed", err)
	}

	fmt.Printf("[release] [8/8] %s: remove containers\n", currentColor)
	if err := downColor(cfg, currentColor); err != nil {
		return err
	}
	return nil
}

func downColor(cfg deployConfig, color string) error {
	downArgs := []string{"compose", "-f", composeFile(cfg, color), "--env-file", cfg.EnvFile, "down"}
	if err := runCmd("docker", downArgs...); err != nil {
		return fmt.Errorf("command failed: docker %s: %w", strings.Join(downArgs, " "), err)
	}
	return nil
}

// confirmCutover 轮询网关 /health 接口，确认连续 N 次返回目标颜色后视为切流成功。
// 任何一次探测失败会重置计数器。
func confirmCutover(cfg deployConfig, targetColor string) error {
	confirmed := 0
	deadline := time.Now().Add(cfg.CutoverTimeout)
	for time.Now().Before(deadline) {
		health, err := gatewayHealth(cfg.GatewayHostPort)
		if err == nil && health.Color == targetColor && health.Version == cfg.Version {
			confirmed++
			fmt.Printf("[release] gateway -> %s version=%s (%d/%d)\n", targetColor, health.Version, confirmed, cfg.CutoverConfirmations)
			if confirmed >= cfg.CutoverConfirmations {
				fmt.Printf("[release] cutover confirmed: all sampled traffic is on %s\n", targetColor)
				return nil
			}
		} else if confirmed > 0 {
			fmt.Printf("[release] probe reset (was %d)\n", confirmed)
			confirmed = 0
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("gateway did not route to %s version=%s for %d consecutive probes before timeout", targetColor, cfg.Version, cfg.CutoverConfirmations)
}

// waitForDrain 等待旧容器排空存量请求。
// 通过 /health/detail 的 activeRequests 字段判断；容器不可达时视为已排空。
func waitForDrain(containerName, appPort string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		output, err := getOutput("docker", "exec", containerName, "wget", "-q", "-O-", "--timeout=2", fmt.Sprintf("http://127.0.0.1:%s/health/detail", appPort))
		if err != nil {
			fmt.Printf("[release] %s: container unreachable, treating as drained\n", containerName)
			return nil
		}
		var detail struct {
			ActiveRequests int64 `json:"activeRequests"`
		}
		if err := json.Unmarshal([]byte(output), &detail); err != nil {
			return fmt.Errorf("parse %s /health/detail response: %w", containerName, err)
		}
		if detail.ActiveRequests == 0 {
			fmt.Printf("[release] %s: no in-flight requests\n", containerName)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("%s still has in-flight requests after %s", containerName, timeout)
}

// postControl 通过 docker exec 向容器内应用发送控制指令（traffic-shift / reject-new-requests）。
func postControl(containerName, appPort, action string) error {
	url := fmt.Sprintf("http://127.0.0.1:%s/_internal/control/%s", appPort, action)
	_, err := getOutput("docker", "exec", containerName, "wget", "-q", "-O-", "--timeout=5", "--post-data=", url)
	return err
}

// cleanupOldImages 按创建时间降序排列本地镜像，保留最新的 N 个 tag，删除其余。
// 通过完整 tag 引用（registry/name:tag）删除，避免多 tag 共享 image ID 导致的删除失败。
func cleanupOldImages(cfg deployConfig) {
	if cfg.KeepImages <= 0 {
		return
	}
	fmt.Printf("\n[cleanup] Removing old images (keeping latest %d versions)...\n", cfg.KeepImages)

	reference := fmt.Sprintf("%s/%s", cfg.Registry, cfg.ImageName)
	output, err := getOutput("docker", "images",
		"--filter", fmt.Sprintf("reference=%s:*", reference),
		"--format", "{{.Repository}}:{{.Tag}}|{{.CreatedAt}}",
		"--no-trunc")
	if err != nil {
		fmt.Printf("[cleanup] Warning: failed to list images: %v\n", err)
		return
	}
	if output == "" {
		fmt.Println("[cleanup] No images found")
		return
	}

	var images []cleanupImageLine
	for _, line := range strings.Split(output, "\n") {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			ref := parts[0]
			images = append(images, cleanupImageLine{
				ref:       ref,
				tag:       extractImageTag(ref),
				createdAt: parts[1],
				createdTS: parseDockerCreatedAt(parts[1]),
			})
		}
	}

	versions := make([]cleanupImageLine, 0, len(images))
	for _, img := range images {
		if img.tag == "" || img.tag == "latest" {
			continue
		}
		versions = append(versions, img)
	}
	if len(versions) == 0 {
		fmt.Println("[cleanup] No versioned images found")
		return
	}

	sort.Slice(versions, func(i, j int) bool {
		// 优先按解析后的时间降序；任一条目解析失败时退化为字符串比较，
		// 同时区下字符串顺序与时间顺序一致，跨时区/夏令时切换会被时间路径捕获。
		ti, tj := versions[i].createdTS, versions[j].createdTS
		if !ti.IsZero() && !tj.IsZero() {
			return ti.After(tj)
		}
		return versions[i].createdAt > versions[j].createdAt
	})

	if len(versions) <= cfg.KeepImages {
		fmt.Printf("[cleanup] Found %d versioned images, no cleanup needed\n", len(versions))
		return
	}

	deleted := 0
	for _, img := range versions[cfg.KeepImages:] {
		if err := runCmd("docker", "rmi", img.ref); err == nil {
			fmt.Printf("[cleanup] Removed: %s\n", img.ref)
			deleted++
		} else {
			fmt.Printf("[cleanup] Warning: failed to remove %s\n", img.ref)
		}
	}
	fmt.Printf("[cleanup] Cleanup complete: removed %d old images\n", deleted)
}

// parseDockerCreatedAt 解析 docker images --format {{.CreatedAt}} 输出的时间。
// 典型格式：'2024-01-02 15:04:05 +0800 CST'，部分版本可能省略时区名。
// 解析失败返回零值，由调用方决定退化策略。
// extractImageTag returns the tag portion of a docker image reference.
func extractImageTag(ref string) string {
	if i := strings.LastIndex(ref, ":"); i >= 0 {
		return ref[i+1:]
	}
	return ""
}

// parseDockerCreatedAt parses the CreatedAt field from docker images output.
func parseDockerCreatedAt(s string) time.Time {
	s = strings.TrimSpace(s)
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 MST",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// ============================================================================
// File Lock
// ============================================================================

// lockFilePath 返回指定环境的锁文件路径（位于项目根目录）。
func lockFilePath(env string) string {
	return projectPath(fmt.Sprintf(".deploy.%s.lock", env))
}

// acquireDeployLock 获取部署文件锁。
func acquireDeployLock(env string) {
	lockFile := lockFilePath(env)

	for {
		startedAt := time.Now()
		err := createDeployLockFile(lockFile, startedAt)
		if err == nil {
			activeLockFileMu.Lock()
			activeLockFile = lockFile
			activeLockFileMu.Unlock()
			fmt.Printf("[lock] Acquired deploy lock for environment '%s' (pid=%d)\n", env, os.Getpid())
			return
		}

		if os.IsExist(err) {
			handleExistingDeployLock(env, lockFile)
			continue
		}

		panic(fmt.Errorf("failed to create lock file %s: %w", lockFile, err))
	}
}

func createDeployLockFile(lockFile string, startedAt time.Time) error {
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lockContent := fmt.Sprintf("pid=%d\nstarted=%s\n", os.Getpid(), startedAt.Format(time.DateTime))
	if _, err := file.WriteString(lockContent); err != nil {
		os.Remove(lockFile)
		return err
	}
	return nil
}

// handleExistingDeployLock 处理已存在的锁文件：stale 则移除后返回（调用方会重试），
// 否则 panic 退出。
func handleExistingDeployLock(env, lockFile string) {
	info, err := os.Stat(lockFile)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(fmt.Errorf("failed to inspect lock file %s: %w", lockFile, err))
	}

	content, readErr := os.ReadFile(lockFile)
	if readErr == nil {
		ownerPID, startTime := parseLockContent(string(content))
		if !startTime.IsZero() && time.Since(startTime) > time.Duration(defaultLockTimeoutMinutes)*time.Minute {
			fmt.Printf("[lock] Stale lock detected (pid=%s, started=%s), removing...\n", ownerPID, startTime.Format(time.DateTime))
			if err := os.Remove(lockFile); err == nil || os.IsNotExist(err) {
				return
			}
			panic(fmt.Errorf("failed to remove stale lock file %s: %w", lockFile, err))
		}
		if isProcessAlive(ownerPID) {
			panic(fmt.Errorf("another deployment is in progress for environment '%s'\n  Lock holder: PID %s (started %s ago)\n  Lock file: %s\n\nIf you believe this is stale, delete the lock file manually:\n  rm %s",
				env, ownerPID, time.Since(startTime).Truncate(time.Second), lockFile, lockFile))
		}

		fmt.Printf("[lock] Lock holder (pid=%s) is no longer running, removing stale lock...\n", ownerPID)
		if err := os.Remove(lockFile); err == nil || os.IsNotExist(err) {
			return
		}
		panic(fmt.Errorf("failed to remove stale lock file %s: %w", lockFile, err))
	}

	if time.Since(info.ModTime()) > time.Duration(defaultLockTimeoutMinutes)*time.Minute {
		fmt.Println("[lock] Stale lock detected (unreadable, expired), removing...")
		if err := os.Remove(lockFile); err == nil || os.IsNotExist(err) {
			return
		}
		panic(fmt.Errorf("failed to remove stale lock file %s: %w", lockFile, err))
	}

	panic(fmt.Errorf("lock file exists but cannot be read: %s\nAnother deployment may be in progress for '%s'", lockFile, env))
}

// releaseDeployLock 释放当前持有的部署锁。
func releaseDeployLock() {
	activeLockFileMu.Lock()
	lockFile := activeLockFile
	activeLockFile = ""
	activeLockFileMu.Unlock()
	if lockFile == "" {
		return
	}
	if err := os.Remove(lockFile); err != nil {
		fmt.Printf("[lock] Warning: Failed to remove lock file %s: %v\n", lockFile, err)
		return
	}
	fmt.Printf("[lock] Released deploy lock: %s\n", lockFile)
}

// cleanupDeployTmpDir 清理 .deploy.tmp 目录，用于信号退出时兜底清理临时文件。
func cleanupDeployTmpDir() {
	tmpDir := projectPath(".deploy.tmp")
	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Printf("[cleanup] Warning: failed to remove %s: %v\n", tmpDir, err)
	}
}

// setupSignalHandler 注册信号处理器，确保进程被中断时优雅终止子进程并释放锁文件。
// 同时监听 os.Interrupt 和 SIGTERM(后者覆盖 docker stop / k8s preStop / kill 默认信号)。
// 信号到达时先转发给当前正在运行的子进程，给它一个清理窗口，再释放锁退出。
func setupSignalHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("\n[lock] Caught signal %v, terminating child process and releasing lock...\n", sig)

		currentChildMu.Lock()
		child := currentChild
		currentChildMu.Unlock()
		if child != nil && child.Process != nil {
			// Windows 不支持 syscall.SIGTERM 投递给子进程，统一用 Kill 兜底；
			// 类 Unix 系统先尝试转发原信号，让 docker compose 等子进程有机会清理。
			if runtime.GOOS == "windows" {
				_ = child.Process.Kill()
			} else {
				_ = child.Process.Signal(sig)
			}
		}

		// 等待短暂窗口让子进程退出，避免锁释放后立即 exit 导致子进程仍占用资源。
		time.Sleep(500 * time.Millisecond)

		cleanupDeployTmpDir()
		releaseDeployLock()
		os.Exit(1)
	}()
}

// parseLockContent 解析锁文件内容，提取 PID 和启动时间。
func parseLockContent(content string) (string, time.Time) {
	pid := "unknown"
	started := time.Time{}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "pid=") {
			pid = strings.TrimPrefix(line, "pid=")
		} else if strings.HasPrefix(line, "started=") {
			if t, err := time.ParseInLocation(time.DateTime, strings.TrimPrefix(line, "started="), time.Local); err == nil {
				started = t
			}
		}
	}
	return pid, started
}

// isProcessAlive 检查指定 PID 的进程是否仍在运行。
// Windows 通过 tasklist 命令查询，Unix 通过 kill -0 探测。
func isProcessAlive(pidStr string) bool {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// /FO CSV /NH 输出格式：每行 "ImageName","PID","Session","Session#","Mem"
		// 进程不存在时输出 "INFO: No tasks are running ..." 或空内容。
		// 通过精确匹配第二列的引号包裹 PID，避免 strings.Contains 的子串误判
		// （例如 PID=123 时不应命中 1234）。
		output, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH", "/FO", "CSV").Output()
		if err != nil || len(output) == 0 {
			return false
		}
		for _, line := range strings.Split(string(output), "\n") {
			fields := strings.Split(line, ",")
			if len(fields) >= 2 {
				pidField := strings.Trim(strings.TrimSpace(fields[1]), `"`)
				if pidField == strconv.Itoa(pid) {
					return true
				}
			}
		}
		return false
	}

	// Unix: 发送信号 0 探测进程是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}

// ============================================================================
// Configuration
// ============================================================================

// loadDeployConfig 从 .env.<env> 文件、系统环境变量和命令行选项中加载部署配置。
// 配置优先级：version/-f 来自命令行；其他配置为 .env 文件 > 系统环境变量 > 默认值。
// 非 local 环境必须显式指定 version。
func loadDeployConfig(env string, options map[string]string) deployConfig {
	if env == "" {
		panic(usageError{message: "Error: environment required"})
	}

	envFile := projectPath(fmt.Sprintf(".env.%s", env))
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		panic(fmt.Errorf("environment file not found: %s", envFile))
	} else if err != nil {
		panic(fmt.Errorf("inspect environment file %s: %w", envFile, err))
	}

	envVars := loadEnv(envFile)
	imageSource := imageSourceConfig(envVars)

	return deployConfig{
		Env:                     env,
		EnvFile:                 envFile,
		AppName:                 resolveConfig(envVars, "APP_NAME", defaultAppName),
		ImageName:               resolveConfig(envVars, "IMAGE_NAME", defaultImageName),
		ImageSource:             imageSource,
		Registry:                resolveConfig(envVars, "IMAGE_REGISTRY", defaultRegistry(env)),
		Version:                 getVersion(options, env),
		GatewayHostPort:         resolveConfig(envVars, "HOST_GATEWAY_PORT", defaultGatewayHostPort),
		GatewayInternalPort:     resolveConfig(envVars, "GATEWAY_INTERNAL_PORT", defaultGatewayInternalPort),
		AppPort:                 resolveConfig(envVars, "APP_PORT", defaultAppPort),
		DashboardPort:           resolveConfig(envVars, "TRAEFIK_DASHBOARD_PORT", defaultDashboardPort),
		DashboardEnabled:        boolConfig(envVars, "TRAEFIK_DASHBOARD_ENABLED", env == "local"),
		HealthTimeout:           secondsConfig(envVars, "DEPLOY_HEALTH_TIMEOUT_SECONDS", defaultHealthTimeoutSeconds),
		CutoverTimeout:          secondsConfig(envVars, "DEPLOY_CUTOVER_TIMEOUT_SECONDS", defaultCutoverTimeoutSeconds),
		CutoverConfirmations:    intConfig(envVars, "DEPLOY_CUTOVER_CONFIRMATIONS", defaultCutoverConfirmations),
		DrainTimeout:            secondsConfig(envVars, "DEPLOY_DRAIN_TIMEOUT_SECONDS", defaultDrainTimeoutSeconds),
		KeepImages:              intConfig(envVars, "DEPLOY_KEEP_IMAGES", defaultKeepImages),
		TraefikComposeFile:      projectPath(resolveConfig(envVars, "TRAEFIK_COMPOSE_FILE", defaultTraefikComposeFile)),
		TraefikDashboardFile:    projectPath(resolveConfig(envVars, "TRAEFIK_DASHBOARD_COMPOSE_FILE", defaultTraefikDashboardFile)),
		Dockerfile:              projectPath(resolveConfig(envVars, "DOCKERFILE", defaultDockerfile)),
		ComposeDir:              projectPath(resolveConfig(envVars, "COMPOSE_DIR", defaultComposeDir)),
		LocalDBComposeFile:      projectPath(resolveConfig(envVars, "LOCAL_DB_COMPOSE_FILE", defaultLocalDBComposeFile)),
		ForceGatewayReplacement: options["force"] == "true",
	}
}

// defaultRegistry 根据环境名返回对应的腾讯云镜像仓库地址。
func defaultRegistry(env string) string {
	registry, ok := registryByEnv[env]
	if !ok {
		panic(fmt.Errorf("unknown environment: %s", env))
	}
	return registry
}

func getVersion(options map[string]string, env string) string {
	if version, ok := options["version"]; ok && version != "" {
		return version
	}
	if env == "local" {
		return defaultLocalVersion
	}
	panic(fmt.Errorf("error: version parameter is required for %s", env))
}

func imageSourceConfig(fileVars map[string]string) string {
	value := resolveConfig(fileVars, "IMAGE_SOURCE", defaultImageSource)
	switch value {
	case "local", "remote":
		return value
	default:
		panic(fmt.Errorf("invalid IMAGE_SOURCE=%q, expected local or remote", value))
	}
}

// resolveConfig 按优先级获取配置值：env 文件 > 系统环境变量 > 默认值。
func resolveConfig(fileVars map[string]string, key, defaultVal string) string {
	if val, ok := fileVars[key]; ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	if val := os.Getenv(key); strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

// getEnvWithDefault 从系统环境读取值，未设置时返回默认值。
func getEnvWithDefault(key, defaultVal string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return defaultVal
}

func secondsConfig(fileVars map[string]string, key string, defaultVal int) time.Duration {
	return time.Duration(intConfig(fileVars, key, defaultVal)) * time.Second
}

func intConfig(fileVars map[string]string, key string, defaultVal int) int {
	value := resolveConfig(fileVars, key, "")
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		panic(fmt.Errorf("invalid %s=%q, expected positive integer", key, value))
	}
	return parsed
}

func boolConfig(fileVars map[string]string, key string, defaultVal bool) bool {
	value := resolveConfig(fileVars, key, "")
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(value)
	if err == nil {
		return parsed
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes", "on":
		return true
	case "no", "off":
		return false
	default:
		panic(fmt.Errorf("invalid %s=%q, expected boolean", key, value))
	}
}

func loadEnv(envFile string) map[string]string {
	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		panic(fmt.Errorf("open env file %s: %w", envFile, err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[strings.TrimSpace(parts[0])] = cleanEnvValue(parts[1])
		}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("failed to read env file %s: %w", envFile, err))
	}
	return env
}

// cleanEnvValue 清理 env 值：去除首尾引号、行内注释（未被引号包裹时）。
func cleanEnvValue(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		quote := value[0]
		if (quote == '\'' || quote == '"') && value[len(value)-1] == quote {
			return value[1 : len(value)-1]
		}
	}
	if idx := strings.Index(value, " #"); idx >= 0 {
		value = value[:idx]
	}
	return strings.TrimSpace(value)
}

// getEnvVar 从指定 env 文件中读取单个配置值的便捷方法。
func getEnvVar(envFile, key, defaultVal string) string {
	info, err := os.Stat(envFile)
	if os.IsNotExist(err) {
		return defaultVal
	}
	if err != nil {
		panic(fmt.Errorf("inspect env file %s: %w", envFile, err))
	}
	if info.IsDir() {
		panic(fmt.Errorf("env file path is a directory: %s", envFile))
	}
	env := loadEnv(envFile)
	return resolveConfig(env, key, defaultVal)
}

// getAppName 从 env 文件中获取应用名称，未指定环境时依次尝试 local/test/production。
func getAppName(envFile string) string {
	if envFile != "" {
		return getEnvVar(envFile, "APP_NAME", defaultAppName)
	}
	for _, file := range []string{
		projectPath(".env.local"),
		projectPath(".env.test"),
		projectPath(".env.production"),
	} {
		if appName := getEnvVar(file, "APP_NAME", ""); appName != "" {
			return appName
		}
	}
	return defaultAppName
}

// ============================================================================
// Helpers
// ============================================================================

// initProjectRoot 探测并缓存项目根目录，允许在项目内移动脚本或从不同目录执行。
func initProjectRoot() {
	if projectRootDir != "" {
		return
	}

	candidates := make([]string, 0, 3)
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, wd)
	}
	if exePath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Dir(exePath))
	}
	if _, currentFile, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Dir(currentFile))
	}

	for _, start := range candidates {
		if root, ok := findProjectRoot(start); ok {
			projectRootDir = root
			return
		}
	}

	panic(fmt.Errorf("failed to locate project root from current working directory, executable path, or source file path"))
}

// projectRoot 返回已缓存的项目根目录。
func projectRoot() string {
	return projectRootDir
}

// findProjectRoot 从起始目录向上查找项目根目录。
func findProjectRoot(start string) (string, bool) {
	current := start
	for {
		if fileExists(filepath.Join(current, "go.mod")) && fileExists(filepath.Join(current, "manifest", "docker", "Dockerfile")) {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

// projectPath 将项目相对路径解析为基于项目根目录的绝对路径。
func projectPath(parts ...string) string {
	segments := append([]string{projectRoot()}, parts...)
	return filepath.Join(segments...)
}

// fileExists 判断文件或目录是否存在。
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// formatImageName 拼接完整镜像名：registry/image:tag。
func formatImageName(cfg deployConfig, tag string) string {
	return fmt.Sprintf("%s/%s:%s", cfg.Registry, cfg.ImageName, tag)
}

// composeFile 返回指定颜色的 compose 文件路径（如 manifest/docker/compose/blue.yml）。
func composeFile(cfg deployConfig, color string) string {
	return fmt.Sprintf("%s/%s.yml", strings.TrimRight(cfg.ComposeDir, "/\\"), color)
}

// traefikComposeArgs 返回 Traefik compose 基础参数；Dashboard 仅在显式开启时挂载端口和路由。
func traefikComposeArgs(cfg deployConfig) []string {
	args := []string{"compose", "-f", cfg.TraefikComposeFile}
	if cfg.DashboardEnabled {
		args = append(args, "-f", cfg.TraefikDashboardFile)
	}
	return append(args, "--env-file", cfg.EnvFile)
}

// runCmd 执行外部命令，stdout/stderr 直接输出到终端。
// 同时把当前 *exec.Cmd 注册为全局 currentChild，便于信号处理器中断时转发信号给子进程。
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectRoot()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	currentChildMu.Lock()
	currentChild = cmd
	currentChildMu.Unlock()
	defer func() {
		currentChildMu.Lock()
		currentChild = nil
		currentChildMu.Unlock()
	}()

	return cmd.Run()
}

// mustRun 执行外部命令，失败时 panic。
func mustRun(name string, args ...string) {
	if err := runCmd(name, args...); err != nil {
		panic(fmt.Errorf("command failed: %s %s: %w", name, strings.Join(args, " "), err))
	}
}

// getOutput 执行外部命令并捕获 stdout 输出，带默认超时避免 docker 偶发挂死阻塞部署流程。
// 调用方需要更长超时时改用 getOutputWithTimeout。
func getOutput(name string, args ...string) (string, error) {
	return getOutputWithTimeout(defaultGetOutputTimeout, name, args...)
}

// getOutputWithTimeout 执行外部命令并捕获 stdout 输出，使用指定超时。
func getOutputWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = projectRoot()
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s %s: timed out after %s", name, strings.Join(args, " "), timeout)
		}
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// healthResponse 解析 /health 接口返回的 JSON。
type healthResponse struct {
	Color   string `json:"color"`
	Version string `json:"version"`
}

// gatewayActiveColor 通过网关 /health 接口解析当前活跃的部署颜色。
func gatewayActiveColor(port string) (string, error) {
	health, err := gatewayHealth(port)
	if err != nil {
		return "", err
	}
	return health.Color, nil
}

// gatewayHealth 通过指定端口访问网关 /health 接口，返回解析后的健康信息。
// 使用 encoding/json 反序列化，避免字符串匹配误判；通过带超时的 httpClient 防止
// 网关挂死时无限阻塞 confirmCutover 的轮询循环。
func gatewayHealth(port string) (healthResponse, error) {
	resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/health", port))
	if err != nil {
		return healthResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return healthResponse{}, fmt.Errorf("/health returned HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return healthResponse{}, err
	}
	var hr healthResponse
	if err := json.Unmarshal(body, &hr); err != nil {
		return healthResponse{}, fmt.Errorf("parse /health response: %w", err)
	}
	if hr.Color == "" {
		return healthResponse{}, fmt.Errorf("/health response missing color field")
	}
	return hr, nil
}

// containerExists 检查指定名称的容器是否正在运行。
func containerExists(name string) (bool, error) {
	output, err := getOutput("docker", "ps", "--format", "{{.Names}}")
	if err != nil {
		return false, err
	}
	return hasLine(output, name), nil
}

// hasLine 检查多行输出中是否包含精确匹配的行。
func hasLine(output, expected string) bool {
	for line := range strings.SplitSeq(output, "\n") {
		if strings.TrimSpace(line) == expected {
			return true
		}
	}
	return false
}

// getContainerNames 返回匹配 appName 的运行中容器名称列表。
func getContainerNames(appName string) []string {
	output, err := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=%s", appName), "--format", "{{.Names}}")
	if err != nil || output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}

// getContainerEnv 从容器环境变量中读取指定 key 的值。
func getContainerEnv(containerName, key string) string {
	output, err := getOutput("docker", "exec", containerName, "printenv", key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(output)
}

// gatewayContainerName 返回网关容器名称：{appName}-gateway。
func gatewayContainerName(appName string) string {
	return fmt.Sprintf("%s-gateway", appName)
}

// appContainerName 返回应用容器名称：{appName}-{color}。
func appContainerName(appName, color string) string {
	return fmt.Sprintf("%s-%s", appName, color)
}

// oppositeColor 返回对立颜色：blue -> green, green -> blue。
func oppositeColor(color string) string {
	switch color {
	case "blue":
		return "green"
	case "green":
		return "blue"
	default:
		panic(fmt.Errorf("unknown deployment color: %s", color))
	}
}

// getGatewayHostPort 通过 docker port 查询网关容器实际映射的宿主机端口。
func getGatewayHostPort(cfg deployConfig) string {
	return getContainerHostPort(gatewayContainerName(cfg.AppName), cfg.GatewayInternalPort)
}

func getContainerHostPort(containerName, internalPort string) string {
	output, err := getOutput("docker", "port", containerName, fmt.Sprintf("%s/tcp", internalPort))
	if err != nil || output == "" {
		return ""
	}
	firstLine := strings.SplitN(strings.TrimSpace(output), "\n", 2)[0]
	parts := strings.Split(strings.TrimSpace(firstLine), ":")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[len(parts)-1])
}

// ============================================================================
// Error Handling
// ============================================================================

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return e.message
}
