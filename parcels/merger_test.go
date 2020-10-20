package parcels

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMerger(t *testing.T) {
	fmt.Fprintln(os.Stdout, "Started Merger tests...")

	time1 := time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC)
	time2 := time.Date(2018, time.November, 10, 23, 10, 0, 0, time.UTC)
	time3 := time.Date(2019, time.March, 7, 12, 0, 0, 0, time.UTC)
	time4 := time.Date(2019, time.April, 8, 8, 0, 0, 0, time.UTC)
	time5 := time.Date(2019, time.May, 6, 5, 10, 10, 10, time.UTC)
	time6 := time.Date(2020, time.June, 22, 22, 0, 0, 0, time.UTC)
	time7 := time.Date(2020, time.April, 12, 23, 0, 0, 0, time.UTC)
	time8 := time.Date(2020, time.May, 9, 23, 0, 0, 0, time.UTC)
	time9 := time.Date(2020, time.June, 29, 21, 0, 0, 0, time.UTC)
	time10 := time.Date(2020, time.June, 30, 14, 0, 0, 0, time.UTC)

	pd1 := ParcelData{
		Provider:          "GlobalCanaio",
		To:                &Address{Country: "Croatia"},
		LastUpdated:       "2020-02-02",
		From:              &Address{Country: "China"},
		Status:            "IN Transport",
		StatusDescription: "in transport to some location",
		TrackingNumber:    "123456789",
		Timeline: &Timeline{
			{Time: time2, Status: "processed"},
			{Time: time3, Status: "shipped"},
			{Time: time7, Status: "arrived at 1"},
			{Time: time8, Status: "arrived at destination"},
			{Time: time6, Status: "out for delivery"},
		},
	}

	pd2 := ParcelData{
		Provider:          "PostHr",
		To:                &Address{Country: "Croatia"},
		LastUpdated:       "2019-02-09",
		From:              &Address{Country: "China"},
		Status:            "SHIPPED",
		StatusDescription: "IN Transport",
		TrackingNumber:    "123456789",
		Timeline: &Timeline{
			{Time: time1, Status: "processed"},
			{Time: time4, Status: "added to box"},
			{Time: time5, Status: "arrived at 1", Location: &Address{Country: "Russia"}},
			{Time: time5, Status: "arrived at this desk", Location: &Address{Country: "Russia"}},
			{Time: time6, Status: "arrived at 1", Location: &Address{Country: "Russia"}},
			{Time: time9, Status: "arrived at destination", Location: &Address{Country: "Italy"}},
			{Time: time10, Status: "received", Location: &Address{Country: "Croatia"}},
		},
	}

	t.Logf("Parcel data 1 value: %+v", pd1)
	t.Logf("Parcel data 1 timeline value: %+v", *pd1.Timeline)

	t.Logf("Parcel data 2 value: %+v", pd2)
	t.Logf("Parcel data 2 timeline value: %+v", *pd2.Timeline)

	dataMerger := NewMerger()
	dataMerger.AddData(&pd1)
	dataMerger.AddData(&pd2)

	final, _ := dataMerger.GetMergedData()
	t.Logf("Merged data value: %+v", *final)
	t.Logf("Merged data timeline value: %+v", *final.Timeline)

	finalTimelineLength := (*final.Timeline).Len()

	if final.Timeline.Len() != (pd1.Timeline.Len() + pd2.Timeline.Len()) {
		t.Fatalf("Expected length of final timeline data: %d, actual: %d\n", (pd1.Timeline.Len() + pd2.Timeline.Len()), final.Timeline.Len())
	}

	if (*final.Timeline)[0].Index != int8(final.Timeline.Len()-1) {
		t.Fatalf("Expected timeline entry index: %d, actual: %d\n", (*final.Timeline)[0].Index, final.Timeline.Len()-1)
	}

	if (*final.Timeline)[0].Time != time10 {
		t.Fatalf("Expected first timeline entry with time: %v, actual: %v\n", time10, (*final.Timeline)[0].Time)
	}

	if (*final.Timeline)[finalTimelineLength-1].Time != time1 {
		t.Fatalf("Expected last timeline entry with time: %v, actual: %v\n", time1, (*final.Timeline)[finalTimelineLength-1].Time)
	}

	fmt.Fprintln(os.Stdout, "Finished Merger tests...")
}
