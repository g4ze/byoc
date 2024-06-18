package core

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
	log.Println("Tesing single delete")
	DeleteLoadBalancer(elbSvc, "test")
	log.Printf("Checking for deletion ALL of load balancer")
	CreateLoadBalancer(elbSvc, "test1")
	CreateLoadBalancer(elbSvc, "test2")
	CreateLoadBalancer(elbSvc, "test3")
	log.Println("Tesing ALL delete")
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

func TestCluster(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	log.Printf("Creating cluster: test")
	CreateCluster(svc, "test")
	DeleteCluster(svc, "test")
	status := ClusterStatus(svc, "test")
	if status != "INACTIVE" {
		t.Errorf("Cluster should be inactive but got %s", status)
	}
}

func TestService(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	sess := session.Must(session.NewSession())
	elbSvc := elbv2.New(sess)
	img := "docker.io/g4ze/cattodb:latest"
	CreateCluster(svc, "test")
	CreateTaskDefinition(svc, "test", img, 80, nil)
	log.Printf("Creating service")
	CreateService(svc, elbSvc, "test", img, int32(80), []types.KeyValuePair{{Name: aws.String("test"), Value: aws.String("test")}})

	time.Sleep(10 * time.Second)
	log.Printf("Deleting service")
	DeleteService(elbSvc, svc, generateName("test", img, "service"), "test", "test")
	time.Sleep(10 * time.Second)
	DeleteCluster(svc, "test")
}
