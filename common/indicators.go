package common

import (
	"fmt"
	"time"
)

// Indicators describe historical market data
type Indicators struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func (i Indicators) String() string {
	return fmt.Sprintf("Open: %f\nHigh: %f\nLow: %f\nClose: %f\nVolume: %d\n", i.Open, i.High, i.Low, i.Close, i.Volume)
}

// DataPoint capture indicators for a particular period of time
type DataPoint struct {
	Timestamp  time.Time
	Indicators Indicators
}

func (dp DataPoint) String() string {
	return fmt.Sprintf("%s\n%v", dp.Timestamp.Format("2006-01-02"), dp.Indicators)
}

// Series is a collection of non-overlapping data points
type Series struct {
	Data []DataPoint
}

// Oracle is a source of market related information
type Oracle interface {
	DailySeries(symbol string) (Series, error)
}
