package core

import (
	"context"
	"log"
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
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	containerName := UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-container"
	serviceName := UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-service"
	loadBalancerName := UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-container" + string(Port)

	serviceInput := &ecs.CreateServiceInput{
		ServiceName:  &serviceName,
		Cluster:      &UserName,
		DesiredCount: aws.Int32(1),
		LaunchType:   types.LaunchTypeFargate,
		LoadBalancers: []types.LoadBalancer{
			{
				ContainerName:    &containerName,
				ContainerPort:    &Port,
				LoadBalancerName: &loadBalancerName,
				TargetGroupArn: func() *string {

					targetGroupName := UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-container-tg"
					createInput := elbv2.CreateTargetGroupInput{
						Name:                &targetGroupName,
						Protocol:            aws.String("HTTP"),
						HealthCheckProtocol: aws.String("HTTP"),
						// Matcher: &elbv2.Matcher{
						// 	HttpCode: aws.String("200"),
						// },
						TargetType: aws.String("ip"),
					}
					createResp, err := elbv2.New(session.Must(session.NewSession())).CreateTargetGroup(&createInput)

					if err != nil {
						log.Fatalf("Couldn't create new TG for LB")
					}
					return createResp.TargetGroups[0].TargetGroupArn
				}(),
			},
		},
	}

	resp, err := svc.CreateService(context.TODO(), serviceInput)
	if err != nil {
		log.Printf("Service created %v", resp.Service.Status)
	}
}
