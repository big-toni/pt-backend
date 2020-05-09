package services

import (
	"encoding/json"
	"log"
	"pt-server/couriers"
)

// GetParcelData func
func GetParcelData(parcelNumber string) ([]byte, bool) {
	c1 := make(chan *couriers.ParcelData)
	quit := make(chan bool)

	log.Println("Searching data for parcelNumber:", parcelNumber)

	go func() {
		log.Println("GlobalCanaio data scaper started")
		gcParcelData, _ := couriers.GetGlobalCanaioData(parcelNumber)
		log.Println("GlobalCanaio data scaper finished")
		if gcParcelData != nil {
			log.Println("GlobalCanaio data found")
			c1 <- gcParcelData
			quit <- true
		} else {
			log.Println("GlobalCanaio data not found")
		}
	}()
	go func() {
		log.Println("OrangeConnex data scaper started")
		ocParcelData, _ := couriers.GetOrangeConnexData(parcelNumber)
		log.Println("OrangeConnex data scaper finished")
		if ocParcelData != nil {
			log.Println("OrangeConnex data found")
			c1 <- ocParcelData
			quit <- true
		} else {
			log.Println("OrangeConnex data not found")
		}
	}()

	var result []byte

	for {
		select {
		case <-quit:
			log.Println("Detected quit signal!")
			return result, true
		case msg1 := <-c1:
			// TODO: Maybe needs to be done before sending to channel
			result, _ = json.Marshal(msg1)
		}
	}
}

// ResolveCourier func
func ResolveCourier(parcelNumber string) ([]string, bool) {
	return couriers.ResolveCourier(parcelNumber)
}
