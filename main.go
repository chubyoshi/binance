package main

import (
	"binance/controller"
	"binance/usecase"
	"log"
)

func main() {
	log.Println("[Main] Starting Task")

	year := 2018
	interval := "1d"
	assets := 1000.0

	uc := usecase.InitUsecase()
	ctrl := controller.InitController(uc)
	ctrl.ProcessBinanceCandleStick(interval, year, assets)
}
