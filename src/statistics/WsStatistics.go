package statistics

import "go/ast"

var WsCount int = 0

func (c chan int <-int )  {
	for ; ;  {
		
	}


}

func GetWsCount() int{
	return WsCount
}

func InitWsStatictics(ch chan int){
	for  {
		WsCount+= <-ch
	}
}