package routes

import (
	"net/http"
	"project_todo/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getAllTodos(context *gin.Context) {
	userId := context.GetInt64("userId")
	todos, err := models.GetAllTodos(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todos"})
		return
	}
	context.JSON(http.StatusOK, todos)
}

func createTodo(context *gin.Context) {
	var todo models.Todo
	err := context.ShouldBindJSON(&todo)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse request data"})
		return
	}
	todo.UserID = context.GetInt64("userId")
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.IsActive = true

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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse todo id"})
		return
	}
	userId := context.GetInt64("userId")
	todo, err := models.GetTodoById(todoId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todo"})
		return
	}
	if todo.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized access"})
		return
	}
	context.JSON(http.StatusOK, todo)
}

func updateTodoById(context *gin.Context) {
	todoId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse todo id"})
		return
	}

	userId := context.GetInt64("userId")
	todo, err := models.GetTodoById(todoId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch todo to update"})
		return
	}
	if todo.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized update"})
		return
	}

	var modifiedTodo models.Todo
	err = context.ShouldBindJSON(&modifiedTodo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse request data"})
		return
	}
	modifiedTodo.ID = todoId
	modifiedTodo.UpdatedAt = time.Now()
	err = modifiedTodo.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to update todo"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func deleteTodoById(context *gin.Context) {
	todoId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse todo id"})
		return
	}
	userId := context.GetInt64("userId")
	todo, err := models.GetTodoById(todoId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch the todo"})
		return
	}
	if todo.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized delete"})
		return
	}
	err = todo.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete the todo"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
