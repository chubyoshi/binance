package main

import (
	"binance/controller"
	"binance/usecase"
	"log"
)

func main() {
	uc := usecase.InitUsecase()
	ctrl := controller.InitController(uc)

	log.Println("[Main] Starting Task")
	interval := "1d"
	year := 2018
	ctrl.ProcessBinanceCandleStick(interval, year)
}
