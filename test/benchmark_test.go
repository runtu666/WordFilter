package test

import (
	"go-wordfilter/ac"
	"go-wordfilter/common"
	"go-wordfilter/dfa"
	"testing"
)

var a = ac.NewAc()
var d = dfa.NewDfa()

func init() {
	a.LoadWords(common.GetWords())
	d.LoadWords(common.GetWords())
}

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
