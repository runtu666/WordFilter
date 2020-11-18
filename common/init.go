package common

var wordCells = make(map[rune]bool)

func IsWordCell(s rune) bool {
	_, ok := wordCells[s]
	return ok
}

func init() {
	//#预先生成好组成单词的字符
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-" {
		wordCells[c] = true
	}

}
