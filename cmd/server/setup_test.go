package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application
var expiredToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2ODkwNDM3NjcsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.IYE7W19-RIo5RZ3AwHtrTJB6ExFUQPxaprAFEwChto4" 


func TestMain(m *testing.M) {
    app.DB = &dbrepo.TestDBRepo{}
    app.Domain = "example.com"
    app.JWTSecret = "adfskjl123408asdfljlk;123408asdflkj1234081234081234-9asdfljasdf"
    os.Exit(m.Run())
}
