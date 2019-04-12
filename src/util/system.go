package util

import (
	"log"
	"os"
	"runtime/pprof"
)

func GetMemoryFile() {
	fm, err := os.OpenFile("./thor_mem.out", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(fm)
	fm.Close()
}

func AOrB(f func() bool, a interface{}, b interface{}) interface{} {
	if f() {
		return a
	} else {
		return b
	}
}
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
