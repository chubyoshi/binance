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
	GetAnnualDataMomentum(interval string, start time.Time) []float64
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

//GetAnnualDataMomentum Get Data from start with interval. return Monthly Momentum for that year
func (uc *UsecaseStruct) GetAnnualDataMomentum(interval string, start time.Time) []float64 {
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

	report := []float64{}
	idx := 12 //First 12 elements are used only for Momentum
	idxf := 12
	assetDollar := 1000.0
	startYearAsset := assetDollar

	//Loop every new month until the before the start of next year
	for currentTime := start; currentTime.Before(endTime); {
		newMonth := currentTime.AddDate(0, 1, 0).Unix()
		startMonthDollar := assetDollar
		buyPrice := 0.0
		coin := ""

		for BTCUSDT[idx].OpenTimestamp/1000 < newMonth {
			//Hold and sell at the End of Following Day
			switch coin {
			case "ETH":
				assetCoin := assetDollar / buyPrice //Buy at highest Momentum
				buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Close, 64)
				assetDollar = assetCoin * buyPrice //Sell at the end of following day
			case "BTC":
				assetCoin := assetDollar / buyPrice
				buyPrice, _ = strconv.ParseFloat(BTCUSDT[idx].Close, 64)
				assetDollar = assetCoin * buyPrice
			case "ETH US":
				assetCoin := assetDollar / buyPrice
				buyPrice, _ = strconv.ParseFloat(ETHUSDC[idxf].Close, 64)
				assetDollar = assetCoin * buyPrice
			case "BTC US":
				assetCoin := assetDollar / buyPrice
				buyPrice, _ = strconv.ParseFloat(BTCUSDC[idxf].Close, 64)
				assetDollar = assetCoin * buyPrice
			}

			//Calculate Offensive
			BTC := utility.CalculateMomentum(BTCUSDT[idx].Close, BTCUSDT[idx-1].Close, BTCUSDT[idx-3].Close, BTCUSDT[idx-6].Close, BTCUSDT[idx-12].Close)
			ETH := utility.CalculateMomentum(ETHUSDT[idx].Close, ETHUSDT[idx-1].Close, ETHUSDT[idx-3].Close, ETHUSDT[idx-6].Close, ETHUSDT[idx-12].Close)
			momentum := 0.0

			if BTC > ETH {
				momentum = BTC
				buyPrice, _ = strconv.ParseFloat(BTCUSDT[idx].Close, 64)
				coin = "BTC"
			} else {
				momentum = ETH
				buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Close, 64)
				coin = "ETH"
			}

			//If negative use highest defensive
			flag := BTCUSDC[idxf].OpenTimestamp == BTCUSDT[idx].OpenTimestamp
			if momentum < 0 {
				//If No Defensive Option
				if flag {
					// fmt.Println(time.Unix(BTCUSDC[idxf].OpenTimestamp/1000, 0), time.Unix(BTCUSDC[idxf].OpenTimestamp/1000, 0), idxf, idx)
					idx++
					idxf++
					continue
				}

				//Calculate Max Defensive
				BTCUS := utility.CalculateMomentum(BTCUSDC[idxf].Close, BTCUSDC[idxf-1].Close, BTCUSDC[idxf-3].Close, BTCUSDC[idxf-6].Close, BTCUSDC[idxf-12].Close)
				ETHUS := utility.CalculateMomentum(ETHUSDC[idxf].Close, ETHUSDC[idxf-1].Close, ETHUSDC[idxf-3].Close, ETHUSDC[idxf-6].Close, ETHUSDC[idxf-12].Close)

				if BTCUS > ETHUS {
					momentum = BTCUS
					buyPrice, _ = strconv.ParseFloat(BTCUSDC[idxf].Close, 64)
					coin = "BTC US"
				} else {
					momentum = ETHUS
					buyPrice, _ = strconv.ParseFloat(ETHUSDC[idxf].Close, 64)
					coin = "ETH US"
				}
			}
			// fmt.Printf("%d. Time:%s Momentum:%f\n", idx, time.Unix(BTCUSDT[idx].OpenTimestamp/1000, 0), total)
			idx++
			if flag {
				idxf++
			}
		}

		returns := ((assetDollar - startMonthDollar) / startMonthDollar) * 100
		report = append(report, returns)
		currentTime = currentTime.AddDate(0, 1, 0)
	}
	report = append(report, ((assetDollar-startYearAsset)/startYearAsset)*100)

	return report
}
