package services

import (
	"encoding/json"
	"log"
	"pt-server/couriers"
)

// GetParcelData func
func GetParcelData(parcelNumber string) ([]byte, bool) {
	c1 := make(chan *couriers.ParcelData)
	c2 := make(chan *couriers.ParcelData)

	log.Println("Searching data for parcelNumber:", parcelNumber)

	go func() {
		gcParcelData, _ := couriers.GetGlobalCanaioData(parcelNumber)
		c1 <- gcParcelData
	}()
	go func() {
		ocParcelData, _ := couriers.GetOrangeconnexData(parcelNumber)
		c2 <- ocParcelData
	}()

	var result []byte

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			log.Println("Received gcParcelData:", msg1)
			result, _ = json.Marshal(msg1)
		case msg2 := <-c2:
			log.Println("Received ocParcelData:", msg2)
			result, _ = json.Marshal(msg2)
		}
	}

	return result, true
}

// ResolveCourier func
func ResolveCourier(parcelNumber string) ([]string, bool) {
	return couriers.ResolveCourier(parcelNumber)
}
