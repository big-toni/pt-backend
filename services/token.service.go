package services

import (
	"crypto/sha256"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"pt-backend/database/models"
)

// TokenDAO interface
type TokenDAO interface {
	Save(token models.Token) string
	GetUserID(tokenHash string) *primitive.ObjectID
}

// TokenService struct
type TokenService struct {
	dao TokenDAO
}

// NewTokenService creates a new TokenService with the given token DAO.
func NewTokenService(dao TokenDAO) *TokenService {
	return &TokenService{dao}
}

// CreateToken func
func (s *TokenService) CreateToken(email string, userID primitive.ObjectID) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")

	tknID := primitive.NewObjectID()

	claims := jwt.StandardClaims{
		Issuer:   jwtIssuer,
		Subject:  email,
		IssuedAt: time.Now().UTC().Unix(),
		Id:       tknID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	sha256 := sha256.Sum256([]byte(tokenString))

	dbToken := models.Token{
		Model: models.Model{
			ID:        tknID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		AppID:     "appid_placeholder",
		TokenHash: string(sha256[:]),
		Source:    "iOS",
		UserID:    userID,
	}

	s.dao.Save(dbToken)

	return tokenString, nil
}

// GetUserID func
func (s *TokenService) GetUserID(token string) *primitive.ObjectID {

	sha256 := sha256.Sum256([]byte(token))
	tokenHash := string(sha256[:])

	userID := s.dao.GetUserID(tokenHash)

	return userID
}

// CreatePasswordResetToken func
func (s *TokenService) CreatePasswordResetToken(email string, userID primitive.ObjectID) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtExpiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))

	tknID := primitive.NewObjectID()

	claims := jwt.StandardClaims{
		Issuer:    jwtIssuer,
		Subject:   email,
		IssuedAt:  time.Now().UTC().Unix(),
		Id:        tknID.Hex(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(jwtExpiration)).UTC().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	// sha256 := sha256.Sum256([]byte(tokenString))

	// dbToken := models.Token{
	// 	Model: models.Model{
	// 		ID:        tknID,
	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 	},
	// 	AppID:     "appid_placeholder",
	// 	TokenHash: string(sha256[:]),
	// 	Source:    "iOS",
	// 	UserID:    userID,
	// }

	// s.dao.Save(dbToken)

	return tokenString, nil
}

// ValidateToken func
func (s *TokenService) ValidateToken(tknStr string) (bool, *jwt.StandardClaims) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tknStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println(err)
		return false, nil
	}

	if !token.Valid {
		log.Println("Token not valid: ", tknStr)
		return false, nil
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		log.Println("Claims not valid for token: ", tknStr)
		return false, nil
	}

	if err := claims.Valid(); err != nil {
		log.Println("Claims not valid")
		return false, nil
	}

	return true, claims
}
