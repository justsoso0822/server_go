package main

import (
	"bufio"
	"fmt"
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
	case "rollback":
		rollback()
	case "status":
		status()
	case "cleanup":
		cleanup()
	case "start-local":
		startLocal()
	case "stop-local":
		stopLocal()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: go run cmd/deploy/main.go <command> [args]

Commands:
  build <env>              构建镜像 (local|test|production)
  push <env>               推送镜像 (local|test|production)
  deploy <env> <color>     蓝绿部署 (env: local|test|production, color: blue|green)
  rollback <env>           回滚 (local|test|production)
  status <env>             查看状态 (local|test|production)
  cleanup <env>            清理环境 (local|test|production)
  start-local              启动本地数据库
  stop-local               停止本地数据库

Examples:
  go run cmd/deploy/main.go build local
  go run cmd/deploy/main.go deploy local blue
  go run cmd/deploy/main.go rollback local`)
}

func getArg(index int, defaultVal string) string {
	if len(os.Args) > index {
		return os.Args[index]
	}
	return defaultVal
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

func getTag() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "latest"
	}
	tag := strings.TrimSpace(string(output))

	cmd = exec.Command("git", "status", "--porcelain")
	output, _ = cmd.Output()
	if len(output) > 0 {
		tag += ".dirty"
	}
	return tag
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	env := getArg(2, "local")
	tag := getTag()
	registry := getRegistry(env)
	image := fmt.Sprintf("%s/server-go:%s", registry, tag)
	imageLatest := fmt.Sprintf("%s/server-go:latest", registry)

	fmt.Printf("Building for environment: %s with tag: %s\n", env, tag)
	fmt.Printf("Building image: %s\n", image)

	if err := runCmd("docker", "build", "-t", image, "-t", imageLatest, "-f", "docker/Dockerfile", "."); err != nil {
		fmt.Println("Build failed")
		os.Exit(1)
	}
	fmt.Printf("Build completed: %s\n", image)
}

func push() {
	env := getArg(2, "local")
	tag := getTag()
	registry := getRegistry(env)
	image := fmt.Sprintf("%s/server-go:%s", registry, tag)
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
	fmt.Printf("Push completed: %s\n", image)
}

func deploy() {
	env := getArg(2, "local")
	color := getArg(3, "")
	if color == "" || (color != "blue" && color != "green") {
		fmt.Println("Color must be 'blue' or 'green'")
		os.Exit(1)
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("Environment file not found: %s\n", envFile)
		os.Exit(1)
	}

	fmt.Printf("Deploying %s environment for %s\n", color, env)
	appName := getEnvVar(envFile, "APP_NAME", "server-go")

	// 启动 traefik
	if !strings.Contains(getOutput("docker", "ps", "--format", "{{.Names}}"), appName+"-gateway") {
		fmt.Println("Starting Traefik gateway...")
		runCmd("docker", "compose", "-f", "docker/compose/traefik.yml", "--env-file", envFile, "up", "-d")
		time.Sleep(5 * time.Second)
	}

	// 部署新颜色
	fmt.Printf("Starting %s service...\n", color)
	composeFile := fmt.Sprintf("docker/compose/%s.yml", color)
	if err := runCmd("docker", "compose", "-f", composeFile, "--env-file", envFile, "up", "-d", "--build"); err != nil {
		fmt.Println("Deployment failed")
		os.Exit(1)
	}

	// 等待健康检查
	fmt.Printf("Waiting for %s service to be healthy...\n", color)
	maxWait := 60
	for i := 0; i < maxWait; i += 3 {
		output := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=%s-%s", appName, color), "--filter", "health=healthy", "--format", "{{.Names}}")
		if strings.Contains(output, color) {
			fmt.Printf("%s service is healthy\n", color)
			break
		}
		if i+3 >= maxWait {
			fmt.Printf("ERROR: %s service failed to become healthy\n", color)
			os.Exit(1)
		}
		fmt.Printf("Waiting... (%d/%d seconds)\n", i, maxWait)
		time.Sleep(3 * time.Second)
	}

	// 停止旧颜色
	oldColor := "green"
	if color == "green" {
		oldColor = "blue"
	}
	if strings.Contains(getOutput("docker", "ps", "--format", "{{.Names}}"), appName+"-"+oldColor) {
		fmt.Printf("Stopping old %s service...\n", oldColor)
		oldComposeFile := fmt.Sprintf("docker/compose/%s.yml", oldColor)
		runCmd("docker", "compose", "-f", oldComposeFile, "--env-file", envFile, "down")
	}

	fmt.Println("Deployment completed successfully!")
	fmt.Printf("Active service: %s\n", color)
	hostPort := getEnvVar(envFile, "HOST_GATEWAY_PORT", "7001")
	dashboardPort := getEnvVar(envFile, "TRAEFIK_DASHBOARD_PORT", "18080")
	fmt.Printf("Gateway: http://localhost:%s\n", hostPort)
	fmt.Printf("Traefik Dashboard: http://localhost:%s/dashboard/\n", dashboardPort)
}

func rollback() {
	env := getArg(2, "local")
	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("Environment file not found: %s\n", envFile)
		os.Exit(1)
	}

	appName := getEnvVar(envFile, "APP_NAME", "server-go")
	output := getOutput("docker", "ps", "--format", "{{.Names}}")

	var currentColor, targetColor string
	if strings.Contains(output, appName+"-blue") {
		currentColor = "blue"
		targetColor = "green"
	} else if strings.Contains(output, appName+"-green") {
		currentColor = "green"
		targetColor = "blue"
	} else {
		fmt.Println("No active deployment found")
		os.Exit(1)
	}

	fmt.Printf("Current active: %s\n", currentColor)
	fmt.Printf("Rolling back to: %s\n", targetColor)

	// 启动目标颜色
	fmt.Printf("Starting %s service...\n", targetColor)
	targetFile := fmt.Sprintf("docker/compose/%s.yml", targetColor)
	runCmd("docker", "compose", "-f", targetFile, "--env-file", envFile, "up", "-d")

	// 等待健康检查
	fmt.Printf("Waiting for %s service to be healthy...\n", targetColor)
	maxWait := 60
	for i := 0; i < maxWait; i += 3 {
		output := getOutput("docker", "ps", "--filter", fmt.Sprintf("name=%s-%s", appName, targetColor), "--filter", "health=healthy", "--format", "{{.Names}}")
		if strings.Contains(output, targetColor) {
			fmt.Printf("%s service is healthy\n", targetColor)
			break
		}
		if i+3 >= maxWait {
			fmt.Printf("ERROR: %s service failed to become healthy\n", targetColor)
			os.Exit(1)
		}
		fmt.Printf("Waiting... (%d/%d seconds)\n", i, maxWait)
		time.Sleep(3 * time.Second)
	}

	// 停止当前颜色
	fmt.Printf("Stopping %s service...\n", currentColor)
	currentFile := fmt.Sprintf("docker/compose/%s.yml", currentColor)
	runCmd("docker", "compose", "-f", currentFile, "--env-file", envFile, "down")

	fmt.Println("Rollback completed successfully!")
	fmt.Printf("Active service: %s\n", targetColor)
}

func status() {
	env := getArg(2, "local")
	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("Environment file not found: %s\n", envFile)
		os.Exit(1)
	}

	appName := getEnvVar(envFile, "APP_NAME", "server-go")
	fmt.Printf("=== Environment: %s ===\n\n", env)
	fmt.Println("Running containers:")
	runCmd("docker", "ps", "--filter", fmt.Sprintf("name=%s", appName), "--format", "table {{.Names}}\t{{.Status}}\t{{.Ports}}")
	fmt.Println("\nNetworks:")
	runCmd("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", appName))
	fmt.Println("\nVolumes:")
	runCmd("docker", "volume", "ls", "--filter", fmt.Sprintf("name=%s", appName))
}

func cleanup() {
	env := getArg(2, "local")
	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		fmt.Printf("Environment file not found: %s\n", envFile)
		os.Exit(1)
	}

	fmt.Printf("Stopping all services for environment: %s\n", env)
	runCmd("docker", "compose", "-f", "docker/compose/blue.yml", "--env-file", envFile, "down")
	runCmd("docker", "compose", "-f", "docker/compose/green.yml", "--env-file", envFile, "down")
	runCmd("docker", "compose", "-f", "docker/compose/traefik.yml", "--env-file", envFile, "down")

	if env == "local" {
		runCmd("docker", "compose", "-f", "docker/compose/local.yml", "--env-file", envFile, "down")
	}

	fmt.Println("Cleanup completed")
}

func startLocal() {
	fmt.Println("Starting local development environment...")
	if err := runCmd("docker", "compose", "-f", "docker/compose/local.yml", "--env-file", ".env.local", "up", "-d"); err != nil {
		fmt.Println("Failed to start local services")
		os.Exit(1)
	}

	fmt.Println("Waiting for services to be ready...")
	time.Sleep(10 * time.Second)

	fmt.Println("Local services started:")
	fmt.Println("  MySQL: 127.0.0.1:330")
	fmt.Println("  Redis: 127.0.0.1:637")
	fmt.Println("")
	fmt.Println("You can now run the application with:")
	fmt.Println("  GF_GCFG_FILE=config.local.yaml go run main.go")
}

func stopLocal() {
	fmt.Println("Stopping local development environment...")
	runCmd("docker", "compose", "-f", "docker/compose/local.yml", "--env-file", ".env.local", "down")
	fmt.Println("Local services stopped")
}