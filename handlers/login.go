package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"froi/go-templ-poc/types"
	"froi/go-templ-poc/ui/pages"

	"golang.org/x/crypto/bcrypt"
)

func PostLoginHandler(db *sql.DB, title string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		formValues := types.LoginInputFormParams{
			Email:                r.Form.Get("email"),
			Password:             r.Form.Get("password"),
			ShowFailedLoginFlag:  false,
			SubmitButtonDisabled: false,
		}
		slog.Info("Form Values", "formvalues", formValues)
		if formValues.Email == "" || formValues.Password == "" {
			slog.Info("Email or password can't be empty")

			formValues.ShowFailedLoginFlag = true
			pages.LoginPage(title, formValues).Render(r.Context(), w)
			return
		}
		user := types.User{}
		err = user.GetUserByEmail(db, formValues.Email)
		if err != nil {
			slog.Error(err.Error())
			formValues.ShowFailedLoginFlag = true
			pages.LoginPage(title, formValues).Render(r.Context(), w)
			return
		}
		slog.Debug("Found user", "user", user)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formValues.Password))
		if err != nil {
			slog.Error(err.Error())
			formValues.ShowFailedLoginFlag = true
			pages.LoginPage(title, formValues).Render(r.Context(), w)
			return
		}
		slog.Info("User authenticated")
		pages.Homepage(title, user.Email, true).Render(r.Context(), w)
	}
}
