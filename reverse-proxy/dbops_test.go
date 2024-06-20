package main

import (
	"log"
	"testing"
)

func Test_dbops(T *testing.T) {
	resp, err := GetLB_DNS("determined-hoatzin-of-advertising-test_user")
	if err != nil {
		log.Fatalf("Error in GetLB_DNS: %v", err)
	}
	log.Printf("resp: %+v", resp)

}
