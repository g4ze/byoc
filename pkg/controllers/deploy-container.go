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
	"github.com/g4ze/byoc/pkg/database"
	byocTypes "github.com/g4ze/byoc/pkg/types"
	"github.com/joho/godotenv"
)

func Deploy_container(UserName string, Image string, Port int32, Environment map[string]string) (*byocTypes.Service, error) {

	// KeyValuePair
	Environment2 := func() []types.KeyValuePair {
		var Environment2 []types.KeyValuePair
		for key, value := range Environment {
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
	err = core.CreateCluster(svc, UserName)
	if err != nil {
		return nil, err
	}

	err = core.CreateTaskDefinition(svc, UserName, Image, Port, Environment2)
	if err != nil {
		return nil, err
	}

	service, err := core.CreateService(svc, elbSvc, UserName, Image, Port, Environment2)
	if err != nil {
		return nil, err
	}

	err = database.InsertService(service, UserName)
	if err != nil {
		return nil, err
	}
	return service, nil

}
