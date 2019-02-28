package statistics

var WsCount int = 0


func GetWsCount() int{
	return WsCount
}

func InitWsStatictics(ch chan int){
	for  {
		WsCount+= <-ch
	}
}