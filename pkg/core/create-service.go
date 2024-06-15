package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/joho/godotenv"
)

func CreateService(UserName string, Image string, Port int32, Environment []types.KeyValuePair) {
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

	// Generate names
	containerName := generateName(UserName, Image, "container")
	serviceName := generateName(UserName, Image, "service")
	loadBalancerName := generateName(UserName, "Image", "lb")
	taskName := generateName(UserName, Image, "task")
	// Create target group
	// targetGroupArn, err := createTargetGroup(UserName, Image)
	if err != nil {
		log.Fatalf("Failed to create target group: %v", err)
	}

	serviceInput := &ecs.CreateServiceInput{
		ServiceName:  &serviceName,
		Cluster:      &UserName,
		DesiredCount: aws.Int32(0),
		LaunchType:   types.LaunchTypeFargate,
		LoadBalancers: []types.LoadBalancer{
			{

				LoadBalancerName: &loadBalancerName,
				ContainerName:    &containerName,
				ContainerPort:    &Port,
			},
		},
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets: []string{
					os.Getenv("SUBNET1"),
					os.Getenv("SUBNET2"),
					os.Getenv("SUBNET3"),
				},
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{"sg-0a9330b2203aac089"},
			},
		},
		SchedulingStrategy: types.SchedulingStrategyReplica,
		DeploymentConfiguration: &types.DeploymentConfiguration{
			DeploymentCircuitBreaker: &types.DeploymentCircuitBreaker{
				Enable:   true,
				Rollback: true,
			},
		},
		TaskDefinition: &taskName,
	}

	resp, err := svc.CreateService(context.TODO(), serviceInput)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	log.Printf("Service created with status: %v", *resp.Service.Status)
}

func generateName(UserName string, Image, suffix string) string {
	return UserName + "-" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + suffix
}
func createTargetGroup(UserName string, Image string) (*string, error) {
	targetGroupName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + "tg"

	// Create ELBv2 client
	sess := session.Must(session.NewSession())
	elbSvc := elbv2.New(sess)

	// Define target group input
	createInput := &elbv2.CreateTargetGroupInput{
		Name:                aws.String(targetGroupName),
		Protocol:            aws.String("HTTP"),
		Port:                aws.Int64(80),
		VpcId:               aws.String(os.Getenv("ECS_VPC")),
		HealthCheckProtocol: aws.String("HTTP"),
		TargetType:          aws.String("ip"),
	}

	// Create the target group
	createResp, err := elbSvc.CreateTargetGroup(createInput)
	if err != nil {
		return nil, err
	}

	// Check if the target group was created successfully
	if len(createResp.TargetGroups) == 0 {
		return nil, fmt.Errorf("no target groups created")
	}

	return createResp.TargetGroups[0].TargetGroupArn, nil
}
