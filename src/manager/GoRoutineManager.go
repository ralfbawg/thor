package manager

type GoRoutineItem struct {
	id int
	chanl chan<-string
}

const (
	RoutineSliceDefaultCap = 100000
	RoutineSliceDefaultLen = 50000
)

//var GoRoutineSlice = make([]GoRoutineItem, RoutineSliceDefaultLen, RoutineSliceDefaultCap)
//var GoRoutineMap = make(map[string]GoRoutineItem,RoutineSliceDefaultCap)

//func AddGoFunc(f func()) {
//	t := make(chan GoRoutineItem,1)
//	GoRoutineMap[uuid.Generate().String()] = GoRoutineItem{
//		id:1,
//	}
//	sliceLen := len(GoRoutineSlice)
//
//}

func DeleteGoFunc()  {

}
