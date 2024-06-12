package controllers

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/joho/godotenv"
)

// CreateCluster creates a cluster if clustername doesnt already exists.
// ClusterName is the unique username of the user
func CreateCluster(clusterName string) {
	// create cluster
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	if checkUserCluster(clusterName, *svc) {
		log.Println("Cluster already exists")
		return
	}
	respCluster, err := svc.CreateCluster(context.TODO(), &ecs.CreateClusterInput{
		CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		ClusterName:       &clusterName,
	})
	if err != nil {
		log.Printf("unable to create cluster, %v\n", err)
	}
	log.Println(respCluster.Cluster.ClusterName, " created")
}
func checkUserCluster(ClusterName string, svc ecs.Client) bool {
	// check if user has a cluster
	resp, err := svc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("unable to list clusters, %v", err)
	}
	log.Println(resp.ClusterArns)
	for _, cluster := range resp.ClusterArns {
		if cluster == ClusterName {
			return true
		}
	}
	return true
}
