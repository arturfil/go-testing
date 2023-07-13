// go:build server
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type test struct {
	name               string
	requestBody        string
	expectedStatusCode int
}

func Test_app_authentication(t *testing.T) {
	var tests = []test{
		{"valid user", `{"email":"admin@example.com", "password":"secret"}`, http.StatusOK},
		{"not json", `not json`, http.StatusUnauthorized},
		{"empty json", `{}`, http.StatusUnauthorized},
		{"empty email", `{"email":"", "password":"secret"}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@example.com", "password":""}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"admin@other.com", "password":"secret"}`, http.StatusUnauthorized},
	}
    for _, test := range tests {
        var reader io.Reader
        reader = strings.NewReader(test.requestBody)
        req, _ := http.NewRequest("POST", "/auth", reader)
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(app.authenticate)

        handler.ServeHTTP(rr, req)

        if test.expectedStatusCode != rr.Code {
            t.Errorf("%s: returned wrong status code; expected %d but got %d", test.name, test.expectedStatusCode, rr.Code)
        }
    }

}
