package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/g4ze/byoc/pkg/database"
	"github.com/g4ze/byoc/pkg/middleware"
	"github.com/g4ze/byoc/pkg/types"
	"github.com/gin-gonic/gin"
)

// would return user doesnt exist if err
func Login(c *gin.Context) {
	// create user

	var user types.Login
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	userReceived, err := database.GetUser(user.UserName, user.Password)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if userReceived == nil {
		http.Error(c.Writer, "User does not exist", http.StatusUnauthorized)
		return
	}
	token, err := middleware.GenerateJWT(user.Password, user.UserName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
		"token":   token,
	})
}
