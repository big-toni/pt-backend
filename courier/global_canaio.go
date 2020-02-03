package courier

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetGlobalCanaioData(parcelNumber string) (string, bool) {
	/* test with:
	LE571379316CN,
	S00000111674824,
	S00000111794186,
	S00000111603104,
	S00000110926105,
	LE711833812CN,
	UW309135048CN,
	LVS10060000009698541,
	S00000095458023,
	S00000089577386,
	LE590547286CN,
	S00000075933024,
	*/
	urlString := fmt.Sprintf("http://global.cainiao.com/detail.htm?mailNoList=%s", parcelNumber)
	resp, err := http.Get(urlString)
	if err != nil {
		return "", false
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

		parcelData, _ := json.Marshal(result["data"].([]interface{})[0])

		return string(parcelData), true
	}

	return "", false
}
