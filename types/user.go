package types

import (
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}

type User struct {
	ID                uuid.UUID              `bun:"id,pk,notnull,type:uuid,default:uuid_generate_v4()" json:"id"`
	Email             string                 `bun:"email,unique,notnull" json:"email"`
	Activated         bool                   `bun:"activated,notnull,default:false" json:"activated"`
	EncryptedPassword string                 `json:"-"`
	ConfirmationToken string                 `json:"confirmation_token"`
	ConfirmedAt       time.Time              `bun:"confirmed_at" json:"confirmed_at"`
	RecoveryToken     string                 `json:"recovery_token"`
	RecoverySentAt    time.Time              `json:"recovery_sent_at"`
	RawAppMetaData    map[string]interface{} `json:"raw_app_meta_data"`
	RawUserMetaData   map[string]interface{} `bun:"raw_user_meta_data,default:'{}'" json:"raw_user_meta_data"`
	Version           int                    `bun:"version,notnull,default:1" json:"version"`
	CreatedAt         time.Time              `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt         time.Time              `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}
