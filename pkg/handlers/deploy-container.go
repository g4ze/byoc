package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/g4ze/byoc/pkg/controllers"
)

func Deploy_container(w http.ResponseWriter, r *http.Request) {
	// define a struct to verify data integrity
	// get request from client
	// deploy container
	// return response to client
	type payload struct {
		Image       string            `json:"image"`
		UserName    string            `json:"userName"`
		Port        int               `json:"port"`
		Environment map[string]string `json:"environment"`
	}
	// check if request payload matches the required payload
	var reqPayload payload

	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request payload
	if reqPayload.Image == "" || reqPayload.UserName == "" || reqPayload.Port == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	controllers.Deploy_container(reqPayload.UserName, reqPayload.Image, reqPayload.Port, reqPayload.Environment)

	// deploy container
	// ...

	// return response to client
	// ...
}
