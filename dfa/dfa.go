package dfa

import (
	"go-wordfilter/common"
	"strings"
)

type (
	DfaNode struct {
		Children map[rune]*DfaNode
		Rank     uint8
		End      bool
	}
	Dfa struct {
		Root *DfaNode
	}
)

func NewDfa() *Dfa {
	return &Dfa{
		Root: NewDfaNode(),
	}
}

func NewDfaNode() *DfaNode {
	return &DfaNode{
		Children: make(map[rune]*DfaNode),
	}
}

func (n *Dfa) LoadWords(words []*common.SensitiveWords) {
	for _, word := range words {
		n.add(word.Word, word.Rank)
	}
}

func (n *Dfa) add(word string, rank uint8) {
	chars := []rune(strings.ToLower(word))
	if len(chars) == 0 {
		return
	}
	nd := n.Root
	for _, char := range chars {
		if _, ok := nd.Children[char]; !ok {
			nd.Children[char] = NewDfaNode()
		}
		nd = nd.Children[char]
	}
	nd.Rank = rank
	nd.End = true
}

func (n *Dfa) Search(contentStr string) []*common.SearchItem {
	result := make([]*common.SearchItem, 0)
	chars := []rune(strings.ToLower(contentStr))
	size := len(chars)
	currentNode := n.Root
	for start, char := range chars {
		child, ok := currentNode.Children[char]
		if !ok {
			continue
		}
		if child.End {
			//if size < start-1 && common.IsWordCell(char) && common.IsWordCell(chars[start+1]) {
			//	continue
			//}
			result = append(result, &common.SearchItem{
				StartP: start,
				EndP:   start,
				Word:   string(chars[start : start+1]),
				Rank:   child.Rank,
			})
		}
		for end := start + 1; end < size; end++ {
			if _, ok := child.Children[chars[end]]; !ok {
				break
			}
			child = child.Children[chars[end]]
			if child.End {
				//if size < End-1 && common.IsWordCell(char) && common.IsWordCell(chars[End+1]) {
				//	continue
				//}
				//if start > 0 && common.IsWordCell(char) && common.IsWordCell(chars[start-1]) {
				//	continue
				//}
				result = append(result, &common.SearchItem{
					StartP: start,
					EndP:   end,
					Word:   string(chars[start : end+1]),
					Rank:   child.Rank,
				})
			}
		}
	}

	return result
}

func (n *Dfa) Replace(content string, rank uint8) *common.FindResponse {
	var res = new(common.FindResponse)
	res.BadWords = make(map[uint8][]*common.SearchItem)
	result := n.Search(content)
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
