package controllers

func Deploy_container(UserName string, Image string, Port int, Environment map[string]string) {
	CreateCluster(UserName)
}
