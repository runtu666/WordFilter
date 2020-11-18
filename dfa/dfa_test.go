package dfa

import (
	"fmt"
	"go-wordfilter/common"
	"log"
	"testing"
)

func TestDfa(t *testing.T) {
	dfa := NewTire()
	dfa.LoadWords(common.GetWords())
	result := dfa.Search("hello av hava 毛泽东 sm 气枪 测试, 支付宝 ")
	log.Println("len:", len(result))
	for _, item := range result {
		fmt.Printf("%+v \n", item)
	}
}
