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
			fmt.Printf("Panic in PostaHrScraper, jsGetTimeline %s", r)
		}
	}()
	buf, _ := ioutil.ReadFile("helpers/postaHr.js")
	funcJS := string(buf)
	invokeFuncJS := `var a = getItems("` + sel + `"); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// GetData func
func (s *PostaHrScraper) GetData(trackingNumber string) (*parcels.ParcelData, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in PostaHrScraper, GetData %s", r)
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

	parcelData := parcels.ParcelData{
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
		chromedp.Evaluate(jsTimeline, &parcelData.Timeline),
	)

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
