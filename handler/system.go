package handler

import (
	"net/http"
	"os"
	"time"

	"server_gin/runtime"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func HealthReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"color":     getenv("APP_COLOR", "unknown"),
		"version":   getenv("APP_VERSION", "unknown"),
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
	})
}

func HealthDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":         "ok",
		"color":          getenv("APP_COLOR", "unknown"),
		"pid":            os.Getpid(),
		"uptime":         int(time.Since(startTime).Seconds()),
		"timestamp":      time.Now().Format("2006/01/02 15:04:05"),
		"draining":       runtime.IsTrafficShift(),
		"rejecting":      runtime.IsRejecting(),
		"activeRequests": runtime.GetActiveRequests(),
	})
}

func HealthLb(c *gin.Context) {
	if runtime.IsTrafficShift() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "draining"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func TrafficShift(c *gin.Context) {
	runtime.StartTrafficShift()
	c.JSON(http.StatusOK, gin.H{"ok": true, "state": "traffic-shift"})
}

func RejectNew(c *gin.Context) {
	runtime.StartRejectNew()
	c.JSON(http.StatusOK, gin.H{"ok": true, "state": "reject-new-requests"})
}

func ResumeTraffic(c *gin.Context) {
	runtime.Resume()
	c.JSON(http.StatusOK, gin.H{"ok": true, "state": "resume-traffic"})
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
