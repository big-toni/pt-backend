package couriers

import "time"

type address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

type timelineEntry struct {
	Description string    `json:"description"`
	Index       int8      `json:"index"`
	Location    *address  `json:"location"`
	Status      string    `json:"status"`
	Time        time.Time `json:"time"`
}

// ParcelData struct
type ParcelData struct {
	Courier           string           `json:"courier"`
	From              *address         `json:"from"`
	LastUpdated       string           `json:"lastUpdated"`
	Provider          string           `json:"provider"`
	ShippingDaysCount float64          `json:"shippingDaysCount"`
	Status            string           `json:"status"`
	StatusDescription string           `json:"statusDescription"`
	Timeline          *[]timelineEntry `json:"timeline"`
	To                *address         `json:"to"`
	TrackingNumber    string           `json:"trackingNumber"`
}
