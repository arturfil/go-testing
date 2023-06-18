package main

import (
	"log"
	"net/http"
	"os"
	"webapp/pkg/db"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

type application struct {
	DB     db.PostgresConn 
	DSN     string
	Session *scs.SessionManager
}

func main() {
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
    
    app.DB = db.PostgresConn{DB: conn}

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
