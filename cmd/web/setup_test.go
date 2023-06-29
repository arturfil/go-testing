package main

import (
	"log"
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"

	"github.com/joho/godotenv"
)

var app application

func TestMain(m *testing.M) {

	// add this line to get the DSN for testing purposes
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	pathToTemplates = "./../../templates/"

	app.Session = getSession()
	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
