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
		{"valid user", `{"email":"admin@example.com", "password":"secret"}`, 200},
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
