package usecase

import (
	"binance/constants"
	"binance/usecase/btc"
	"binance/usecase/eth"
)

//UsecaseInterface Abstract Object
type UsecaseInterface interface {
}

//UsecaseStruct Object Field
type UsecaseStruct struct {
	BTCUSDC btc.BTCUSDCInterface
	BTCUSDT btc.BTCUSDTInterface
	ETHUSDC eth.ETHUSDCInterface
	ETHUSDT eth.ETHUSDTInterface
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
