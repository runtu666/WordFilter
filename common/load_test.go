package common

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	fi, err := os.Open("../words.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	var words = make([]*SensitiveWords, 0)
	br := bufio.NewReader(fi)
	i := 1
	for {
		w, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		word := string(w)
		if word != "" {
			words = append(words, &SensitiveWords{
				Word: word,
				Rank: i%5 + 1,
			})
			i++
		}
	}
	marshal, err := json.Marshal(words)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("../bad_words.json", marshal, 0644)
}
