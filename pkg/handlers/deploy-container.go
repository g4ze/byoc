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
		User        string            `json:"user"`
		Port        int32             `json:"port"`
		Environment map[string]string `json:"environment"`
	}
	// check if request payload matches the required payload
	var reqPayload payload
	log.Printf("Request body: %v", c.Request.Body)
	err := json.NewDecoder(c.Request.Body).Decode(&reqPayload)
	if err != nil {
		http.Error(c.Writer, err.Error()+"invalid boady", http.StatusBadRequest)
		return
	}
	log.Printf("Received request to deploy container: %v", reqPayload)

	// validate the request payload
	if reqPayload.Image == "" || reqPayload.User == "" || (reqPayload.Port) == 0 {
		http.Error(c.Writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := controllers.Deploy_container(reqPayload.User, reqPayload.Image, reqPayload.Port, reqPayload.Environment)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Service already exists, deployed with new image"})
		return
	}
	log.Printf("Inserting service into database: %v", resp.Name)
	err = database.InsertService(resp, reqPayload.User)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, resp)
}
