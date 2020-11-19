package dfa

import (
	"fmt"
	"go-wordfilter/common"
	"log"
	"testing"
)

func TestDfa(t *testing.T) {
	dfa := NewDfa()
	dfa.LoadWords(common.GetWords())
	result := dfa.Search("hello av java 毛泽东 sm 气枪 测试, 支付宝 ")
	log.Println("len:", len(result))
	for _, item := range result {
		fmt.Printf("%+v \n", item)
	}
}

func TestDfaReplace(t *testing.T) {
	dfa := NewDfa()
	dfa.LoadWords(common.GetWords())
	result := dfa.Replace("hello av java 毛泽东 sm 气枪 测试, 支付宝 ", 0)
	fmt.Println(result)
}

func TestB(t *testing.T) {
	ac := NewDfa()
	ac.LoadWords([]*common.SensitiveWords{
		&common.SensitiveWords{
			Word: "一二三四",
			Rank: 1,
		},
		&common.SensitiveWords{
			Word: "一二三",
			Rank: 1,
		},
		&common.SensitiveWords{
			Word: "四五六七",
			Rank: 1,
		},
		&common.SensitiveWords{
			Word: "三四",
			Rank: 1,
		},
	})
	result := ac.Search("一二三四五六七")
	log.Println("len:", len(result))
	for _, item := range result {
		fmt.Printf("%+v \n", item)
	}
}
