package routes

import (
	"fmt"
	"net/http"
	"project_todo/models"
	"project_todo/utils"

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

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request body"})
		return
	}
	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to authenticate the user"})
		return
	}
	jwtToken, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to authenticate the user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
}
