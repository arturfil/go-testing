package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type test struct {
    name string
    url string
    expectedStatusCode int
}

type sessionTest struct {
    name string
    putInSession string
    expectedHTML string 
}

func Test_application_handlers(t *testing.T) {
    var theTests = []test {
        {"home", "/", http.StatusOK},
        {"404", "/fish", http.StatusNotFound},
    }

    routes := app.routes()

    // create a test server
    ts := httptest.NewTLSServer(routes)
    defer ts.Close()

    pathToTemplates = "./../../templates/"

    for _, test := range theTests {
        resp, err := ts.Client().Get(ts.URL + test.url)
        if err != nil {
            t.Log(err)
            t.Fatal(err)
        }

        if resp.StatusCode != test.expectedStatusCode {
            t.Errorf("for %s: expected status %d, bt got %d", test.name, test.expectedStatusCode, resp.StatusCode)
        }
    }
}

func TestAppHome(t *testing.T) {
    var tests = []sessionTest {
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

func getCtx(req *http.Request) context.Context {
    ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
    return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
    req = req.WithContext(getCtx(req))
    ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
    return req.WithContext(ctx)
}
