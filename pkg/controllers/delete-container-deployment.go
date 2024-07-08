package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/g4ze/byoc/pkg/core"
	"github.com/g4ze/byoc/pkg/database"
	"github.com/g4ze/byoc/pkg/types"
)

// deletes the container deployed
func DeleteContainerDeployment(service *types.Service) error {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	// Create ECS client
	svc := ecs.NewFromConfig(cfg)
	sess := session.Must(session.NewSession())
	elbSvc := elbv2.New(sess)
	// #1-remove service from AWS
	err = core.DeleteService(elbSvc, svc, service)
	if err != nil {
		return fmt.Errorf("error deleting task definition: %v", err)
	}
	log.Println("Service deletion successful")
	// #2-remove cluster from AWS if no services
	err = core.DeleteCluster(svc, service.Cluster)
	if err != nil {
		log.Println("Can't delete cluster")
	}
	// #3-remove service from database
	err = database.DeleteService(types.DeleteContainer{
		Image:    service.Image,
		UserName: service.Cluster,
	})
	if err != nil {
		return fmt.Errorf("error deleting service from database: %v", err)
	}

	return nil
}
