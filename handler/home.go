package handler

import (
	"net/http"

	"go-auth-template/view/home"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, home.Index())
}
