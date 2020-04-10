package routes

import (
	"fmt"
	"net/http"
	"strings"

	"pt-server/services"

	"github.com/gorilla/mux"
)

// Parcel func
func Parcel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, _ := services.GetParcelData(vars["trackingNumber"])
	w.Write(data)
}

// Courier func
func Courier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	couriers, _ := services.ResolveCourier(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"couriers\": %q }", vars["trackingNumber"], strings.Join(couriers, ", "))
}
