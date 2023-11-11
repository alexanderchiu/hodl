package cmd

import (
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

var avUrl = "https://www.alphavantage.co/query?"

// NewTickerProvider creates an instace of a TickerProvider
// AlphaVantage API: https://www.alphavantage.co/documentation/
func NewTickerProvider(client *resty.Client, apiKey string) TickerProvider {
	return alphaVantageClient{client, apiKey}
}

type alphaVantageClient struct {
	client *resty.Client
	apiKey string
}

func avFunction(f string) string {
	return avUrl + "function=" + f
}

func param(base string, paramType string, paramValue string) string {
	return base + "&" + paramType + "=" + paramValue
}

// DailySeries returns daily indicators for the provided ticker symbol
func (mo alphaVantageClient) DailySeries(symbol string) (Series, error) {
	daily := avFunction("TIME_SERIES_DAILY")
	daily = param(daily, "symbol", symbol)
	daily = param(daily, "apikey", mo.apiKey)
	resp, err := mo.client.R().Get(daily)
	if err != nil {
		return Series{}, err
	}
	var result map[string]interface{}
	json.Unmarshal(resp.Body(), &result)
	timeSeries := result["Time Series (Daily)"].(map[string]interface{})
	data, err := unmarshal(timeSeries)
	if err != nil {
		return Series{}, err
	}
	return Series{data}, nil
}

func parseFloat(fields map[string]interface{}, key string) (float64, error) {
	return strconv.ParseFloat(fields[key].(string), 64)
}

func parseInt(fields map[string]interface{}, key string) (int64, error) {
	return strconv.ParseInt(fields[key].(string), 0, 64)
}

func unmarshal(timeSeries map[string]interface{}) ([]DataPoint, error) {
	var data []DataPoint
	for key, value := range timeSeries {
		d := value.(map[string]interface{})
		open, err := parseFloat(d, "1. open")
		if err != nil {
			return nil, err
		}
		high, err := parseFloat(d, "2. high")
		if err != nil {
			return nil, err
		}
		low, err := parseFloat(d, "3. low")
		if err != nil {
			return nil, err
		}
		close, err := parseFloat(d, "4. close")
		if err != nil {
			return nil, err
		}
		volume, err := parseInt(d, "5. volume")
		if err != nil {
			return nil, err
		}
		indicators := Indicators{open, high, low, close, volume}
		date, err := time.Parse("2006-01-02", key)
		if err != nil {
			return nil, err
		}
		data = append(data, DataPoint{date, indicators})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.After(data[j].Timestamp)
	})
	return data, nil
}
