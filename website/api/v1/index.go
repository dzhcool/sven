package v1

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title": "This is index",
	})
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"ContentType": "json",
		"title":       "This is index",
	})
}
