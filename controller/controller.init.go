package controller

import (
	"binance/usecase"
	"log"
	"time"
)

//CoinInterface Abstract Obj.
type CoinInterface interface {
	ProcessBinanceCandleStick(interval time.Duration, year int)
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
func (ctrl *CoinController) ProcessBinanceCandleStick(interval time.Duration, year int) {
	//Acceptable Interval: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
	dateStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC", 0))

	//loop every year once
	for i := year; i < 2020; i++ { //TODO magic variable
		//go routine for offensive and defensive, getting data and churn data into yearly outcome
		ctrl.Usecase.GetAnnualDataMargin("1d", dateStart) //TODO change using input scanned and validate

	}

	//Process into Spreadsheet -> write into folder

	log.Println("[ProcessBinanceCandleStick]Done")
}
