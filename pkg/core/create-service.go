package core

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateService(svc *ecs.Client, elbSvc *elbv2.ELBV2, UserName string, Image string, Port int32, Environment []types.KeyValuePair) *string {

	containerName := generateName(UserName, Image, "container")
	serviceName := generateName(UserName, Image, "service")
	taskName := generateName(UserName, Image, "task")
	isService, serviceStatus, err := serviceExists(svc, serviceName, UserName)
	if err != nil {
		log.Fatalf("Error checking if service exists: %v", err)
	}
	if isService {
		log.Printf("Service %s already exists. Skipping creation.", serviceName)
		if serviceStatus == "ACTIVE" {
			_, err := svc.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
				Service:            &serviceName,
				Cluster:            &UserName,
				ForceNewDeployment: true,
				DesiredCount:       aws.Int32(2),
			})
			if err != nil {
				log.Fatalf("Unable to update service: %v", err)
			}
		}

		return aws.String("OK")
	}
	// this needs a change, we cant keep creating load balancers
	// based on the image name, we need to create a unique name
	loadBalancerArn, lbdns, err := CreateLoadBalancer(elbSvc, Image)
	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}
	targetGroupArn, err := CreateTargetGroup(elbSvc, Image)
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
func serviceExists(svc *ecs.Client, serviceName string, clusterName string) (bool, any, error) {
	input := &ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Cluster:  aws.String(clusterName),
	}

	result, err := svc.DescribeServices(context.TODO(), input)
	if err != nil {
		log.Printf("Error describing service: %v", err)
		return false, nil, err
	}

	if len(result.Services) == 0 {
		return false, nil, nil
	}
	log.Printf("service is %v", *result.Services[0].ServiceName)

	return true, *result.Services[0].Status, nil
}
