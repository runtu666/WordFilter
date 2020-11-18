package common

import (
	"encoding/json"
	"io/ioutil"
)

func GetWords() []*SensitiveWords {
	var wordList []*SensitiveWords
	f, err := ioutil.ReadFile("../bad_words.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &wordList)
	if err != nil {
		panic(err)
	}
	return wordList
}
