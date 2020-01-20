package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	aa  "github.com/gorilla/http"
)

// Define our struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func main() {
	r := mux.NewRouter()

	amw := authenticationMiddleware{}
	a := make(map[string]string)
	amw.tokenUsers = a
	amw.Populate()

	// r.Use(amw.Middleware)
	r.Use(loggingMiddleware)

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/user/{id:[0-9]+}/", UserHandler)
	r.HandleFunc("/parcel/", ParcelHandler)

	http.ListenAndServe(":3000", r)
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers["00000000"] = "user0"
	amw.tokenUsers["aaaaaaaa"] = "userA"
	amw.tokenUsers["05f717e5"] = "randomUser"
	amw.tokenUsers["deadbeef"] = "user0"
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		log.Println("toni")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Parcels")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User: %v\n", vars["id"])
}

func ParcelHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := aa.Get(os.Stdout, "http://www.gorillatoolkit.org/"); err != nil {
		log.Fatalf("could not fetch: %v", err)
	}

	fmt.Fprintf(w, "Parcel data")
}
