package routes

import (
	"project_todo/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.GET("/todos", getAllTodos)
	authenticated.POST("/todos", createTodo)
	authenticated.GET("/todos/:id", getTodoById)
	authenticated.PUT("/todos/:id", updateTodoById)
	authenticated.DELETE("/todos/:id", deleteTodoById)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
