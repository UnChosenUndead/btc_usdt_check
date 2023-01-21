package queries

import (
	"context"
	"github.com/StanDenisov/btc_usdt_check/db"
	"github.com/StanDenisov/btc_usdt_check/entyties"
	"github.com/StanDenisov/btc_usdt_check/service"
	"github.com/StanDenisov/btc_usdt_check/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"golang.org/x/exp/slices"
	"log"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

func InsertBTCUSD(kucionResponse service.KucoinResponse) {
	btcUsdDataId := checkUniqueBTCValueBeforeInsert(kucionResponse)
	parseAverage, err := strconv.ParseFloat(kucionResponse.Data.AveragePrice, 64)
	if err != nil {
		log.Fatal(err)
	}
	roundedAveragePrice := utils.Round(parseAverage, 4)
	if btcUsdDataId != 0 {
		row := db.Conn.QueryRow(ctx,
			`INSERT INTO btc_usd(average_price, btc_date_id) VALUES ($1, $2) RETURNING id`,
			roundedAveragePrice, btcUsdDataId)
		var id int64
		err := row.Scan(&id)
		if err != nil {
			log.Printf("Unable to INSERT: %v\n", err)
			return
		}
	}
	utils.KucionReady = btcUsdDataId
	log.Printf("From insert btc_usd KucionReady now is %d", utils.KucionReady)
}

func checkUniqueBTCValueBeforeInsert(kucionResponse service.KucoinResponse) int64 {
	parseAverage, err := strconv.ParseFloat(kucionResponse.Data.AveragePrice, 64)
	if err != nil {
		log.Fatal(err)
	}
	roundedAveragePrice := utils.Round(parseAverage, 4)
	row := db.Conn.QueryRow(ctx,
		`SELECT id FROM btc_usd WHERE average_price = $1`, roundedAveragePrice)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return insertBtcUsdData(kucionResponse.Data.Time)
	}
	return 0
}

func insertBtcUsdData(data int64) int64 {
	row := db.Conn.QueryRow(ctx,
		`INSERT INTO btc_usd_date(btc_date_value) VALUES ($1) RETURNING ID`, data)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("cant insert btc_usd_date : %s", err)
	}
	return id
}

func InsertFiat(valCurs service.ValCurs) {
	cbrDate, err := time.Parse("02.01.2006", valCurs.Date)
	if err != nil {
		log.Println(err.Error())
	}
	fiatDateId := checkUniqueFiatDataBeforeInsert(cbrDate)
	if fiatDateId != 0 {
		var fiatArray []entyties.Fiat
		for _, fiat := range valCurs.Valute {
			nominal, err := strconv.Atoi(fiat.Nominal)
			if err != nil {
				log.Fatal(err.Error())
			}
			s := strings.Replace(fiat.Value, ",", ".", -1)
			log.Printf("this is fiat value after parse: %s", s)
			if err != nil {
				log.Fatal(err.Error())
			}
			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatal("german fail", err)
			}
			x := entyties.Fiat{Date: fiatDateId, Name: fiat.Name, CharCode: fiat.CharCode, Nominal: nominal, Value: value}
			fiatArray = append(fiatArray, x)
		}
		for _, fiat := range fiatArray {
			row := db.Conn.QueryRow(ctx,
				`INSERT INTO fiat(fiat_date_id, fiat_name, fiat_char_code, fiat_nominal, fiat_value) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				fiat.Date, fiat.Name, fiat.CharCode, fiat.Nominal, fiat.Value)
			var id int
			err := row.Scan(&id)
			if err != nil {
				log.Fatalf("Error insert fiat %s", err)
			}
		}
		utils.CbrReady = fiatDateId
		log.Printf("From insert Fiat CbrReady now is %d", utils.CbrReady)
	}
}

func checkUniqueFiatDataBeforeInsert(date time.Time) int64 {
	var fiatTimeStamp = date.Unix()
	row := db.Conn.QueryRow(ctx,
		`SELECT id FROM fiat_date WHERE fiat_date_value = $1`, fiatTimeStamp)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return insertFiatDate(fiatTimeStamp)
	}
	return 0
}

func insertFiatDate(data int64) int64 {
	row := db.Conn.QueryRow(ctx,
		`INSERT INTO fiat_date(fiat_date_value) VALUES ($1) RETURNING ID`, data)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("cant insert fiat_date : %s", err)
	}
	return id
}

func InsertFiatBTC() {
	var fiatArray []*entyties.Fiat
	var btcUsd []*entyties.BTCUSD
	if utils.CbrReady == 0 {
		fiatArray = SelectFiatByLastDate()
	} else {
		fiatArray = SelectFiatByDate(utils.CbrReady)
	}
	if utils.CbrReady == 0 {
		btcUsd = SelectBtcUsdByDate(utils.KucionReady)
	} else {
		btcUsd = SelectBtcUsdByLastDate()
	}
	log.Printf("btcUsd is: %+v\\n", btcUsd)
	log.Printf("Fiat array len is: %d", len(fiatArray))
	rubFiatNum := slices.IndexFunc(fiatArray, func(c *entyties.Fiat) bool { return c.CharCode == "USD" })
	btcUsdU := btcUsd[0]
	rubBtcSum := InsertRubBtcAndReturnSum(fiatArray[rubFiatNum], btcUsdU)
	fiatArray = RemoveIndex(fiatArray, rubFiatNum)
	for _, fiat := range fiatArray {
		btcUsdU := btcUsd[0]
		fiatBtcRoundedSum := utils.Round(rubBtcSum/fiat.Value*float64(fiat.Nominal), 4)
		log.Printf("fiatBtcSum is %f", fiatBtcRoundedSum)
		row := db.Conn.QueryRow(ctx,
			`INSERT INTO fiat_btc(fiat_date_id, btc_usd_date_id, fiat_char_code, fiat_btc_sum_value) VALUES ($1, $2, $3, $4) RETURNING ID`,
			fiat.Date, btcUsdU.BtcDateId, fiat.CharCode, fiatBtcRoundedSum)
		var id int64
		err := row.Scan(&id)
		if err != nil {
			log.Fatalf("cant insert fiat_btc : %s", err)
		}
	}
	utils.CbrReady = 0
	utils.KucionReady = 0
}

func RemoveIndex(s []*entyties.Fiat, index int) []*entyties.Fiat {
	return append(s[:index], s[index+1:]...)
}

func InsertRubBtcAndReturnSum(fiat *entyties.Fiat, btcusd *entyties.BTCUSD) float64 {
	fiatBtcSum := utils.Round(fiat.Value*btcusd.AveragePrice, 4)
	row := db.Conn.QueryRow(ctx,
		`INSERT INTO fiat_btc(fiat_date_id, btc_usd_date_id, fiat_char_code, fiat_btc_sum_value) VALUES ($1, $2, $3, $4) RETURNING ID`,
		fiat.Date, btcusd.BtcDateId, "RUB", fiatBtcSum)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("cant insert fiat_btc : %s", err)
	}
	return fiatBtcSum
}

func SelectFiatByDate(fiatDateId int64) []*entyties.Fiat {
	var fiat []*entyties.Fiat
	err := pgxscan.Select(ctx, db.Conn, &fiat, `SELECT fiat_date_id, fiat_char_code, fiat_nominal, fiat_name, fiat_value FROM fiat WHERE fiat_date_id = $1;`, fiatDateId)
	if err != nil {
		log.Fatalf("selecting err is Fiat %s", err)
	}
	return fiat
}

func SelectFiatByLastDate() []*entyties.Fiat {
	var fiat []*entyties.Fiat
	err := pgxscan.Select(ctx, db.Conn, &fiat, `SELECT fiat_date_id, fiat_char_code, fiat_nominal, fiat_name, fiat_value FROM fiat WHERE fiat_date_id =(select max(id) from fiat_date);`)
	if err != nil {
		log.Fatalf("selecting err is Fiat %s", err)
	}
	return fiat
}

func SelectBtcUsdByDate(btcDate int64) []*entyties.BTCUSD {
	var btcUsd []*entyties.BTCUSD
	err := pgxscan.Select(ctx, db.Conn, &btcUsd, `SELECT btc_date_id, average_price FROM btc_usd WHERE btc_date_id = $1;`, btcDate)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return btcUsd
}

func SelectBtcUsdByLastDate() []*entyties.BTCUSD {
	var btcUsd []*entyties.BTCUSD
	err := pgxscan.Select(ctx, db.Conn, &btcUsd, `SELECT btc_date_id, average_price FROM btc_usd WHERE btc_date_id = (select max(id) from btc_usd_date);`)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return btcUsd
}

func SelectLastBtcUsdCourse() entyties.BTCUSDCourseAll {
	var btcUsdCount entyties.BtcUsdCount
	countQuery := `
		SELECT count(*)
		FROM btc_usd
		INNER JOIN btc_usd_date as bud
					ON  btc_usd.btc_date_id = bud.id
		 and bud.id = (
                            SELECT MAX(id)
                            FROM btc_usd_date
        );					
`
	row, _ := db.Conn.Query(ctx, countQuery)
	err := pgxscan.ScanOne(&btcUsdCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}
	var btcUsdtAll []*entyties.BTCUSDCourse
	query := `
		SELECT bud.btc_date_value, average_price
		FROM btc_usd
         INNER JOIN btc_usd_date as bud
                    ON  btc_usd.btc_date_id = bud.id
         and bud.id = (
                            SELECT MAX(id)
                            FROM btc_usd_date
        );
	`
	err = pgxscan.Select(ctx, db.Conn, &btcUsdtAll, query)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return entyties.BTCUSDCourseAll{Count: btcUsdCount, History: btcUsdtAll}
}

func SelectFilteredBtcUsdCourse(filter utils.UsdBtcFilter) entyties.BTCUSDCourseAll {
	var btcUsdCount entyties.BtcUsdCount
	countQuery := `
		SELECT count(*)
		FROM btc_usd
		INNER JOIN btc_usd_date as bud
					ON  btc_usd.btc_date_id = bud.id
		 and bud.btc_date_value = $1;				
`
	row, _ := db.Conn.Query(ctx, countQuery, filter.Date)
	err := pgxscan.ScanOne(&btcUsdCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}
	var btcUsdtAll []*entyties.BTCUSDCourse
	if filter.Id == 0 {
		filter.Id = 1
	}
	query := `
		SELECT bud.btc_date_value, average_price
		FROM btc_usd
         INNER JOIN btc_usd_date as bud
                    ON  btc_usd.btc_date_id = bud.id
                    and bud.btc_date_value = $1
         WHERE btc_usd.id >= $2
         ORDER BY btc_usd.id ASC
		 LIMIT $3;
	`
	err = pgxscan.Select(ctx, db.Conn, &btcUsdtAll, query, filter.Date, filter.Id, filter.Count)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return entyties.BTCUSDCourseAll{Count: btcUsdCount, History: btcUsdtAll}
}

func SelectLastFiatRubCourse() entyties.FiatRubAll {
	var fiatRubCount entyties.FiatRubCount
	countQuery := `
		SELECT count(*)
		FROM fiat
		INNER JOIN fiat_date as fd
					ON  fiat.fiat_date_id = fd.id
		and fd.id = (
                            SELECT MAX(id)
                            FROM fiat_date
        );			
`
	row, _ := db.Conn.Query(ctx, countQuery)
	err := pgxscan.ScanOne(&fiatRubCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}
	var fiatRub []*entyties.FiatRub
	query := `
		SELECT fd.fiat_date_value, fiat_value, fiat_char_code, fiat_nominal
		FROM fiat
         INNER JOIN fiat_date as fd
                    ON  fiat.fiat_date_id = fd.id
         and fd.id = (
                            SELECT MAX(id)
                            FROM fiat_date
        );
	`
	err = pgxscan.Select(ctx, db.Conn, &fiatRub, query)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return entyties.FiatRubAll{Count: fiatRubCount, FiatRubSelect: fiatRub}
}

func SelectFilteredFiatRubCourse(filter utils.FiatFilter) entyties.FiatRubAll {
	var fiatRubCount entyties.FiatRubCount
	countQuery := `
		SELECT count(*)
		FROM fiat
		INNER JOIN fiat_date as fd
					ON  fiat.fiat_date_id = fd.id
		 and fd.fiat_date_value = $1;				
`
	row, _ := db.Conn.Query(ctx, countQuery, filter.FiatDate)
	err := pgxscan.ScanOne(&fiatRubCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}

	var fiatRub []*entyties.FiatRub
	query := `
		SELECT fd.fiat_date_value, fiat_value, fiat_char_code, fiat_nominal
		FROM fiat
         INNER JOIN fiat_date as fd
                    ON  fiat.fiat_date_id = fd.id
         			AND fd.fiat_date_value = $1
		 WHERE fiat.id >= $2
         ORDER BY fiat.id ASC
		 LIMIT $3;
	`
	err = pgxscan.Select(ctx, db.Conn, &fiatRub, query, filter.FiatDate, filter.Id, filter.Count)
	if err != nil {
		log.Fatalf("selecting err is btcUsdByDate %s", err)
	}
	return entyties.FiatRubAll{Count: fiatRubCount, FiatRubSelect: fiatRub}
}

func SelectLastFiatBtcCourse() entyties.FiatBTCAll {
	var fiatBtcCount entyties.FiatBTCCount
	countQuery := `
		SELECT count(*)
		FROM fiat_btc
         INNER JOIN btc_usd_date as bud
                    ON  fiat_btc.btc_usd_date_id = bud.id
         and bud.id = (
                            SELECT MAX(id)
                            FROM btc_usd_date
        )
		INNER JOIN fiat_date as fd
					ON  fiat_btc.fiat_date_id = fd.id
		 and fd.id = (
                            SELECT MAX(id)
                            FROM fiat_date
        );				
`
	row, _ := db.Conn.Query(ctx, countQuery)
	err := pgxscan.ScanOne(&fiatBtcCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}

	var fiatBtcSelect []*entyties.FiatBTCCourse
	query := `
		SELECT bud.btc_date_value, fd.fiat_date_value, fiat_btc_sum_value, fiat_char_code
		FROM fiat_btc
         INNER JOIN btc_usd_date as bud
                    ON  fiat_btc.btc_usd_date_id = bud.id
         and bud.id = (
                            SELECT MAX(id)
                            FROM btc_usd_date
        )
		INNER JOIN fiat_date as fd
					ON  fiat_btc.fiat_date_id = fd.id
		 and fd.id = (
                            SELECT MAX(id)
                            FROM fiat_date
        );
	`
	err = pgxscan.Select(ctx, db.Conn, &fiatBtcSelect, query)
	if err != nil {
		log.Fatalf("selecting err is btcFiatAll %s", err)
	}
	return entyties.FiatBTCAll{Count: fiatBtcCount, FiatSelect: fiatBtcSelect}
}

func SelectFilteredFiatBtcCourse(filter utils.FiatBtcFilter) entyties.FiatBTCAll {
	log.Println(filter)
	var fiatBtcCount entyties.FiatBTCCount
	countQuery := `
		SELECT count(*)
		FROM fiat_btc
         INNER JOIN btc_usd_date as bud
                    ON  fiat_btc.btc_usd_date_id = bud.id
         and bud.btc_date_value = $1
		INNER JOIN fiat_date as fd
					ON  fiat_btc.fiat_date_id = fd.id
		 and fd.fiat_date_value = $2;				
`
	row, _ := db.Conn.Query(ctx, countQuery, filter.UsdBtcDate, filter.FiatDate)
	err := pgxscan.ScanOne(&fiatBtcCount, row)
	if err != nil {
		log.Fatalf("selecting err is btcFiatCount %s", err)
	}

	var fiatBtcSelect []*entyties.FiatBTCCourse
	query := `
		SELECT bud.btc_date_value, fd.fiat_date_value, fiat_btc_sum_value, fiat_char_code
		FROM fiat_btc
         INNER JOIN btc_usd_date as bud
                    ON  fiat_btc.btc_usd_date_id = bud.id
         and bud.btc_date_value = $1
		INNER JOIN fiat_date as fd
					ON  fiat_btc.fiat_date_id = fd.id
		 and fd.fiat_date_value = $2
		 WHERE fiat_btc.id >= $3
		 LIMIT $4;
	`
	err = pgxscan.Select(ctx, db.Conn, &fiatBtcSelect, query, filter.UsdBtcDate, filter.FiatDate, filter.Id, filter.Count)
	if err != nil {
		log.Fatalf("selecting err is btcFiatAll %s", err)
	}
	return entyties.FiatBTCAll{Count: fiatBtcCount, FiatSelect: fiatBtcSelect}
}

func SelectFiatBtcByCharCode(charCode string) []*entyties.FiatBTCCourse {
	log.Printf("charCode is %s", charCode)
	var fiatBtcSelect []*entyties.FiatBTCCourse
	query := `
		SELECT bud.btc_date_value, fd.fiat_date_value, fiat_btc_sum_value, fiat_char_code
		FROM fiat_btc
         INNER JOIN btc_usd_date as bud
                    ON  fiat_btc.btc_usd_date_id = bud.id
         and bud.id = (
                            SELECT MAX(id)
                            FROM btc_usd_date
        )
		INNER JOIN fiat_date as fd
					ON  fiat_btc.fiat_date_id = fd.id
		 and fd.id = (
                            SELECT MAX(id)
                            FROM fiat_date
        )
		WHERE fiat_char_code = $1;
	`
	err := pgxscan.Select(ctx, db.Conn, &fiatBtcSelect, query, charCode)
	if err != nil {
		log.Fatalf("selecting err is btcFiatAll %s", err)
	}
	log.Printf("%s", fiatBtcSelect)
	return fiatBtcSelect
}
