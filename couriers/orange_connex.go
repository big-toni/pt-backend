package couriers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func jsGetDetails() (js string) {
	const funcJS = `function getDetails() {
				var x = {};
				var header = document.body.querySelector("div[class='el-collapse-item__header']");

				items = header.querySelectorAll("p");
				var details = {};
				details.description = items[1].textContent.trim();
				details.date = items[2] && items[2].textContent.trim();
				details.location = items[3] && items[3].textContent.trim();

				x.details = details;

				var from = {}
				from.city = [...header.querySelectorAll("ul[class='fl']>li>b")].reduce((sum, x) => sum.concat(x.textContent.trim(), ", ") ,"");
				const foundPostCode1 = header.querySelector("li[class='fl']>b[class='postCode']")
				from.postCode = foundPostCode1 && foundPostCode1.textContent.trim();
				x.from  = from;

				var to = {}
				to.city = [...header.querySelectorAll("ul[class='fr']>li>b")].reduce((sum, x) => sum.concat(x.textContent.trim(), ", ") ,"");
				const foundPostCode2 = header.querySelector("ul[class='fr']>li>b[class='postCode']")
				to.postCode = foundPostCode2 && foundPostCode2.textContent.trim();
				x.to = to;

				return x
			 };`

	invokeFuncJS := `var a = getDetails(); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

func jsGetOrangeConnexTimeline(sel string) (js string) {
	buf, _ := ioutil.ReadFile("jsHelpers/orangeConnex.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetOrangeConnexData func
func GetOrangeConnexData(parcelNumber string) (*ParcelData, bool) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	urlString := fmt.Sprintf(`https://www.orangeconnex.com/tracking?language=en&trackingnumber=%s`, parcelNumber)

	timelineEvaluate := jsGetOrangeConnexTimeline("ul[class='timeline']")

	details := jsGetDetails()

	var timeline []timelineEntry
	parcelData := ParcelData{
		Provider: "OrangeConnex",
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		// chromedp.WaitVisible("ul[data-v-41aef011] > h3[data-v-41aef011] > div"),
		chromedp.WaitVisible("div > div[data-v-41aef011][class='part1']"),
		chromedp.Evaluate(timelineEvaluate, &timeline),
		chromedp.Evaluate(details, &parcelData),
	)
	parcelData.Timeline = &timeline

	if err != nil {
		log.Fatal(err)
		return nil, true
	}

	// TODO: need to refactor
	if len(*parcelData.Timeline) == 0 {
		return nil, true
	}

	return &parcelData, true
}