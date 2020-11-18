package wordfilter

import (
	"strings"
)

type (
	SearchItem struct {
		StartP int    `json:"start_p"`
		EndP   int    `json:"end_p"`
		Words  string `json:"words"`
		Rank   uint8  `json:"rank"`
	}
	AcNode struct {
		Next   map[rune]*AcNode `json:"next"`
		Fail   *AcNode          `json:"fail"`
		IsWord bool             `json:"isWord"`
		Rank   uint8            `json:"Rank"`
	}
	FindResponse struct {
		Status     uint8                   `json:"status"`
		NewContent string                  `json:"new_content"`
		ErrMsg     string                  `json:"err_msg"`
		BadWords   map[uint8][]*SearchItem `json:"bad_words"`
	}
	Ac struct {
		Root *AcNode `json:"root"`
	}
)

func NewAc() *Ac {
	return &Ac{
		Root: newAcNode(),
	}
}
func newAcNode() *AcNode {
	return &AcNode{
		Next: make(map[rune]*AcNode),
	}
}

func (ac *Ac) AddWord(word string, rank uint8) {
	chars := []rune(strings.ToLower(word))
	tmp := ac.Root
	for _, c := range chars {
		if _, ok := tmp.Next[c]; !ok {
			tmp.Next[c] = newAcNode()
		}
		tmp = tmp.Next[c]
	}
	tmp.IsWord = true
	tmp.Rank = rank
}

func (ac *Ac) Make() {
	queue := make([]*AcNode, 0)
	queue = append(queue, ac.Root)
	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]
		for k, child := range parent.Next {
			if parent == ac.Root {
				child.Fail = ac.Root
			} else {
				FailNode := parent.Fail
				for FailNode != nil {
					if _, ok := FailNode.Next[k]; ok {
						child.Fail = FailNode.Next[k]
						break
					}
					FailNode = FailNode.Fail
				}
				if FailNode == nil {
					child.Fail = ac.Root
				}
			}
			queue = append(queue, child)
		}
	}
}

func (ac *Ac) Search(contentStr string) []*SearchItem {
	rawContent := contentStr
	content := []rune(strings.ToLower(contentStr))
	p := ac.Root
	result := make([]*SearchItem, 0)
	startWordIndex := 0
	contentLen := len(content)
	for currentPosition, word := range content {
		// 检索状态机，直到匹配
		for _, ok := p.Next[word]; !ok && p != ac.Root; {
			//直到找到失败节点，或者找到根节点
			p = p.Fail
		}
		if _, ok := p.Next[word]; ok {
			if p == ac.Root {
				//# 若当前节点是根且存在转移状态，则说明是匹配词的开头，记录词的起始位置
				startWordIndex = currentPosition
			}
			//# 转移状态机的状态
			p = p.Next[word]
			if p.IsWord {
				//# 若状态为词的结尾，则把词放进结果集
				//#判断当前这些位置是否为单词的边界
				if startWordIndex > 0 && isWordCell(content[startWordIndex-1]) && isWordCell(content[startWordIndex]) {
					//#当前字符和前面的字符都是字母,那么它是连续单词
					continue
				}
				if currentPosition < contentLen-1 && isWordCell(content[currentPosition+1]) && isWordCell(content[currentPosition]) {
					//#print '后面不是单词边界'
					continue
				}
				result = append(result, &SearchItem{
					StartP: startWordIndex,
					EndP:   currentPosition,
					Words:  string([]rune(rawContent)[startWordIndex : currentPosition+1]),
					Rank:   p.Rank,
				})
			}
		}

	}
	return result

}

func (ac *Ac) Replace(content string, rank uint8) *FindResponse {
	var res = new(FindResponse)
	res.BadWords = make(map[uint8][]*SearchItem)
	if ac == nil {
		res.ErrMsg = "ac is nil"
		return res
	}
	result := ac.Search(content)
	contentBuff := []rune(content)
	for _, item := range result {
		if item.Rank > rank && rank != 0 {
			continue
		}
		for i := item.StartP; i <= item.EndP; i++ {
			contentBuff[i] = '*'
		}
		res.BadWords[item.Rank] = append(res.BadWords[item.Rank], item)
	}
	res.Status = 0
	res.NewContent = string(contentBuff)
	return res
}

var wordCells = make(map[rune]bool)

func isWordCell(s rune) bool {
	_, ok := wordCells[s]
	return ok
}

func init() {
	//#预先生成好组成单词的字符
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-" {
		wordCells[c] = true
	}

}
