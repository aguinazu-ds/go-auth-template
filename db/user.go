package db

import (
	"context"
	"crypto/sha256"
	"go-auth-template/types"
)

func CreateUser(user *types.User) error {
	// Create a new user in the database.
	_, err := Bun.NewInsert().Model(user).Returning("ID").Exec(context.Background(), user)
	// print the id of the new user
	return err
}

func UpdateUser(user *types.User) error {
	_, err := Bun.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(context.Background())
	return err
}

func GetUserByToken(token string, scope string) (types.User, error) {
	tokenHash := sha256.Sum256([]byte(token))

	user := new(types.User)
	err := Bun.NewSelect().
		Model(user).
		Join("join tokens as t on t.user_id = id").
		Where("t.hash = ?", tokenHash[:]).
		Where("t.scope = ?", scope).
		Where("t.expires_at > now()").
		Scan(context.Background())
	return *user, err
}
