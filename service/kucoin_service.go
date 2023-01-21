package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type KucoinResponse struct {
	Data struct {
		Time         int64  `json:"time"`
		AveragePrice string `json:"averagePrice"`
	} `json:"data"`
}

func GetBtcUsdCurrency() (KucoinResponse, error) {
	resp, err := http.Get(os.Getenv("BTCUSDUrl"))
	kucoinResponse := KucoinResponse{}
	if err != nil {
		log.Printf("http get: %s", err)
		return kucoinResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("body close: %s", err)
		return kucoinResponse, err
	}
	err = json.NewDecoder(resp.Body).Decode(&kucoinResponse)
	if err != nil {
		log.Printf("json decoder: %s", err)
		return kucoinResponse, err
	}
	return kucoinResponse, nil
}
