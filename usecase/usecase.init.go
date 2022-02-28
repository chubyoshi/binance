package usecase

import (
	"binance/constants"
	"binance/usecase/btc"
	"binance/usecase/eth"
	"binance/utility"
	"fmt"
	"strconv"
	"time"
)

//UsecaseInterface Abstract Object
type UsecaseInterface interface {
	GetBinanceReport(interval string, start time.Time, initialMoney float64) ([]float64, error)
}

//InitUsecase Initialize Usecase
func InitUsecase() UsecaseInterface {
	return &UsecaseStruct{
		BTCUSDT: btc.InitBTCUSDT(constants.BTCUSDT),
		ETHUSDT: eth.InitETHUSDT(constants.ETHUSDT),
	}
}

//GetBinanceReport Get Data from start until now and return the returns in montly slices pf float64. It takes the interval and when to start in time obj
func (uc *UsecaseStruct) GetBinanceReport(interval string, start time.Time, initialMoney float64) ([]float64, error) {
	var duration time.Duration
	var durationInt int
	var err error
	intervalFormat := interval[len(interval)-1:]

	//Validate Interval Format
	switch intervalFormat {
	case "h", "m", "s":
		duration, err = time.ParseDuration(interval)
		if err != nil {
			return nil, fmt.Errorf("error interval format %s. Err: %v\n", interval, err)
		}
	case "d", "w", "M", "Y":
		durationInt, err = strconv.Atoi(interval[:len(interval)-1])
		if err != nil {
			return nil, fmt.Errorf("error interval format %s. Err: %v\n", interval, err)
		}
	default:
		return nil, fmt.Errorf("error interval format %s\n", interval)
	}

	//Get Offensive Coin 12 interval data prior for momentum counting
	BTCUSDT := uc.BTCUSDT.GetInitialData(interval, start.Unix())
	ETHUSDT := uc.ETHUSDT.GetInitialData(interval, start.Unix())

	//Validate Initial data
	if BTCUSDT == nil || ETHUSDT == nil {
		return nil, fmt.Errorf("Data from Binance Error")
	}

	monthly, yearly := start.AddDate(0, 1, 0), start.AddDate(1, 0, 0)
	initialMonthly, initialYearly := initialMoney, initialMoney
	report := []float64{}
	var assetMoney, assetCoin, buyPrice float64
	var coinName string
	idx := 12 //First 12 elements are used only for Momentum which is the size of each coin slices
	for i := start; i.Before(time.Now()); {
		//Get Data from Binance for the offensive coin goroutine until all coins data get

		if i.After(monthly) {
			percentage := ((assetMoney - initialMonthly) / initialMonthly) * 100
			report = append(report, percentage)
			initialMonthly = assetMoney
			monthly = monthly.AddDate(0, 1, 0)
		}

		if i.After(yearly) {
			percentage := ((assetMoney - initialYearly) / initialYearly) * 100
			report = append(report, percentage)
			initialYearly = assetMoney
			yearly = yearly.AddDate(1, 0, 0)
		}

		//Selling the Coin at Open Price the next interval, if empty don't = momentum < 0 so just hold until momentum is positive
		switch coinName {
		case "BTCUSDT":
			buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Open, 64)
			assetMoney = assetCoin * buyPrice
		case "ETHUSDT":
			buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Open, 64)
			assetMoney = assetCoin * buyPrice
		case "": //Skip selling since, previous interval didn't buy coin
		}

		//Finding Momentum of each offensive coin
		btcusdtMomentum := utility.CalculateMomentum(BTCUSDT[idx].Open, BTCUSDT[idx-1].Open, BTCUSDT[idx-3].Open, BTCUSDT[idx-6].Open, BTCUSDT[idx-12].Open)
		ethusdtMomentum := utility.CalculateMomentum(ETHUSDT[idx].Open, ETHUSDT[idx-1].Open, ETHUSDT[idx-3].Open, ETHUSDT[idx-6].Open, ETHUSDT[idx-12].Open)

		//Finding which coin to buy, by the maximum momentum
		momentum := 0.0
		if btcusdtMomentum > ethusdtMomentum {
			momentum = btcusdtMomentum
			buyPrice, _ = strconv.ParseFloat(BTCUSDT[idx].Open, 64)
			coinName = "BTCUSDT"
		} else {
			momentum = ethusdtMomentum
			buyPrice, _ = strconv.ParseFloat(ETHUSDT[idx].Open, 64)
			coinName = "ETHUSDT"
		}

		//If max momentum is negative, don't buy.
		if momentum < 0 {
			coinName = ""
		} else {
			assetCoin = assetMoney / buyPrice
		}

		//Append time based on interval
		switch intervalFormat {
		case "m", "h", "s":
			i.Add(duration)
		case "d":
			i.AddDate(0, 0, durationInt)
		case "w":
			i.AddDate(0, 0, durationInt*7)
		case "M":
			i.AddDate(0, durationInt, 0)
		case "Y":
			i.AddDate(durationInt, 0, 0)
		}
	}

	return report, nil
}
