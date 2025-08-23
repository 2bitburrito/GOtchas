package main

import (
	"fmt"
	"math/rand"
	"runtime"
)

func main() {
	n := 1_000_000
	m := make(map[int]*[128]byte)
	printAlloc()

	for i := range n { // Adds 1 million elements
		m[i] = &[128]byte{}
	}
	printAlloc()

	for i := range n { // Deletes 1 million elements
		delete(m, i)
	}

	runtime.GC() // Triggers a manual GC
	printAlloc()
	runtime.KeepAlive(m) // Keeps a reference to m so that the map isn’t collected
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/(1024*1024))
}

func other() {
	n := 1_000_000
	m := make(map[int]Foo)

	for i := range n { // Adds 1 million elements
		m[i] = Foo{}
	}

	for i, foo := range m { // Deletes 50% of elements
		if foo.isExpired() {
			delete(m, i)
		}
	}
}

type Foo struct{}

func (f Foo) isExpired() bool {
	flip := rand.Intn(2)
	return flip == 0
}
