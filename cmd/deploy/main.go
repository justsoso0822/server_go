package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

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

func printUsage() {
	fmt.Println(`Usage: go run cmd/deploy/main.go <command> <env> [options]

Commands:
  build <env> [version=xxx]         构建镜像
  push <env> [version=xxx]          推送镜像
  deploy <env> [version=xxx] [-f]   部署（自动检测颜色并切换，失败自动回滚）
  status                            查看容器状态
  start-local-db                    启动本地数据库
  stop-local-db                     停止本地数据库

Environment:
  local       本地环境
  test        测试环境
  production  生产环境

Options:
  version=xxx  指定版本标签
  -f           网关端口冲突时强制替换网关

Examples:
  go run main.go                                         # 本地开发，自动使用 config.yaml
  go run cmd/deploy/main.go start-local-db               # 启动本地数据库
  go run cmd/deploy/main.go build local                  # 构建本地镜像
  go run cmd/deploy/main.go build test version=v1.2.3    # 构建测试镜像并指定版本
  go run cmd/deploy/main.go push local                   # 推送本地镜像，默认 version=1.0.0
  go run cmd/deploy/main.go push test version=v1.2.3     # 推送测试镜像并指定版本
  go run cmd/deploy/main.go deploy test version=v1.2.3   # 部署到测试环境
  go run cmd/deploy/main.go deploy local -f              # 强制替换本地网关后部署
  go run cmd/deploy/main.go status                       # 查看容器状态`)
}

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

func getRegistry(env string) string {
	switch env {
	case "local":
		return "ccr.ccs.tencentyun.com/justsoso-local"
	case "test":
		return "ccr.ccs.tencentyun.com/justsoso-test"
	case "production":
		return "ccr.ccs.tencentyun.com/justsoso-production"
	default:
		fmt.Printf("Unknown environment: %s\n", env)
		os.Exit(1)
		return ""
	}
}

func getBuildVersion(options map[string]string, env string) string {
	if version, ok := options["version"]; ok && version != "" {
		return version
	}

	if env == "local" {
		return "1.0.0"
	}

	fmt.Println("Error: version parameter is required for test/production build/push")
	fmt.Println("Usage: version=v1.2.3")
	os.Exit(1)
	return ""
}

func getDeployVersion(options map[string]string) string {
	if version, ok := options["version"]; ok && version != "" {
		return version
	}

	return "latest"
}

func getGitVersion() string {
	version := strings.TrimSpace(getOutput("git", "rev-parse", "--short", "HEAD"))
	if version == "" {
		return ""
	}

	dirty := strings.TrimSpace(getOutput("git", "status", "--porcelain")) != ""
	if dirty {
		version += ".dirty"
	}
	return version
}

func getAppName(envFile string) string {
	if appName := os.Getenv("APP_NAME"); appName != "" {
		return appName
	}

	for _, file := range []string{envFile, ".env.local", ".env.test", ".env.production"} {
		if appName := getEnvVar(file, "APP_NAME", ""); appName != "" {
			return appName
		}
	}

	return "server-go"
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gatewayNeedsRecreate(appName, desiredHostPort string) bool {
	containerName := fmt.Sprintf("%s-gateway", appName)
	runningContainers := getOutput("docker", "ps", "--format", "{{.Names}}")
	if !strings.Contains(runningContainers, containerName) {
		return true
	}

	currentHostPort := getGatewayHostPort(appName, "")
	if currentHostPort == "" {
		return true
	}

	return desiredHostPort != "" && currentHostPort != desiredHostPort
}

func gatewayExists(appName string) bool {
	containerName := fmt.Sprintf("%s-gateway", appName)
	runningContainers := getOutput("docker", "ps", "--format", "{{.Names}}")
	return strings.Contains(runningContainers, containerName)
}

func forceGatewayReplace(options map[string]string) bool {
	return options["force"] == "true"
}

func getGatewayHostPort(appName, defaultPort string) string {
	containerName := fmt.Sprintf("%s-gateway", appName)
	output := getOutput("docker", "port", containerName, "7001/tcp")
	if output == "" {
		return defaultPort
	}

	parts := strings.Split(strings.TrimSpace(output), ":")
	if len(parts) == 0 {
		return defaultPort
	}

	port := strings.TrimSpace(parts[len(parts)-1])
	if port == "" {
		return defaultPort
	}
	return port
}

func getOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output))
}

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
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}

func getEnvVar(envFile, key, defaultVal string) string {
	env := loadEnv(envFile)
	if val, ok := env[key]; ok {
		return val
	}
	return defaultVal
}

func build() {
	env, options := parseArgs()
	if env == "" {
		fmt.Println("Error: environment required")
		printUsage()
		os.Exit(1)
	}

	// local 默认 1.0.0；test/production 必须显式指定 version。
	version := getBuildVersion(options, env)
	registry := getRegistry(env)
	image := fmt.Sprintf("%s/server-go:%s", registry, version)
	imageLatest := fmt.Sprintf("%s/server-go:latest", registry)

	fmt.Printf("Building for environment: %s with version: %s\n", env, version)
	fmt.Printf("Building image: %s\n", image)

	if err := runCmd("docker", "build", "-t", image, "-t", imageLatest, "-f", "manifest/docker/Dockerfile", "."); err != nil {
		fmt.Println("Build failed")
		os.Exit(1)
	}
	fmt.Printf("Build completed: %s\n", image)
}

func push() {
	env, options := parseArgs()
	if env == "" {
		fmt.Println("Error: environment required")
		printUsage()
		os.Exit(1)
	}

	version := getBuildVersion(options, env)
	registry := getRegistry(env)
	image := fmt.Sprintf("%s/server-go:%s", registry, version)
	imageLatest := fmt.Sprintf("%s/server-go:latest", registry)

	fmt.Printf("Pushing image: %s\n", image)
	if err := runCmd("docker", "push", image); err != nil {
		fmt.Println("Push failed")
		os.Exit(1)
	}
	if err := runCmd("docker", "push", imageLatest); err != nil {
		fmt.Println("Push failed")
		os.Exit(1)
	}
	fmt.Printf("Push completed: %s and latest\n", image)
}

func deploy() {
	env, options := parseArgs()
	if env == "" {
		fmt.Println("Error: environment required")
		printUsage()
		os.Exit(1)
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("Environment file not found: %s\n", envFile)
		os.Exit(1)
	}

	appName := getEnvVar(envFile, "APP_NAME", "server-go")
	imageSource := getEnvVar(envFile, "IMAGE_SOURCE", "remote")

	version := getDeployVersion(options)
	if env == "local" && imageSource == "local" {
		version = getBuildVersion(options, env)
	}

	// 检测当前运行的颜色
	currentColor := ""
	targetColor := ""
	blueRunning := strings.Contains(getOutput("docker", "ps", "--format", "{{.Names}}"), appName+"-blue")
	greenRunning := strings.Contains(getOutput("docker", "ps", "--format", "{{.Names}}"), appName+"-green")

	if blueRunning {
		currentColor = "blue"
		targetColor = "green"
	} else if greenRunning {
		currentColor = "green"
		targetColor = "blue"
	} else {
		// 首次部署，默认部署到 blue
		targetColor = "blue"
		fmt.Println("No active deployment found, deploying to blue")
	}

	if currentColor != "" {
		fmt.Printf("Current active: %s, deploying to: %s\n", currentColor, targetColor)
	}

	// 启动或校验 traefik 网关
	fmt.Println("[release] [1/8] ensure traefik gateway")
	desiredGatewayHostPort := getEnvVar(envFile, "HOST_GATEWAY_PORT", "7001")
	gatewayRunning := gatewayExists(appName)
	currentGatewayHostPort := getGatewayHostPort(appName, "")
	forceReplaceGateway := forceGatewayReplace(options)

	switch {
	case !gatewayRunning:
		if err := runCmd("docker", "compose", "-f", "manifest/docker/compose/traefik.yml", "--env-file", envFile, "up", "-d"); err != nil {
			fmt.Println("Failed to start traefik gateway")
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	case currentGatewayHostPort == desiredGatewayHostPort:
		fmt.Printf("[release] gateway already aligned on host port %s\n", desiredGatewayHostPort)
	case forceReplaceGateway:
		fmt.Printf("[release] gateway host port mismatch: current=%s desired=%s, force replacing gateway\n", currentGatewayHostPort, desiredGatewayHostPort)
		if err := runCmd("docker", "compose", "-f", "manifest/docker/compose/traefik.yml", "--env-file", envFile, "up", "-d", "--force-recreate"); err != nil {
			fmt.Println("Failed to force replace traefik gateway")
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	default:
		fmt.Printf("ERROR: gateway host port mismatch: current=%s desired=%s\n", currentGatewayHostPort, desiredGatewayHostPort)
		fmt.Println("Refusing to replace gateway automatically. Re-run with -f to force replace the gateway.")
		os.Exit(1)
	}

	// 本地环境且镜像来源是 local，在网关校验通过后再本地构建
	if env == "local" && imageSource == "local" {
		fmt.Println("Local environment detected, building image after gateway check...")
		registry := getRegistry(env)
		image := fmt.Sprintf("%s/server-go:%s", registry, version)
		imageLatest := fmt.Sprintf("%s/server-go:latest", registry)

		if err := runCmd("docker", "build", "-t", image, "-t", imageLatest, "-f", "manifest/docker/Dockerfile", "."); err != nil {
			fmt.Println("Build failed")
			os.Exit(1)
		}
		fmt.Printf("Build completed: %s\n", image)
	}

	// 部署新颜色
	fmt.Printf("[release] [3/8] start %s (version=%s)\n", targetColor, version)
	composeFile := fmt.Sprintf("manifest/docker/compose/%s.yml", targetColor)

	// 设置镜像环境变量
	registry := getRegistry(env)
	os.Setenv("APP_IMAGE", fmt.Sprintf("%s/server-go:%s", registry, version))

	composeArgs := []string{"compose", "-f", composeFile, "--env-file", envFile, "up", "-d"}
	if env == "local" && imageSource == "local" {
		composeArgs = append(composeArgs, "--build")
	}
	if err := runCmd("docker", composeArgs...); err != nil {
		fmt.Println("Deployment failed")
		os.Exit(1)
	}

	// 等待健康检查
	fmt.Printf("[release] [4/8] wait for %s to be healthy (timeout=60s)\n", targetColor)
	maxWait := 60
	healthy := false
	for i := 0; i < maxWait; i++ {
		output := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=%s-%s", appName, targetColor), "--filter", "health=healthy", "--format", "{{.Names}}")
		if strings.Contains(output, targetColor) {
			fmt.Printf("[release] %s healthy\n", targetColor)
			healthy = true
			break
		}
		if i+1 >= maxWait {
			break
		}
		fmt.Printf("Waiting... (%d/%d seconds)\n", i, maxWait)
		time.Sleep(1 * time.Second)
	}

	if !healthy {
		fmt.Printf("ERROR: %s service failed to become healthy, rolling back...\n", targetColor)
		// 自动回滚：停止新部署的服务
		runCmd("docker", "compose", "-f", composeFile, "--env-file", envFile, "down")
		fmt.Println("Rollback completed, deployment failed")
		os.Exit(1)
	}

	// 如果有旧容器，执行流量切换
	if currentColor != "" {
		oldContainerName := fmt.Sprintf("%s-%s-1", appName, currentColor)

		// 步骤 5: 调用 traffic-shift
		fmt.Printf("[release] [5/8] http control -> %s: trigger traffic-shift, /health/lb now returns 503\n", currentColor)
		script := `wget -q -O- --timeout=5 --post-data="" http://127.0.0.1:7001/_internal/control/traffic-shift`
		if err := runCmd("docker", "exec", oldContainerName, "sh", "-c", script); err != nil {
			fmt.Printf("[release] Warning: failed to call traffic-shift: %v\n", err)
		}

		// 步骤 6: 确认切流完成
		fmt.Printf("[release] [6/8] confirm gateway routes to %s (9 consecutive, timeout=30s)\n", targetColor)
		gatewayHostPort := getGatewayHostPort(appName, getEnvVar(envFile, "HOST_GATEWAY_PORT", "7001"))
		confirmed := false
		consecutive := 0
		deadline := time.Now().Add(30 * time.Second)

		for time.Now().Before(deadline) {
			healthURL := fmt.Sprintf("http://localhost:%s/health", gatewayHostPort)
			resp, err := http.Get(healthURL)
			if err == nil {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				if strings.Contains(string(body), fmt.Sprintf(`"color":"%s"`, targetColor)) {
					consecutive++
					fmt.Printf("[release] gateway -> %s (%d/9)\n", targetColor, consecutive)
					if consecutive >= 9 {
						fmt.Printf("[release] cutover confirmed: all traffic is on %s, safe to stop %s\n", targetColor, currentColor)
						confirmed = true
						break
					}
				} else {
					if consecutive > 0 {
						fmt.Printf("[release] probe reset (was %d)\n", consecutive)
						consecutive = 0
					}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}

		if !confirmed {
			fmt.Println("[release] Warning: cutover confirmation timeout")
		}

		// 步骤 7: 排水旧容器
		fmt.Printf("[release] [7/8] http control -> %s: reject any remaining new requests\n", currentColor)
		script = `wget -q -O- --timeout=5 --post-data="" http://127.0.0.1:7001/_internal/control/reject-new-requests`
		runCmd("docker", "exec", oldContainerName, "sh", "-c", script)

		fmt.Printf("[release] waiting %s in-flight requests (timeout=15s)\n", currentColor)
		drainDeadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(drainDeadline) {
			output := getOutput("docker", "exec", oldContainerName, "wget", "-q", "-O-", "--timeout=2", "http://127.0.0.1:7001/health/detail")
			if strings.Contains(output, `"activeRequests":0`) {
				fmt.Printf("[release] %s: no in-flight requests\n", currentColor)
				break
			}
			time.Sleep(1 * time.Second)
		}

		// 步骤 8: 停止旧容器
		fmt.Printf("[release] [8/8] %s: remove containers\n", currentColor)
		oldComposeFile := fmt.Sprintf("manifest/docker/compose/%s.yml", currentColor)
		runCmd("docker", "compose", "-f", oldComposeFile, "--env-file", envFile, "down")
	}

	fmt.Printf("\n[release] SUCCESS: %s now served by %s (version=%s)\n", env, targetColor, version)
	gatewayHostPort := getGatewayHostPort(appName, getEnvVar(envFile, "HOST_GATEWAY_PORT", "7001"))
	dashboardPort := getEnvVar(envFile, "TRAEFIK_DASHBOARD_PORT", "18080")
	fmt.Printf("Gateway: http://localhost:%s\n", gatewayHostPort)
	fmt.Printf("Traefik Dashboard: http://localhost:%s/dashboard/\n", dashboardPort)

	// 清理旧镜像，只保留最近 10 个
	cleanupOldImages()
}

func status() {
	env, _ := parseArgs()
	envFile := ""
	if env != "" {
		envFile = fmt.Sprintf(".env.%s", env)
	}
	appName := getAppName(envFile)
	fmt.Printf("=== Container Status ===\n\n")
	fmt.Println("Running containers:")
	runCmd("docker", "ps", "--filter", fmt.Sprintf("name=%s", appName), "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}")
	fmt.Println("\nNetworks:")
	runCmd("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", appName))
	fmt.Println("\nVolumes:")
	runCmd("docker", "volume", "ls", "--filter", fmt.Sprintf("name=%s", appName))
}

func cleanupOldImages() {
	fmt.Println("\n[cleanup] Removing old images (keeping latest 10)...")

	output := getOutput("docker", "images", "--filter", "reference=*/server-go", "--format", "{{.ID}}|{{.Repository}}:{{.Tag}}|{{.CreatedAt}}", "--no-trunc")
	if output == "" {
		fmt.Println("[cleanup] No images found")
		return
	}

	lines := strings.Split(output, "\n")
	if len(lines) <= 10 {
		fmt.Printf("[cleanup] Found %d images, no cleanup needed\n", len(lines))
		return
	}

	toDelete := lines[10:]
	deleted := 0
	for _, line := range toDelete {
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			imageID := parts[0]
			imageName := parts[1]
			if err := runCmd("docker", "rmi", imageID); err == nil {
				fmt.Printf("[cleanup] Removed: %s\n", imageName)
				deleted++
			}
		}
	}
	fmt.Printf("[cleanup] Cleanup complete: removed %d old images\n", deleted)
}

func startLocalDB() {
	fmt.Println("Starting local database services...")
	if err := runCmd("docker", "compose", "-f", "manifest/docker/compose/local.yml", "--env-file", ".env.local", "up", "-d"); err != nil {
		fmt.Println("Failed to start local database services")
		os.Exit(1)
	}

	fmt.Println("Local database services started:")
	fmt.Println("  MySQL: 127.0.0.1:330")
	fmt.Println("  Redis: 127.0.0.1:637")
	fmt.Println("")
	fmt.Println("You can now run the application with:")
	fmt.Println("  go run main.go")
}

func stopLocalDB() {
	fmt.Println("Stopping local database services...")
	runCmd("docker", "compose", "-f", "manifest/docker/compose/local.yml", "--env-file", ".env.local", "down")
	fmt.Println("Local database services stopped")
}
