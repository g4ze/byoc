package database

import (
	"context"
	"log"

	"github.com/g4ze/byoc/pkg/database/db"
	"github.com/g4ze/byoc/pkg/types"
)

// pushes service params to the database
func InsertService(Service *types.Service, userName string) error {
	// create service
	log.Printf("Creating service %s", Service.Name)
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	_, err := client.Service.CreateOne(
		db.Service.Name.Set(Service.Name),
		db.Service.Arn.Set(Service.Arn),
		db.Service.TaskFamily.Set(Service.TaskFamily),
		db.Service.LoadBalancerARN.Set(Service.LoadBalancerARN),
		db.Service.TargetGroupARN.Set(Service.TargetGroupARN),
		db.Service.LoadbalancerDNS.Set(Service.LoadbalancerDNS),
		db.Service.DesiredCount.Set(int(Service.DesiredCount)),
		db.Service.User.Link(
			db.User.UserName.Equals(userName),
		),
	).Exec(ctx)
	if err != nil {
		return err
	}
	log.Printf("Service %s created", Service.Name)
	return nil
}
