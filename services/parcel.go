package services

import (
	"encoding/json"
	"log"
	"pt-server/couriers"
	"sync"
	"time"
)

// GetParcelData func
func GetParcelData(parcelNumber string) ([]byte, bool) {
	ch := make(chan *couriers.ParcelData)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	var result []byte

	log.Println("Searching data for parcelNumber:", parcelNumber)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "GlobalCanaio data scraper")
		log.Println("GlobalCanaio data scraper started")
		gcParcelData, _ := couriers.GetGlobalCanaioData(parcelNumber)
		ch <- gcParcelData
	}(ch, wg)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "OrangeConnex data scraper")
		log.Println("OrangeConnex data scraper started")
		ocParcelData, _ := couriers.GetOrangeConnexData(parcelNumber)
		ch <- ocParcelData
	}(ch, wg)

	go func(ch chan<- *couriers.ParcelData, wg *sync.WaitGroup) {
		defer timeTrack(time.Now(), "PostaHr data scraper")
		log.Println("PostaHr data scraper started")
		phParcelData, _ := couriers.GetPostaHrData(parcelNumber)
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
func ResolveCourier(parcelNumber string) ([]string, bool) {
	return couriers.ResolveCourier(parcelNumber)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s finished with execution time %s", name, elapsed)
}
