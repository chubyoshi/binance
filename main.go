package main

import (
	"binance/controller"
	"binance/usecase"
	"fmt"
	"log"
	"strconv"
)

func main() {
	log.Println("[Main] Starting Task")

	scan := ""
	fmt.Print("Enter start Year: ")
	fmt.Scanln(&scan)
	year, err := strconv.Atoi(scan)
	if err != nil {
		fmt.Printf("Error Parsing int: %+v", err)
	}

	fmt.Print(`Enter Interval Period ("m", "h", "s", "d", "w", "M", "Y"): `)
	interval := "1d"
	fmt.Scanln(&interval)

	fmt.Print("Enter starting Asset: ")
	fmt.Scanln(&scan)
	assets, err := strconv.ParseFloat(scan, 64)
	if err != nil {
		fmt.Printf("Error Parsing float64: %+v", err)
	}

	uc := usecase.InitUsecase()
	ctrl := controller.InitController(uc)
	ctrl.ProcessBinanceCandleStick(interval, year, assets)
}
