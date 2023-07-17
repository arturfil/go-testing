package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    // register middleware
    mux.Use(middleware.Recoverer)
    mux.Use(app.enableCORS)
    // auth routes - auth handler, refresh
    mux.Post("/auth", app.authenticate)
    mux.Post("/refresh-token", app.refresh)
    // protected routes
    mux.Route("/users", func(mux chi.Router) {
        mux.Use(app.authRequired)
        mux.Get("/", app.allUsers)
        mux.Get("/{userID}", app.getUser)
        mux.Post("/", app.insertUser)
        mux.Patch("/{userID}", app.updateUser)
        mux.Delete("/{userID}", app.deleteUser)
    })
    return mux
}
