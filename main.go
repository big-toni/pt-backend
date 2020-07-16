package main

import (
	"log"
	"net/http"

	"os"
	"strings"

	"pt-server/routes"
	"pt-server/services"

	"pt-server/database"

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
	router := mux.NewRouter()

	router.HandleFunc("/", routes.Root)
	router.Use(loggingMiddleware)

	web := router.PathPrefix("/").Subrouter()
	web.HandleFunc("/", routes.Root)
	// admin := router.PathPrefix("/admin").Subrouter()
	api := router.PathPrefix("/api").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(loggingMiddleware)
	auth.HandleFunc("/login", routes.Login)
	auth.HandleFunc("/signup", routes.SignUp)

	users := api.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware)

	account := users.PathPrefix("/account").Subrouter()
	account.HandleFunc("/data", routes.Account)

	parcels := api.PathPrefix("/parcels").Subrouter()
	parcels.Use(authMiddleware)

	parcels.HandleFunc("/data/{trackingNumber:[a-zA-Z0-9]+}/", routes.Parcel)
	parcels.HandleFunc("/courier/{trackingNumber:[a-zA-Z0-9]+}/", routes.Courier)

	parcels.HandleFunc("/add/{userId:[a-zA-Z0-9]+}/", routes.AddParcels)
	parcels.HandleFunc("/get/{userId:[a-zA-Z0-9]+}/", routes.GetParcels)

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
		authorised := s.AuthenticateUser(jwtToken)

		if !authorised {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
