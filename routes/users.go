package routes

import (
	"net/http"
	"project_todo/models"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request body"})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to save the user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
