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
	ProcessBinanceCandleStick(interval string, year int, asset float64)
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
func (ctrl *CoinController) ProcessBinanceCandleStick(interval string, year int, asset float64) {
	//Acceptable Interval: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
	switch interval {
	case "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M":
	default:
		log.Printf("[ProcessBinanceCandleStick] Interval Invalid: %s\n", interval)
		return
	}

	//Get Yearly Data Report
	dateStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))
	report := []float64{}
	data := []float64{}
	for i := dateStart; i.Before(time.Now()) && !i.Equal(time.Now()); i = i.AddDate(1, 0, 0) {
		data, asset = ctrl.Usecase.GetAnnualDataMomentum(interval, dateStart, asset)
		report = append(report, data...)
		dateStart = dateStart.AddDate(1, 0, 0)
	}

	//Process into Spreadsheet -> write into folder
	utility.FormatToSpreadsheet(report, interval, year)
	log.Printf("[ProcessBinanceCandleStick] Done. File %s Created/Updated\n", fmt.Sprintf("Kline%dInterval%s.csv", year, interval))
	log.Printf("Final Assets = %f", asset)
}
