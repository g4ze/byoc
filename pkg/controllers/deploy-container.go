package controllers

import (
	"github.com/g4ze/byoc/pkg/core"
)

func Deploy_container(UserName string, Image string, Port int, Environment map[string]string) {
	core.CreateCluster(UserName)
	core.CreateTaskDefinition(UserName, Image, Port, Environment)
	// CreateTaskDefinition(UserName, Image, Port, Environment)
}
