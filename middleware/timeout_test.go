package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRequestTimeoutSetsDeadline(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestTimeout(50 * time.Millisecond))
	r.GET("/check", func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		if !ok {
			t.Fatal("请求 context 没有 deadline")
		}
		if wait := time.Until(deadline); wait > 60*time.Millisecond {
			t.Fatalf("deadline 距离现在 %v，期望不超过 50ms 左右", wait)
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/check", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("状态码 = %d", rec.Code)
	}
}

func TestRequestTimeoutExpires(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var ctxErr error
	r := gin.New()
	r.Use(RequestTimeout(20 * time.Millisecond))
	r.GET("/slow", func(c *gin.Context) {
		<-c.Request.Context().Done()
		ctxErr = c.Request.Context().Err()
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/slow", nil))
	if !errors.Is(ctxErr, context.DeadlineExceeded) {
		t.Fatalf("ctx.Err() = %v，期望 context.DeadlineExceeded", ctxErr)
	}
}

func TestRequestTimeoutFallsBackToDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestTimeout(0))
	r.GET("/check", func(c *gin.Context) {
		deadline, ok := c.Request.Context().Deadline()
		if !ok {
			t.Fatal("请求 context 没有 deadline")
		}
		if wait := time.Until(deadline); wait < time.Second {
			t.Fatalf("fallback 后 deadline 只有 %v，期望接近 2s", wait)
		}
		c.Status(http.StatusOK)
	})
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/check", nil))
}
