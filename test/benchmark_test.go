package test

import (
	"testing"
)

func BenchmarkAC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Search(Text2)
	}
}

func BenchmarkDfa(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		d.Search(Text2)
	}
}
