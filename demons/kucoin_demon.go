package demons

import (
	"github.com/StanDenisov/btc_usdt_check/queries"
	"github.com/StanDenisov/btc_usdt_check/service"
	"log"
)

func RunKucoinDaemon() {
	log.Println("RunKucionDeamon")
	kucionCurs, err := service.GetBtcUsdCurrency()
	if err != nil {
		log.Println(err)
	}
	queries.InsertBTCUSD(kucionCurs)
}
