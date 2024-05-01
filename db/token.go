package db

import (
	"context"
	"go-auth-template/types"
	"time"

	"github.com/google/uuid"
)

func CreateActivationToken(token []byte, user *types.User) error {
	activationToken := types.Token{
		Hash:      token,
		UserID:    user.ID,
		User:      user,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		Scope:     types.ScopeActivation,
		CreatedAt: time.Now(),
	}
	// Create a new token in the database.
	_, err := Bun.NewInsert().Model(&activationToken).Exec(context.Background())
	return err
}

func DeleteAllTokensByUserIDAndScope(userID uuid.UUID, scope string) error {
	_, err := Bun.NewDelete().
		Model(&types.Token{}).
		Where("user_id = ?", userID).
		Where("scope = ?", scope).
		Exec(context.Background())
	return err
}
