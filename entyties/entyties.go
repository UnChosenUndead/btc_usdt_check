package entyties

type Fiat struct {
	Date      int64   `db:"fiat_date_id"`
	Name      string  `db:"fiat_name"`
	CharCode  string  `db:"fiat_char_code"`
	Nominal   int     `db:"fiat_nominal"`
	Value     float64 `db:"fiat_value"`
	OneRubBtc int     `db:"one_rub_btc_id"`
}

type BTCUSD struct {
	BtcDateId    int64   `db:"btc_date_id"`
	AveragePrice float64 `db:"average_price"`
}
type BtcUsdCount int32

type BTCUSDCourse struct {
	BtcDate      float64 `db:"btc_date_value"`
	AveragePrice float64 `db:"average_price"`
}

type BTCUSDCourseAll struct {
	Count   BtcUsdCount     `json:"count"`
	History []*BTCUSDCourse `json:"history"`
}

type FiatRubAll struct {
	Count         FiatRubCount `json:"total"`
	FiatRubSelect []*FiatRub   `json:"history"`
}

type FiatRub struct {
	FiatNominal  int64   `db:"fiat_nominal"`
	FiatValue    float64 `db:"fiat_value"`
	FiatDate     int64   `db:"fiat_date_value"`
	FiatCharCode string  `db:"fiat_char_code"`
}

type FiatRubCount int32

type FiatBTCAll struct {
	Count      FiatBTCCount     `db:"count" json:"total"`
	FiatSelect []*FiatBTCCourse `json:"history"`
}

type FiatBTCCount int32

type FiatBTCCourse struct {
	BtcDate         int64   `db:"btc_date_value"`
	FiatDate        int64   `db:"fiat_date_value"`
	FiatBtcSumValue float64 `db:"fiat_btc_sum_value"`
	FiatCharCode    string  `db:"fiat_char_code"`
}
