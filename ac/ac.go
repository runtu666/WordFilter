package ac

import (
	"log"
	"strings"
	"time"

	"go-wordfilter/common"
)

type (
	AcNode struct {
		Children map[rune]*AcNode `json:"next"`
		//Prev *AcNode          `json:"prev"`
		Position int     `json:"position"`
		Fail     *AcNode `json:"fail"`
		End      bool    `json:"isWord"`
		Rank     int     `json:"Rank"`
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
		Children: make(map[rune]*AcNode),
	}
}

func (ac *Ac) LoadWords(words []*common.SensitiveWords) {
	t1 := time.Now()
	for _, row := range words {
		ac.AddWord(row.Word, row.Rank)
	}
	log.Println("load Word:", len(words), "sec:", time.Now().Sub(t1).Seconds())
	ac.Make()
}

func (ac *Ac) AddWord(word string, rank int) {
	chars := []rune(strings.ToLower(word))
	nd := ac.Root
	if len(chars) == 0 {
		return
	}
	for i, c := range chars {
		if _, ok := nd.Children[c]; !ok {
			nd.Children[c] = newAcNode()
		}
		nd.Children[c].Position = i
		nd = nd.Children[c]
	}
	nd.End = true
	nd.Rank = rank
}

func (ac *Ac) Make() {
	queue := make([]*AcNode, 0)
	queue = append(queue, ac.Root)
	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]
		for k, child := range parent.Children {
			if parent == ac.Root {
				child.Fail = ac.Root
			} else {
				FailNode := parent.Fail
				for FailNode != nil {
					if _, ok := FailNode.Children[k]; ok {
						child.Fail = FailNode.Children[k]
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
	//contentLen := len(content)
	for currentPosition, word := range content {
		// 检索状态机，直到匹配
		for p.Children[word] == nil && p != ac.Root {
			//直到找到失败节点，或者找到根节点
			p = p.Fail
		}
		if _, ok := p.Children[word]; ok {
			//# 转移状态机的状态
			p = p.Children[word]
			if p.End {
				startWordIndex := currentPosition - p.Position

				//#当前字符和前面的字符都是字母,那么它是连续单词
				if startWordIndex > 0 && common.IsWordCell(content[startWordIndex-1]) && common.IsWordCell(content[startWordIndex]) {
					continue
				}
				//if currentPosition < contentLen-1 && common.IsWordCell(content[currentPosition+1]) && common.IsWordCell(content[currentPosition]) {
				//	//#print '后面不是单词边界'
				//	continue
				//}
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

func (ac *Ac) Replace(content string, rank int) *common.FindResponse {
	var res = new(common.FindResponse)
	res.BadWords = make(map[int][]*common.SearchItem)
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
