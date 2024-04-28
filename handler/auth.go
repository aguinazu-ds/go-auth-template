package handler

import (
	"net/http"

	"go-auth-template/pkg/validator"
	"go-auth-template/view/auth"

	"golang.org/x/crypto/bcrypt"
)

func HandleAuthLogin(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleAuthSignup(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func HandleAuthSignupPost(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	// Validate the form
	v := validator.New()
	v.ValidateEmail(params.Email)
	v.ValidatePassword(params.Password)
	v.Check(params.Password == params.ConfirmPassword, "confirmPassword", "Las contraseñas no coinciden.")

	if !v.Valid() {
		return render(r, w, auth.SignupForm(params, auth.SignupErrors{
			InvalidEmail:     v.Errors["email"],
			InvalidPassword:  v.Errors["password"],
			PasswordMismatch: v.Errors["confirmPassword"],
		}))
	}

	// Create the user
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return err
	}

	println("Hash: ", string(hash))

	compare := bcrypt.CompareHashAndPassword(hash, []byte(params.Password))
	if compare != nil {
		println("Error: ", compare)
	}

	// user := types.User{
	// 	Email: params.Email,

	// }
	// if err := db.CreateUser(params.Email, params.Password); err != nil {
	// 	return err
	// }

	// Redirect to the login page
	hxRedirect(w, r, "/login")
	return nil
}
