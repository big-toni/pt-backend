package handlers

import (
	services "../services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func ParcelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, _ := services.GetParcelData(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"data\": %v }", vars["trackingNumber"], data)
}
