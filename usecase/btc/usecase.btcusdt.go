package btc

//BTCUSDTInterface Abstract Obj.
type BTCUSDTInterface interface {
}

//InitBTCUSDT Init Obj
func InitBTCUSDT(name string) BTCUSDTInterface {
	return &BTCUSDT{
		Name: name,
	}
}
