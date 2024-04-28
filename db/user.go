package db

import (
	"context"
	"go-auth-template/types"
)

func CreateUser(user types.User) error {
	// Create a new user in the database.
	_, err := Bun.NewInsert().Model(&user).Exec(context.Background())
	return err
}
