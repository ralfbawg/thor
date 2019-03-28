package pojo

import "time"

type App struct {
	Id int64
	Name string
	Key string
	ManagerUid string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}
