package types

import (
	"time"

	"github.com/google/uuid"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Hash      []byte    `bun:"hash,notnull" json:"hash"`
	UserID    uuid.UUID `bun:"user_id,notnull" json:"user_id"`
	User      *User     `bun:"rel:belongs-to" json:"user"`
	ExpiresAt time.Time `bun:"expires_at,notnull" json:"expires_at"`
	Scope     string    `bun:"scope,notnull" json:"scope"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
}
