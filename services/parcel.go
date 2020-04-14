package services

import (
	"encoding/json"
	"strconv"

	"pt-server/couriers"

	"github.com/mitchellh/mapstructure"
)

type trackingInfo struct {
	Description string `json:"description" mapstructure:"desc"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	Time        string `json:"time"`
	TimeZone    string `json:"timeZone"`
}

// parcelData struct
type parcelData struct {
	DestCountry        string         `json:"destCountry"`
	LastUpdated        string         `json:"lastUpdated"`
	LatestTrackingInfo trackingInfo   `json:"latestTrackingInfo"`
	OriginCountry      string         `json:"originCountry"`
	ShippingDaysCount  float64        `json:"shippingDaysCount"`
	Status             string         `json:"status"`
	StatusDescription  string         `json:"statusDescription"`
	TrackingHistory    []trackingInfo `json:"trackingHistory"`
	TrackingNumber     string         `json:"trackingNumber"`
}

// GetParcelData func
func GetParcelData(parcelNumber string) ([]byte, bool) {
	response, _ := couriers.GetGlobalCanaioData(parcelNumber)

	if response == nil {
		return nil, false
	}

	var result map[string]interface{}
	json.Unmarshal(response, &result)

	var parcelData parcelData

	parcelData.DestCountry = result["destCountry"].(string)
	parcelData.LastUpdated = result["cachedTime"].(string)
	// mapstructure.Decode(result["latestTrackingInfo"].(map[string]interface{}), &parcelData.LatestTrackingInfo)
	parcelData.OriginCountry = result["originCountry"].(string)
	parcelData.ShippingDaysCount = result["shippingTime"].(float64)
	parcelData.Status = result["status"].(string)
	parcelData.StatusDescription = result["statusDesc"].(string)
	mapstructure.Decode(result["section2"].(map[string]interface{})["detailList"], &parcelData.TrackingHistory)
	parcelData.TrackingNumber = result["mailNo"].(string)

	latestTrackingInfo := result["latestTrackingInfo"].(map[string]interface{})
	parcelData.LatestTrackingInfo.Description = latestTrackingInfo["desc"].(string)
	parcelData.LatestTrackingInfo.Status = latestTrackingInfo["status"].(string)
	parcelData.LatestTrackingInfo.Time = latestTrackingInfo["time"].(string)
	parcelData.LatestTrackingInfo.TimeZone = latestTrackingInfo["timeZone"].(string)

	historyLen := len(parcelData.TrackingHistory)
	parcelData.LatestTrackingInfo.ID = strconv.Itoa(historyLen - 1)

	// Add indexes in reversed order
	for i := range parcelData.TrackingHistory {
		parcelData.TrackingHistory[historyLen-1-i].ID = strconv.Itoa(i)
	}

	json, _ := json.Marshal(parcelData)
	return json, true
}

// ResolveCourier func
func ResolveCourier(parcelNumber string) ([]string, bool) {
	return couriers.ResolveCourier(parcelNumber)
}
