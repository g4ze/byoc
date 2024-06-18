package core

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateListener(elbSvc *elbv2.ELBV2, loadBalancerArn *string, targetGroupArn *string) error {
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

	_, err := elbSvc.CreateListener(listenerInput)
	return err
}

// start using ars you stupid human
// func deleteListener(elbSvc *elbv2.ELBV2, listenerArn *string) error {
// 	_, err := elbSvc.DeleteListener(&elbv2.DeleteListenerInput{
// 		ListenerArn: listenerArn,
// 	})
// 	return err
// }
