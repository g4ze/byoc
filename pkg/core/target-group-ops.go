package core

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateTargetGroup(elbSvc *elbv2.ELBV2, Image string) (*string, error) {
	targetGroupName := generateName("", Image, "tg")

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

// should be instable
func DeleteTargetGroup(elbSvc *elbv2.ELBV2, targetGroupArn *string) error {
	_, err := elbSvc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: targetGroupArn,
	})
	return err
}
