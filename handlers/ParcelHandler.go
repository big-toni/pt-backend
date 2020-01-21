package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	gorillaHttp  "github.com/gorilla/http"
	services "../services"
)

func ParcelHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := gorillaHttp.Get(os.Stdout, "http://www.gorillatoolkit.org/"); err != nil {
		log.Fatalf("could not fetch: %v", err)
	}

	data, _ := services.GetParcelData("aaaa")

	fmt.Fprintf(w, "Parcel data: %v", data)
}
