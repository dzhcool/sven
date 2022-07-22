package middlewares

import (
	"github.com/gin-gonic/gin"
)

func InitMiddleware(r *gin.Engine) {
	r.Use(Cors()) // 设置跨域header
}
