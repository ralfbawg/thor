package util

import (
	"sync"
)

const (
	defaultMapSize  = 10000
	deletethreshold = 100
)

type ConcurrentMap struct {
	sync.RWMutex
	Map         map[string]interface{}
	keys        []string
	deleteCount int
}

func NewConcurrentMap() *ConcurrentMap {
	//sm := new(ConcurrentMap)
	//sm.Map = make(map[string]interface{})
	return &ConcurrentMap{
		Map:  make(map[string]interface{}, defaultMapSize),
		keys: make([]string, defaultMapSize),
	}

}

func (cm *ConcurrentMap) Get(key string) interface{} {
	cm.RLock()
	value := cm.Map[key]
	cm.RUnlock()
	return value
}

func (cm *ConcurrentMap) Len() int {
	cm.RLock()
	value := len(cm.Map)
	cm.RUnlock()
	return value
}

func (cm *ConcurrentMap) Put(key string, value interface{}) {
	cm.Lock()
	cm.Map[key] = value
	cm.keys = append(cm.keys, key)
	cm.Unlock()
}
func (cm *ConcurrentMap) Del(id string) {
	cm.Lock()
	delete(cm.Map, id)
	cm.deleteCount++
	if cm.deleteCount > deletethreshold {
		factor := 0
		if len(cm.Map)%defaultMapSize > 0 {
			factor = 1
		}
		tmpSize := (factor + len(cm.Map)/defaultMapSize) * defaultMapSize
		tmp := cm.Map
		cm = &ConcurrentMap{
			Map:  make(map[string]interface{}, tmpSize),
			keys: make([]string, tmpSize),
		}
		for k, v := range tmp {
			cm.Map[k] = v
		}
		tmp = nil
	}
	cm.Unlock()
}
func (cm *ConcurrentMap) Foreach(f func(string, interface{})) {
	//cm.RLock()
	//for k, v := range cm.Map {
	//	f(k, v)
	//}
	//cm.RUnlock()
	cm.ForeachN(0, -1, f)
}

//todo 需要测试
func (cm *ConcurrentMap) ForeachN(start int, offset int, f func(string, interface{})) {

	for k, v := range cm.keys {
		if k >= start && (k < start+offset || offset == -1) {
			cm.RLock()
			if cm.Map[v] == nil {
				cm.keys = append(cm.keys[:k], cm.keys[k+1:]...)
			} else {
				f(v, cm.Map[v])
			}
			cm.RUnlock()
		}
	}

}

type Tuple struct {
	Key string
	Val interface{}
}

func (cm *ConcurrentMap) Iter(f func(key string, value interface{})) {
	for item := range cm.iter() {
		f(item.Key, item.Val)
	}
}
func (cm *ConcurrentMap) iter() <-chan Tuple {
	length := cm.Len()
	ch := make(chan Tuple, length)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(length)
		for k, v := range cm.Map {
			cm.RLock()
			ch <- Tuple{k, v}
			cm.RUnlock()
			wg.Done()
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}
