package parcels

import (
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

// GetMergedData func
func (m *Merger) GetMergedData() (*ParcelData, error) {
	var mergedTimeline Timeline

	for _, x := range m.dataPool {
		m.addInfo(x)

		for _, y := range *x.Timeline {
			mergedTimeline = append(mergedTimeline, y)
		}
	}

	sort.Sort(sort.Reverse(mergedTimeline))

	// add indices to timeline
	length := mergedTimeline.Len()
	for i := range mergedTimeline {
		mergedTimeline[i].Index = int8(length - i - 1)
	}

	m.mergedData.Timeline = &mergedTimeline

	return &m.mergedData, nil
}

func (m *Merger) addInfo(x *ParcelData) {
	// TODO: refactor when go introduces generics in 1.17 version
	if x.Courier != "" {
		if m.mergedData.Courier != "" {
			m.mergedData.Courier = m.mergedData.Courier + ", " + x.Courier
		} else {
			m.mergedData.Courier = x.Courier
		}
	}

	if x.From != nil && m.mergedData.From == nil {
		m.mergedData.From = x.From
	}

	if x.LastUpdated != "" {
		if m.mergedData.LastUpdated != "" {
			m.mergedData.LastUpdated = m.mergedData.LastUpdated + ", " + x.LastUpdated
		} else {
			m.mergedData.LastUpdated = x.LastUpdated
		}
	}

	if x.Provider != "" {
		if m.mergedData.Provider != "" {
			m.mergedData.Provider = m.mergedData.Provider + ", " + x.Provider
		} else {
			m.mergedData.Provider = x.Provider
		}
	}

	if x.ShippingDaysCount != 0 && m.mergedData.ShippingDaysCount == 0 {
		m.mergedData.ShippingDaysCount = x.ShippingDaysCount
	}

	if x.Status != "" {
		if m.mergedData.Status != "" {
			m.mergedData.Status = m.mergedData.Status + ", " + x.Status
		} else {
			m.mergedData.Status = x.Status
		}
	}

	if x.StatusDescription != "" {
		if m.mergedData.StatusDescription != "" {
			m.mergedData.StatusDescription = m.mergedData.StatusDescription + ", " + x.StatusDescription
		} else {
			m.mergedData.StatusDescription = x.StatusDescription
		}
	}

	if x.To != nil && m.mergedData.To == nil {
		m.mergedData.To = x.To
	}

	if x.TrackingNumber != "" {
		if m.mergedData.TrackingNumber != "" {
			m.mergedData.TrackingNumber = m.mergedData.TrackingNumber + ", " + x.TrackingNumber
		} else {
			m.mergedData.TrackingNumber = x.TrackingNumber
		}
	}
}
