package usecase

import (
	"binance/constants"
	"binance/usecase/btc"
	"binance/usecase/eth"
	"binance/utility"
	"fmt"
	"log"
	"strconv"
	"time"
)

//UsecaseInterface Abstract Object
type UsecaseInterface interface {
	GetAnnualDataMargin(interval string, start time.Time) []float64
}

//InitUsecase Initialize Usecase
func InitUsecase() UsecaseInterface {
	return &UsecaseStruct{
		BTCUSDC: btc.InitBTCUSDC(constants.BTCUSDC),
		BTCUSDT: btc.InitBTCUSDT(constants.BTCUSDT),
		ETHUSDC: eth.InitETHUSDC(constants.ETHUSDC),
		ETHUSDT: eth.InitETHUSDT(constants.ETHUSDT),
	}
}

//GetAnnualDataMargin Get Data from start with interval. return Monthly Margin for that year
func (uc *UsecaseStruct) GetAnnualDataMargin(interval string, start time.Time) []float64 {
	//end = current year however there's a limit of 1000. TODO validation limit for other interval
	//Acceptable Interval: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
	var startTime time.Time
	switch interval[len(interval)-1:] {
	case "m", "h":
		duration, err := time.ParseDuration(interval)
		if err != nil {
			log.Printf("Error Interval %s, Err: %+v\n", interval, err)
			return nil
		}
		startTime.Add(-13 * duration)
	case "d":
		duration, err := strconv.Atoi(interval[:len(interval)-1])
		if err != nil {
			log.Printf("Error Atoi %s, Err: %+v\n", interval[:len(interval)-1], err)
			return nil
		}
		startTime = start.AddDate(0, 0, duration-13)
	case "w":
		startTime = start.AddDate(0, 0, -91) //91 = 7d * 13
	case "M":
		startTime = start.AddDate(0, -13, 0)
	default:
		log.Printf("Error Interval %s\n", interval)
		return nil
	}
	endTime := start.AddDate(1, 0, 0)

	//Get Yearly Data
	BTCUSDT := uc.BTCUSDT.GetAnnualData(interval, startTime.Unix(), endTime.Unix())
	BTCUSDC := uc.BTCUSDC.GetAnnualData(interval, startTime.Unix(), endTime.Unix())
	ETHUSDT := uc.ETHUSDT.GetAnnualData(interval, startTime.Unix(), endTime.Unix())
	ETHUSDC := uc.ETHUSDC.GetAnnualData(interval, startTime.Unix(), endTime.Unix())

	fmt.Println(len(BTCUSDT), len(ETHUSDT), len(ETHUSDC), len(BTCUSDC))

	report := []float64{}
	endTime = start.AddDate(0, 2, 0)
	annual := 0.0
	idx := 12 //First 12 elements are used only for margin
	//Loop every new month until the before the start of next year
	for currentTime := start; currentTime.Before(endTime); {
		newMonth := currentTime.AddDate(0, 1, 0).Unix()

		fmt.Println(currentTime.Month())
		total := 0.0
		for BTCUSDT[idx].OpenTimestamp/1000 < newMonth {

			//Calculate Offensive
			BTC := utility.CalculateMargin(BTCUSDT[idx].Close, BTCUSDT[idx-1].Close, BTCUSDT[idx-3].Close, BTCUSDT[idx-6].Close, BTCUSDT[idx-12].Close)
			ETH := utility.CalculateMargin(ETHUSDT[idx].Close, ETHUSDT[idx-1].Close, ETHUSDT[idx-3].Close, ETHUSDT[idx-6].Close, ETHUSDT[idx-12].Close)
			offenseMargin := utility.Max(BTC, ETH)
			//If negative use highest defensive
			if offenseMargin < 0 {
				//If No Defensive Option
				if BTCUSDC[idx].OpenTimestamp != BTCUSDT[idx].OpenTimestamp {
					total += offenseMargin
					idx++
					continue
				}

				//Calculate Max Defensive
				BTCUS := utility.CalculateMargin(BTCUSDC[idx].Close, BTCUSDC[idx-1].Close, BTCUSDC[idx-3].Close, BTCUSDC[idx-6].Close, BTCUSDC[idx-12].Close)
				ETHUS := utility.CalculateMargin(ETHUSDC[idx].Close, ETHUSDC[idx-1].Close, ETHUSDC[idx-3].Close, ETHUSDC[idx-6].Close, ETHUSDC[idx-12].Close)
				defenseMargin := utility.Max(BTCUS, ETHUS)

				total += defenseMargin
			} else {
				total += offenseMargin
			}
			fmt.Printf("%d. Time:%s Margin:%f\n", idx, time.Unix(BTCUSDT[idx].OpenTimestamp/1000, 0), total)
			idx++
		}
		report = append(report, total)
		annual += total
		currentTime = currentTime.AddDate(0, 1, 0)
	}
	report = append(report, annual)

	return report
}
