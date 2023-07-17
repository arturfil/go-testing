package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
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
		{"get", "GET", false},
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

func Test_app_authRequired(t *testing.T) {
    type test struct {
        name string
        token string
        expectedAuthorization bool
        setHeader bool
    }

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
    testUser := data.User {
        ID: 1,
        FirstName: "Admin",
        LastName: "User",
        Email: "admin@example.com",
    } 

    tokens, _ := app.generateTokenPair(&testUser)
    
    var tests = []test {
        {"valid token", fmt.Sprintf("Bearer %s", tokens.Token), true, true},
        {"no token", fmt.Sprintf("Bearer %s", ""), false, false},
        {"invalid token", fmt.Sprintf("Bearer %s", expiredToken), false, true},
    }

    for _, e := range tests {
        req, _ := http.NewRequest("GET", "/", nil)
        if e.setHeader {
            req.Header.Set("Authorization", e.token)
        }
        rr := httptest.NewRecorder()

        handlerToTest := app.authRequired(nextHandler)
        handlerToTest.ServeHTTP(rr, req)

        if e.expectedAuthorization && rr.Code == http.StatusUnauthorized {
            t.Errorf("%s: got code 402, and should not have", e.name)
        }

        if !e.expectedAuthorization && rr.Code != http.StatusUnauthorized {
            t.Errorf("%s: did not get code 401, and should have", e.name)
        }
    }
    
}
