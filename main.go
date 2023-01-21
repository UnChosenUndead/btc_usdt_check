package main

import (
	"fmt"
	"github.com/StanDenisov/btc_usdt_check/db"
	"github.com/StanDenisov/btc_usdt_check/demons"
	"github.com/StanDenisov/btc_usdt_check/router"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)
	//init before start actions
	godotenv.Load()
	db.InitConnection()

	//start calculate before server run
	demons.RunCbrDemon()
	demons.RunKucoinDaemon()
	demons.RunSumDemon()

	//start cron scheduler
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().Do(func() { demons.RunCbrDemon() })
	s.Every(10).Second().Do(func() { demons.RunKucoinDaemon() })
	s.Every(12).Second().Do(func() { demons.RunSumDemon() })
	s.StartAsync()
	//gin
	r := gin.Default()
	router.Router(r)
	address := fmt.Sprintf("%s:%s", os.Getenv("APPHOST"), os.Getenv("APPPORT"))

	r.Run(address) // listen and serve on 0.0.0.0:8080
}
