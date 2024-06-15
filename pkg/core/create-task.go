package core

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/joho/godotenv"
)

func CreateTaskDefinition(UserName string, Image string, Ports []int32, Environment []types.KeyValuePair) {
	// create task definition
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	type Tasks struct {
		containerName *string
		image         *string
		ports         *[]int32
		environment   *[]types.KeyValuePair
		family        string
		essential     bool
		compatibility types.Compatibility
		cpu           string
		memory        string
	}

	task := func() *Tasks {
		return &Tasks{
			essential: true,
			cpu:       "1024",
			memory:    "2048",
		}
	}()

	containerName := UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-container"
	task.containerName = &containerName
	task.image = &Image
	task.ports = &Ports
	task.environment = &Environment
	task.family = UserName + "-" + strings.Replace(strings.Replace(strings.Replace(Image, "/", "-", -1), ".", "-", -1), ":", "-", -1) + "-task"

	log.Printf("Creating task definition for %v", *task.containerName)

	respTask, err := svc.RegisterTaskDefinition(context.TODO(), &ecs.RegisterTaskDefinitionInput{
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     &task.cpu,
		Memory:                  &task.memory,
		RuntimePlatform: &types.RuntimePlatform{
			CpuArchitecture:       types.CPUArchitectureX8664,
			OperatingSystemFamily: types.OSFamilyLinux,
		},
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Name:      task.containerName,
				Image:     task.image,
				Essential: &task.essential,
				PortMappings: func() []types.PortMapping {
					var portMappings []types.PortMapping
					for _, port := range *task.ports {
						portName := *task.containerName + "-" + string(port)
						portMappings = append(portMappings, types.PortMapping{
							ContainerPort: &port,
							Name:          &portName,
						})
					}
					return portMappings
				}(),
				Environment: *task.environment,
			},
		},
		Family: &task.family,
	})
	if err != nil {
		log.Fatalf("Error creating task definition, %v", err)
	}
	log.Printf("Task definition created: %v", *respTask.TaskDefinition.TaskDefinitionArn)

}
