package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 临时验证：预检请求应在 CORS 中间件内直接响应，普通请求应带 Access-Control-Allow-Origin。
func TestCORSVerify(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("APP_ENV", "local")

	newEngine := func(origins []string) *gin.Engine {
		r := gin.New()
		r.Use(CORS(origins))
		r.GET("/api/user/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 0})
		})
		return r
	}

	// 1. 简单请求：local 环境白名单为空时回显 Origin
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/user/login?sign=1", nil)
	req.Header.Set("Origin", "http://localhost:18090")
	newEngine(nil).ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:18090" {
		t.Fatalf("简单请求 Allow-Origin = %q", got)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("简单请求状态码 = %d", rec.Code)
	}

	// 2. 预检请求：应返回 204 并带允许头，且不进业务处理
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodOptions, "/api/user/login", nil)
	req.Header.Set("Origin", "http://localhost:18090")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type,x-custom")
	newEngine(nil).ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("预检状态码 = %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != "content-type,x-custom" {
		t.Fatalf("预检 Allow-Headers = %q", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("Vary = %q", got)
	}

	// 3. 生产白名单：名单内的源放行，名单外的源不返回 CORS 头
	t.Setenv("APP_ENV", "production")
	engine := newEngine([]string{"https://h5.example.com"})

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/user/login", nil)
	req.Header.Set("Origin", "https://h5.example.com")
	engine.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://h5.example.com" {
		t.Fatalf("白名单内 Allow-Origin = %q", got)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/user/login", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	engine.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("白名单外不应返回 Allow-Origin，实际 = %q", got)
	}
}
