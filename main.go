package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"froi/go-templ-poc/handlers"
	"froi/go-templ-poc/types"
	"froi/go-templ-poc/ui/pages"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
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
var sessionManager *scs.SessionManager

func main() {
	title := "GO Templ PoC"
	db, err := sql.Open("sqlite3", "./app_store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionManager = scs.New()
	sessionManager.Lifetime = 10 * time.Minute
	sessionManager.Store = sqlite3store.NewWithCleanupInterval(db, 20*time.Minute)

	pageTitle := "GO Templ PoC"
	mux := http.NewServeMux()
	// http.Handle("GET /{$}", templ.Handler(pages.Homepage(pageTitle)))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		user := sessionManager.GetString(r.Context(), "user")
		loggedIn := false
		if len(user) > 0 {
			loggedIn = true
		}
		pages.Homepage(title, user, loggedIn).Render(r.Context(), w)
	})
	subFS, err := fs.Sub(staticFS, "ui/static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(subFS))))
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		loginFormParams := types.LoginInputFormParams{
			Email:                "",
			Password:             "",
			ShowFailedLoginFlag:  false,
			SubmitButtonDisabled: false,
		}
		pages.LoginPage(title, loginFormParams, false).Render(r.Context(), w)
	})
	mux.HandleFunc("POST /login", handlers.PostLoginHandler(db, sessionManager, pageTitle))
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Destroy(r.Context())
		pages.Homepage(title, "", false).Render(r.Context(), w)
	})
	mux.Handle("GET /signup", templ.Handler(pages.SignupPage(pageTitle, types.SignupInputFormParams{})))
	mux.HandleFunc("POST /signup", handlers.GetSignupHandler(db))
	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", sessionManager.LoadAndSave(mux)); err != nil {
		slog.Error(err.Error())
	}
}
