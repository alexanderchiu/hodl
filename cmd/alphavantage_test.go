package cmd

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const threshold = 1e-9

func TestSearch(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	var response = `{
    "bestMatches": [
        {
            "1. symbol": "TSCO.LON",
            "2. name": "Tesco PLC",
            "3. type": "Equity",
            "4. region": "United Kingdom",
            "5. marketOpen": "08:00",
            "6. marketClose": "16:30",
            "7. timezone": "UTC+01",
            "8. currency": "GBX",
            "9. matchScore": "0.7273"
        },
        {
            "1. symbol": "TSCDF",
            "2. name": "Tesco plc",
            "3. type": "Equity",
            "4. region": "United States",
            "5. marketOpen": "09:30",
            "6. marketClose": "16:00",
            "7. timezone": "UTC-04",
            "8. currency": "USD",
            "9. matchScore": "0.7143"
        }
    ]
}`
	url := "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=tesco&apikey=test"
	addJsonResponse("GET", url, response, http.StatusOK)

	tickerProvider := NewTickerProvider(client, "test")
	searchResults, _ := tickerProvider.Search("tesco")
	assert.Equal(t, 2, len(searchResults), "The number of search results should be correct")

	expected := Ticker{"TSCO.LON", "Tesco PLC", "Equity", "United Kingdom", "GBX"}
	assert.Equal(t, expected, searchResults[0])
}

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
	addJsonResponse("GET", `=~^https://www\.alphavantage\.co\/query.*$`, data, http.StatusOK)
	ticker := NewTickerProvider(client, "test")
	series, _ := ticker.DailySeries("test")
	assert.Equal(t, 2, len(series.Data))
	assert.Equal(t, time.Date(2020, 05, 29, 0, 0, 0, 0, time.UTC), series.Data[0].Timestamp)
	assert.EqualValues(t, 101, series.Data[1].Indicators.Open)
	assert.EqualValues(t, 108, series.Data[1].Indicators.High)
	assert.EqualValues(t, 104, series.Data[1].Indicators.Low)
	assert.EqualValues(t, 107, series.Data[1].Indicators.Close)
	assert.EqualValues(t, 2270535, series.Data[1].Indicators.Volume)
}

func addJsonResponse(method string, url string, json string, status int) {
	httpmock.RegisterResponder(method, url, func(r *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, json)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})
}
