package couriers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/html"
)

type gcTimelineEntry struct {
	Date        string   `json:"date"`
	Description string   `json:"description" mapstructure:"desc"`
	Index       int8     `json:"index"`
	Location    *address `json:"location"`
	Status      string   `json:"status"`
	Time        string   `json:"time"`
}

// GetGlobalCanaioData func
func GetGlobalCanaioData(trackingNumber string) (*ParcelData, error) {
	urlString := fmt.Sprintf("http://global.cainiao.com/detail.htm?mailNoList=%s", trackingNumber)
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
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
			return nil, errors.New("RESULT_EMPTY")
		}

		data2, _ := json.Marshal(data)

		parcelDataPointer, _ := mapData(data2)

		return parcelDataPointer, nil
	}

	return nil, errors.New("Error in GetGlobalCanaioData func")
}

func mapData(data []byte) (parcelData *ParcelData, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic in global_canaio mapData %s", r)
			parcelData = nil
		}
	}()

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	var timeline []gcTimelineEntry

	pd := ParcelData{
		Provider: "GlobalCanaio",
	}
	pd.To = &address{Country: result["destCountry"].(string)}
	pd.LastUpdated = result["cachedTime"].(string)
	pd.ShippingDaysCount = result["shippingTime"].(float64)
	pd.From = &address{Country: result["originCountry"].(string)}
	pd.Status = result["status"].(string)
	pd.StatusDescription = result["statusDesc"].(string)
	mapstructure.Decode(result["section2"].(map[string]interface{})["detailList"], &timeline)
	pd.TrackingNumber = result["mailNo"].(string)

	pd.Timeline = getGCTimelineData(timeline)

	historyLen := len(*pd.Timeline)

	//Add indexes in reversed order
	for i := range *pd.Timeline {
		p := &pd
		t := *p.Timeline
		t[historyLen-1-i].Index = int8(i)
	}

	return &pd, err
}

func getGCTimelineData(gcTimeline []gcTimelineEntry) *[]timelineEntry {
	var parsedTimeline []timelineEntry
	timelineLen := len(gcTimeline)

	for i, item := range gcTimeline {
		entry := timelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.Index = int8(timelineLen - 1 - i)
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
