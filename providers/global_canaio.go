package providers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"pt-backend/parcels"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type gcTimelineEntry struct {
	Date        string           `json:"date"`
	Description string           `json:"description" mapstructure:"desc"`
	Index       int8             `json:"index"`
	Location    *parcels.Address `json:"location"`
	Status      string           `json:"status"`
	Time        string           `json:"time"`
}

// GlobalCanaioScraper struct
type GlobalCanaioScraper struct {
}

// NewGlobalCanaioScraper creates a new GlobalCanaioScraper.
func NewGlobalCanaioScraper() *GlobalCanaioScraper {
	return &GlobalCanaioScraper{}
}

func (s *GlobalCanaioScraper) jsGetDetails() (js string) {
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

func (s *GlobalCanaioScraper) jsGetTimeline(sel string) (js string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in GlobalCanaioScraper, jsGetTimeline %s", r)
		}
	}()
	buf, _ := ioutil.ReadFile("helpers/globalCanaio.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetData func
func (s *GlobalCanaioScraper) GetData(trackingNumber string) (*parcels.ParcelData, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in GlobalCanaioScraper, GetData %s", r)
		}
	}()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		// chromedp.Flag("disable-gpu", false),
		// chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	urlString := fmt.Sprintf(`http://global.cainiao.com/detail.htm?mailNoList=%s`, trackingNumber)

	timelineEvaluate := s.jsGetTimeline("div[class='TrackingDetail--shipSteps--kvxVRO1']")

	// details := s.jsGetDetails()

	var timeline []gcTimelineEntry
	parcelData := parcels.ParcelData{
		Provider: "GlobalCanaio",
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		chromedp.WaitVisible("div[class='TrackingDetail--shipSteps--kvxVRO1"),
		chromedp.Evaluate(timelineEvaluate, &timeline),
		// chromedp.Evaluate(details, &parcelData),
	)
	parcelData.Timeline = s.getTimelineData(timeline)

	if err != nil {
		// log.Fatal(err)
		log.Println("GlobalCanaio GetData", err)
		return nil, true
	}

	// TODO: need to refactor
	if len(*parcelData.Timeline) == 0 {
		return nil, true
	}

	return &parcelData, true
}

func (s *GlobalCanaioScraper) getTimelineData(gcTimeline []gcTimelineEntry) *parcels.Timeline {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in GlobalCanaioScraper, getTimelineData %s", r)
		}
	}()
	var parsedTimeline parcels.Timeline
	timelineLen := len(gcTimeline)

	for i, item := range gcTimeline {
		entry := parcels.TimelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.Index = int8(timelineLen - 1 - i)
		entry.Location = item.Location
		entry.Status = item.Status

		// 2022-10-10 16:20:20 GMT+8'

		layout := "2006-01-02 15:04:05 MST"
		t, err := time.Parse(layout, item.Time)
		entry.Time = t
		if err != nil {
			log.Println(err)
		}
		parsedTimeline = append(parsedTimeline, entry)
	}

	return &parsedTimeline
}