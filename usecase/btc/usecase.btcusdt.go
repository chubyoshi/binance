package btc

import (
	"binance/constants"
	"binance/utility"
	"fmt"
)

//BTCUSDTInterface Abstract Obj.
type BTCUSDTInterface interface {
	GetInitialData(interval string, endDate int64) []utility.CandleStickData
	GetData(start int64) utility.CandleStickData
}

//InitBTCUSDT Init Obj
func InitBTCUSDT(name string) BTCUSDTInterface {
	return &BTCUSDT{
		Name: name,
	}
}

func (btc *BTCUSDT) GetInitialData(interval string, endDate int64) []utility.CandleStickData {
	//Get Data from Binance. Binance LIMIT, Default = 500, Max = 1000
	url := fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&endTime=%d", constants.GET_CANDLESTICK_URL, btc.Name, interval, start*1000, end*1000)

	//Get data prior to endDate
	result := []utility.CandleStickData{}
	for i := 0; i < 30; i++ {

	}
	data := utility.GetFromURL(url)
	if data == nil {
		return nil
	}
	if len(data) < 13 {
		return nil
	}

}

//GetAnnualData Get and Data on a yearly basis. Return data base on interval
func (btc *BTCUSDT) GetData(start int64) utility.CandleStickData {
	url := fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&limit=%d", constants.GET_CANDLESTICK_URL, btc.Name, interval, start*1000, 5)

	//Get data until end of year period
	data := utility.GetFromURL(url)
	if data == nil {
		return nil
	}

	for data[len(data)-1].OpenTimestamp/1000 != end {
		starting := data[len(data)-1]

		url = fmt.Sprintf("%ssymbol=%s&interval=%s&startTime=%d&endTime=%d", constants.GET_CANDLESTICK_URL, btc.Name, interval, starting.OpenTimestamp, end*1000)
		data = append(data, utility.GetFromURL(url)[1:]...) //delete the start otherwise duplicate of starting point
	}
	return data[0]
}
