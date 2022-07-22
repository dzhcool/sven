package api

import (
	v1 "website/api/v1"

	"github.com/gin-gonic/gin"
)

func V1(r *gin.RouterGroup) {
	r.GET("/index", v1.Index)
}
