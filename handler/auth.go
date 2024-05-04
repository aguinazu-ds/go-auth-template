package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"go-auth-template/db"
	"go-auth-template/pkg/argon2id"
	"go-auth-template/pkg/authsession"
	"go-auth-template/pkg/mailer"
	"go-auth-template/pkg/utils"
	"go-auth-template/pkg/validator"
	"go-auth-template/types"
	"go-auth-template/view/auth"
)

func HandleAuthLogin(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleAuthLoginPost(w http.ResponseWriter, r *http.Request) error {
	params := auth.LoginParams{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	v := validator.New()
	v.ValidateEmailPasswordForLogin(params.Email, params.Password)

	if !v.Valid() {
		return render(r, w, auth.LoginForm(params, auth.LoginErrors{
			InvalidCredentials: v.Errors["invalidCredentials"],
		}))
	}

	user, err := db.GetUserByEmail(params.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return render(r, w, auth.LoginForm(params, auth.LoginErrors{
				InvalidCredentials: "Las credenciales ingresadas no son válidas.",
			}))
		}
		return err
	}

	match, err := argon2id.ComparePasswordAndHash(params.Password, user.EncryptedPassword)
	if err != nil {
		return err
	}

	if !match {
		return render(r, w, auth.LoginForm(params, auth.LoginErrors{
			InvalidCredentials: "Las credenciales ingresadas no son válidas.",
		}))
	}

	store := authsession.GetStore()
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["user"] = types.AuthenticatedUser{
		ID:       user.ID,
		Email:    user.Email,
		LoggedIn: true,
	}

	if err := session.Save(r, w); err != nil {
		return err
	}

	hxRedirect(w, r, "/")
	return nil
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

	// Hash the password
	hash, err := argon2id.CreateHash(params.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	RawAppMetaData := map[string]interface{}{
		"provider":  "email",
		"providers": []string{"email"},
	}

	newUser := &types.User{
		Email:             params.Email,
		Activated:         false,
		EncryptedPassword: string(hash),
		RawAppMetaData:    RawAppMetaData,
	}

	// Start a transaction
	tx, err := db.Bun.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Create the user
	if err := db.CreateUser(newUser); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			_ = tx.Rollback()
			return render(r, w, auth.SignupForm(params, auth.SignupErrors{
				EmailInUse: "El email ya está en uso.",
			}))
		}
		_ = tx.Rollback()
		return err
	}

	// Create activation token
	token, plaintext, err := utils.GenerateRandomTokenHash()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Save the token in the database
	if err := db.CreateActivationToken(token, newUser); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	data := mailer.WelcomeEmailData{
		ActivationUrl: os.Getenv("APP_HOST") + "/activate?token=" + plaintext,
	}

	// Send welcome email with activation link
	if err := mailer.SendEmailUsingTemplate(params.Email, "Bienvenido a la aplicación", "aguinazu-dev", "noreply@aguinazu-dev.xyz", "user_welcome.tmpl", data); err != nil {
		return err
	}

	return render(r, w, auth.SignupSuccess(params.Email))
}

func HandleAuthActivate(w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	if token == "" {
		slog.Error("token is empty")
		return render(r, w, auth.ActivationError(auth.ActivationErrors{
			InvalidToken: "Hubo un error al activar la cuenta. Por favor, intenta nuevamente. Si el problema persiste, contacta al soporte.",
		}))
	}

	user, err := db.GetUserByToken(token, types.ScopeActivation)
	if err != nil {
		slog.Error("error getting user by token: ", err)
		return render(r, w, auth.ActivationError(auth.ActivationErrors{
			InvalidToken: "Hubo un error al activar la cuenta. Por favor, intenta nuevamente. Si el problema persiste, contacta al soporte.",
		}))
	}

	user.Activated = true

	if err := db.UpdateUser(&user); err != nil {
		slog.Error("error updating user: ", err)
		return render(r, w, auth.ActivationError(auth.ActivationErrors{
			InvalidToken: "Hubo un error al activar la cuenta. Por favor, intenta nuevamente. Si el problema persiste, contacta al soporte.",
		}))
	}

	// Delete all activation tokens for this user
	if err := db.DeleteAllTokensByUserIDAndScope(user.ID, types.ScopeActivation); err != nil {
		slog.Error("error deleting all tokens by user id and scope: ", err)
		return render(r, w, auth.ActivationError(auth.ActivationErrors{
			InvalidToken: "Hubo un error al activar la cuenta. Por favor, intenta nuevamente. Si el problema persiste, contacta al soporte.",
		}))
	}

	return render(r, w, auth.ActivationSuccess())
}
