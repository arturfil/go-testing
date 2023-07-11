package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
    DSN string
    DB repository.DatabaseRepo
    Domain string
    JWTSecret string
}

func main() {
    var app application
    err := godotenv.Load()
    if err != nil {
        fmt.Println(err)
    }

    flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application")
    flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 password=secret user=root dbname=unit_testing_db sslmode=disable timezone=UTC connect_timeout=5", "dsn string")
    flag.StringVar(&app.JWTSecret, "jwt-secret", "adfskjl123408asdfljlk;123408asdflkj1234081234081234-9asdfljasdf", "signging secret")
    flag.Parse()

    err = godotenv.Load()
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
    log.Printf("Starting api on port %d\n", port)
    err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
    if err != nil {
        log.Fatal(err)
    }

}
