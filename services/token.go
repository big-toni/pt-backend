package services

import (
	"crypto/sha256"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// TODO: Move to database
var users = map[string]string{
	"test@parcel.tracker": "123TVtv",
	"user2":               "password2",
}

// CreateToken func
func CreateToken(email string, password string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")

	tknID := uuid.NewV4()
	expectedPassword, ok := users[email]

	if !ok || expectedPassword != password {
		return "", errors.New("Wrong password")
	}

	claims := jwt.StandardClaims{
		Issuer:   jwtIssuer,
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
	jwtSecret := os.Getenv("JWT_SECRET")
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

// SaveToken func
func SaveToken(token string) (bool, error) {
	sha256 := sha256.Sum256([]byte(token))

	log.Printf("sha256: %x\n", sha256)

	return true, nil
}
