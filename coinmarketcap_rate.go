package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type CoinMarketCapRate struct {
	currency string
	rate     float64
}

func FetchUSDRate(apiKey string) (*CoinMarketCapRate, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sandbox-api.coinmarketcap.com/v2/tools/price-conversion", nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := url.Values{}

	q.Add("amount", "1")
	q.Add("symbol", "NGN")
	q.Add("convert", "USD,NGN")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to server")
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		log.Println("error fetching USD rate")
		return nil, errors.New("failed to fetch USD rate")
	}
	respBody, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}

	err = json.Unmarshal(respBody, &result)

	if err != nil {
		return nil, err
	}

	data := result["data"].(map[string]interface{})
	currencyRate := data["NGN"].(map[string]interface{})
	quote := currencyRate["quote"].(map[string]interface{})
	nGN := quote["NGN"].(map[string]interface{})
	nGPrice := nGN["price"].(float64)

	return &CoinMarketCapRate{
		rate:     nGPrice,
		currency: "NGN",
	}, nil

}