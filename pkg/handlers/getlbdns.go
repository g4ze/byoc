package handlers

import (
	"log"

	"github.com/g4ze/byoc/pkg/database"
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
	dns, err := database.GetLB_DNS(slug)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"dns": dns, "slug": slug})

}
