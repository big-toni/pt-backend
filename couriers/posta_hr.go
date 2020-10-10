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

type postaHrTimelineEntry struct {
	Date        string   `json:"date"`
	Description string   `json:"description"`
	Index       int8     `json:"index"`
	Location    *address `json:"location"`
	Status      string   `json:"status"`
	Time        string   `json:"time"`
}

// PostaHrScraper struct
type PostaHrScraper struct {
}

// NewPostaHrScraper creates a new PostaHrScraper.
func NewPostaHrScraper() *PostaHrScraper {
	return &PostaHrScraper{}
}

func (s *PostaHrScraper) jsGetTimeline(sel string) (js string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in PostaHrScraper, jsGetTimeline %s", r)
		}
	}()
	buf, _ := ioutil.ReadFile("helpers/postaHr.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetData func
func (s *PostaHrScraper) GetData(trackingNumber string) (*ParcelData, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in PostaHrScraper, GetData %s", r)
		}
	}()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
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

	urlString := fmt.Sprintf(`https://posiljka.posta.hr/Tracking/Info`)

	var timeline []postaHrTimelineEntry
	parcelData := ParcelData{
		Provider: "PostaHr",
	}
	jsTimeline := s.jsGetTimeline("div[class='styles__table___5Ule6']")

	var foundData string

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlString),
		chromedp.WaitVisible("input[class='__c-form-field__text']"),
		chromedp.Click("input[class='__c-form-field__text']"),
		chromedp.SetValue("input[class='__c-form-field__text']", trackingNumber),
		chromedp.Click("input[class='__c-form-field__text']"),
		chromedp.Sleep(1*time.Second),
		chromedp.Focus("input[class='__c-form-field__text']"),
		chromedp.Sleep(1*time.Second),
		chromedp.Submit("input[class='__c-form-field__text']"),

		chromedp.Click("button[class='__c-btn form-submit']"),

		chromedp.InnerHTML("div[class='__c-shipment__details']", &foundData),

		chromedp.Text("sham-shipment-origin-date", &parcelData.StatusDescription),
		chromedp.Text("div[class='__c-heading __c-heading--h4 __c-heading--bold __u-mb--none']", &parcelData.TrackingNumber),
		chromedp.Evaluate(jsTimeline, &timeline),
	)

	parcelData.Timeline = s.getTimelineData(timeline)

	if foundData == "" {
		chromedp.Stop()
		return nil, true
	}

	if err != nil {
		log.Fatal(err)
		return nil, true
	}

	return &parcelData, true
}

func (s *PostaHrScraper) getTimelineData(phrTimeline []postaHrTimelineEntry) *[]timelineEntry {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in PostaHrScraper, getTimelineData %s", r)
		}
	}()
	var parsedTimeline []timelineEntry

	for i, item := range phrTimeline {
		entry := timelineEntry{}

		entry.Description = item.Description
		//Add indices in reversed order
		entry.Index = int8(i)
		entry.Location = item.Location
		entry.Status = item.Status

		layout := "1/2/2006T15:04:05 PM"
		t, err := time.Parse(layout, item.Date+"T"+item.Time)
		entry.Time = t
		if err != nil {
			log.Println(err)
		}
		// parsedTimeline = append(parsedTimeline, entry)
		parsedTimeline = append([]timelineEntry{entry}, parsedTimeline...)
	}

	return &parsedTimeline
}
