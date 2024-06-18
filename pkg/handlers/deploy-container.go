package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/g4ze/byoc/pkg/controllers"
	"github.com/g4ze/byoc/pkg/database"
	"github.com/gin-gonic/gin"
)

func Deploy_Container(c *gin.Context) {
	// define a struct to verify data integrity
	// get request from client
	// deploy container
	// return response to client
	type payload struct {
		Image       string            `json:"image"`
		UserName    string            `json:"userName"`
		Port        int32             `json:"port"`
		Environment map[string]string `json:"environment"`
	}
	// check if request payload matches the required payload
	var reqPayload payload

	err := json.NewDecoder(c.Request.Body).Decode(&reqPayload)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Received request to deploy container: %v", reqPayload)

	// validate the request payload
	if reqPayload.Image == "" || reqPayload.UserName == "" || (reqPayload.Port) == 0 {
		http.Error(c.Writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := controllers.Deploy_container(reqPayload.UserName, reqPayload.Image, reqPayload.Port, reqPayload.Environment)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	database.InsertService(resp, reqPayload.UserName)
	c.JSON(http.StatusOK, resp)
}
