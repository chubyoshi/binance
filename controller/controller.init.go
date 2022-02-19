package controller

import (
	"binance/usecase"
)

//CoinInterface Abstract Obj.
type CoinInterface interface {
}

//CoinController Obj.
type CoinController struct {
	Usecase usecase.UsecaseInterface
}

//InitController Init Obj
func InitController(usecase usecase.UsecaseInterface) CoinInterface {
	return &CoinController{
		Usecase: usecase,
	}
}
