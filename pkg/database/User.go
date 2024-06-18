package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/g4ze/byoc/pkg/database/db"
)

func CreateUser(UserName string, Email string, Password string) error {
	// create user
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

	if isunique, err := ValidateUniqueUserName(client, ctx, UserName); !isunique && err == nil {
		log.Printf("User already exists %s", UserName)
		return fmt.Errorf("Exists")
	} else if err != nil {
		return err
	}
	user, err := client.User.CreateOne(
		db.User.UserName.Set(UserName),
		db.User.Email.Set(Email),
		db.User.Password.Set(Password),
		db.User.Cluster.Set(UserName),
	).Exec(ctx)
	if err != nil {
		return err
	}
	result, _ := json.MarshalIndent(user, "", "  ")
	log.Printf("created post: %s\n", result)
	log.Printf("post: %+v", user)
	return nil
}
func ValidateUniqueUserName(client *db.PrismaClient, ctx context.Context, UserName string) (bool, error) {
	log.Printf("Validating UserName: %s", UserName)
	resp, err := client.User.FindMany(
		db.User.UserName.Equals(UserName),
	).Exec(ctx)
	if err != nil {
		return false, err
	}
	log.Printf("resp: %+v", resp)
	if len(resp) > 0 {
		return false, nil
	}
	return true, nil
}
func DeleteUser(UserName string) error {
	// delete user
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
	if userexists, err := ValidateUniqueUserName(client, ctx, UserName); !userexists && err == nil {
		return fmt.Errorf("user does not exists")
	} else if err != nil {
		return err
	}
	_, err := client.User.FindUnique(
		db.User.UserName.Equals(UserName),
	).Delete().Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}
