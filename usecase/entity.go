package usecase

import (
	"binance/usecase/btc"
	"binance/usecase/eth"
)

//UsecaseStruct Object Field
type UsecaseStruct struct {
	BTCUSDC btc.BTCUSDCInterface
	BTCUSDT btc.BTCUSDTInterface
	ETHUSDC eth.ETHUSDCInterface
	ETHUSDT eth.ETHUSDTInterface
}
