package controllers

import (
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/g4ze/byoc/pkg/core"
)

func Deploy_container(UserName string, Image string, Ports int32, Environment map[string]string) {
	core.CreateCluster(UserName)
	// KeyValuePair
	Environment2 := func() []types.KeyValuePair {
		var Environment2 []types.KeyValuePair
		for key, value := range Environment {
			Environment2 = append(Environment2, types.KeyValuePair{
				Name:  &key,
				Value: &value,
			})
		}
		return Environment2
	}()
	core.CreateTaskDefinition(UserName, Image, Ports, Environment2)
	core.CreateService(UserName, Image, Ports, Environment2)
	// CreateTaskDefinition(UserName, Image, Port, Environment)
}
