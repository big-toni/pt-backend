package couriers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
)

type gcTimelineEntry struct {
	Date        string   `json:"date"`
	Description string   `json:"description" mapstructure:"desc"`
	ID          string   `json:"id"`
	Location    *address `json:"location"`
	Status      string   `json:"status"`
	Time        string   `json:"time"`
}

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
		success := data.(map[string]interface{})["success"]
		if errCode == "RESULT_EMPTY" || success == false {
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

	var timeline []gcTimelineEntry

	parcelData := ParcelData{
		Provider: "GlobalCanaio",
	}
	parcelData.To = &address{Country: result["destCountry"].(string)}
	parcelData.LastUpdated = result["cachedTime"].(string)
	parcelData.From = &address{Country: result["originCountry"].(string)}
	parcelData.ShippingDaysCount = result["shippingTime"].(float64)
	parcelData.Status = result["status"].(string)
	parcelData.StatusDescription = result["statusDesc"].(string)
	mapstructure.Decode(result["section2"].(map[string]interface{})["detailList"], &timeline)
	parcelData.TrackingNumber = result["mailNo"].(string)

	parcelData.Timeline = getGCTimelineData(timeline)

	historyLen := len(*parcelData.Timeline)

	//Add indexes in reversed order
	for i := range *parcelData.Timeline {
		p := &parcelData
		t := *p.Timeline
		t[historyLen-1-i].ID = strconv.Itoa(i)
	}

	return &parcelData, true

}

func getGCTimelineData(gcTimeline []gcTimelineEntry) *[]timelineEntry {
	var parsedTimeline []timelineEntry
	timelineLen := len(gcTimeline)

	for i, item := range gcTimeline {
		entry := timelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.ID = strconv.Itoa(timelineLen - 1 - i)
		entry.Location = item.Location
		entry.Status = item.Status

		layout := "2006-01-02 15:04:05"
		t, err := time.Parse(layout, item.Time)
		entry.Time = t
		if err != nil {
			log.Println(err)
		}
		parsedTimeline = append(parsedTimeline, entry)
	}
	return &parsedTimeline
}
