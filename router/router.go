package router

import (
	"github.com/StanDenisov/btc_usdt_check/controller"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	api := r.Group("api")
	api.GET("/btcusdt", controller.GetLastCourseOfBtcUsdt)
	api.POST("/btcusdt", controller.GetCourseOfBtcUsdtWithFilter)
	api.GET("/currencies", controller.GetLastCourseFiatRub)
	api.POST("/currencies", controller.GetCourseFilteredFiatRub)
	api.GET("/latest", controller.GetLastCoursesFiatBtc)
	api.POST("/latest", controller.GetCourseFilteredFiatBtc)
	api.GET("/latest/", controller.GetFiatBtcByCharCode)
}
