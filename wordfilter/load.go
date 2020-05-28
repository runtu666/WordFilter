package wordfilter

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var serviceAC *Ac

type SensitiveWords struct {
	Id   int64
	Word string
	Rank uint8
}

func getAc() *Ac {
	return serviceAC
}
func setAc(ac *Ac) {
	serviceAC = ac
	log.Println("set ac:", ac)
}

func LoadWords() {
	log.Println("start load ")
	ac, err := loadWordDb()
	if err != nil {
		log.Fatal("load err:", err)
	}
	setAc(ac)
	log.Println("start load end  ")
}

func loadWordDb() (*Ac, error) {
	var wordList []*SensitiveWords
	NewConfig().MysqlConn.Find(&wordList)
	ac := NewAc()
	c := 0
	t1 := time.Now()
	for _, row := range wordList {
		c++
		ac.AddWord(row.Word, row.Rank)
	}
	t2 := time.Now()
	log.Println("load Words:", c, "sec:", t2.Sub(t1).Seconds())
	ac.Make()
	return ac, nil
}
