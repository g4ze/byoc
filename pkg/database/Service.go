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
		db.Service.Cluster.Set(Service.Cluster),
		db.Service.Image.Set(Service.Image),
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
func GetService(request types.DeleteContainer) ([]db.ServiceModel, error) {
	// get service
	log.Printf("Getting service %s", request.Image)
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	resp, err := client.Service.FindMany(
		db.Service.Image.Equals(request.Image),
		db.Service.UserName.Equals(request.UserName),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("resp: %+v", resp)
	return resp, nil
}
func DeleteService(request types.DeleteContainer) error {
	// delete service
	log.Printf("Deleting service %s", request.Image)
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
	_, err := client.Service.FindMany(
		db.Service.Image.Equals(request.Image),
		db.Service.UserName.Equals(request.UserName),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}
	log.Printf("Service %s deleted", request.Image)
	return nil
}
