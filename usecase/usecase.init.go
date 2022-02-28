package usecase

import (
	"binance/constants"
	"binance/usecase/btc"
	"binance/usecase/eth"
	"binance/utility"
	"log"
	"strconv"
	"time"
)

//UsecaseInterface Abstract Object
type UsecaseInterface interface {
	GetAnnualDataMomentum(interval string, start time.Time, assetDollar float64) ([]float64, float64)
}

//InitUsecase Initialize Usecase
func InitUsecase() UsecaseInterface {
	return &UsecaseStruct{
		BTCUSDT: btc.InitBTCUSDT(constants.BTCUSDT),
		ETHUSDT: eth.InitETHUSDT(constants.ETHUSDT),
	}
}

//GetAnnualDataMomentum Get Data from start with interval. return Monthly Momentum for that year
func (uc *UsecaseStruct) GetAnnualDataMomentum(interval string, start time.Time, assetDollar float64) ([]float64, float64) {
	//Acceptable Interval: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
	var startTime time.Time
	switchInterval := interval[len(interval)-1:]
	var duration time.Duration
	var err error
	var durationInt int
	// set startTime = start - 13 interval time (12 for findingthe  momentum, 13 = data at now)
	switch switchInterval {
	case "m", "h":
		duration, err = time.ParseDuration(interval)
		if err != nil {
			log.Printf("Error Interval %s, Err: %+v\n", interval, err)
			return nil, assetDollar
		}
		startTime.Add(-12 * duration)
	case "d":
		durationInt, err = strconv.Atoi(interval[:len(interval)-1])
		if err != nil {
			log.Printf("Error Atoi %s, Err: %+v\n", interval[:len(interval)-1], err)
			return nil, assetDollar
		}
		startTime = start.AddDate(0, 0, -12*durationInt)
	case "w":
		durationInt, err = strconv.Atoi(interval[:len(interval)-1])
		if err != nil {
			log.Printf("Error Atoi %s, Err: %+v\n", interval[:len(interval)-1], err)
			return nil, assetDollar
		}
		startTime = start.AddDate(0, 0, -12*durationInt*7)
	case "M":
		startTime = start.AddDate(0, -12, 0)
	default:
		log.Printf("Error Interval %s\n", interval)
		return nil, assetDollar
	}
	endTime := start.AddDate(1, 0, 0)
	if endTime.After(time.Now()) {
		y, m, d := time.Now().Date()
		endTime = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}

	//Get Offensive Coin
	BTCUSDT := uc.BTCUSDT.GetAnnualData(interval, startTime.Unix(), endTime.Unix())
	ETHUSDT := uc.ETHUSDT.GetAnnualData(interval, startTime.Unix(), endTime.Unix())

	//Loop every new month and append monthly return
	idx := 12 //First 12 elements are used only for Momentum
	report := []float64{}
	assetCoin := 0.0
	coinName := ""
	buyPrice := 0.0
	reportIdx := 0

	startYearAsset := assetDollar
	currentTime := start
	for currentTime.Before(endTime) && !currentTime.Equal(endTime) {
		newMonth := currentTime.AddDate(0, 1, 0)
		startMonthDollar := assetDollar
		// fmt.Println(idx, currentTime, time.Unix(BTCUSDT[idx].OpenTimestamp/1000, 0), endTime, BTCUSDT[idx].OpenTimestamp)

		for currentTime.Before(newMonth) {
			//Sell at the End of Following Day
			switch coinName {
			case "ETHUSDT":
				buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Close, 64)
				assetDollar = assetCoin * buyPrice //Sell at the end of following day
			case "BTCUSDT":
				buyPrice, _ = strconv.ParseFloat(BTCUSDT[idx].Close, 64)
				assetDollar = assetCoin * buyPrice
			case "": //Skip
			}

			//Calculate Offensive Momentum
			btcusdtMomentum := utility.CalculateMomentum(BTCUSDT[idx].Open, BTCUSDT[idx-1].Open, BTCUSDT[idx-3].Open, BTCUSDT[idx-6].Open, BTCUSDT[idx-12].Open)
			ethusdtMomentum := utility.CalculateMomentum(ETHUSDT[idx].Open, ETHUSDT[idx-1].Open, ETHUSDT[idx-3].Open, ETHUSDT[idx-6].Open, ETHUSDT[idx-12].Open)
			momentum := 0.0

			//Can be turned into a function
			if btcusdtMomentum > ethusdtMomentum {
				momentum = btcusdtMomentum
				buyPrice, _ = strconv.ParseFloat(BTCUSDT[idx].Open, 64)
				coinName = "BTCUSDT"
			} else if ethusdtMomentum > btcusdtMomentum {
				momentum = ethusdtMomentum
				buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Open, 64)
				coinName = "ETHUSDT"
			}

			//If momentum is negative don't buy
			if momentum < 0 {
				coinName = ""
			} else {
				assetCoin = assetDollar / buyPrice
			}
			idx++

			switch switchInterval {
			case "m", "h":
				currentTime = currentTime.Add(duration)
			case "d":
				currentTime = currentTime.AddDate(0, 0, durationInt)
			case "w":
				currentTime = currentTime.AddDate(0, 0, durationInt*7)
			case "M":
				currentTime = currentTime.AddDate(0, durationInt, 0)
			}
		}
		// fmt.Println(currentTime.Month(), currentTime, currentTime.Equal(endTime))
		//Monthly Returns
		returns := ((assetDollar - startMonthDollar) / startMonthDollar) * 100
		report = append(report, returns)
		reportIdx++
	}

	//Yearly Returns
	for i := reportIdx; i < 12; i++ {
		report = append(report, 0.0)
	}
	report = append(report, ((assetDollar-startYearAsset)/startYearAsset)*100)
	return report, assetDollar
}
