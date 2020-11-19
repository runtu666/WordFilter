package ac

import (
	"fmt"
	"go-wordfilter/common"
	"log"
	"testing"
)

func TestAC(t *testing.T) {
	ac := NewAc()
	ac.LoadWords(common.GetWords())
	result := ac.Search("hello av java 毛泽东 sm 气枪 测试, 支付宝 ")
	log.Println("len:", len(result))
	for _, item := range result {
		fmt.Printf("%+v \n", item)
	}
}

func TestB(t *testing.T) {
	ac := NewAc()
	ac.LoadWords([]*common.SensitiveWords{
		&common.SensitiveWords{
			Word: "三四",
			Rank: 1,
		},
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
	})
	result := ac.Search("一二三四五六七")
	for _, item := range result {
		fmt.Printf("%+v \n", item)
	}
}

func TestAc_Replace(t *testing.T) {
	ac := NewAc()
	ac.LoadWords(common.GetWords())
	result := ac.Replace("hello av java 毛泽东 sm 气枪 测试, 支付宝 ", 0)
	fmt.Printf("%+v\n", result.NewContent)
}

func BenchmarkTire(b *testing.B) {
	b.ReportAllocs()
	ac := NewAc()
	ac.LoadWords(common.GetWords())
	for i := 0; i < b.N; i++ {
		ac.Replace("hello av java 毛泽东 sm 气枪 测试, 支付宝 ", 0)
	}
}
