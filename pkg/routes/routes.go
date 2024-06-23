package routes

import (
	"github.com/g4ze/byoc/pkg/handlers"
	"github.com/g4ze/byoc/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	r := gin.Default()

	// Enable CORS
	r.Use(cors.Default())

	authRoutes := r.Group("/v1")
	authRoutes.Use(middleware.JwtMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	// 208 = user already exists
	// 200 = fed
	// 429 = too many requests
	r.POST("/create-user", handlers.Create_User)
	r.POST("/login", handlers.Login)
	authRoutes.Use(cors.Default())
	// Create a new group for routes that require JWT middleware
	authRoutes.POST("/whoami", handlers.WhoAMI)
	authRoutes.POST("/make-cluster", handlers.Make_Cluster)
	authRoutes.DELETE("/delete-cluster", handlers.Delete_Cluster)
	authRoutes.POST("/deploy-container", handlers.Deploy_Container)
	authRoutes.DELETE("/delete-container", handlers.Delete_Container)
	authRoutes.GET("/get-services", handlers.Get_Services)

	r.Run(":2001") // listen and serve on 0.0.0.0:2001
}
