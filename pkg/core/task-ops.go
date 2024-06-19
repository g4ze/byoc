package core

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// creates a new task defination
// if the task defination with the same faily exists,
// it makes a new revision on it.
func CreateTaskDefinition(svc *ecs.Client, UserName string, Image string, Port int32, Environment []types.KeyValuePair) error {
	// create task definition
	Essential := true
	// cpu and memory values are sensitive to weach other
	Cpu := "1024"
	Memory := "2048"
	containerName := generateName(UserName, Image, "container")
	Family := generateName(UserName, Image, "task")
	log.Printf("Family: %v", Family)
	log.Printf("Creating task definition for %v", containerName)
	portName := containerName + "-" + strconv.Itoa(int(Port))
	// Ensure port mapping name matches the expected pattern
	if !isValidPortMappingName(portName) {
		return fmt.Errorf("invalid port mapping name: %v", portName)
	}
	log.Printf("Port mapping name: %v", portName)
	respTask, err := svc.RegisterTaskDefinition(context.TODO(), &ecs.RegisterTaskDefinitionInput{
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		NetworkMode:             types.NetworkModeAwsvpc,
		Cpu:                     &Cpu,
		Memory:                  &Memory,
		RuntimePlatform: &types.RuntimePlatform{
			CpuArchitecture:       types.CPUArchitectureX8664,
			OperatingSystemFamily: types.OSFamilyLinux,
		},
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Name:      &containerName,
				Image:     &Image,
				Essential: &Essential,
				PortMappings: []types.PortMapping{
					{
						ContainerPort: &Port,
						Name:          &portName,
					},
				},
				Environment: Environment,
			},
		},
		Family: &Family,
	})
	if err != nil {
		return fmt.Errorf("error creating task definition, %v", err)
	}
	log.Printf("Task definition created: %v", *respTask.TaskDefinition.TaskDefinitionArn)
	return nil
}

// Check if the port mapping name is valid
func isValidPortMappingName(name string) bool {
	// Regex pattern for valid port mapping name
	pattern := `^[a-z0-9_][a-z0-9_-]{0,63}$`
	match, _ := regexp.MatchString(pattern, name)
	return match
}

// function deregisters all the revisions of a task defination
// and then deletes the task defination
func DeleteTaskDefination(svc *ecs.Client, UserName string, Image string) error {

	Family := generateName(UserName, Image, "task")
	log.Printf("Deleting task definition for %v", Family)

	revisions, err := FindRevisions(svc, Family)
	if err != nil {
		return fmt.Errorf("error finding task definition revisions, %v", err)
	}

	for _, taskARN := range revisions {
		_, err := svc.DeregisterTaskDefinition(context.TODO(), &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: &taskARN,
		})
		if err != nil {
			return fmt.Errorf("error deregistering task definition, %v", err)
		} else {
			log.Printf("Task definition deregistered: %v", taskARN)
			_, err := svc.DeleteTaskDefinitions(context.TODO(), &ecs.DeleteTaskDefinitionsInput{
				TaskDefinitions: []string{taskARN},
			})
			if err != nil {
				return fmt.Errorf("error deleting task definition, %v", err)
			}
			log.Printf("Task definition deleted: %v", taskARN)

		}

	}

	_, err2 := svc.DeleteTaskDefinitions(context.TODO(), &ecs.DeleteTaskDefinitionsInput{
		TaskDefinitions: []string{Family},
	})
	if err2 != nil {
		return fmt.Errorf("error deleting task definition, %v", err2)
	}
	log.Printf("Task definition deleted: %v", Family)
	return nil
}
func FindRevisions(svc *ecs.Client, Family string) ([]string, error) {
	// find task definition revisions

	log.Printf("Finding task definition revisions for %v", Family)
	resp, err := svc.ListTaskDefinitions(context.TODO(), &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: &Family,
	})
	if err != nil {
		return nil, fmt.Errorf("error finding task definition revisions, %v", err)
	}
	log.Printf("Task definition revisions found: %v", len(resp.TaskDefinitionArns))
	return (resp.TaskDefinitionArns), nil
}
