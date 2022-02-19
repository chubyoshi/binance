package btc

//BTCUSDCInterface Abstract Obj.
type BTCUSDCInterface interface {
}

//InitBTCUSDC Init Obj
func InitBTCUSDC(name string) BTCUSDCInterface {
	return &BTCUSDC{
		Name: name,
	}
}
