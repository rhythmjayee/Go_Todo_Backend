package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var todos map[int] string
var counter int

type Todo struct {
	Text   string `json:"text" binding:"required"`
}

func testAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is working",
	})
}

func addTodo(c *gin.Context) {
	var todoJson Todo
	if err := c.ShouldBindJSON(&todoJson); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(todoJson)
	if todoJson.Text != "" {
		var todo = todoJson.Text
		counter++
		todos[counter] = todo
		c.JSON(200, gin.H{
			"id": counter,
			"message": "Todo : {"+ todo +"} has been saved.",
		})
		return
	}
	c.JSON(400, gin.H{
		"error": "Something went wrong",
	}) 
}

func main() {
	route := gin.Default()
	todos = make(map[int] string)
	counter = 0
	route.GET("/test", testAPI)
	route.POST("/add", addTodo)
	route.Run() // listen and serve on 0.0.0.0:8080
}