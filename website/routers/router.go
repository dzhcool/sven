package routers

import (
	"html/template"
	"website/api"

	"github.com/dzhcool/sven/setting"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 加载自定义函数
	r.SetFuncMap(template.FuncMap{
		"add":    Add,
		"strCut": StrCut,
	})

	web_inf := setting.Config.MustString("app.web_inf", "./")
	r.Static("/static", web_inf+"/static")
	r.LoadHTMLGlob(web_inf + "/views/*.html")

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(200, "")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "uri:/")
	})

	_v1 := r.Group("/api/v1")
	{
		api.V1(_v1)
	}
}
