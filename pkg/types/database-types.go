package types

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Service struct {
	Name            string `json:"name"`
	Arn             string `json:"arn"`
	TaskFamily      string `json:"taskFamily"`
	LoadBalancerARN string `json:"loadBalancerARN"`
	TargetGroupARN  string `json:"targetGroupARN"`
	LoadbalancerDNS string `json:"loadbalancerDNS"`
	DesiredCount    int32  `json:"desiredCount"`
	Cluster         string `json:"cluster"`
	Image           string `json:"image"`
	Slug            string `json:"slug"`
}
