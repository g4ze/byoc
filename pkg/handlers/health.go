package handlers

import (
	"os"

	"github.com/g4ze/byoc/pkg/database/db"
	"github.com/gin-gonic/gin"
)

// health should check the health of the service
// database activity, tables
// aws keys and connection
// other related env variables
// and return errors accordingly
func Health(c *gin.Context) {

	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ECS_VPC := os.Getenv("ECS_VPC")
	SUBNET1 := os.Getenv("SUBNET1")
	SUBNET2 := os.Getenv("SUBNET2")
	SUBNET3 := os.Getenv("SUBNET3")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if AWS_REGION == "" || AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" || ECS_VPC == "" || SUBNET1 == "" || SUBNET2 == "" || SUBNET3 == "" || JWT_SECRET == "" {
		c.JSON(500, gin.H{"error": "Missing env variables"})
	}

	// databse pinging
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		c.JSON(500, gin.H{"error": "Error connecting to database", "message": err.Error()})
		return
	}
	func() {
		if err := client.Prisma.Disconnect(); err != nil {
			c.JSON(500, gin.H{"error": "Error disconnecting from database", "message": err.Error()})
		}
	}()
	c.JSON(200, gin.H{"message": "Service is healthy"})

}
