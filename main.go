package main

import (
	"fmt"
	"project_todo/db"
	"project_todo/routes"

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

	// Serve static files
	server.Static("/static", "./static")

	// Serve the HTML files based on the URL path
	server.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	server.GET("/app-signup", func(c *gin.Context) {
		c.File("./static/signup.html")
	})

	server.GET("/app-login", func(c *gin.Context) {
		c.File("./static/login.html")
	})

	server.GET("/todolist", func(c *gin.Context) {
		c.File("./static/todos.html")
	})

	routes.RegisterRoutes(server)

	// Serve index.html as the default route
	server.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	server.Run(":8080") // localhost:8080
}
