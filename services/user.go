package services

import (
	"errors"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"pt-server/database"
	"pt-server/models"
)

// UserDAO interface
type UserDAO interface {
	GetUserForEmail(email string) *models.User
	GetUserByID(ID primitive.ObjectID) models.User
	Save(user models.User) primitive.ObjectID
}

// UserService struct
type UserService struct {
	dao UserDAO
}

// NewUserService creates a new UserService with the given user DAO.
func NewUserService(dao UserDAO) *UserService {
	return &UserService{dao}
}

// AuthenticateUser func
func (s *UserService) AuthenticateUser(tknStr string) bool {
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
		log.Println("Token not valid: ", tknStr)
		return false
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		log.Println("Claims not valid for token: ", tknStr)
		return false
	}

	if err := claims.Valid(); err != nil {
		log.Println("Claims not valid")
		return false
	}

	ts := NewTokenService(database.NewTokenDAO())
	userID := ts.GetUserID(tknStr)

	if userID == nil {
		log.Println("No token in database: ", tknStr)
		return false
	}

	user := s.dao.GetUserByID(*userID)

	if &user == nil {
		log.Println("No user with token: ", tknStr)
		return false
	}

	return true
}

// CreateUser func
func (s *UserService) CreateUser(email string, password string) primitive.ObjectID {

	dbUser := models.User{
		Email:        email,
		PasswordHash: password,
	}

	id := s.dao.Save(dbUser)

	return id
}

// GetUserForEmail func
func (s *UserService) GetUserForEmail(email string) *models.User {
	return s.dao.GetUserForEmail(email)
}
