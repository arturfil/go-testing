package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_enableCORS(t *testing.T) {

	type test struct {
		name           string
		method         string
		expectedHeader bool
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	var tests = []test{
		{"preflight", "OPTIONS", true},
	}

    for _, e := range tests {
        handlerToTest := app.enableCORS(nextHandler)
        req := httptest.NewRequest(e.method, "http://testing", nil)
        rr := httptest.NewRecorder()

        handlerToTest.ServeHTTP(rr, req)

        if e.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") == "" {
            t.Errorf("%s: expected header, but did not find it", e.name)
        }

        if !e.expectedHeader && rr.Header().Get("Access-Control-Allow-Credentials") != "" {
            t.Errorf("%s: expected no heder but got one", e.name)
        }
    }
}
