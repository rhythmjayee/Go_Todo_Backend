package main

import (
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var todos map[int] string
var counter int

type Todo struct {
	Id int
	Text   string `json:"text" binding:"required"`
}
type TodoList struct {
	Todos   []Todo 
}

func testAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is working",
	})
}

func getTodos(c *gin.Context) {
	len := len(todos)
	fmt.Println(len)
	list := TodoList{}
	(list).Todos = make([]Todo, len)
	count := 0
	for k,v := range todos {
		(list).Todos[count] = Todo{
			Id: k,
			Text: v,
		}
		count = count + 1
	}
	res, err := json.Marshal(list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(res)
	c.JSON(200, gin.H{
		"Todos": list.Todos,
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
	route.GET("/todos", getTodos)
	route.POST("/add", addTodo)
	route.Run() // listen and serve on 0.0.0.0:8080
}