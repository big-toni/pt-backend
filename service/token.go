package service

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

const (
	jwtSecret = "tonis_stupid_app"
)

var users = map[string]string{
	"test@parcel.tracker": "123TVtv",
	"user2":               "password2",
}

// CreateToken func
func CreateToken(email string, password string) (string, error) {
	tknID, _ := uuid.NewV4()
	expectedPassword, ok := users[email]

	if !ok || expectedPassword != password {
		return "", errors.New("Wrong password")
	}

	claims := jwt.StandardClaims{
		Issuer:   "ParcelTrackr",
		Subject:  email,
		IssuedAt: time.Now().UTC().Unix(),
		Id:       tknID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

// AuthenticateUser func
func AuthenticateUser(tknStr string) bool {
	token, err := jwt.ParseWithClaims(tknStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	if !token.Valid {
		log.Println("Token not valid")
		return false
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		log.Println("Claims not valid")
		return false
	}

	if err := claims.Valid(); err != nil {
		log.Println("Claims not valid")
		return false
	}

	return true
}
