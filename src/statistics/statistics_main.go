package statistics

func InitStatistics()  {
	InitWsStatictics(make(chan int,10))
}
