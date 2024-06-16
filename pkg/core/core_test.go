package core

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/joho/godotenv"
)

func TestLoadBalancerFunctions(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	sess := session.Must(session.NewSession())
	elbSvc := elbv2.New(sess)
	log.Printf("Creating load balancer")
	CreateLoadBalancer(elbSvc, "test")
	log.Printf("Deleting load balancer")
	DeleteLoadBalancer(elbSvc, "test")
	log.Printf("Checking for deletion of load balancer")
	CreateLoadBalancer(elbSvc, "test1")
	CreateLoadBalancer(elbSvc, "test2")
	CreateLoadBalancer(elbSvc, "test3")
	DeleteAllLoadBalancers(elbSvc)
}
func TestTaskFunctions(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	log.Printf("Creating task definition")
	// CreateTaskDefinition(svc, "test", "test", 80, nil)
	CreateTaskDefinition(svc, "test", "test", 80, nil)
	CreateTaskDefinition(svc, "test", "test", 80, nil)
	DeleteTaskDefination(svc, "test", "test")
}
