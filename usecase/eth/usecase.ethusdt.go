package eth

import (
	"binance/constants"
	"binance/utility"
	"fmt"
	"strconv"
	"time"
)

//ETHUSDTInterface Abstract Obj.
type ETHUSDTInterface interface {
	GetInitialData(interval string, endDate time.Time) []utility.CandleStickData
	GetData(start int64) utility.CandleStickData
}

//InitETHUSDT Init Obj
func InitETHUSDT(name string) ETHUSDTInterface {
	return &ETHUSDT{
		Name: name,
	}
}

//GetInitialData Get the initial data for momentum counting
func (eth *ETHUSDT) GetInitialData(interval string, endDate time.Time) []utility.CandleStickData {
	var duration time.Duration
	var durationInt int
	var err error
	intervalFormat := interval[len(interval)-1:]

	switch intervalFormat {
	case "h", "m", "s":
		duration, err = time.ParseDuration(interval)
		if err != nil {
			return nil
		}
	case "d", "w", "M":
		durationInt, err = strconv.Atoi(interval[:len(interval)-1])
		if err != nil {
			return nil
		}
	default:
		return nil
	}

	var startDate time.Time
	switch intervalFormat {
	case "m", "h", "s":
		startDate = endDate.Add(-12 * duration)
	case "d":
		startDate = endDate.AddDate(0, 0, -12*durationInt)
	case "w":
		startDate = endDate.AddDate(0, 0, -12*durationInt*7)
	case "M":
		startDate = endDate.AddDate(0, -12*durationInt, 0)
	}

	//Get data prior to endDate
	result := []utility.CandleStickData{}
	for i := 0; i < 13; i++ {
		result = append(result, eth.GetData(startDate.Unix()))

		switch intervalFormat {
		case "m", "h", "s":
			startDate = startDate.Add(duration)
		case "d":
			startDate = startDate.AddDate(0, 0, durationInt)
		case "w":
			startDate = startDate.AddDate(0, 0, durationInt*7)
		case "M":
			startDate = startDate.AddDate(0, durationInt, 0)
		}

	}
	return result
}

//GetData Get Data on that timestamp
func (eth *ETHUSDT) GetData(start int64) utility.CandleStickData {
	url := fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&limit=%d", constants.GET_CANDLESTICK_URL, eth.Name, constants.CANDLESTICK_LOWEST_INTERVAL, start*1000, 5)

	//Get data until end of year period
	data := utility.GetFromURL(url)

	return data[0]
}
