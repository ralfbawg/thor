package filter

import (
	"common/logging"
)

type BaseFilter struct {
	filterI
}

func (b *BaseFilter) before() { //空方法
	logging.Debug("do before")
}
func (b *BaseFilter) after() { //空方法
	logging.Debug("do after")
}
