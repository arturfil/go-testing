package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"webapp/pkg/data"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

type application struct {
	DB     repository.DatabaseRepo
	DSN     string
	Session *scs.SessionManager
}

func main() {
    gob.Register(data.User{})
	// setup an app config
	app := application{}

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.DSN = os.Getenv("DSN")
	conn, err := app.connectToDB()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    app.DB = &dbrepo.PostgresDBRepo{DB: conn} 

	// get a session manager
	app.Session = getSession()

	// print oot a message
	log.Println("Starting server on port: 8080...")

	// start the server
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
