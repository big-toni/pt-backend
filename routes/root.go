package routes

import (
	"fmt"
	"net/http"
)

// Root func
func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Parcels home page")
}
