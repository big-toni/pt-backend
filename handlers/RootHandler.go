package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Parcels home page")
}
