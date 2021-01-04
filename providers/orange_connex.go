package providers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"pt-server/parcels"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type orangeConnexTimelineEntry struct {
	Date        string           `json:"date"`
	Description string           `json:"description"`
	Index       int8             `json:"index"`
	Location    *parcels.Address `json:"location"`
	Status      string           `json:"status"`
	Time        string           `json:"time"`
}

// OrangeConnexScraper struct
type OrangeConnexScraper struct {
}

// NewOrangeConnexScraper creates a new OrangeConnexScraper.
func NewOrangeConnexScraper() *OrangeConnexScraper {
	return &OrangeConnexScraper{}
}

func (s *OrangeConnexScraper) jsGetDetails() (js string) {
	const funcJS = `function getDetails() {
				var x = {};
				var header = document.body.querySelector("div[class=el-collapse-item__header]");

				items = header && header.querySelectorAll("td[data-v-41aef011]");

				var from = {}
				from.country = items[1]
				.textContent
				.trim();

				from.city = items[2]
				.textContent
				.trim();

				from.zip = items[3]
				.textContent
				.trim();

				x.from = from;

				var to = {}
				to.country = items[5]
				.textContent
				.trim();

				to.city = items[6]
				.textContent
				.trim();

				to.zip = items[7]
				.textContent
				.trim();

				x.to = to;

				return x
			 };`

	invokeFuncJS := `var a = getDetails(); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

func (s *OrangeConnexScraper) jsGetTimeline(sel string) (js string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in OrangeConnexScraper, jsGetTimeline %s", r)
		}
	}()
	buf, _ := ioutil.ReadFile("helpers/orangeConnex.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetData func
func (s *OrangeConnexScraper) GetData(trackingNumber string) (*parcels.ParcelData, bool) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Panic in OrangeConnexScraper, GetData %s", r)
	// 	}
	// }()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.Flag("disable-gpu", false),
		// chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-extensions", false),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()

	// create chrome instance
	ctx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	// defer cancel()

	// create a timeout
	ctx, _ = context.WithTimeout(ctx, 50*time.Second)
	// defer cancel()

	urlString := fmt.Sprintf(`https://www.orangeconnex.com/tracking?language=en&trackingnumber=%s`, trackingNumber)

	timelineEvaluate := s.jsGetTimeline("placeholder")

	details := s.jsGetDetails()

	var timeline []orangeConnexTimelineEntry
	parcelData := parcels.ParcelData{
		Provider: "OrangeConnex",
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		chromedp.Click("div[class='el-collapse-item__header'][role='button']"),
		chromedp.Evaluate(timelineEvaluate, &timeline),
		chromedp.Evaluate(details, &parcelData),
	)
	parcelData.Timeline = s.getTimelineData(timeline)

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

func (s *OrangeConnexScraper) getTimelineData(ocTimeline []orangeConnexTimelineEntry) *parcels.Timeline {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in OrangeConnexScraper, getTimelineData %s", r)
		}
	}()
	var parsedTimeline parcels.Timeline
	timelineLen := len(ocTimeline)

	for i, item := range ocTimeline {
		entry := parcels.TimelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.Index = int8(timelineLen - 1 - i)
		entry.Location = item.Location
		entry.Status = item.Status

		layout := "2006/01/02 15:04"
		t, err := time.Parse(layout, item.Time)
		entry.Time = t
		if err != nil {
			log.Println(err)
		}
		parsedTimeline = append(parsedTimeline, entry)
	}

	return &parsedTimeline
}