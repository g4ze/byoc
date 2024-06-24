package types

type DeleteContainer struct {
	Image    string `json:"image"`
	UserName string `json:"userName"`
}
type DeployContainerPayload struct {
	DeploymentName string            `json:"deploymentName"`
	Image          string            `json:"image"`
	UserName       string            `json:"userName"`
	Port           int               `json:"port"`
	Environment    map[string]string `json:"environment"`
}
