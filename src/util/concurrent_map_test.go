package util

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkCm(b *testing.B) {
	b.ReportAllocs()
	m := NewConcurrentMap()

	for i := 0; i < 100000; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	m.Foreach(func(s string, i interface{}) {
		fmt.Printf("get k=%s & v=%s\n", s, i.(Animal).name)
	})

}
func BenchmarkCm2(b *testing.B) {
	b.ReportAllocs()
	m := NewConcurrentMap()

	for i := 0; i < 100000; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	m.Iter(func(s string, i interface{}) {
		fmt.Printf("get k=%s & v=%s\n", s, i.(Animal).name)
	})
}

type Animal struct {
	name string
}
