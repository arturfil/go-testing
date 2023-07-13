package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
    app.DB = &dbrepo.TestDBRepo{}
    app.Domain = "example.com"
    app.JWTSecret = "adfskjl123408asdfljlk;123408asdflkj1234081234081234-9asdfljasdf"
    os.Exit(m.Run())
}
