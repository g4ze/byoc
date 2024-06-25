package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Get_LBDNS(c *gin.Context) {
	// get request from client
	// get lb dns
	// return response to client
	slug := c.Request.URL.Query().Get("slug")
	log.Printf("Received request to get lb dns: %v", slug)
	if slug == "" {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	// get lb dns
	GetLB_DNS(slug)

}
