package parcels

import (
	"fmt"
	"sort"
)

// Merger struct
type Merger struct {
	dataPool   []*ParcelData
	mergedData ParcelData
}

// NewMerger creates a new Merger.
func NewMerger() *Merger {
	dataPool := []*ParcelData{}
	mergedData := ParcelData{}
	return &Merger{dataPool, mergedData}
}

// AddData func
func (m *Merger) AddData(data *ParcelData) (bool, error) {
	m.dataPool = append(m.dataPool, data)
	return true, nil
}

// GetFinalData func
func (m *Merger) GetFinalData() (*ParcelData, error) {
	var parsedTimeline Timeline

	for _, x := range m.dataPool {
		m.addInfo(x)

		for _, y := range *x.Timeline {
			parsedTimeline = append(parsedTimeline, y)
		}
	}

	sort.Sort(sort.Reverse(parsedTimeline))
	m.mergedData.Timeline = &parsedTimeline

	// add indices to timeline
	length := m.mergedData.Timeline.Len()
	for i := range *m.mergedData.Timeline {
		(*m.mergedData.Timeline)[i].Index = int8(length - i - 1)
	}

	return &m.mergedData, nil
}

func (m *Merger) addInfo(x *ParcelData) {
	// TODO: refactor when go introduces generics in 1.17 version
	if x.Courier != "" {
		m.mergedData.Courier = m.mergedData.Courier + ", " + x.Courier
	}
	if x.From != nil {
		m.mergedData.From = x.From
	}

	if x.LastUpdated != "" {
		m.mergedData.LastUpdated = m.mergedData.LastUpdated + ", " + x.LastUpdated
	}

	if x.Provider != "" {
		m.mergedData.Provider = fmt.Sprintf("%v, %v", m.mergedData.Provider, x.Provider)
		// m.mergedData.Provider = m.mergedData.Provider + ", " + x.Provider
	}

	if x.ShippingDaysCount != 0 {
		m.mergedData.ShippingDaysCount = x.ShippingDaysCount
	}

	if x.Status != "" {
		m.mergedData.Status = m.mergedData.Status + ", " + x.Status
	}

	if x.StatusDescription != "" {
		m.mergedData.StatusDescription = m.mergedData.StatusDescription + ", " + x.StatusDescription
	}

	if x.To != nil {
		m.mergedData.To = x.To
	}

	if x.TrackingNumber != "" {
		m.mergedData.TrackingNumber = x.TrackingNumber
	}
}
