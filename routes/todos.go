package routes

import (
	"fmt"
	"net/http"
	"project_todo/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func getTodoById(context *gin.Context) {
	todoId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todo id"})
		return
	}
	todo, err := models.GetTodoById(todoId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todo"})
		return
	}
	context.JSON(http.StatusOK, todo)
}
