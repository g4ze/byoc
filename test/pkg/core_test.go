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
	"github.com/g4ze/byoc/pkg/core"
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
	lbarn, _, err := core.CreateLoadBalancer(elbSvc, "test")
	if err != nil {
		log.Fatalf("Error creating load balancer: %v", err)
	}
	log.Printf("Deleting load balancer")
	log.Println("Tesing single delete")
	core.DeleteLoadBalancerARN(elbSvc, lbarn)
	log.Printf("Checking for deletion ALL of load balancer")
	core.CreateLoadBalancer(elbSvc, "test1")
	core.CreateLoadBalancer(elbSvc, "test2")
	core.CreateLoadBalancer(elbSvc, "test3")
	log.Println("Tesing ALL delete")
	core.DeleteAllLoadBalancers(elbSvc)
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
	err = core.CreateTaskDefinition(svc, "test", "test", 80, nil)
	if err != nil {
		log.Fatal(err)
	}
	core.CreateTaskDefinition(svc, "test", "test", 80, nil)
	core.DeleteTaskDefinition(svc, "test", "test")
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
	err = core.CreateCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error creating cluster: %v", err)
	}
	err = core.DeleteCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error deleting cluster: %v", err)
	}
	status, err := core.ClusterStatus(svc, "test")
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
	err = core.CreateCluster(svc, "test")
	if err != nil {
		log.Fatalf("Error creating cluster: %v", err)
	}

	core.CreateTaskDefinition(svc, "test", img, 80, nil)
	log.Printf("Creating service")
	service, err := core.CreateService(svc, elbSvc, "test", img, int32(80), []types.KeyValuePair{{Name: aws.String("test"), Value: aws.String("test")}}, "test")
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}
	time.Sleep(10 * time.Second)
	log.Printf("Deleting service")
	err = core.DeleteService(elbSvc, svc, service)
	if err != nil {
		log.Fatalf("Error deleting service: %v", err)
	}
	time.Sleep(10 * time.Second)
	core.DeleteCluster(svc, "test")
}
