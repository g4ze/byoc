package core

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// CreateCluster creates a cluster if clustername doesnt already exists.
// ClusterName is the unique username of the user
func CreateCluster(svc *ecs.Client, clusterName string) error {
	// check if cluster already exists
	if CheckUserCluster(clusterName, *svc) {
		return fmt.Errorf("cluster already exists")
	}
	respCluster, err := svc.CreateCluster(context.TODO(), &ecs.CreateClusterInput{
		CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		ClusterName:       &clusterName,
		DefaultCapacityProviderStrategy: []types.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String("FARGATE"),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("unable to create cluster, %v", err)
	}
	log.Println(respCluster.Cluster.ClusterName, " created")
	return nil
}

func DeleteCluster(svc *ecs.Client, clusterName string) error {

	_, err := svc.DeleteCluster(context.TODO(), &ecs.DeleteClusterInput{
		Cluster: &clusterName,
	})
	if err != nil {
		return fmt.Errorf("unable to delete cluster, %v", err)
	}
	log.Println("Cluster deleted")
	return nil
}

// this also needs to work on cluster arn and not cluster name
func CheckUserCluster(ClusterName string, svc ecs.Client) bool {
	// check if user has a cluster
	resp, err := svc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("unable to list clusters, %v", err)
	}
	log.Println(resp.ClusterArns, len(resp.ClusterArns))
	for _, cluster := range resp.ClusterArns {
		// cluster = arn:aws:ecs:ap-south-1:272197635538:cluster/g4ze
		cluster = strings.Split(cluster, ":")[5]
		cluster = strings.Split(cluster, "/")[1]
		// cluster = g4ze
		log.Println(cluster)
		if cluster == ClusterName {
			return true
		}
	}
	return false
}
func ClusterStatus(svc *ecs.Client, clusterName string) (string, error) {
	resp, err := svc.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: []string{clusterName},
	})
	if err != nil {
		return "", fmt.Errorf("unable to describe cluster, %v", err)
	}
	log.Printf("Cluster status: %v", *resp.Clusters[0].Status)
	return *resp.Clusters[0].Status, nil
}
