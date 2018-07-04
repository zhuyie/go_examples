package cgo

import "testing"

func BenchmarkCGO(b *testing.B) {
	CallCgo(b.N)
}

func BenchmarkGo(b *testing.B) {
	CallGo(b.N)
}

// On my MacBookPro:
// go test -bench . -gcflags '-l'
// BenchmarkCGO-8   	30000000	         52.2 ns/op
// BenchmarkGo-8    	2000000000	         1.40 ns/op

