package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type middlewareTest struct {
        headerame string
        headerValue string
        addr string
        emptyAddr bool
}

func Test_application_addIPtoContext(t *testing.T) {
    tests := []middlewareTest {
        {"", "", "", false},
        {"", "", "", true},
        {"X-Forwarded-For", "192.3.2.1", "", false},
        {"", "", "hello:world", false},
    }

    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        val := r.Context().Value(contextUserKey)
        if val == nil {
            t.Error(contextUserKey, "not present")
        }
        ip, ok := val.(string)
        if !ok {
            t.Error("not string")
        }
        t.Log(ip)
    })

    for _, test := range tests {
        handlerToTest := app.addIPToContext(nextHandler)

        req := httptest.NewRequest("GET", "http://testing", nil)

        if test.emptyAddr {
            req.RemoteAddr = ""
        }

        if len(test.headerame) > 0 {
            req.Header.Add(test.headerame, test.headerValue)
        }

        if len(test.addr) > 0 {
            req.RemoteAddr = test.addr
        }

        handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
    }

    
}

func Test_application_ipFromContext(t *testing.T) {
    // get a context 
    ctx := context.Background()   
    ctx = context.WithValue(ctx, contextUserKey, "test_api")    

    // call the function
    ip := app.ipFromContext(ctx)

    // perform the test
    if !strings.EqualFold("test_api", ip) {
        t.Error("wrong value retunred from context")
    }

}
