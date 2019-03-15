package statistics

var StatArr = []string{"taskSum", "taskClientSum"}
var StatApp []*Statistics

type Statistics struct {
	name  string
	key string
	subKey string
	count int
	s   func (a string)
	step  chan int
}

func (s *Statistics) fun(a string) {
	s.count += 1
}
func InitStatistics() {
	t := make([]*Statistics, 100)
	for _, v := range StatArr {
		t = append(t, &Statistics{
			name:  v,
			step:  make(chan int),
			count: 0,
		})
	}
}
