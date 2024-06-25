package controllers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/g4ze/byoc/pkg/core"
	byocTypes "github.com/g4ze/byoc/pkg/types"
	"github.com/joho/godotenv"
)

// Deploy the container
func Deploy_container(newDeployment *byocTypes.DeployContainerPayload) (*byocTypes.Service, error) {

	// KeyValuePair
	Environment2 := func() []types.KeyValuePair {
		var Environment2 []types.KeyValuePair
		for key, value := range newDeployment.Environment {
			Environment2 = append(Environment2, types.KeyValuePair{
				Name:  &key,
				Value: &value,
			})
		}
		return Environment2
	}()
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
	err = core.CreateCluster(svc, newDeployment.UserName)
	if err != nil {
		return nil, err
	}

	err = core.CreateTaskDefinition(svc, newDeployment.UserName, newDeployment.Image, int32(newDeployment.Port), Environment2)
	if err != nil {
		return nil, err
	}

	service, err := core.CreateService(svc, elbSvc, newDeployment.UserName, newDeployment.Image, int32(newDeployment.Port), Environment2, newDeployment.DeploymentName)
	if err != nil {
		return nil, err
	}
	// means the service was updated
	// with new task and deployed
	if service == nil {
		log.Printf("service=nil")
		return nil, nil
	}

	service.Slug = (newDeployment.DeploymentName)

	log.Printf("returning service: %+v", service)
	return service, nil

}
