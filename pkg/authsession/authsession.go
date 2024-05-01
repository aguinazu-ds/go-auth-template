package authsession

import (
	"encoding/gob"
	"go-auth-template/types"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func Init() {
	key := os.Getenv("SECRET_KEY")
	store = sessions.NewCookieStore([]byte(key))
	gob.Register(&types.AuthenticatedUser{})
}

func GetStore() *sessions.CookieStore {
	return store
}
