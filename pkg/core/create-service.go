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

func CreateService(UserName string, Image string, Port int32, Environment []types.KeyValuePair) *string {
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
	containerName := generateName(UserName, Image, "container")
	serviceName := generateName(UserName, Image, "service")
	// loadBalancerName := generateName(UserName, "Image", "lb")
	taskName := generateName(UserName, Image, "task")
	if serviceExists(svc, serviceName, UserName) {
		log.Printf("Service %s already exists. Skipping creation.", serviceName)
		_, err := svc.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
			Service:            &serviceName,
			Cluster:            &UserName,
			ForceNewDeployment: true,
			DesiredCount:       aws.Int32(2),
		})
		if err != nil {
			log.Fatalf("Unable to update service: %v", err)
		}
		return aws.String("OK")
	}

	loadBalancerArn, lbdns, err := createLoadBalancer(elbSvc, UserName, Image)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	targetGroupArn, err := createTargetGroup(UserName, Image)
	if err != nil {
		log.Fatalf("Failed to create target group: %v", err)
	}
	err = createListener(elbSvc, loadBalancerArn, targetGroupArn)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	// Generate names

	serviceInput := &ecs.CreateServiceInput{
		ServiceName:  &serviceName,
		Cluster:      &UserName,
		DesiredCount: aws.Int32(2),
		LaunchType:   types.LaunchTypeFargate,
		LoadBalancers: []types.LoadBalancer{
			{
				TargetGroupArn: targetGroupArn,
				ContainerName:  &containerName,
				ContainerPort:  &Port,
			},
		},
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets: []string{
					// os.Getenv("SUBNET1"),
					// os.Getenv("SUBNET2"),
					// os.Getenv("SUBNET3"),
					"subnet-03f664a0d4fe40293",
					"subnet-01850c7c6f49dfb7f",
					"subnet-0450df2a14564e3d5",
				},
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{"sg-0d1526e6316cb2abf"},
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
	return lbdns
}

func generateName(UserName string, Image, suffix string) string {
	return UserName + "-" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + suffix
}
func createTargetGroup(_ string, Image string) (*string, error) {
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
	log.Printf("TG CREATED")
	return createResp.TargetGroups[0].TargetGroupArn, nil
}
func createLoadBalancer(elbSvc *elbv2.ELBV2, _ string, Image string) (*string, *string, error) {
	loadBalancerName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + "lb"

	subnets := []*string{
		aws.String(os.Getenv("SUBNET1")),
		aws.String(os.Getenv("SUBNET2")),
		aws.String(os.Getenv("SUBNET3")),
	}

	createInput := &elbv2.CreateLoadBalancerInput{
		Name:           aws.String(loadBalancerName),
		Subnets:        subnets,
		SecurityGroups: []*string{aws.String("sg-0a9330b2203aac089")},
		Scheme:         aws.String("internet-facing"),
		Type:           aws.String("application"),
	}

	createResp, err := elbSvc.CreateLoadBalancer(createInput)
	if err != nil {
		log.Fatal(err)
	}

	if len(createResp.LoadBalancers) == 0 {
		log.Fatal("No LB created")
	}
	log.Println("LB CREATED")
	lbdns := *createResp.LoadBalancers[0].DNSName
	log.Printf("LB public DNS %v", lbdns)
	return createResp.LoadBalancers[0].LoadBalancerArn, &lbdns, nil
}

func createListener(elbSvc *elbv2.ELBV2, loadBalancerArn *string, targetGroupArn *string) error {
	listenerInput := &elbv2.CreateListenerInput{
		LoadBalancerArn: loadBalancerArn,
		Protocol:        aws.String("HTTP"),
		Port:            aws.Int64(80),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: targetGroupArn,
						},
					},
				},
			},
		},
	}

	_, err := elbSvc.CreateListener(listenerInput)
	return err
}
func serviceExists(svc *ecs.Client, serviceName string, clusterName string) bool {
	input := &ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Cluster:  aws.String(clusterName),
	}

	result, err := svc.DescribeServices(context.TODO(), input)
	if err != nil {
		log.Printf("Error describing service: %v", err)
		return false
	}

	if len(result.Services) == 0 {
		return false
	}
	log.Printf("service is %v", *result.Services[0].ServiceName)

	return *result.Services[0].Status != "INACTIVE"
}
