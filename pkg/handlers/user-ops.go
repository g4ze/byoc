package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/g4ze/byoc/pkg/database"
	"github.com/g4ze/byoc/pkg/types"
	"github.com/gin-gonic/gin"
)

func Create_User(c *gin.Context) {
	// create user

	var user types.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("user: %+v", user)
	if err := database.CreateUser(user.UserName, user.Email, user.Password); err != nil {
		log.Printf("Error creating user: %s", err.Error())
		http.Error(c.Writer, err.Error(), http.StatusAlreadyReported)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Created",
	})
}
func WhoAMI(c *gin.Context) {
	// get user info
	user := c.Params.ByName("user")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
