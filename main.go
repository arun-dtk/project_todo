package main

import (
	"fmt"
	"net/http"
	"project_todo/db"
	"project_todo/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	db.InitDB()
	server := gin.Default()

	server.GET("/todos", getAllTodos)

	server.POST("/todos", createTodo)

	server.Run(":8080") // localhost:8080
}

func getAllTodos(context *gin.Context) {
	todos, err := models.GetAllTodos()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todos"})
		return
	}
	context.JSON(http.StatusOK, todos)
}

func createTodo(context *gin.Context) {
	var todo models.Todo
	err := context.ShouldBindBodyWithJSON(&todo)

	// Debugging: Print the entire todo struct after binding JSON
	fmt.Printf("Parsed todo: %+v\n", todo)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse request data"})
		return
	}
	fmt.Println("Parsed success", todo)

	err = todo.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create todo"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Todo created", "todo": todo})
}
