package usecase

import (
	"binance/usecase/btc"
	"binance/usecase/eth"
)

//UsecaseStruct Object Field
type UsecaseStruct struct {
	BTCUSDT btc.BTCUSDTInterface
	ETHUSDT eth.ETHUSDTInterface
}
