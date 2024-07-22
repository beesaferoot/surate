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

func FetchUSDRate(apiKey string, apiUrl string) (*CoinMarketCapRate, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := url.Values{}

	// amount=1&convert_id=2819&id=2781
	q.Add("amount", "1")
	q.Add("id", "2781")
	q.Add("convert_id", "2819")

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
	currencyRate := data["2781"].(map[string]interface{})
	quote := currencyRate["quote"].(map[string]interface{})
	nGN := quote["2819"].(map[string]interface{})
	nGPrice := nGN["price"].(float64)

	return &CoinMarketCapRate{
		rate:     nGPrice,
		currency: "NGN",
	}, nil

}
