package handler

import (
	"context"
	"go-auth-template/pkg/authsession"
	"go-auth-template/types"
	"log/slog"
	"net/http"
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
			next.ServeHTTP(w, r)
			return
		}
		user := session.Values["user"].(*types.AuthenticatedUser)
		ctx := context.WithValue(r.Context(), types.UserContextKey, types.AuthenticatedUser{
			Email:    user.Email,
			LoggedIn: user.LoggedIn,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		store := authsession.GetStore()
		session, err := store.Get(r, "session")
		if err != nil {
			slog.Error("Error getting session", err)
			next.ServeHTTP(w, r)
			return
		}
		if session.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
