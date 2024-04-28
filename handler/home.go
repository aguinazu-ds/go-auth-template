package handler

import (
	"net/http"

	"go-auth-template/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, home.Index())
}
