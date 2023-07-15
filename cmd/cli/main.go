package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type application struct {
    JWTSecret string
    Action string
}

func main() {
    var app application
    flag.StringVar(&app.JWTSecret, "jwt-secret", "adfskjl123408asdfljlk;123408asdflkj1234081234081234-9asdfljasdf", "signging secret")
    flag.StringVar(&app.Action, "action", "valid", "action: valid|expired")
    flag.Parse()

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["name"] = "John Doe"
    claims["sub"] = 1
    claims["admin"] = true
    claims["aud"] = "example.com"
    claims["iss"] = "example.com"
    
    if app.Action == "valid" {
        expires := time.Now().UTC().Add(time.Hour * 72)
        claims["exp"] = expires.Unix()
        fmt.Println("VALID Token:")
    } else {
        expires := time.Now().UTC().Add(time.Hour * 100 * -1)
        claims["exp"] = expires.Unix()
        fmt.Println("EXPIRED Token:")
    }

    signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(signedAccessToken))
}
