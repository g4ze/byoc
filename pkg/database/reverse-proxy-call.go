package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/g4ze/byoc/pkg/database/db"
)

func GetLB_DNS(subdomain string) (string, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return "", err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()

	resp, err := client.Service.FindMany(
		db.Service.Slug.Equals(subdomain),
	).Exec(ctx)
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", sql.ErrNoRows
	}
	log.Print("Got response: ", resp[0].LoadbalancerDNS)
	return resp[0].LoadbalancerDNS, nil

}
