package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User: %v\n", vars["id"])
}
