package main

import (
	"binance/controller"
	"binance/usecase"
	"log"
	"time"
)

func main() {
	uc := usecase.InitUsecase()
	ctrl := controller.InitController(uc)

	log.Println("[Main] Starting Task")
	interval := 24 * time.Hour
	year := 2019
	ctrl.ProcessBinanceCandleStick(interval, year)
}
