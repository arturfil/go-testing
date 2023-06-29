package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type test struct {
	name               string
	url                string
	expectedStatusCode int
    expectedUrl string
	expectedFirstStatusCode int 
}

type sessionTest struct {
	name         string
	putInSession string
	expectedHTML string
	
}

func Test_application_handlers(t *testing.T) {
	var theTests = []test{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound, "/fish", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{
        Transport: tr,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }

	for _, test := range theTests {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s: expected status %d, bt got %d", test.name, test.expectedStatusCode, resp.StatusCode)
		}

        if resp.Request.URL.Path != test.expectedUrl {
            t.Errorf("%s: expected final url of %s but got %s", test.name, test.expectedUrl, resp.Request.URL.Path)
        }

        resp2, _ := client.Get(ts.URL + test.url) 
        if resp2.StatusCode != test.expectedFirstStatusCode{
            t.Errorf("%s: expected first returned status code to be %d but got %d", test.name, test.expectedFirstStatusCode, resp2.StatusCode)
        }
	}
}

func TestAppHome(t *testing.T) {
	var tests = []sessionTest{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello, world!", "<small>From Session:"},
	}

	for _, test := range tests {

		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, app)

		_ = app.Session.Destroy(req.Context())
		if test.putInSession != "" {
			app.Session.Put(req.Context(), "test", test.putInSession)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(rr, req)

		// check status code
		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), test.expectedHTML) {
			t.Errorf("%s: did not find %s in response body", test.name, test.expectedHTML)
		}
	}
}

func TestApp_renderWithBadTemplate(t *testing.T) {
	// set template to a location with a bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("expected error from template but didn't get one")
	}

	pathToTemplates = "./../../templates/"
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}

func Test_app_Login(t *testing.T) {
	type Test struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}

	var tests = []Test{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
        {
			name: "missing form data",
			postedData: url.Values{
				"email":    {""},
				"password": {""},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
        {
			name: "user not found",
			postedData: url.Values{
				"email":    {"bad@example.com"},
				"password": {"badpassword123"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
        {
			name: "bad credentials",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"badpassword123"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(test.postedData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatusCode {
			t.Errorf("%s: retunred wrong status code: expected %d, but got %d", test.name, test.expectedStatusCode, rr.Code)
		}

		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != test.expectedLoc {
				t.Errorf("%s: expected location %s but got %s", test.name, test.expectedLoc, actualLoc.String())
			}
		} else {
            t.Errorf("%s no location header set", test.name)
        }

	}
}
