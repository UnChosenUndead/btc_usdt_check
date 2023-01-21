package controller

import (
	"github.com/StanDenisov/btc_usdt_check/queries"
	"github.com/StanDenisov/btc_usdt_check/utils"
	"github.com/gin-gonic/gin"
)

func GetLastCourseOfBtcUsdt(c *gin.Context) {
	response := queries.SelectLastBtcUsdCourse()
	c.JSON(200, response)
}

func GetCourseOfBtcUsdtWithFilter(c *gin.Context) {
	var f utils.UsdBtcFilter
	c.Bind(&f)
	response := queries.SelectFilteredBtcUsdCourse(f)
	c.JSON(200, response)
}

func GetLastCoursesFiatBtc(c *gin.Context) {
	response := queries.SelectLastFiatBtcCourse()
	c.JSON(200, response)
}

func GetCourseFilteredFiatBtc(c *gin.Context) {
	var f utils.FiatBtcFilter
	c.Bind(&f)
	response := queries.SelectFilteredFiatBtcCourse(f)
	c.JSON(200, response)
}

func GetLastCourseFiatRub(c *gin.Context) {
	responseStruct := struct {
		Count   int32              `json:"total"`
		History map[string]float64 `json:"history"`
	}{}
	responseStruct.History = make(map[string]float64)
	response := queries.SelectLastFiatRubCourse()
	responseStruct.Count = int32(response.Count)
	for _, rub := range response.FiatRubSelect {
		responseStruct.History[rub.FiatCharCode] = rub.FiatValue
	}
	c.JSON(200, responseStruct)
}

func GetCourseFilteredFiatRub(c *gin.Context) {
	var f utils.FiatFilter
	c.Bind(&f)
	responseStruct := struct {
		Count   int32              `json:"total"`
		History map[string]float64 `json:"history"`
	}{}
	responseStruct.History = make(map[string]float64)
	response := queries.SelectFilteredFiatRubCourse(f)
	responseStruct.Count = int32(response.Count)
	for _, rub := range response.FiatRubSelect {
		responseStruct.History[rub.FiatCharCode] = rub.FiatValue
	}
	c.JSON(200, responseStruct)
}

func GetFiatBtcByCharCode(c *gin.Context) {
	var b struct {
		CharCode string `form:"char_code"`
	}
	c.Bind(&b)
	response := queries.SelectFiatBtcByCharCode(b.CharCode)
	c.JSON(200, response)
}
