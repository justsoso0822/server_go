// 蓝绿部署工具
//
// 部署流程：
//   1. 确保 Traefik 网关运行
//   2. 构建镜像（仅 local + IMAGE_SOURCE=local）
//   3. 启动目标颜色容器（blue/green）
//   4. 等待新容器健康检查通过
//   5. 通知旧容器触发 traffic-shift（/health/lb 返回 503，Traefik 自动摘流）
//   6. 轮询网关确认流量已切至新容器
//   7. 通知旧容器拒绝新请求
//   8. 等待旧容器排空存量请求后移除
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAppName               = "server-go"
	defaultImageName             = "server-go"
	defaultLocalVersion          = "1.0.0"
	defaultGatewayHostPort       = "7001"
	defaultGatewayInternalPort   = "7001"
	defaultAppInternalPort       = "7001"
	defaultDashboardPort         = "18080"
	defaultHealthTimeoutSeconds  = 60
	defaultCutoverTimeoutSeconds = 30
	defaultCutoverConfirmations  = 9
	defaultDrainTimeoutSeconds   = 15
	defaultKeepImages            = 10
	defaultImageSource           = "remote"
	defaultTraefikComposeFile    = "manifest/docker/compose/traefik.yml"
	defaultDockerfile            = "manifest/docker/Dockerfile"
	defaultLocalDBComposeFile    = "manifest/docker/compose/local.yml"
	defaultComposeDir            = "manifest/docker/compose"
)

var registryByEnv = map[string]string{
	"local":      "ccr.ccs.tencentyun.com/justsoso-local",
	"test":       "ccr.ccs.tencentyun.com/justsoso-test",
	"production": "ccr.ccs.tencentyun.com/justsoso-production",
}

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
	AppInternalPort         string
	DashboardPort           string
	HealthTimeout           time.Duration
	CutoverTimeout          time.Duration
	CutoverConfirmations    int
	DrainTimeout            time.Duration
	KeepImages              int
	TraefikComposeFile      string
	Dockerfile              string
	ComposeDir              string
	LocalDBComposeFile      string
	ForceGatewayReplacement bool
}

// main 解析子命令并分发到对应处理函数。
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

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
		startLocalDB()
	case "stop-local-db":
		stopLocalDB()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// printUsage 输出命令行帮助信息。
func printUsage() {
	fmt.Println(`Usage: go run cmd/deploy/main.go <command> <env> [options]

Commands:
  build <env> [version=xxx]         构建镜像
  push <env> [version=xxx]          推送镜像
  deploy <env> [version=xxx] [-f]   部署（蓝绿切换，关键失败会保留旧实例）
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
  APP_INTERNAL_PORT                应用容器内 HTTP 端口
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

// loadDeployConfig 从 .env.<env> 文件和命令行选项中加载部署配置。
// 配置优先级：命令行选项 > .env 文件 > 系统环境变量 > 默认值。
// 非 local 环境必须显式指定 version。
func loadDeployConfig(env string, options map[string]string) deployConfig {
	if env == "" {
		fatalUsage("Error: environment required")
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fatalf("Environment file not found: %s", envFile)
	}

	envVars := loadEnv(envFile)
	registry := getConfigValue(envVars, "IMAGE_REGISTRY", defaultRegistry(env))
	imageSource := getConfigValue(envVars, "IMAGE_SOURCE", defaultImageSource)

	version := getVersion(options, env, env != "local")
	if env == "local" && imageSource == "local" && version == "" {
		version = defaultLocalVersion
	}
	if version == "" {
		version = "latest"
	}

	return deployConfig{
		Env:                     env,
		EnvFile:                 envFile,
		AppName:                 getConfigValue(envVars, "APP_NAME", defaultAppName),
		ImageName:               getConfigValue(envVars, "IMAGE_NAME", defaultImageName),
		ImageSource:             imageSource,
		Registry:                registry,
		Version:                 version,
		GatewayHostPort:         getConfigValue(envVars, "HOST_GATEWAY_PORT", defaultGatewayHostPort),
		GatewayInternalPort:     getConfigValue(envVars, "GATEWAY_INTERNAL_PORT", defaultGatewayInternalPort),
		AppInternalPort:         getConfigValue(envVars, "APP_INTERNAL_PORT", defaultAppInternalPort),
		DashboardPort:           getConfigValue(envVars, "TRAEFIK_DASHBOARD_PORT", defaultDashboardPort),
		HealthTimeout:           secondsConfig(envVars, "DEPLOY_HEALTH_TIMEOUT_SECONDS", defaultHealthTimeoutSeconds),
		CutoverTimeout:          secondsConfig(envVars, "DEPLOY_CUTOVER_TIMEOUT_SECONDS", defaultCutoverTimeoutSeconds),
		CutoverConfirmations:    intConfig(envVars, "DEPLOY_CUTOVER_CONFIRMATIONS", defaultCutoverConfirmations),
		DrainTimeout:            secondsConfig(envVars, "DEPLOY_DRAIN_TIMEOUT_SECONDS", defaultDrainTimeoutSeconds),
		KeepImages:              intConfig(envVars, "DEPLOY_KEEP_IMAGES", defaultKeepImages),
		TraefikComposeFile:      getConfigValue(envVars, "TRAEFIK_COMPOSE_FILE", defaultTraefikComposeFile),
		Dockerfile:              getConfigValue(envVars, "DOCKERFILE", defaultDockerfile),
		ComposeDir:              getConfigValue(envVars, "COMPOSE_DIR", defaultComposeDir),
		LocalDBComposeFile:      getConfigValue(envVars, "LOCAL_DB_COMPOSE_FILE", defaultLocalDBComposeFile),
		ForceGatewayReplacement: options["force"] == "true",
	}
}

// defaultRegistry 根据环境名返回对应的腾讯云镜像仓库地址。
func defaultRegistry(env string) string {
	registry, ok := registryByEnv[env]
	if !ok {
		fatalf("Unknown environment: %s", env)
	}
	return registry
}

// getVersion 从选项中获取版本号，非 local 环境未指定时 fatal 退出。
func getVersion(options map[string]string, env string, required bool) string {
	if version, ok := options["version"]; ok && version != "" {
		return version
	}
	if required {
		fatalf("Error: version parameter is required for %s", env)
	}
	if env == "local" {
		return defaultLocalVersion
	}
	return ""
}

// getConfigValue 按优先级获取配置值：env 文件 > 系统环境变量 > 默认值。
func getConfigValue(env map[string]string, key, defaultVal string) string {
	if val, ok := env[key]; ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	if val := os.Getenv(key); strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

// secondsConfig 读取整数配置并转换为 time.Duration（秒）。
func secondsConfig(env map[string]string, key string, defaultVal int) time.Duration {
	return time.Duration(intConfig(env, key, defaultVal)) * time.Second
}

// intConfig 读取正整数配置值，无效值时 fatal 退出。
func intConfig(env map[string]string, key string, defaultVal int) int {
	value := getConfigValue(env, key, "")
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		fatalf("Invalid %s=%q, expected positive integer", key, value)
	}
	return parsed
}

// imageRef 拼接完整镜像引用：registry/image:tag。
func imageRef(cfg deployConfig, tag string) string {
	return fmt.Sprintf("%s/%s:%s", cfg.Registry, cfg.ImageName, tag)
}

// composeFile 返回指定颜色的 compose 文件路径（如 manifest/docker/compose/blue.yml）。
func composeFile(cfg deployConfig, color string) string {
	return fmt.Sprintf("%s/%s.yml", strings.TrimRight(cfg.ComposeDir, "/\\"), color)
}

// runCmd 执行外部命令，stdout/stderr 直接输出到终端。
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// mustRun 执行外部命令，失败时 fatal 退出。
func mustRun(name string, args ...string) {
	if err := runCmd(name, args...); err != nil {
		fatalf("Command failed: %s %s: %v", name, strings.Join(args, " "), err)
	}
}

// getOutput 执行外部命令并捕获 stdout 输出。
func getOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// loadEnv 解析 .env 文件，支持 # 注释、引号包裹值、export 前缀。
func loadEnv(envFile string) map[string]string {
	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		return env
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
	return getConfigValue(loadEnv(envFile), key, defaultVal)
}

// build 构建 Docker 镜像，同时打 version 和 latest 两个 tag。
func build() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	fmt.Printf("Building for environment: %s with version: %s\n", env, cfg.Version)
	buildImage(cfg, cfg.Version)
	fmt.Printf("Build completed: %s\n", imageRef(cfg, cfg.Version))
}

// push 将 version 和 latest 两个 tag 的镜像推送到远程仓库。
func push() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	image := imageRef(cfg, cfg.Version)
	imageLatest := imageRef(cfg, "latest")

	fmt.Printf("Pushing image: %s\n", image)
	mustRun("docker", "push", image)
	mustRun("docker", "push", imageLatest)
	fmt.Printf("Push completed: %s and latest\n", image)
}

// deploy 执行蓝绿部署：启动新颜色 -> 健康检查 -> 切流 -> 排水 -> 移除旧颜色。
// 首次部署（无活跃容器）时跳过切流和排水步骤。
func deploy() {
	env, options := parseArgs()
	cfg := loadDeployConfig(env, options)

	currentColor, targetColor := detectDeploymentColors(cfg)
	if currentColor == "" {
		fmt.Printf("No active deployment found, deploying to %s\n", targetColor)
	} else {
		fmt.Printf("Current active: %s, deploying to: %s\n", currentColor, targetColor)
	}

	fmt.Println("[release] [1/8] ensure traefik gateway")
	ensureGateway(cfg)

	if cfg.Env == "local" && cfg.ImageSource == "local" {
		fmt.Println("[release] [2/8] local image source detected, building image")
		buildImage(cfg, cfg.Version)
		fmt.Printf("Build completed: %s\n", imageRef(cfg, cfg.Version))
	}

	fmt.Printf("[release] [3/8] start %s (version=%s)\n", targetColor, cfg.Version)
	startColor(cfg, targetColor)

	fmt.Printf("[release] [4/8] wait for %s to be healthy (timeout=%s)\n", targetColor, cfg.HealthTimeout)
	if err := waitForHealthy(cfg, targetColor); err != nil {
		fmt.Printf("ERROR: %v, rolling back new %s deployment...\n", err, targetColor)
		mustRun("docker", "compose", "-f", composeFile(cfg, targetColor), "--env-file", cfg.EnvFile, "down")
		fatalf("Rollback completed, deployment failed")
	}

	if currentColor != "" {
		cutover(cfg, currentColor, targetColor)
	}

	fmt.Printf("\n[release] SUCCESS: %s now served by %s (version=%s)\n", cfg.Env, targetColor, cfg.Version)
	fmt.Printf("Gateway: http://localhost:%s\n", cfg.GatewayHostPort)
	fmt.Printf("Traefik Dashboard: http://localhost:%s/dashboard/\n", cfg.DashboardPort)

	cleanupOldImages(cfg)
}

// buildImage 执行 docker build，注入 APP_PORT 构建参数，同时打 version 和 latest tag。
func buildImage(cfg deployConfig, version string) {
	mustRun("docker", "build", "--build-arg", fmt.Sprintf("APP_PORT=%s", cfg.AppInternalPort), "-t", imageRef(cfg, version), "-t", imageRef(cfg, "latest"), "-f", cfg.Dockerfile, ".")
}

// startColor 生成包含运行时变量的临时 env 文件，启动指定颜色的 compose 服务。
func startColor(cfg deployConfig, color string) {
	releaseEnvFile := writeReleaseEnvFile(cfg)
	defer os.Remove(releaseEnvFile)
	composeArgs := []string{"compose", "-f", composeFile(cfg, color), "--env-file", releaseEnvFile, "up", "-d"}
	mustRun("docker", composeArgs...)
}

// writeReleaseEnvFile 基于原始 .env 文件生成临时发布 env 文件，
// 追加 APP_IMAGE、APP_VERSION 等运行时变量供 compose 使用。
func writeReleaseEnvFile(cfg deployConfig) string {
	content, err := os.ReadFile(cfg.EnvFile)
	if err != nil {
		fatalf("Failed to read env file %s: %v", cfg.EnvFile, err)
	}

	file, err := os.CreateTemp("", fmt.Sprintf("%s-release-*.env", cfg.AppName))
	if err != nil {
		fatalf("Failed to create release env file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		fatalf("Failed to write release env file: %v", err)
	}
	if _, err := fmt.Fprintf(file, "\nAPP_IMAGE=%s\nAPP_VERSION=%s\nAPP_INTERNAL_PORT=%s\nGATEWAY_INTERNAL_PORT=%s\n", imageRef(cfg, cfg.Version), cfg.Version, cfg.AppInternalPort, cfg.GatewayInternalPort); err != nil {
		fatalf("Failed to append release env values: %v", err)
	}
	return file.Name()
}

// ensureGateway 确保 Traefik 网关容器运行且端口与配置一致。
// 端口不一致时需要 -f 标志才会强制重建，避免意外中断线上流量。
func ensureGateway(cfg deployConfig) {
	gatewayRunning, err := containerExists(gatewayContainerName(cfg.AppName))
	if err != nil {
		fatalf("Failed to inspect gateway: %v", err)
	}
	currentGatewayHostPort := getGatewayHostPort(cfg)

	switch {
	case !gatewayRunning:
		mustRun("docker", "compose", "-f", cfg.TraefikComposeFile, "--env-file", cfg.EnvFile, "up", "-d")
		time.Sleep(2 * time.Second)
	case currentGatewayHostPort == cfg.GatewayHostPort:
		fmt.Printf("[release] gateway already aligned on host port %s\n", cfg.GatewayHostPort)
	case cfg.ForceGatewayReplacement:
		fmt.Printf("[release] gateway host port mismatch: current=%s desired=%s, force replacing gateway\n", currentGatewayHostPort, cfg.GatewayHostPort)
		mustRun("docker", "compose", "-f", cfg.TraefikComposeFile, "--env-file", cfg.EnvFile, "up", "-d", "--force-recreate")
		time.Sleep(2 * time.Second)
	default:
		fatalf("ERROR: gateway host port mismatch: current=%s desired=%s\nRefusing to replace gateway automatically. Re-run with -f to force replace the gateway.", currentGatewayHostPort, cfg.GatewayHostPort)
	}
}

// detectDeploymentColors 检测当前活跃颜色和目标部署颜色。
// 双容器同时运行时通过网关 /health 响应判断活跃方。
func detectDeploymentColors(cfg deployConfig) (string, string) {
	blueRunning, err := containerExists(colorContainerName(cfg.AppName, "blue"))
	if err != nil {
		fatalf("Failed to inspect blue container: %v", err)
	}
	greenRunning, err := containerExists(colorContainerName(cfg.AppName, "green"))
	if err != nil {
		fatalf("Failed to inspect green container: %v", err)
	}

	switch {
	case blueRunning && greenRunning:
		active, err := gatewayActiveColor(cfg)
		if err != nil {
			fatalf("Both blue and green are running, but active color cannot be determined from gateway: %v", err)
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

// cutover 执行流量切换：通知旧容器摘流 -> 确认网关路由到新容器 -> 排水 -> 移除旧容器。
// 任何步骤失败都会保留旧容器，避免服务中断。
func cutover(cfg deployConfig, currentColor, targetColor string) {
	oldContainerName := colorContainerName(cfg.AppName, currentColor)

	fmt.Printf("[release] [5/8] http control -> %s: trigger traffic-shift, /health/lb now returns 503\n", currentColor)
	if err := postControl(oldContainerName, cfg.AppInternalPort, "traffic-shift"); err != nil {
		fatalf("Failed to call traffic-shift on %s: %v. Keeping old container running.", currentColor, err)
	}

	fmt.Printf("[release] [6/8] confirm gateway routes to %s (%d consecutive, timeout=%s)\n", targetColor, cfg.CutoverConfirmations, cfg.CutoverTimeout)
	if err := confirmCutover(cfg, targetColor); err != nil {
		fatalf("Cutover confirmation failed: %v. Keeping old container running.", err)
	}

	fmt.Printf("[release] [7/8] http control -> %s: reject any remaining new requests\n", currentColor)
	if err := postControl(oldContainerName, cfg.AppInternalPort, "reject-new-requests"); err != nil {
		fatalf("Failed to reject new requests on %s: %v. Keeping old container running.", currentColor, err)
	}

	fmt.Printf("[release] waiting %s in-flight requests (timeout=%s)\n", currentColor, cfg.DrainTimeout)
	if err := waitForDrain(oldContainerName, cfg.AppInternalPort, cfg.DrainTimeout); err != nil {
		fatalf("Drain failed: %v. Keeping old container running.", err)
	}

	fmt.Printf("[release] [8/8] %s: remove containers\n", currentColor)
	mustRun("docker", "compose", "-f", composeFile(cfg, currentColor), "--env-file", cfg.EnvFile, "down")
}

// waitForHealthy 轮询 docker ps 等待目标容器健康检查通过，超时返回 error。
func waitForHealthy(cfg deployConfig, color string) error {
	deadline := time.Now().Add(cfg.HealthTimeout)
	for time.Now().Before(deadline) {
		output, err := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=^%s$", colorContainerName(cfg.AppName, color)), "--filter", "health=healthy", "--format", "{{.Names}}")
		if err != nil {
			return err
		}
		if hasLine(output, colorContainerName(cfg.AppName, color)) {
			fmt.Printf("[release] %s healthy\n", color)
			return nil
		}
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	return fmt.Errorf("%s service failed to become healthy", color)
}

// confirmCutover 轮询网关 /health 接口，确认连续 N 次返回目标颜色后视为切流成功。
// 任何一次探测失败会重置计数器。
func confirmCutover(cfg deployConfig, targetColor string) error {
	confirmed := 0
	deadline := time.Now().Add(cfg.CutoverTimeout)
	for time.Now().Before(deadline) {
		active, err := gatewayActiveColor(cfg)
		if err == nil && active == targetColor {
			confirmed++
			fmt.Printf("[release] gateway -> %s (%d/%d)\n", targetColor, confirmed, cfg.CutoverConfirmations)
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
	return fmt.Errorf("gateway did not route to %s for %d consecutive probes before timeout", targetColor, cfg.CutoverConfirmations)
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
		if strings.Contains(output, `"activeRequests":0`) {
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
	return runCmd("docker", "exec", containerName, "wget", "-q", "-O-", "--timeout=5", "--post-data=", url)
}

// gatewayActiveColor 通过网关 /health 接口解析当前活跃的部署颜色。
func gatewayActiveColor(cfg deployConfig) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", cfg.GatewayHostPort))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	content := string(body)
	switch {
	case strings.Contains(content, `"color":"blue"`):
		return "blue", nil
	case strings.Contains(content, `"color":"green"`):
		return "green", nil
	default:
		return "", fmt.Errorf("/health response does not contain deployment color: %s", strings.TrimSpace(content))
	}
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

// gatewayContainerName 返回网关容器名称：{appName}-gateway。
func gatewayContainerName(appName string) string {
	return fmt.Sprintf("%s-gateway", appName)
}

// colorContainerName 返回颜色容器名称：{appName}-{color}。
func colorContainerName(appName, color string) string {
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
		fatalf("Unknown deployment color: %s", color)
		return ""
	}
}

// getGatewayHostPort 通过 docker port 查询网关容器实际映射的宿主机端口。
func getGatewayHostPort(cfg deployConfig) string {
	output, err := getOutput("docker", "port", gatewayContainerName(cfg.AppName), fmt.Sprintf("%s/tcp", cfg.GatewayInternalPort))
	if err != nil || output == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(output), ":")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[len(parts)-1])
}

// status 显示与应用相关的容器、网络、卷的运行状态。
func status() {
	env, _ := parseArgs()
	envFile := ""
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}
	appName := getAppName(envFile)
	fmt.Printf("=== Container Status ===\n\n")
	fmt.Println("Running containers:")
	mustRun("docker", "ps", "--filter", fmt.Sprintf("name=%s", appName), "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}")
	fmt.Println("\nNetworks:")
	mustRun("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", appName))
	fmt.Println("\nVolumes:")
	mustRun("docker", "volume", "ls", "--filter", fmt.Sprintf("name=%s", appName))
}

// getAppName 从 env 文件中获取应用名称，未指定环境时依次尝试 local/test/production。
func getAppName(envFile string) string {
	if envFile != "" {
		return getEnvVar(envFile, "APP_NAME", defaultAppName)
	}
	for _, file := range []string{".env.local", ".env.test", ".env.production"} {
		if appName := getEnvVar(file, "APP_NAME", ""); appName != "" {
			return appName
		}
	}
	return defaultAppName
}

// cleanupOldImages 按创建时间降序排列本地镜像，保留最新的 N 个，删除其余。
func cleanupOldImages(cfg deployConfig) {
	if cfg.KeepImages <= 0 {
		return
	}
	fmt.Printf("\n[cleanup] Removing old images (keeping latest %d)...\n", cfg.KeepImages)

	reference := fmt.Sprintf("%s/%s", cfg.Registry, cfg.ImageName)
	output, err := getOutput("docker", "images", "--filter", fmt.Sprintf("reference=%s:*", reference), "--format", "{{.ID}}|{{.Repository}}:{{.Tag}}|{{.CreatedAt}}", "--no-trunc")
	if err != nil {
		fmt.Printf("[cleanup] Warning: failed to list images: %v\n", err)
		return
	}
	if output == "" {
		fmt.Println("[cleanup] No images found")
		return
	}

	lines := strings.Split(output, "\n")
	sort.Slice(lines, func(i, j int) bool {
		pi := strings.Split(lines[i], "|")
		pj := strings.Split(lines[j], "|")
		if len(pi) < 3 || len(pj) < 3 {
			return false
		}
		return pi[2] > pj[2]
	})
	if len(lines) <= cfg.KeepImages {
		fmt.Printf("[cleanup] Found %d images, no cleanup needed\n", len(lines))
		return
	}

	deleted := 0
	seen := map[string]bool{}
	for _, line := range lines[cfg.KeepImages:] {
		parts := strings.Split(line, "|")
		if len(parts) < 2 || seen[parts[0]] {
			continue
		}
		seen[parts[0]] = true
		if err := runCmd("docker", "rmi", parts[0]); err == nil {
			fmt.Printf("[cleanup] Removed: %s\n", parts[1])
			deleted++
		}
	}
	fmt.Printf("[cleanup] Cleanup complete: removed %d old images\n", deleted)
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

// fatalUsage 输出错误信息和用法后退出。
func fatalUsage(message string) {
	fmt.Println(message)
	printUsage()
	os.Exit(1)
}

// fatalf 格式化输出错误信息后以状态码 1 退出。
func fatalf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	fmt.Print(message)
	os.Exit(1)
}
