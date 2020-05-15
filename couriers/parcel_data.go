package couriers

type address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

type timelineEntry struct {
	Date        string   `json:"date"`
	Description string   `json:"description" mapstructure:"desc"`
	ID          string   `json:"id"`
	Location    *address `json:"location"`
	Status      string   `json:"status"`
	Time        string   `json:"time"`
	TimeZone    string   `json:"timeZone"`
}

// ParcelData struct
type ParcelData struct {
	Courier            string           `json:"courier"`
	From               *address         `json:"from"`
	LastUpdated        string           `json:"lastUpdated"`
	LatestTrackingInfo *timelineEntry   `json:"latestTrackingInfo"`
	Provider           string           `json:"provider"`
	ShippingDaysCount  float64          `json:"shippingDaysCount"`
	Status             string           `json:"status"`
	StatusDescription  string           `json:"statusDescription"`
	Timeline           *[]timelineEntry `json:"timeline"`
	To                 *address         `json:"to"`
	TrackingNumber     string           `json:"trackingNumber"`
}
