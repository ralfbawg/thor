package util

import (
	"sync"
)

type ConcurrentMap struct {
	sync.RWMutex
	Map map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	//sm := new(ConcurrentMap)
	//sm.Map = make(map[string]interface{})
	return &ConcurrentMap{
		Map:make(map[string]interface{}),
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
	cm.Unlock()
}
func (cm *ConcurrentMap) Del(id string) {
	cm.Lock()
	delete(cm.Map, id)
	cm.Unlock()
}
func (cm *ConcurrentMap) Foreach(f func(string,interface{})) {
	cm.RLock()
	for k,v := range cm.Map {
		f(k,v)
	}
	cm.RUnlock()
}
