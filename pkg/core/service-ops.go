package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/service/elbv2"
	byocTypes "github.com/g4ze/byoc/pkg/types"
)

func CreateService(svc *ecs.Client, elbSvc *elbv2.ELBV2, UserName string, Image string, Port int32, Environment []types.KeyValuePair, DeploymentName string) (*byocTypes.Service, error) {
	var desiredCount int32 = 2
	containerName := generateName(UserName, Image, "container")
	serviceName := strings.ReplaceAll(DeploymentName, " ", "-")
	Family := generateName(UserName, Image, "task")
	isService, serviceStatus, err := ServiceExists(svc, serviceName, UserName)
	if err != nil {
		return nil, fmt.Errorf("error checking if service esists %+v", err)
	}
	if isService {

		log.Printf("Service %s already exists. Skipping creation.", serviceName)
		if serviceStatus == "ACTIVE" {
			log.Printf("Service is active, updating task definition")
			err = CreateTaskDefinition(svc, UserName, Image, Port, Environment)
			if err != nil {
				return nil, fmt.Errorf("error creating task definition: %v", err)
			}
			log.Printf("Redeploying service %s", serviceName)
			_, err := svc.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
				Service:            &serviceName,
				Cluster:            &UserName,
				ForceNewDeployment: true,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to update service: %+v", err)

			}
			return nil, nil
		} else {
			log.Printf("Service is %v, deleting service", serviceStatus)
			err := DeleteService(elbSvc, svc, &byocTypes.Service{
				Name:    serviceName,
				Cluster: UserName,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to delete service: %+v", err)
			}

		}
	}
	log.Printf("Creating service %s", serviceName)
	// this needs a change, we cant keep creating load balancers
	// based on the service/deploymeny name, we need to create a unique name
	loadBalancerArn, lbdns, err := CreateLoadBalancer(elbSvc, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %+v", err)
	}
	targetGroupArn, err := CreateTargetGroup(elbSvc, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create target group: %v", err)
	}
	EventListenerARN, err := CreateListener(elbSvc, loadBalancerArn, targetGroupArn)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %v", err)
	}
	// Generate names

	serviceInput := &ecs.CreateServiceInput{
		ServiceName:  &serviceName,
		Cluster:      &UserName,
		DesiredCount: &desiredCount,
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
					*aws.String(os.Getenv("SUBNET1")),
					*aws.String(os.Getenv("SUBNET2")),
					*aws.String(os.Getenv("SUBNET3")),
					// "subnet-03f664a0d4fe40293",
					// "subnet-01850c7c6f49dfb7f",
					// "subnet-0450df2a14564e3d5",
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
		// automatically takes latest rev
		TaskDefinition: &Family,
	}

	resp, err := svc.CreateService(context.TODO(), serviceInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	log.Printf("Service created with status: %v", *resp.Service.Status)
	return &byocTypes.Service{
		Name:             serviceName,
		Arn:              *resp.Service.ServiceArn,
		TaskFamily:       Family,
		LoadBalancerARN:  *loadBalancerArn,
		TargetGroupARN:   *targetGroupArn,
		LoadbalancerDNS:  *lbdns,
		DesiredCount:     desiredCount,
		Cluster:          UserName,
		Image:            Image,
		EventListenerARN: EventListenerARN,
		DeploymentName:   DeploymentName,
	}, nil
}

func ServiceExists(svc *ecs.Client, serviceName string, clusterName string) (bool, string /*string for status*/, error) {
	input := &ecs.DescribeServicesInput{
		Services: []string{serviceName},
		Cluster:  aws.String(clusterName),
	}

	result, err := svc.DescribeServices(context.TODO(), input)
	if err != nil {
		log.Printf("Error describing service: %v", err)
		return false, "", err
	}

	if len(result.Services) == 0 {
		return false, "", nil
	}
	log.Printf("service is %v", *result.Services[0].ServiceName)

	return true, *result.Services[0].Status, nil
}
func DeleteService(elbSvc *elbv2.ELBV2, svc *ecs.Client, service *byocTypes.Service) error {
	_, status, err := ServiceExists(svc, service.Name, service.Cluster)
	if err != nil {
		return err
	}
	if status == "ACTIVE" {
		err := UpdateServiceToZeroCount(svc, service.Name, service.Cluster)
		if err != nil {
			return err
		}
	} // else is inactive/ draining
	_, err = svc.DeleteService(context.TODO(), &ecs.DeleteServiceInput{
		Service: &service.Name,
		Cluster: &service.Cluster,
		Force:   aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("unable to delete service: %v", err)
	}
	log.Printf("Service %v deleted from AWS", service.Name)
	if service.LoadBalancerARN != "" {
		err = DeleteLoadBalancerARN(elbSvc, &service.LoadBalancerARN)
		if err != nil {
			return err
		}
	}
	if service.EventListenerARN != "" {
		err = DeleteListener(elbSvc, &service.EventListenerARN)
		if err != nil {
			return err
		}
	}
	time.Sleep(10 * time.Second)
	if service.TaskFamily != "" {
		err = DeleteTaskDefination(svc, service.Cluster, service.Name)
		if err != nil {
			return err
		}
	}
	time.Sleep(10 * time.Second)
	if service.TargetGroupARN != "" {
		err = DeleteTargetGroup(elbSvc, &service.TargetGroupARN)
		if err != nil {
			return err
		}
	}
	log.Printf("Service %v's dependencies deleted from AWS", service.Name)
	return nil

}
func UpdateServiceToZeroCount(svc *ecs.Client, serviceName string, clusterName string) error {
	_, err := svc.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Service:            &serviceName,
		Cluster:            &clusterName,
		ForceNewDeployment: true,
		DesiredCount:       aws.Int32(0),
	})
	if err != nil {
		return fmt.Errorf("unable to update service to 0: %v", err)
	} else {
		log.Printf("Service updated to 0")
	}
	return nil
}
