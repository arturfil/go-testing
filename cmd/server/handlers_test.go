// go:build server
package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"

	"github.com/go-chi/chi/v5"
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

func Test_app_refresh(t *testing.T) {
	type test struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}

	var tests = []test{
		{"valid", "", http.StatusOK, true},
		{"valid but not ready to expire", "", http.StatusTooEarly, false},
		{"expired token", expiredToken, http.StatusBadRequest, false},
	}

	testUser := data.User{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
	}

	oldRefreshTime := refreshTokenExpiry

	for _, e := range tests {
		var tkn string
		if e.token == "" {
			if e.resetRefreshTime {
				refreshTokenExpiry = time.Second * 1
			}
			tokens, _ := app.generateTokenPair(&testUser)
			tkn = tokens.RefreshToken
		} else {
			tkn = e.token
		}

		postedData := url.Values{
			"refresh_token": {tkn},
		}

		req, _ := http.NewRequest("POST", "/refresh-token", strings.NewReader(postedData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.refresh)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: expected status of %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		refreshTokenExpiry = oldRefreshTime
	}

}

func Test_app_userHandlers(t *testing.T) {
	type test struct {
		name           string
		method         string
		json           string
		paramID        string
		handler        http.HandlerFunc
		expectedStatus int
	}

	var tests = []test{
		{"allUsers", "GET", "", "", app.allUsers, http.StatusOK},
		{"deleteUser", "DELETE", "", "1", app.deleteUser, http.StatusNoContent},
		{"deleteUser bad id", "DELETE", "", "Y", app.deleteUser, http.StatusBadRequest},
		{"getUserValid", "GET", "", "1", app.getUser, http.StatusOK},
		{"getUserValid bad id", "GET", "", "Y", app.getUser, http.StatusBadRequest},
		{"getUserInvalid", "GET", "", "2", app.getUser, http.StatusBadRequest},
		{
			"updateUser valid",
			"PATCH",
			`{
                "id":1, 
                "first_name": "Administrator", 
                "last_name": "User", 
                "email": "admin@example.com"
            }`,
			"",
			app.updateUser,
			http.StatusNoContent,
		},
        {
			"updateUser invalid",
			"PATCH",
			`{
                "id":2, 
                "first_name": "Administrator", 
                "last_name": "User", 
                "email": "admin@example.com"
            }`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
        {
			"updateUser valid",
			"PATCH",
			`{
                "id":1, 
                first_name": "Administrator", 
                "last_name": "User", 
                "email": "admin@example.com"
            }`,
			"",
			app.updateUser,
			http.StatusBadRequest,
		},
        {
			"insert user",
			"POST",
			`{
                "first_name": "Jack", 
                "last_name": "Example", 
                "email": "jack@example.com"
            }`,
			"",
			app.insertUser,
			http.StatusNoContent,
		},
        {
			"insert user invalid attribute",
			"POST",
			`{
                "foo":"bar",
                "first_name": "Arturo", 
                "last_name": "Filio", 
                "email": "arturo@example.com"
            }`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
        {
			"insert user invalid json",
			"POST",
			`{
                first_name": "Arturo", 
                "last_name": "Filio", 
                "email": "arturo@example.com"
            }`,
			"",
			app.insertUser,
			http.StatusBadRequest,
		},
	}

	for _, e := range tests {
		var req *http.Request

		if e.json == "" {
			req, _ = http.NewRequest(e.method, "/", nil)
		} else {
			req, _ = http.NewRequest(e.method, "/", strings.NewReader(e.json))
		}

		if e.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("userID", e.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatus {
			t.Errorf("%s: wrong status returned; exptected %d but got %d", e.name, e.expectedStatus, rr.Code)
		}
	}
}
