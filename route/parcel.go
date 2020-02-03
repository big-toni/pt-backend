package route

import (
	service "../service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Parcel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, _ := service.GetParcelData(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"data\": %v }", vars["trackingNumber"], data)
}

func Courier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	courier, _ := service.ResolveCourier(vars["trackingNumber"])
	fmt.Fprintf(w, "{ \"trackingNumber\": \"%v\",\"courier\": \"%v\" }", vars["trackingNumber"], courier)
}
