package eth

//ETHUSDTInterface Abstract Obj
type ETHUSDTInterface interface {
}

//InitETHUSDT Initialize
func InitETHUSDT(name string) ETHUSDTInterface {
	return &ETHUSDT{
		Name: name,
	}
}
