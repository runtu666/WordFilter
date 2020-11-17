package wordfilter

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"time"
)

var serviceAC *Ac

type SensitiveWords struct {
	Id   int64  `json:"id"`
	Word string `json:"word"`
	Rank uint8  `json:"rank"`
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
	ac, err := loadWords()
	if err != nil {
		log.Fatal("load err:", err)
	}
	setAc(ac)
	log.Println("start load end  ")
}

func loadWords() (*Ac, error) {
	return loadWordFile()
}

func loadWordFile() (*Ac, error) {
	var wordList []*SensitiveWords
	f, err := ioutil.ReadFile("../bad_words.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &wordList)
	if err != nil {
		return nil, err
	}
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
