package validator

import (
	"go-auth-template/pkg/utils"
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define a new Validator type which contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// New is a helper which creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) ValidateEmail(email string) {
	if !EmailRX.MatchString(email) {
		v.Errors["email"] = "El email ingresado no es válido."
	}
}

func (v *Validator) ValidatePassword(password string) {
	if len(password) < 8 {
		v.Errors["password"] = "La contraseña debe tener al menos 8 caracteres."
	}

	if len(password) > 72 {
		v.Errors["password"] = "La contraseña debe tener menos de 72 caracteres."
	}

	if !utils.ContainsLowercase(password) {
		v.Errors["password"] = "La contraseña debe tener al menos una letra minúscula."
	}

	if !utils.ContainsNumber(password) {
		v.Errors["password"] = "La contraseña debe tener al menos un número."
	}

	if !utils.ContainsNoWhitespace(password) {
		v.Errors["password"] = "La contraseña no debe tener espacios en blanco."
	}
}

func (v *Validator) ValidateEmailPasswordForLogin(email, password string) {
	if email == "" {
		v.Errors["invalidCredentials"] = "Las credenciales ingresadas no son válidas."
	}

	if password == "" {
		v.Errors["invalidCredentials"] = "Las credenciales ingresadas no son válidas."
	}
}
