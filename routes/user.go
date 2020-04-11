package routes

import (
	"encoding/json"

	"log"
	"net/http"

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

	// log.Printf("Authenticate user\nemail: %s\npassword: %s", creds.Email, creds.Password)

	tokenString, err := services.CreateToken(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	services.SaveToken(tokenString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{tokenString})
}
