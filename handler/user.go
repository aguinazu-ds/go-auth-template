package handler

import (
	"net/http"

	"go-auth-template/view/user"
)

func HandleUserAccount(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, user.Account())
}
