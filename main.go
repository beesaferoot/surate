package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Exchange struct {
	USD int     `json:"USD"`
	NGN float64 `json:"NGN"`
}

type RateResponse struct {
	CBN               Exchange `json:"cbn,omitempty"`
	CoinMarketCapRate Exchange `json:"coinmarketcap,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey, exists := os.LookupEnv("COIN_MARKET_CAP_API_KEY")
	if !exists {
		log.Println("Environment variable COIN_MARKET_CAP_API_KEY not set")
		os.Exit(1)
	}

	coinMarketCapAPIUrl := os.Getenv("COIN_MARKET_CAP_API_URL")

	if len(coinMarketCapAPIUrl) == 0 {
		coinMarketCapAPIUrl = "https://sandbox-api.coinmarketcap.com/v2/tools/price-conversion"
	}

	http.HandleFunc("/api/myrate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			var resp RateResponse
			resp, err := loadFromDisk("usd_rate.gob")
			if err != nil {

				// fetch cbn rate
				cbnRate, err := FetchCBNRate()
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
					return
				}

				coinCapRate, err := FetchUSDRate(apiKey, coinMarketCapAPIUrl)

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
					return
				}

				resp = RateResponse{
					CBN: Exchange{
						USD: 1,
						NGN: cbnRate.buyRate,
					},
					CoinMarketCapRate: Exchange{
						USD: 1,
						NGN: coinCapRate.rate,
					},
				}

				err = saveToDisk(resp, "usd_rate.gob")

				if err != nil {
					log.Printf("error saving rate data to disk: %s\n", err.Error())
				} else {
					log.Println("saved rate data to disk")
				}

			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Println(err)
				// Handle error if JSON encoding fails
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			}
		} else {
			// If not GET, return a method not allowed status
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// serve client dist
	fs := http.FileServer(http.Dir("./client/dist"))
	http.Handle("/", fs)

	// Start the server on port 5000
	fmt.Println("Starting server on :5008")
	if err := http.ListenAndServe(":5008", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
