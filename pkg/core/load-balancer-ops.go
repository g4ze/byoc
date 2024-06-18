package core

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// major issue here is that we are working on the basis of loadbalancernames
// which are not unique, so we need to work with lb arn
// for that we shall implement db soon
func CreateLoadBalancer(elbSvc *elbv2.ELBV2, Image string) (*string, *string, error) {
	loadBalancerName := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + "lb"

	subnets := []*string{
		aws.String(os.Getenv("SUBNET1")),
		aws.String(os.Getenv("SUBNET2")),
		aws.String(os.Getenv("SUBNET3")),
	}

	createInput := &elbv2.CreateLoadBalancerInput{
		Name:           aws.String(loadBalancerName),
		Subnets:        subnets,
		SecurityGroups: []*string{aws.String("sg-0a9330b2203aac089")},
		Scheme:         aws.String("internet-facing"),
		Type:           aws.String("application"),
	}

	createResp, err := elbSvc.CreateLoadBalancer(createInput)
	if err != nil {
		log.Fatal(err)
	}

	if len(createResp.LoadBalancers) == 0 {
		log.Fatal("No LB created")
	}
	log.Println("LB CREATED")
	lbdns := *createResp.LoadBalancers[0].DNSName
	log.Printf("LB public DNS %v", lbdns)
	return createResp.LoadBalancers[0].LoadBalancerArn, &lbdns, nil
}

func DeleteLoadBalancer(elbSvc *elbv2.ELBV2, Image string) {
	// delete load balancer
	loadBalancerName := generateName("", Image, "lb")
	log.Printf("Deleting load balancer %s", loadBalancerName)
	loadBalancerList, err := elbSvc.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{&loadBalancerName},
	})
	if err != nil {
		log.Printf("Unable to describe load balancer: %v", err)
	}
	for _, lb := range loadBalancerList.LoadBalancers {
		log.Printf("LB name %v", *lb.LoadBalancerName)
		if *lb.LoadBalancerName == loadBalancerName {
			log.Printf("Deleting load balancer %s", *lb.LoadBalancerArn)
			_, err := elbSvc.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: lb.LoadBalancerArn,
			})
			if err != nil {
				log.Printf("Unable to delete load balancer: %v", err)
				return
			} else {
				log.Printf("Load balancer %s deleted", loadBalancerName)
			}
		}
	}
}
func DeleteAllLoadBalancers(elbSvc *elbv2.ELBV2) {
	loadBalancerList, err := elbSvc.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		log.Printf("Unable to describe load balancer: %v", err)
	}
	for _, lb := range loadBalancerList.LoadBalancers {
		_, err := elbSvc.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			log.Printf("Unable to delete load balancer: %v", err)
			return
		} else {
			log.Printf("Load balancer %s deleted", *lb.LoadBalancerName)
		}
	}
}
