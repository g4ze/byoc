package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/joho/godotenv"
)

func Make_cluster(w http.ResponseWriter, r *http.Request) {
	// get request from client
	// create cluster
	// return response to client
	clusterName := r.URL.Query().Get("clusterName")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	respCluster, err := svc.CreateCluster(context.TODO(), &ecs.CreateClusterInput{
		CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		ClusterName:       &clusterName,
	})
	if err != nil {
		log.Printf("unable to create cluster, %v\n", err)
	}
	log.Println(respCluster.Cluster.ClusterName)
	resp, err := svc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("unable to list clusters, %v", err)
	}
	log.Println(resp.ClusterArns)
}
func Delete_cluster(w http.ResponseWriter, r *http.Request) {

	clusterName := r.URL.Query().Get("clusterName")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := ecs.NewFromConfig(cfg)
	resp, err := svc.DeleteCluster(context.TODO(), &ecs.DeleteClusterInput{
		Cluster: &clusterName,
	})
	if err != nil {
		log.Println("can't delete cluster", err)
	} else {
		log.Panicln("Deleted Cluster", resp)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "deleted cluster",
	}
	json.NewEncoder(w).Encode(response)
}
