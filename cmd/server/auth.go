package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"webapp/pkg/data"

	"github.com/golang-jwt/jwt/v4"
)

const jwtTokenExpiry = time.Minute * 15
const refreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
    Token string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type Claims struct {
    UserName string `json:"name"` 
    jwt.RegisteredClaims
}

func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
    // add a header
    w.Header().Add("Vary", "Authorization")
    // get the authorization header
    authHeader := r.Header.Get("Authorization")
    // sanity check
    if authHeader == "" {
        return "", nil, errors.New("no auth header")
    }
    // split the header in spaces
    headerParts := strings.Split(authHeader, " ")
    if len(headerParts) != 2 {
        return "", nil, errors.New("invalid auth header")
    }
    // check if we have "Bearer"
    if headerParts[0] != "Bearer" {
        return "", nil, errors.New("unauthorized: no Bearer")
    }
    token := headerParts[1]
    // declared an empty Claims variable
    claims := &Claims{}
    // parse the token with our claims, using our secret
    _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error){
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(app.JWTSecret), nil
    })

    // check for an error
    if err != nil {
        if strings.HasPrefix(err.Error(), "token is expired by") {
            return "", nil, errors.New("expired token")
        }
        return "", nil, err
    }

    if claims.Issuer != app.Domain {
        return "", nil, errors.New("incorrect issuer")
    }

    return token, claims, nil
}

func (app *application) generateTokenPair(user *data.User) (TokenPairs, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
    claims["sub"] = fmt.Sprint(user.ID)
    claims["aud"] = app.Domain
    claims["iss"] = app.Domain
    if user.IsAdmin == 1 {
        claims["admin"] = true
    } else {
        claims["admin"] = false
    }
    claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()
    
    // created signed token
    signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
    if err != nil {
        return TokenPairs{}, err
    }
    // create the refresh token
    refreshToken := jwt.New(jwt.SigningMethodHS256)
    refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
    refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
    // set expiry
    refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()
    // created signed refreshed token
    signedRefreshedToken, err := refreshToken.SignedString([]byte(app.JWTSecret))
    if err != nil {
        return TokenPairs{}, err
    }

    var tokenPairs = TokenPairs{
        Token: signedAccessToken,
        RefreshToken: signedRefreshedToken,
    }

    return tokenPairs, nil

}
