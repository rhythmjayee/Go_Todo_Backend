package main

import 
(
	"fmt"
	"time"
	// "strconv"
	// "encoding/json"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	// "github.com/satori/go.uuid"
)

var users map[string] User
var todos map[string] UserTodos

var mySigningKey = []byte("AllYourBase")

type User struct {
	Name string `json:"name" binding:"required"`
	Password string	`json:"password" binding:"required"`
}

type UserTodos struct {
	Todos  map[string] Todo
}

type ID struct {
	Id string `uri:"id" binding:"required"`
}

type Todo struct {
	Id 		string 
	Text   string `json:"text" binding:"required"`
}
type TodoList struct {
	Todos   []Todo 
}

type JWTClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func JWTParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"];
		if len(header) == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header is missing",
			})
			return
		}
		tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
	
		if token.Valid {
			fmt.Println("Valid JWT token found")
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		
	}
}

func getJWTokenString(username string) (string) {
	claims := JWTClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    "TodoApp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	return ss
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	users[user.Name] = user
	todos[user.Name] = UserTodos{
		Todos: make(map[string] Todo),
	}	
	
	token := getJWTokenString(user.Name)
	c.JSON(200, gin.H{
		"response": "User has been registerd with ID: "+user.Name,
		"token": token,
	})
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, ok := users[user.Name]
	if ok == false {
		c.JSON(400, gin.H{
			"response": "User not found",
		})
		return
	}
	if res.Password != user.Password {
		c.JSON(400, gin.H{
			"response": "Wrong Password",
		})
		return
	}
	token := getJWTokenString(user.Name)
	c.JSON(200, gin.H{
		"response": "User has loggedIn",
		"token": token,
	})
}

// func formTodoList() (*TodoList) {
// 	len := len(todos)
// 	fmt.Println(len)
// 	list := TodoList{}
// 	(list).Todos = make([]Todo, len)
// 	count := 0
// 	for _,v := range todos {
// 		(list).Todos[count] = v
// 		count = count + 1
// 	}
// 	return &list
// }

func testAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "API is working",
	})
}

// func getTodos(c *gin.Context) {
// 	list := formTodoList()
// 	c.JSON(200, gin.H{
// 		"Todos": (*list).Todos,
// 	})
// }

// func addTodo(c *gin.Context) {
// 	var todoJson Todo
// 	if err := c.ShouldBindJSON(&todoJson); err != nil {
// 		c.JSON(400, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	fmt.Println(todoJson)
// 	if todoJson.Text != "" {
// 		var todo = todoJson.Text
// 		todoId := uuid.NewV4().String()
// 		todoJson.Id = todoId
// 		todos[todoId] = todoJson
// 		c.JSON(200, gin.H{
// 			"id": todoId,
// 			"message": "Todo : {"+ todo +"} has been saved.",
// 		})
// 		return
// 	}
// 	c.JSON(400, gin.H{
// 		"error": "Something went wrong",
// 	}) 
// }

// func deleteTodo(c *gin.Context) {
// 	var todo ID
// 	if err := c.ShouldBindUri(&todo); err != nil {
// 		c.JSON(400, gin.H{"msg": err})
// 		return
// 	}
// 	delete(todos, todo.Id)
// 	list := formTodoList()
// 	c.JSON(200, gin.H{
// 		"Todos": (*list).Todos,
// 	})
// }

func main() {
	route := gin.Default()
	users = make(map[string] User)
	todos = make(map[string] UserTodos)

	route.POST("/register", register)
	route.POST("/login", login)

	authorized := route.Group("/")
	authorized.Use(JWTParser())
	{
		authorized.GET("/test", testAPI)
	}
	// route.GET("/todos", getTodos)
	// route.POST("/add", addTodo)
	// route.DELETE("/todo/:id", deleteTodo)
	route.Run() // listen and serve on 0.0.0.0:8080
}