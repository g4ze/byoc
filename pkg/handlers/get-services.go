package handlers

import (
	"log"

	"github.com/g4ze/byoc/pkg/database"
	"github.com/gin-gonic/gin"
)

func Get_Services(c *gin.Context) {
	// get services
	userName := c.Params.ByName("user")
	log.Print("Getting services for user: ", userName)
	services, err := database.GetServices(userName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, services)

}
