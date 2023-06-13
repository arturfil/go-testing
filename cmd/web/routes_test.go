package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type registeredRoutes struct {
	route  string
	method string
}

func Test_application_routes(t *testing.T) {
	var registered = []registeredRoutes{
		{"/", "GET"},
        {"/login", "POST"},
		{"/static/*", "GET"},
	}

	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		if !routesExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}

	}
}

func routesExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(
        chiRoutes, 
        func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute,
        ) {
			found = true
		}
		return nil
	})
	return found
}
