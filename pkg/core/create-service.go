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
	isService, serviceStatus, err := ServiceExists(svc, serviceName, UserName)
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
	err = CreateListener(elbSvc, loadBalancerArn, targetGroupArn)
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

func ServiceExists(svc *ecs.Client, serviceName string, clusterName string) (bool, any /*string for status*/, error) {
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
func DeleteService(elbSvc *elbv2.ELBV2, svc *ecs.Client, serviceName string, clusterName string, Image string) {
	if _, status, _ := ServiceExists(svc, serviceName, clusterName); status == "ACTIVE" {
		UpdateServiceToZeroCount(svc, serviceName, clusterName)
	}
	_, err := svc.DeleteService(context.TODO(), &ecs.DeleteServiceInput{
		Service: &serviceName,
		Cluster: &clusterName,
	})
	if err != nil {
		log.Fatalf("Unable to delete service: %v", err)
	}
	log.Printf("Service deleted")

	DeleteLoadBalancer(elbSvc, Image)
	DeleteTaskDefination(svc, clusterName, Image)
	// DeleteTargetGroup(_)
}
func UpdateServiceToZeroCount(svc *ecs.Client, serviceName string, clusterName string) {
	_, err := svc.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Service:            &serviceName,
		Cluster:            &clusterName,
		ForceNewDeployment: true,
		DesiredCount:       aws.Int32(0),
	})
	if err != nil {
		log.Fatalf("Unable to update service to 0: %v", err)
	} else {
		log.Printf("Service updated to 0")
	}
}
