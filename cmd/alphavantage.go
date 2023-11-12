package cmd

import (
	"log"
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

// Search implements TickerProvider.
func (ticketProvider alphaVantageClient) Search(keywords string) ([]Ticker, error) {
	search := avFunction("SYMBOL_SEARCH")
	search = param(search, "keywords", keywords)

	var searchResponse AlphaVantageSymbolSearchResponse
	_, err := ticketProvider.get(search, &searchResponse)
	if err != nil {
		log.Println("Could not parse search response", err)
		return []Ticker{}, err
	}

	tickers := make([]Ticker, len(searchResponse.BestMatches))
	for idx, x := range searchResponse.BestMatches {
		tickers[idx] = x.asTicker()
	}
	return tickers, nil
}

func avFunction(f string) string {
	return avUrl + "function=" + f
}

func param(base string, paramType string, paramValue string) string {
	return base + "&" + paramType + "=" + paramValue
}

func (ticketProvider alphaVantageClient) get(url string, result interface{}) (*resty.Response, error) {
	url = param(url, "apikey", ticketProvider.apiKey)
	return ticketProvider.client.R().SetResult(result).Get(url)
}

// DailySeries returns daily indicators for the provided ticker symbol
func (ticketProvider alphaVantageClient) DailySeries(symbol string) (Series, error) {
	daily := avFunction("TIME_SERIES_DAILY")
	daily = param(daily, "symbol", symbol)

	var timeSeriesDaily AlphaVantageTimeSeriesDailyResponse
	_, err := ticketProvider.get(daily, &timeSeriesDaily)
	if err != nil {
		log.Println("Could not parse daily time series response", err)
		return Series{}, err
	}

	dataPoints := make([]DataPoint, len(timeSeriesDaily.TimeSeries))
	idx := 0
	for key, value := range timeSeriesDaily.TimeSeries {
		date, err := time.Parse("2006-01-02", key)
		if err != nil {
			log.Println("Could not parse time series key response", err)
			return Series{}, err
		}
		dataPoints[idx] = DataPoint{date, value.asCandle()}
		idx++
	}
	return Series{dataPoints}, nil
}

type AlphaVantageTimeSeriesDailyResponse struct {
	TimeSeries map[string]AlphaVantageCandle `json:"Time Series (Daily)"`
}

type AlphaVantageCandle struct {
	Open   float64 `json:"1. open,string"`
	High   float64 `json:"2. high,string"`
	Low    float64 `json:"3. low,string"`
	Close  float64 `json:"4. close,string"`
	Volume int64   `json:"5. volume,string"`
}

func (avc AlphaVantageCandle) asCandle() Candle {
	return Candle{
		Open:   avc.Open,
		High:   avc.High,
		Low:    avc.Low,
		Close:  avc.Close,
		Volume: avc.Volume,
	}
}

type AlphaVantageSymbolSearchResponse struct {
	BestMatches []AlphaVantageSymbol
}

type AlphaVantageSymbol struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	MarketOpen  string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
	MatchScore  string `json:"9. matchScore"`
}

func (avs AlphaVantageSymbol) asTicker() Ticker {
	return Ticker{
		Symbol:   avs.Symbol,
		Name:     avs.Name,
		Type:     avs.Type,
		Region:   avs.Region,
		Currency: avs.Currency,
	}
}
