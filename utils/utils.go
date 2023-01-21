package utils

import (
	"math"
)

var CbrReady int64 = 0
var KucionReady int64 = 0

type UsdBtcFilter struct {
	Id    int
	Date  int64
	Count int64
}

type FiatBtcFilter struct {
	Id         int   `json:"id"`
	UsdBtcDate int64 `json:"usd_btc_date"`
	FiatDate   int64 `json:"fiat_date"`
	Count      int64 `json:"count"`
}

type FiatFilter struct {
	Id       int   `json:"id"`
	FiatDate int64 `json:"fiat_date"`
	Count    int64 `json:"count"`
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / pow
}
