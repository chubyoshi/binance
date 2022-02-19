package eth

//ETHUSDCInterface Abstract Obj.
type ETHUSDCInterface interface {
}

//InitETHUSDC Initialize Obj
func InitETHUSDC(name string) ETHUSDCInterface {
	return &ETHUSDC{
		Name: name,
	}
}
