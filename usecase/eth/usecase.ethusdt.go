package eth

import (
	"binance/constants"
	"binance/utility"
	"fmt"
)

//ETHUSDTInterface Abstract Obj
type ETHUSDTInterface interface {
	GetAnnualData(interval string, start, end int64) []utility.CandleStickData
}

//InitETHUSDT Initialize
func InitETHUSDT(name string) ETHUSDTInterface {
	return &ETHUSDT{
		Name: name,
	}
}

//GetAnnualData Get and Data on a yearly basis. Return data base on interval
func (eth *ETHUSDT) GetAnnualData(interval string, start, end int64) []utility.CandleStickData {
	//Get Data from Binance. Binance LIMIT, Default = 500, Max = 1000
	url := fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&endTime=%d", constants.GET_CANDLESTICK_URL, eth.Name, interval, start*1000, end*1000)

	//Get data until end of year period
	data := utility.GetFromURL(url)
	if data == nil {
		//Error
		return nil
	}

	for data[len(data)-1].OpenTimestamp/1000 != end {
		starting := data[len(data)-1]

		url = fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&endTime=%d", constants.GET_CANDLESTICK_URL, eth.Name, interval, starting.OpenTimestamp, end*1000)
		data = append(data, utility.GetFromURL(url)[1:]...) //delete the start otherwise duplicate of starting point
	}

	return data
}
