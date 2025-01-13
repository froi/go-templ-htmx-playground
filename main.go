package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"froi/go-templ-poc/handlers"
	"froi/go-templ-poc/types"
	"froi/go-templ-poc/ui/pages"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

func getCurrentDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

//go:embed ui/static
var staticFS embed.FS

func main() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
	pageTitle := "GO Templ PoC"
	// http.Handle("GET /{$}", templ.Handler(pages.Homepage(pageTitle)))
	http.Handle("GET /{$}", templ.Handler(pages.Homepage(pageTitle, "", false)))
	subFS, err := fs.Sub(staticFS, "ui/static")
	if err != nil {
		panic(err)
	}
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(subFS))))
	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		title := "GO Templ PoC"
		loginFormParams := types.LoginInputFormParams{
			Email:                "",
			Password:             "",
			ShowFailedLoginFlag:  false,
			SubmitButtonDisabled: false,
		}
		pages.LoginPage(title, loginFormParams).Render(r.Context(), w)
	})
	http.HandleFunc("POST /login", handlers.PostLoginHandler(db, pageTitle))
	http.Handle("GET /signup", templ.Handler(pages.SignupPage(pageTitle, types.SignupInputFormParams{})))
	http.HandleFunc("POST /signup", handlers.GetSignupHandler(db))
	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error(err.Error())
	}
}
