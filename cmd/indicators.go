package cmd

import (
	"fmt"
	"time"
)

// Candle describe historical market data
type Candle struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func (i Candle) String() string {
	return fmt.Sprintf("Open: %f\nHigh: %f\nLow: %f\nClose: %f\nVolume: %d\n", i.Open, i.High, i.Low, i.Close, i.Volume)
}

// DataPoint capture indicators for a particular period of time
type DataPoint struct {
	Timestamp time.Time
	Candle    Candle
}

func (dp DataPoint) String() string {
	return fmt.Sprintf("%s\n%v", dp.Timestamp.Format("2006-01-02"), dp.Candle)
}

// Series is a collection of non-overlapping data points
type Series struct {
	Data []DataPoint
}

type TickerProvider interface {
	DailySeries(symbol string) (Series, error)
	Search(keywords string) ([]Ticker, error)
}

type Ticker struct {
	Symbol   string
	Name     string
	Type     string
	Region   string
	Currency string
}
