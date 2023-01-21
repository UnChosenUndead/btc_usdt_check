package demons

import (
	"github.com/StanDenisov/btc_usdt_check/queries"
	"github.com/StanDenisov/btc_usdt_check/utils"
	"log"
)

func RunSumDemon() {
	log.Println("RunSumDemon")
	log.Printf("Cbr ready now is: %d", utils.CbrReady)
	log.Printf("Kucion ready now is: %d", utils.KucionReady)
	if utils.CbrReady != 0 || utils.KucionReady != 0 {
		queries.InsertFiatBTC()
		log.Println("CbrReady and KucionReady now 0")
	}
}
