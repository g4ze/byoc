package handlers

import (
	"encoding/json"
	"log"

	"github.com/g4ze/byoc/pkg/database"

	"github.com/g4ze/byoc/pkg/controllers"
	"github.com/g4ze/byoc/pkg/types"
	"github.com/gin-gonic/gin"
)

func Delete_Container(c *gin.Context) {
	var request types.DeleteContainer
	err := json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Request: %+v", request)
	if request.Image == "" {
		c.JSON(400, gin.H{"error": "Image name is required"})
		return
	}

	service, err := database.GetService(request)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(service) == 0 {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}
	// if len(service) != 1 {
	// 	c.JSON(404, gin.H{"error": "More than one service with same name received. Please contact support"})
	// 	return
	// }
	// delete container

	err = controllers.DeleteContainerDeployment(&types.Service{
		Name:            service[0].Name,
		Arn:             service[0].Arn,
		TaskFamily:      service[0].TaskFamily,
		LoadBalancerARN: service[0].LoadBalancerARN,
		TargetGroupARN:  service[0].TargetGroupARN,
		LoadbalancerDNS: service[0].LoadbalancerDNS,
		DesiredCount:    int32(service[0].DesiredCount),
		Cluster:         service[0].Cluster,
		Image:           service[0].Image,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Container deleted successfully"})

}
