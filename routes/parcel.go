package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"pt-server/database"
	"pt-server/database/models"
	"pt-server/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// AddParcels func
func AddParcels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var pInfos []services.ParcelInfo
	err := json.NewDecoder(r.Body).Decode(&pInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ps := services.NewParcelService(database.NewParcelDAO())

	dbUserID, _ := primitive.ObjectIDFromHex(userID)

	ps.AddParcels(pInfos, dbUserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetParcels func
func GetParcels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	dbUserID, _ := primitive.ObjectIDFromHex(userID)
	ps := services.NewParcelService(database.NewParcelDAO())
	parcels := ps.GetParcels(dbUserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parcels)
}

// EditParcel func
func EditParcel(w http.ResponseWriter, r *http.Request) {
	var dbParcel models.Parcel
	err := json.NewDecoder(r.Body).Decode(&dbParcel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ps := services.NewParcelService(database.NewParcelDAO())
	id, err := ps.UpdateParcel(dbParcel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}

// DeleteParcels func
func DeleteParcels(w http.ResponseWriter, r *http.Request) {
	var dbParcels []models.Parcel
	err := json.NewDecoder(r.Body).Decode(&dbParcels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ps := services.NewParcelService(database.NewParcelDAO())
	id, err := ps.DeleteParcels(dbParcels)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}
