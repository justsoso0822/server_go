package router

import "github.com/gin-gonic/gin"

// Handle 在路由组上同时注册 GET 和 POST
func Handle(g *gin.RouterGroup, path string, handlers ...gin.HandlerFunc) {
	g.GET(path, handlers...)
	g.POST(path, handlers...)
}
