package parcels

import "time"

// Address struct
type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

// TimelineEntry struct
type TimelineEntry struct {
	Description string    `json:"description"`
	Index       int8      `json:"index"`
	Location    *Address  `json:"location"`
	Status      string    `json:"status"`
	Time        time.Time `json:"time"`
}

// Timeline type
type Timeline []TimelineEntry

func (slice Timeline) Len() int {
	return len(slice)
}

func (slice Timeline) Less(i, j int) bool {
	return slice[i].Time.Before(slice[j].Time)
}

func (slice Timeline) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ParcelData struct
type ParcelData struct {
	Courier           string    `json:"courier"`
	From              *Address  `json:"from"`
	LastUpdated       string    `json:"lastUpdated"`
	Provider          string    `json:"provider"`
	ShippingDaysCount float64   `json:"shippingDaysCount"`
	Status            string    `json:"status"`
	StatusDescription string    `json:"statusDescription"`
	Timeline          *Timeline `json:"timeline"`
	To                *Address  `json:"to"`
	TrackingNumber    string    `json:"trackingNumber"`
}
