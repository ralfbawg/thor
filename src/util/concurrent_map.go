package util

import (
	"sync"
)

type ConcurrentMap struct {
	sync.RWMutex
	Map map[string]interface{}
	keys []string
}

func NewConcurrentMap() *ConcurrentMap {
	//sm := new(ConcurrentMap)
	//sm.Map = make(map[string]interface{})
	return &ConcurrentMap{
		Map: make(map[string]interface{}),
		keys:make([]string,20),
	}

}

func (cm *ConcurrentMap) Get(key string) interface{} {
	cm.RLock()
	value := cm.Map[key]
	cm.RUnlock()
	return value
}

func (cm *ConcurrentMap) Put(key string, value interface{}) {
	cm.Lock()
	cm.Map[key] = value
	cm.keys = append(cm.keys,key)
	cm.Unlock()
}
func (cm *ConcurrentMap) Del(id string) {
	cm.Lock()
	delete(cm.Map, id)
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
	cm.RLock()
	for  k,v := range cm.keys {
		if k >= start && (k < start+offset || offset == -1) {
			if cm.Map[v] == nil {
				cm.keys=append(cm.keys[:k],cm.keys[k+1:]...)
			}else{
				f(v, cm.Map[v])
			}

		}
	}
	cm.RUnlock()
}
