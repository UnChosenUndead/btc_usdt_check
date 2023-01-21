package demons

import (
	"github.com/StanDenisov/btc_usdt_check/queries"
	"github.com/StanDenisov/btc_usdt_check/service"
	"log"
)

func RunCbrDemon() {
	log.Println("RunCbrDemon")
	valCurs, err := service.GetCbrCurrency()
	if err != nil {
		log.Println(err)
	}
	queries.InsertFiat(valCurs)
}
