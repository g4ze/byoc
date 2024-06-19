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
	"github.com/joho/godotenv"
)

func DeleteContainerDeployment(service *types.Service) error {
	// delete container
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	// Create ECS client
	svc := ecs.NewFromConfig(cfg)
	sess := session.Must(session.NewSession())
	elbSvc := elbv2.New(sess)
	err = core.DeleteService(elbSvc, svc, service)
	if err != nil {
		return fmt.Errorf("error deleting task definition: %v", err)
	}
	log.Printf("Service %s deleted", service.Name)
	err = core.DeleteCluster(svc, service.Cluster)
	if err != nil {
		log.Println("Can't delete cluster")
	}
	database.DeleteService(types.DeleteContainer{
		Image:    service.Image,
		UserName: service.Cluster,
	})
	return nil
}
