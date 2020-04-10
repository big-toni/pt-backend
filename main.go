package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"pt-server/routes"
	"pt-server/services"

	"github.com/gorilla/mux"
)

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

	account := api.PathPrefix("/account").Subrouter()
	account.Use(authMiddleware)
	account.HandleFunc("/data", routes.Account)
	account.HandleFunc("/user/{id:[0-9]+}/", routes.User)

	parcels := api.PathPrefix("/parcels").Subrouter()
	parcels.Use(authMiddleware)

	parcels.HandleFunc("/data/{trackingNumber:[a-zA-Z0-9]+}/", routes.Parcel)
	parcels.HandleFunc("/courier/{trackingNumber:[a-zA-Z0-9]+}/", routes.Courier)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.ListenAndServe(":"+port, router)

	// if port == "" {
	// 	log.Fatal("$PORT must be set")
	// }

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	// router.Static("/static", "static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	// router.Run(":" + port)

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
		authorised := services.AuthenticateUser(jwtToken)

		if !authorised {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
