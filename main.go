package main

import "github.com/gin-gonic/gin"

func testAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is working",
	})
}

func main() {
	route := gin.Default()
	route.GET("/test", testAPI)
	route.Run() // listen and serve on 0.0.0.0:8080
}