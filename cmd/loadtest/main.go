package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Stats struct {
	total      atomic.Int64
	failed     atomic.Int64
	instances  sync.Map
	errorTypes sync.Map
}

type HealthResponse struct {
	Status    string `json:"status"`
	Color     string `json:"color"`
	PID       int    `json:"pid"`
	Uptime    int    `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

func main() {
	concurrency := flag.Int("c", 10, "并发数")
	duration := flag.Int("d", 20, "测试时间（秒）")
	url := flag.String("url", "http://localhost:17001/health", "请求URL")
	flag.Parse()

	fmt.Printf("压测配置: 并发=%d, 时长=%ds, URL=%s\n\n", *concurrency, *duration, *url)

	stats := &Stats{}
	done := make(chan struct{})
	var wg sync.WaitGroup

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(*duration) * time.Second)

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 100,
					IdleConnTimeout:     90 * time.Second,
				},
			}
			
			for time.Now().Before(endTime) {
				stats.total.Add(1)
				
				resp, err := client.Get(*url)
				if err != nil {
					stats.failed.Add(1)
					errType := fmt.Sprintf("请求失败: %v", err)
					val, _ := stats.errorTypes.LoadOrStore(errType, &atomic.Int64{})
					val.(*atomic.Int64).Add(1)
					continue
				}
				
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				
				if err != nil {
					stats.failed.Add(1)
					errType := fmt.Sprintf("读取响应失败: %v", err)
					val, _ := stats.errorTypes.LoadOrStore(errType, &atomic.Int64{})
					val.(*atomic.Int64).Add(1)
					continue
				}
				
				if resp.StatusCode != http.StatusOK {
					stats.failed.Add(1)
					errType := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body[:min(100, len(body))]))
					val, _ := stats.errorTypes.LoadOrStore(errType, &atomic.Int64{})
					val.(*atomic.Int64).Add(1)
					continue
				}
				
				var health HealthResponse
				if err := json.Unmarshal(body, &health); err != nil {
					stats.failed.Add(1)
					errType := fmt.Sprintf("解析JSON失败 (前100字节: %s)", string(body[:min(100, len(body))]))
					val, _ := stats.errorTypes.LoadOrStore(errType, &atomic.Int64{})
					val.(*atomic.Int64).Add(1)
					continue
				}
				
				instanceID := fmt.Sprintf("%s (PID-%d)", health.Color, health.PID)
				val, _ := stats.instances.LoadOrStore(instanceID, &atomic.Int64{})
				val.(*atomic.Int64).Add(1)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	elapsed := time.Since(startTime)

	fmt.Printf("\n=== 压测结果 ===\n")
	fmt.Printf("总请求数: %d\n", stats.total.Load())
	fmt.Printf("失败次数: %d\n", stats.failed.Load())
	fmt.Printf("成功次数: %d\n", stats.total.Load()-stats.failed.Load())
	fmt.Printf("耗时: %.2fs\n", elapsed.Seconds())
	fmt.Printf("QPS: %.2f\n\n", float64(stats.total.Load())/elapsed.Seconds())

	fmt.Println("=== 实例分布 ===")
	stats.instances.Range(func(key, value interface{}) bool {
		count := value.(*atomic.Int64).Load()
		fmt.Printf("实例 [%s]: %d 次 (%.2f%%)\n", 
			key.(string), 
			count, 
			float64(count)/float64(stats.total.Load()-stats.failed.Load())*100)
		return true
	})

	if stats.failed.Load() > 0 {
		fmt.Println("\n=== 失败原因统计 ===")
		stats.errorTypes.Range(func(key, value interface{}) bool {
			count := value.(*atomic.Int64).Load()
			fmt.Printf("%s: %d 次\n", key.(string), count)
			return true
		})
	}
}