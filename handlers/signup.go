package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"froi/go-templ-poc/auth"
	"froi/go-templ-poc/types"

	"github.com/google/uuid"
)

func parseForm(r *http.Request) (types.SignupInputFormParams, error) {
	if err := r.ParseForm(); err != nil {
		return types.SignupInputFormParams{}, err
	}
	formValues := types.SignupInputFormParams{
		Email:     r.Form.Get("email"),
		Password1: r.Form.Get("password1"),
		Password2: r.Form.Get("password2"),
	}
	return formValues, nil
}

func validateRegistrationValues(values types.SignupInputFormParams) error {
	if values.Password1 == "" || values.Password2 == "" {
		return errors.New("passwords cannot be empty")
	}
	if values.Password1 != values.Password2 {
		return errors.New("passwords do not match")
	}
	return nil
}

func createUser(db *sql.DB, email, password string) error {
	userID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	user := types.User{
		ID:       userID.String(),
		Email:    email,
		Password: string(hashedPassword),
	}
	insertSQL := "INSERT INTO users (id, email, password) VALUES (?, ?,?)"
	_, err = db.Exec(insertSQL, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetSignupHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		formValues, err := parseForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = validateRegistrationValues(formValues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = createUser(db, formValues.Email, formValues.Password1)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("User created")
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	}
}
