package main

import (
	"context"
	"log"
	"net/http"

	"os"
	"strings"

	"pt-backend/routes"
	"pt-backend/services"

	"pt-backend/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	database.Connect()
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "job1" {
		log.Println("Executing:" + os.Args[1])
	} else {
		server()
	}
}

func server() {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	// web := router.PathPrefix("/").Subrouter()
	// web.HandleFunc("/", routes.Root).Methods("GET")
	// admin := router.PathPrefix("/admin").Subrouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(loggingMiddleware)
	auth.HandleFunc("/login/", routes.Login).Methods("POST")
	auth.HandleFunc("/signup/", routes.SignUp).Methods("POST")
	auth.HandleFunc("/reset/{key}/", routes.Reset).Methods("GET", "POST")
	auth.HandleFunc("/forgot/", routes.Forgot).Methods("POST")

	account := api.PathPrefix("/me").Subrouter()
	account.HandleFunc("/", routes.Account).Methods("GET")

	parcels := api.PathPrefix("/parcels").Subrouter()
	parcels.Use(authMiddleware)

	parcels.HandleFunc("/", routes.AddParcels).Methods("POST")
	parcels.HandleFunc("/", routes.GetParcels).Methods("GET")
	parcels.HandleFunc("/", routes.DeleteParcels).Methods("DELETE")
	parcels.HandleFunc("/", routes.EditParcel).Methods("PATCH")

	parcels.HandleFunc("/data/{trackingNumber:[a-zA-Z0-9]+}/", routes.Parcel).Methods("GET")
	parcels.HandleFunc("/courier/{trackingNumber:[a-zA-Z0-9]+}/", routes.Courier).Methods("GET")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.ListenAndServe(":"+port, router)

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Middleware function, which will be called for each request
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		authArr := strings.Split(authToken, " ")

		if len(authArr) != 2 {
			log.Println("Authentication header is invalid: " + authToken)
			w.WriteHeader(http.StatusUnauthorized)
		}

		jwtToken := authArr[1]
		s := services.NewUserService(database.NewUserDAO())
		user, err := s.AuthenticateUser(jwtToken)

		userID := user.ID
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			ctx := context.WithValue(r.Context(), "userID", userID.Hex()) // nolint
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
