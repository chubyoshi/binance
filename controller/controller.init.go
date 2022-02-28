package controller

import (
	"binance/usecase"
	"binance/utility"
	"fmt"
	"log"
	"time"
)

//CoinInterface Abstract Obj.
type CoinInterface interface {
	ProcessBinanceCandleStick(interval string, year int, startingAsset float64)
}

//CoinController Obj.
type CoinController struct {
	Usecase usecase.UsecaseInterface
}

//InitController Init Obj
func InitController(usecase usecase.UsecaseInterface) CoinInterface {
	return &CoinController{
		Usecase: usecase,
	}
}

//ProcessBinanceCandleStick Take Interval and start year. Get Data & process into spreadsheet from date to now
func (ctrl *CoinController) ProcessBinanceCandleStick(interval string, year int, startingAsset float64) {

	//Get Data report from binance
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0)) //January 1 of selected year
	report, err := ctrl.Usecase.GetBinanceReport(interval, startDate, startingAsset)
	if err != nil {
		log.Printf("[ProcessBinanceCandleStick][GetBinanceReport] Error: %+v\n", err)
		return
	}

	//Process into Spreadsheet -> write into folder
	utility.FormatToSpreadsheet(report, interval, year)
	log.Printf("[ProcessBinanceCandleStick] Done. File %s Created/Updated\n", fmt.Sprintf("Kline%dInterval%s.csv", year, interval))
}
