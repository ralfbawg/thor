package util

import "sync"

var (
	BroadcastIns = &Broadcast{

	}
)
var BroadcastTaskPool = sync.Pool{
	New: func() interface{} {
		return &BroadcastTask{}
	}
}

type Broadcast struct {
}
type BroadcastTask struct {
}

func (b *Broadcast) Run() {
}

func newBroadcastTask() {

}
