package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.GET("/todos", getAllTodos)
	server.POST("/todos", createTodo)
	server.GET("/todos/:id", getTodoById)
	server.PUT("/todos/:id", updateTodoById)
}
