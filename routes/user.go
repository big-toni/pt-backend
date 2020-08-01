package routes

import (
	"crypto/sha256"
	"encoding/json"
	"strings"

	"log"
	"net/http"

	"pt-server/database"
	"pt-server/services"

	"github.com/gorilla/mux"
)

var userService = services.NewUserService(database.NewUserDAO())

// User func
func User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	json, err := json.Marshal(struct {
		UserID string `json:"userId"`
	}{
		vars["id"],
	})

	if err != nil {
		log.Println(err)
	}

	w.Write(json)
}

// Login func
func Login(w http.ResponseWriter, r *http.Request) {
	var creds services.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser := userService.GetUserForEmail(creds.Email)

	if existingUser == nil {
		http.Error(w, ("Wrong password or username!"), http.StatusBadRequest)
		return
	}

	passwordHash := sha256.Sum256([]byte(creds.Password))
	passwordHashString := string(passwordHash[:])

	expectedPassword := existingUser.PasswordHash

	if expectedPassword != passwordHashString {
		http.Error(w, ("Wrong password!"), http.StatusBadRequest)
		return
	}

	s := services.NewTokenService(database.NewTokenDAO())

	tokenString, err := s.CreateToken(creds.Email, existingUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		tokenString,
	})
}

// SignUp func
func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds services.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUser := userService.GetUserForEmail(creds.Email)

	if existingUser != nil {
		http.Error(w, ("User already exists!"), http.StatusBadRequest)
		return
	}

	passwordHash := sha256.Sum256([]byte(creds.Password))
	passwordHashString := string(passwordHash[:])

	userID := userService.CreateUser(creds.Email, passwordHashString)

	s := services.NewTokenService(database.NewTokenDAO())

	tokenString, err := s.CreateToken(creds.Email, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		tokenString,
	})
}

// Account func
func Account(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	jwtToken := authArr[1]

	ts := services.NewTokenService(database.NewTokenDAO())
	userID := ts.GetUserID(jwtToken)
	if userID == nil {
		http.Error(w, ("Found no user id with token"), http.StatusBadRequest)
		return
	}

	user := userService.GetUser(*userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
