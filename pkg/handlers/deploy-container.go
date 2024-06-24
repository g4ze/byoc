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
	// get request from client
	// deploy container
	// return response to client

	// check if request payload matches the required payload
	var reqPayload DeployContainerPayload

	err := json.NewDecoder(c.Request.Body).Decode(&reqPayload)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	// extracting user name from the auth add on param
	reqPayload.UserName = c.Params.ByName("user")
	log.Printf("Received request to deploy container: %v", reqPayload)

	// validate the request payload
	if reqPayload.Image == "" || reqPayload.UserName == "" || (reqPayload.Port) == 0 {
		http.Error(c.Writer, "Invalid request payload", http.StatusBadRequest)
		return
	}
	Int32Port := int32(reqPayload.Port)
	resp, err := controllers.Deploy_container(reqPayload.UserName, reqPayload.Image, Int32Port, reqPayload.Environment)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Service already exists, deployed with new image"})
		return
	}
	log.Printf("Inserting service into database: %v", resp.Name)
	err = database.InsertService(resp, reqPayload.UserName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, resp)
}
