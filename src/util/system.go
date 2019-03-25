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
