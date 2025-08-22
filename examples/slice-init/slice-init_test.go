package main

import "testing"

const n = 1_000_000

var global = make([]Bar, 0)

func BenchmarkEmptySlice(b *testing.B) {
	foos := make([]Foo, n)
	local := make([]Bar, 0)
	for b.Loop() {
		local = convertEmptySlice(foos)
	}
	global = local
}

func BenchmarkCapSlice(b *testing.B) {
	foos := make([]Foo, n)
	local := make([]Bar, 0)
	for b.Loop() {
		convertGivenCapacity(foos)
	}

	global = local
}

func BenchmarkLenSlice(b *testing.B) {
	foos := make([]Foo, n)
	local := make([]Bar, 0)
	for b.Loop() {
		convertGivenLength(foos)
	}
	global = local
}
