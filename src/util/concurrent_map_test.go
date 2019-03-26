package util

import (
	"fmt"
	"strconv"
	"testing"
	"sync"
)

const countStep = 1000000

func BenchmarkCm(b *testing.B) { //0.60 ns/op
	b.ReportAllocs()
	m := NewConcurrentMap()

	for i := 0; i < countStep; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	count := 0
	m.Foreach(func(s string, i interface{}) {
		tcount, _ := strconv.Atoi(i.(Animal).name)
		count += tcount
	})
	fmt.Printf("get cm count = %d\n", count)
}
func BenchmarkCm2(b *testing.B) { //0.02 ns/op
	b.ReportAllocs()
	m := NewConcurrentMap()

	for i := 0; i < countStep; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})

	}
	count := int64(0)
	m.Iter(func(s string, i interface{}) {
		tcount, _ := strconv.Atoi(i.(Animal).name)
		count += int64(tcount)
	})
	fmt.Printf("get cm2 count = %d\n", count)
}
func BenchmarkSyncmap(b *testing.B) { //0.04 ns/op
	b.ReportAllocs()
	m := &sync.Map{}

	for i := 0; i < countStep; i++ {
		m.Store(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	count := 0
	m.Range(func(key, value interface{}) bool {
		tcount, _ := strconv.Atoi(value.(Animal).name)
		count += tcount
		return true
	})
	fmt.Printf("get cm3 count = %d\n", count)
}
func BenchmarkSyncmap2(b *testing.B) { //0.04 ns/op
	b.ReportAllocs()
	m := New()

	for i := 0; i < countStep; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	count := 0
	for t:= range m.IterBuffered(){
		tcount, _ := strconv.Atoi(t.Val.(Animal).name)
		count += tcount
	}

	//m.IterBuffered(func(key string, v interface{}) {
	//	tcount, _ := strconv.Atoi(v.(Animal).name)
	//	count += tcount
	//})
	fmt.Printf("get cm4 count = %d\n", count)
}
func BenchmarkSyncmap3(b *testing.B) { //0.04 ns/op
	b.ReportAllocs()
	m := New()

	for i := 0; i < countStep; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	count := 0
	//for t:= range m.IterBuffered(){
	//	tcount, _ := strconv.Atoi(t.Val.(Animal).name)
	//	count += tcount
	//}

	m.IterCb(func(key string, v interface{}) {
		tcount, _ := strconv.Atoi(v.(Animal).name)
		count += tcount
	})
	fmt.Printf("get cm5 count = %d\n", count)
}
//func BenchmarkGood(b *testing.B) { //0.04 ns/op
//	b.ReportAllocs()
//
//	count := 0
//	for i := 0; i < countStep; i++ {
//		count += i
//	}
//
//	fmt.Printf("get cm4 count = %d\n", count)
//}

type Animal struct {
	name string
}
