package routes

import (
	"log"
	"os"

	"github.com/g4ze/byoc/pkg/handlers"
	"github.com/g4ze/byoc/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	r := gin.Default()

	// Enable CORS

	// With this custom configuration:
	config := cors.DefaultConfig()
	FE_URL, exists := os.LookupEnv("FE_URL")
	log.Println("FE_URL: " + FE_URL)
	if !exists {
		FE_URL = "http://localhost:3000"
	}
	config.AllowOrigins = []string{FE_URL} // Replace with your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	authRoutes := r.Group("/v1")
	authRoutes.Use(cors.Default(), middleware.JwtMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	// 208 = user already exists
	// 200 = fed
	// 429 = too many requests
	r.POST("/create-user", handlers.Create_User)
	r.POST("/login", handlers.Login)
	r.GET("/get-lbdns", handlers.Get_LBDNS)
	r.GET("/health", handlers.Health)
	// Create a new group for routes that require JWT middleware
	authRoutes.POST("/whoami", handlers.WhoAMI)
	authRoutes.POST("/make-cluster", handlers.Make_Cluster)
	authRoutes.DELETE("/delete-cluster", handlers.Delete_Cluster)
	authRoutes.POST("/deploy-container", handlers.Deploy_Container)
	authRoutes.DELETE("/delete-container", handlers.Delete_Container)
	authRoutes.GET("/get-services", handlers.Get_Services)

	r.Run(":2001")
}
