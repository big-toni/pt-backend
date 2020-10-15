package couriers

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

type dhlHrTimelineEntry struct {
	Date        string           `json:"date"`
	Description string           `json:"description"`
	Index       int8             `json:"index"`
	Location    *parcels.Address `json:"location"`
	Status      string           `json:"status"`
	Time        string           `json:"time"`
}

// DhlHrScraper struct
type DhlHrScraper struct {
}

// NewDhlHrScraper creates a new DhlHrScraper.
func NewDhlHrScraper() *DhlHrScraper {
	return &DhlHrScraper{}
}

func (s *DhlHrScraper) jsGetTimeline(sel string) (js string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in DhlHrScraper jsGetTimeline %s", r)
		}
	}()
	buf, _ := ioutil.ReadFile("helpers/dhlHr.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetData func
func (s *DhlHrScraper) GetData(trackingNumber string) (*parcels.ParcelData, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in DhlHrScraper GetData %s", r)
		}
	}()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), ///- not working without head
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
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	urlString := fmt.Sprintf(`https://www.dhl.com/hr-en/home/tracking/tracking-parcel.html?submit=1&tracking-id=%s`, trackingNumber)

	var timeline []dhlHrTimelineEntry
	parcelData := parcels.ParcelData{
		Provider: "DHL_hr",
	}
	jsTimeline := s.jsGetTimeline("div[class='l-grid l-grid--w-100pc-s l-grid--w-auto-m']")

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		chromedp.WaitVisible("h3[class='c-tracking-result--status-copy-message']"),
		// chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(jsTimeline, &timeline),
	)

	parcelData.Timeline = s.getTimelineData(timeline)

	if err != nil {
		log.Fatal(err)
		return nil, true
	}

	return &parcelData, true
}

func (s *DhlHrScraper) getTimelineData(dhrTimeline []dhlHrTimelineEntry) *parcels.Timeline {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in DhlHrScraper getTimelineData %s", r)
		}
	}()
	var parsedTimeline parcels.Timeline

	for i, item := range dhrTimeline {
		entry := parcels.TimelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.Index = int8(i)
		entry.Location = item.Location
		entry.Status = item.Status

		layout := "2. January 2006 15:04 "
		t, err := time.Parse(layout, "28. August 2020 23:38")
		entry.Time = t
		if err != nil {
			log.Println(err)
		}
		// parsedTimeline = append(parsedTimeline, entry)
		parsedTimeline = append([]parcels.TimelineEntry{entry}, parsedTimeline...)
	}

	return &parsedTimeline
}
