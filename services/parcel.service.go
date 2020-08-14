package services

import (
	"encoding/json"
	"log"
	"pt-server/couriers"
	"pt-server/database/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParcelDAO interface
type ParcelDAO interface {
	Save(parcel []models.Parcel) []primitive.ObjectID
	GetParcelsForUserID(userID primitive.ObjectID) []*models.Parcel
	Update(parcel models.Parcel) primitive.ObjectID
	Delete(parcel []models.Parcel) []primitive.ObjectID
}

// ParcelService struct
type ParcelService struct {
	dao ParcelDAO
}

// ParcelInfo type
type ParcelInfo struct {
	Description    string `json:"description"`
	Name           string `json:"name"`
	TrackingNumber string `json:"trackingNumber"`
}

// NewParcelService creates a new ParcelService with the given parcel DAO.
func NewParcelService(dao ParcelDAO) *ParcelService {
	return &ParcelService{dao}
}

// GetParcelData func
func (s *ParcelService) GetParcelData(trackingNumber string) ([]byte, bool) {

	ch := make(chan *couriers.ParcelData)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	var result []byte

	log.Println("Searching data for trackingNumber:", trackingNumber)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "GlobalCanaio data scraper")
		log.Println("GlobalCanaio data scraper started")
		gcParcelData, _ := couriers.GetGlobalCanaioData(trackingNumber)
		ch <- gcParcelData
	}(ch, wg)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "OrangeConnex data scraper")
		log.Println("OrangeConnex data scraper started")
		ocParcelData, _ := couriers.GetOrangeConnexData(trackingNumber)
		ch <- ocParcelData
	}(ch, wg)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "PostaHr data scraper")
		log.Println("PostaHr data scraper started")
		phParcelData, _ := couriers.GetPostaHrData(trackingNumber)
		ch <- phParcelData
	}(ch, wg)

	go func(ch <-chan *couriers.ParcelData, wg *sync.WaitGroup) {
		for msg := range ch {
			if msg != nil {
				result, _ = json.Marshal(msg)
			}
			wg.Done()
		}

	}(ch, wg)

	wg.Wait()
	// close(ch)
	return result, true

}

// ResolveCourier func
func (s *ParcelService) ResolveCourier(trackingNumber string) ([]string, bool) {
	return couriers.ResolveCourier(trackingNumber)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s finished with execution time %s", name, elapsed)
}

// AddParcels func
func (s *ParcelService) AddParcels(parcelInfos []ParcelInfo, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	var parcels []models.Parcel

	for _, item := range parcelInfos {

		dbParcel := models.Parcel{
			Model: models.Model{
				ID:        primitive.ObjectID{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Description:    item.Description,
			Name:           item.Name,
			TrackingNumber: item.TrackingNumber,
			UserID:         userID,
		}
		parcels = append(parcels, dbParcel)
	}

	savedIDs := s.dao.Save(parcels)

	return savedIDs, nil
}

// GetParcels func
func (s *ParcelService) GetParcels(userID primitive.ObjectID) []*models.Parcel {
	return s.dao.GetParcelsForUserID(userID)
}

// UpdateParcel func
func (s *ParcelService) UpdateParcel(parcel models.Parcel) (primitive.ObjectID, error) {
	updatedID := s.dao.Update(parcel)
	return updatedID, nil
}

// DeleteParcels func
func (s *ParcelService) DeleteParcels(parcels []models.Parcel) ([]primitive.ObjectID, error) {
	deletedIDs := s.dao.Delete(parcels)
	return deletedIDs, nil
}