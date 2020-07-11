package routes

import (
	"crypto/sha256"
	"encoding/json"

	"log"
	"net/http"

	"pt-server/database"
	"pt-server/services"

	"github.com/gorilla/mux"
)

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
	// fmt.Fprintf(w, "User: %v\n", vars["id"])
}

// Account func
func Account(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	json, err := json.Marshal(struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}{
		"testid",
		"test@email.com",
	})

	if err != nil {
		log.Println(err)
	}

	w.Write(json)
}

// Credentials type
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Token string `json:"token"`
}

// Login func
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us := services.NewUserService(database.NewUserDAO())
	existingUser := us.GetUserForEmail(creds.Email)

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
	json.NewEncoder(w).Encode(response{tokenString})
}

// SignUp func
func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us := services.NewUserService(database.NewUserDAO())
	existingUser := us.GetUserForEmail(creds.Email)

	if existingUser != nil {
		http.Error(w, ("User already exists!"), http.StatusBadRequest)
		return
	}

	passwordHash := sha256.Sum256([]byte(creds.Password))
	passwordHashString := string(passwordHash[:])

	userID := us.CreateUser(creds.Email, passwordHashString)

	s := services.NewTokenService(database.NewTokenDAO())

	tokenString, err := s.CreateToken(creds.Email, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{tokenString})
}
