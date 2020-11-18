package ac

import (
	"go-wordfilter/common"
	"log"
	"strings"
	"time"
)

type (
	AcNode struct {
		Next map[rune]*AcNode `json:"next"`
		Fail *AcNode          `json:"fail"`
		End  bool             `json:"isWord"`
		Rank uint8            `json:"Rank"`
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

func (ac *Ac) LoadWords(words []*common.SensitiveWords) {
	t1 := time.Now()
	for _, row := range words {
		ac.AddWord(row.Word, row.Rank)
	}
	t2 := time.Now()
	log.Println("load Word:", len(words), "sec:", t2.Sub(t1).Seconds())
	ac.Make()
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
	tmp.End = true
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

func (ac *Ac) Search(contentStr string) []*common.SearchItem {
	rawContent := contentStr
	content := []rune(strings.ToLower(contentStr))
	p := ac.Root
	result := make([]*common.SearchItem, 0)
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
			if p.End {
				//# 若状态为词的结尾，则把词放进结果集
				//#判断当前这些位置是否为单词的边界
				if startWordIndex > 0 && common.IsWordCell(content[startWordIndex-1]) && common.IsWordCell(content[startWordIndex]) {
					//#当前字符和前面的字符都是字母,那么它是连续单词
					continue
				}
				if currentPosition < contentLen-1 && common.IsWordCell(content[currentPosition+1]) && common.IsWordCell(content[currentPosition]) {
					//#print '后面不是单词边界'
					continue
				}
				result = append(result, &common.SearchItem{
					StartP: startWordIndex,
					EndP:   currentPosition,
					Word:   string([]rune(rawContent)[startWordIndex : currentPosition+1]),
					Rank:   p.Rank,
				})
			}
		}

	}
	return result

}

func (ac *Ac) Replace(content string, rank uint8) *common.FindResponse {
	var res = new(common.FindResponse)
	res.BadWords = make(map[uint8][]*common.SearchItem)
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
