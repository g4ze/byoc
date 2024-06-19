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
	lbarn, _, err := CreateLoadBalancer(elbSvc, "test")
	if err != nil {
		log.Fatalf("Error creating load balancer: %v", err)
	}
	log.Printf("Deleting load balancer")
	log.Println("Tesing single delete")
	DeleteLoadBalancerARN(elbSvc, lbarn)
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
	err = CreateTaskDefinition(svc, "test", "test", 80, nil)
	if err != nil {
		log.Fatal(err)
	}
	// CreateTaskDefinition(svc, "test", "test", 80, nil)
	// DeleteTaskDefination(svc, "test", "test")
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
	err = CreateCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error creating cluster: %v", err)
	}
	err = DeleteCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error deleting cluster: %v", err)
	}
	status, err := ClusterStatus(svc, "test")
	if err != nil {
		log.Fatalf("Error getting cluster status: %v", err)
	}
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
	img := "test"
	err = CreateCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error creating cluster: %v", err)
	}

	CreateTaskDefinition(svc, "test", img, 80, nil)
	log.Printf("Creating service")
	service, err := CreateService(svc, elbSvc, "test", img, int32(80), []types.KeyValuePair{{Name: aws.String("test"), Value: aws.String("test")}})
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}
	time.Sleep(10 * time.Second)
	log.Printf("Deleting service")
	DeleteService(elbSvc, svc, service)
	time.Sleep(10 * time.Second)
	DeleteCluster(svc, "test")
}
