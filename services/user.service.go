package services

import (
	"crypto/sha256"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"pt-server/database"
	"pt-server/database/models"
)

// UserDAO interface
type UserDAO interface {
	GetUserForEmail(email string) *models.User
	GetUserByID(ID primitive.ObjectID) *models.User
	Save(user models.User) primitive.ObjectID
	Update(user models.User) primitive.ObjectID
}

// UserService struct
type UserService struct {
	dao UserDAO
}

// Credentials type
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewUserService creates a new UserService with the given user DAO.
func NewUserService(dao UserDAO) *UserService {
	return &UserService{dao}
}

// AuthenticateUser func
func (s *UserService) AuthenticateUser(tknStr string) bool {
	ts := NewTokenService(database.NewTokenDAO())
	valid, _ := ts.ValidateToken(tknStr)

	if !valid {
		return false
	}

	// TODO: find better way
	userID := ts.GetUserID(tknStr)
	// userID := &s.dao.GetUserForEmail(claims.Subject).ID

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
		Model: models.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
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

// GetUser func
func (s *UserService) GetUser(id primitive.ObjectID) *models.User {
	return s.dao.GetUserByID(id)
}

// UpdatePassword func
func (s *UserService) UpdatePassword(email string, password string) (primitive.ObjectID, error) {

	user := *s.dao.GetUserForEmail(email)

	passwordHash := sha256.Sum256([]byte(password))
	passwordHashString := string(passwordHash[:])

	user.PasswordHash = passwordHashString

	updatedID := s.dao.Update(user)
	return updatedID, nil
}
