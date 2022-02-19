package utility

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/bitly/go-simplejson"
)

//FormatToExcel Take data and return it into Spreadsheet
func FormatToSpreadsheet() {}

//GetFromURL Get Data from URL return []CandleStickData
func GetFromURL(url string) []CandleStickData {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GetFromURL: %s Error: %+v\n", url, err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("[GetFromURL][Do] Error: %+v\n", err)
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[GetFromURL][ReadAll] Error: %+v\n", err)
		return nil
	}

	j, err := newJSON(body)
	if err != nil {
		log.Printf("[GetFromURL][newJSON] Error: %+v\n", err)
		return nil
	}

	num := len(j.MustArray())
	data := make([]CandleStickData, num)
	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		if len(item.MustArray()) <= 11 {
			err = fmt.Errorf("invalid kline response")
			log.Printf("[GetFromURL][MustArray] Error: %+v\n", err)
		}

		data[i] = CandleStickData{
			OpenTimestamp:  item.GetIndex(0).MustInt64(),
			Open:           item.GetIndex(1).MustString(),
			High:           item.GetIndex(2).MustString(),
			Low:            item.GetIndex(3).MustString(),
			Close:          item.GetIndex(4).MustString(),
			Volume:         item.GetIndex(5).MustString(),
			CloseTimestamp: item.GetIndex(6).MustInt64(),
			Quote:          item.GetIndex(7).MustString(),
			NumberOfTrades: item.GetIndex(8).MustInt64(),
			TakerBuyBase:   item.GetIndex(9).MustString(),
			TakerBuyQuote:  item.GetIndex(10).MustString(),
			Ignore:         item.GetIndex(11).MustString(),
		}
	}

	sortDATA(data)
	return data
}

//NewJSON Seperate 2d Arrays from response
func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

//sortDATA Sort Slices for CandleStickData
func sortDATA(data []CandleStickData) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].OpenTimestamp < data[j].OpenTimestamp
	})
}

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

//CalculateMargin Calculate the margin
func CalculateMargin(today, one, three, six, twelve string) float64 {
	//Convert String to float
	todayPrice, err := strconv.ParseFloat(today, 64)
	if err != nil {
		log.Printf("Error Convert today's Price: %v", err)
	}
	onePrice, err := strconv.ParseFloat(one, 64)
	if err != nil {
		log.Printf("Error Convert interval -1 Price: %v", err)
	}
	threePrice, err := strconv.ParseFloat(three, 64)
	if err != nil {
		log.Printf("Error Convert interval -3 Price: %v", err)
	}
	sixPrice, err := strconv.ParseFloat(six, 64)
	if err != nil {
		log.Printf("Error Convert interval -6 Price: %v", err)
	}
	twelvePrice, err := strconv.ParseFloat(twelve, 64)
	if err != nil {
		log.Printf("Error Convert interval -12 Price: %v", err)
	}

	//Formula = (12 * (todayPrice/interval1 - 1)) + (4 * (todayPrice/interval3 - 1)) + (2 * (todayPrice/interval6 - 1)) + (todayPrice/interval12 - 1)
	return (12 * (todayPrice/onePrice - 1)) + (4 * (todayPrice/threePrice - 1)) + (2 * (todayPrice/sixPrice - 1)) + (todayPrice/twelvePrice - 1)
}
