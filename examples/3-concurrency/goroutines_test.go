package main

import (
	"math/rand"
	"testing"
	"time"
)

var global []int

func Benchmark_mergesortV1(b *testing.B) {
	var local []int
	for b.Loop() {
		b.StopTimer()
		input := getRandomElements()
		b.StartTimer()

		sequentialMergesort(input)
		local = input
	}
	global = local
}

func Benchmark_mergesortV2(b *testing.B) {
	var local []int
	for b.Loop() {
		b.StopTimer()
		input := getRandomElements()
		b.StartTimer()

		parallelMergesortV1(input)
		local = input
	}
	global = local
}

func Benchmark_finalMergesortV3(b *testing.B) {
	var local []int
	for b.Loop() {
		b.StopTimer()
		input := getRandomElements()
		b.StartTimer()
		parallelMergesortV2(input)
		local = input
	}
	global = local
}

func getRandomElements() []int {
	n := 10_000
	res := make([]int, n)
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	for i := range n {
		res[i] = rnd.Int()
	}
	return res
}
