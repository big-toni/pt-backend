package couriers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

// GetGlobalCanaioData func
func GetGlobalCanaioData(parcelNumber string) ([]byte, bool) {
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

		parcelData, _ := json.Marshal(data)

		return parcelData, true
	}

	return nil, false
}
