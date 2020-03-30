package route

import (
	"fmt"
	"net/http"
	"strings"

	service "../service"
	"github.com/gorilla/mux"
)

// Parcel func
func Parcel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, _ := service.GetParcelData(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"data\": %v }", vars["trackingNumber"], data)
}

// Courier func
func Courier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	couriers, _ := service.ResolveCourier(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"couriers\": %q }", vars["trackingNumber"], strings.Join(couriers, ", "))
}
