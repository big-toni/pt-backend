package couriers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
)

// GetGlobalCanaioData func
func GetGlobalCanaioData(parcelNumber string) (*ParcelData, bool) {
	urlString := fmt.Sprintf("http://global.cainiao.com/detail.htm?mailNoList=%s", parcelNumber)
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := string(body)

	re := regexp.MustCompile(`id="waybill_list_val_box">(.*)</textarea>`)
	submatchall := re.FindAllStringSubmatch(result, -1)
	for _, element := range submatchall {
		unescaped := html.UnescapeString(element[1])

		// Declared an empty interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(unescaped), &result)

		data := result["data"].([]interface{})[0]

		errCode := data.(map[string]interface{})["errorCode"]
		if errCode == "RESULT_EMPTY" {
			return nil, false
		}

		data2, _ := json.Marshal(data)

		parcelDataPointer, _ := mapData(data2)

		return parcelDataPointer, true
	}

	return nil, false
}

func mapData(data []byte) (*ParcelData, bool) {

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	var parcelData ParcelData
	parcelData.To = &address{Country: result["destCountry"].(string)}
	parcelData.LastUpdated = result["cachedTime"].(string)
	// mapstructure.Decode(result["latestTrackingInfo"].(map[string]interface{}), &parcelData.LatestTrackingInfo)
	parcelData.From = &address{Country: result["originCountry"].(string)}
	parcelData.ShippingDaysCount = result["shippingTime"].(float64)
	parcelData.Status = result["status"].(string)
	parcelData.StatusDescription = result["statusDesc"].(string)
	mapstructure.Decode(result["section2"].(map[string]interface{})["detailList"], &parcelData.Timeline)
	parcelData.TrackingNumber = result["mailNo"].(string)

	latestTrackingInfo := result["latestTrackingInfo"].(map[string]interface{})

	parcelData.LatestTrackingInfo = &timelineEntry{
		// Date        string `json:"date"`
		Description: latestTrackingInfo["desc"].(string),
		// ID          string `json:"id"`
		Status:   latestTrackingInfo["status"].(string),
		Time:     latestTrackingInfo["time"].(string),
		TimeZone: latestTrackingInfo["timeZone"].(string),
	}
	historyLen := len(*parcelData.Timeline)
	parcelData.LatestTrackingInfo.ID = strconv.Itoa(historyLen - 1)

	//Add indexes in reversed order
	for i := range *parcelData.Timeline {
		p := &parcelData
		t := *p.Timeline
		t[historyLen-1-i].ID = strconv.Itoa(i)
	}

	return &parcelData, true

}
