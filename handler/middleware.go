package handler

import (
	"context"
	"go-auth-template/pkg/authsession"
	"go-auth-template/types"
	"net/http"
	"strconv"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		store := authsession.GetStore()
		session, err := store.Get(r, "session")
		if err != nil {
			println(err)
			next.ServeHTTP(w, r)
			return
		}
		if session.Values["user"] == nil {
			println("User is not authenticated")
			next.ServeHTTP(w, r)
			return
		}
		user := session.Values["user"].(*types.AuthenticatedUser)
		println("email: " + user.Email)
		println("logged in: " + strconv.FormatBool(user.LoggedIn))
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
