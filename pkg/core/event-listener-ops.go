package core

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateListener(elbSvc *elbv2.ELBV2, loadBalancerArn *string, targetGroupArn *string) (string, error) {
	listenerInput := &elbv2.CreateListenerInput{
		LoadBalancerArn: loadBalancerArn,
		Protocol:        aws.String("HTTP"),
		Port:            aws.Int64(80),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				ForwardConfig: &elbv2.ForwardActionConfig{
					TargetGroups: []*elbv2.TargetGroupTuple{
						{
							TargetGroupArn: targetGroupArn,
						},
					},
				},
			},
		},
	}

	listener, err := elbSvc.CreateListener(listenerInput)
	return *listener.Listeners[0].ListenerArn, err
}

// start using ars you stupid human
func DeleteListener(elbSvc *elbv2.ELBV2, listenerArn *string) error {
	log.Printf("Deleting listener %s", *listenerArn)
	_, err := elbSvc.DeleteListener(&elbv2.DeleteListenerInput{
		ListenerArn: listenerArn,
	})
	return err
}
