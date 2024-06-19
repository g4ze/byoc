package handlers

import (
	"encoding/json"

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
	service, err := database.GetService(request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// delete container\
	controllers.DeleteContainerDeployment(&types.Service{
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

}
