package cmd

import (
	"math"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const threshold = 1e-9

func TestDailySeries(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()
	data :=
		`{
			"Meta Data": {
				"1. Information": "Daily Prices (open, high, low, close) and Volumes",
				"2. Symbol": "SNH.JOH",
				"3. Last Refreshed": "2020-05-29",
				"4. Output Size": "Compact",
				"5. Time Zone": "US/Eastern"
			},
			"Time Series (Daily)": {
				"2020-05-29": {
					"1. open": "109.0000",
					"2. high": "110.0000",
					"3. low": "101.0000",
					"4. close": "101.0000",
					"5. volume": "4583994"
				},
				"2020-05-28": {
					"1. open": "101.0000",
					"2. high": "108.0000",
					"3. low": "104.0000",
					"4. close": "107.0000",
					"5. volume": "2270535"
				}
			}
		}`
	httpmock.RegisterResponder("GET", `=~^https://www\.alphavantage\.co\/query.*$`, httpmock.NewStringResponder(200, data))
	oracle := NewTickerProvider(client, "test")
	series, _ := oracle.DailySeries("test")
	assertEquals(t, 2, len(series.Data))
	assertEquals(t, time.Date(2020, 05, 29, 0, 0, 0, 0, time.UTC), series.Data[0].Timestamp)
	assertClose(t, 101, series.Data[1].Indicators.Open)
	assertClose(t, 108, series.Data[1].Indicators.High)
	assertClose(t, 104, series.Data[1].Indicators.Low)
	assertClose(t, 107, series.Data[1].Indicators.Close)
	assertEquals(t, int64(2270535), series.Data[1].Indicators.Volume)
}

func assertClose(t *testing.T, expected, actual float64) {
	diff := math.Abs(actual - expected)
	if diff >= threshold {
		t.Errorf("Actual = %v, and Expected = %v difference %f greater than threshold ", actual, expected, diff)
	}
}

func assertEquals(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		t.Errorf("Error actual = %v (%T), and Expected = %v (%T).", actual, actual, expected, expected)
	}
}
