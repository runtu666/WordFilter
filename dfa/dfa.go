package dfa

import (
	"go-wordfilter/common"
	"strings"
)

type (
	trieNode struct {
		children map[rune]*trieNode
		rank     uint8
		end      bool
	}

	scope struct {
		start int
		stop  int
	}
)

func NewTire() *trieNode {
	n := new(trieNode)
	return n
}

func (n *trieNode) LoadWords(words []*common.SensitiveWords) {
	for _, word := range words {
		n.add(word.Word, word.Rank)
	}
}

func (n *trieNode) add(word string, rank uint8) {
	chars := []rune(strings.ToLower(word))
	if len(chars) == 0 {
		return
	}
	nd := n
	for _, char := range chars {
		if nd.children == nil {
			child := new(trieNode)
			nd.children = map[rune]*trieNode{
				char: child,
			}
			nd = child
		} else if child, ok := nd.children[char]; ok {
			nd = child
		} else {
			child := new(trieNode)
			nd.children[char] = child
			nd = child
		}
	}
	nd.rank = rank
	nd.end = true
}

func (n *trieNode) Search(contentStr string) []*common.SearchItem {
	result := make([]*common.SearchItem, 0)
	chars := []rune(strings.ToLower(contentStr))
	size := len(chars)
	for start, char := range chars {
		child, ok := n.children[char]
		if !ok {
			continue
		}
		if child.end {
			if size < start-1 && common.IsWordCell(char) && common.IsWordCell(chars[start+1]) {
				continue
			}
			result = append(result, &common.SearchItem{
				StartP: start,
				EndP:   start,
				Word:   string(chars[start : start+1]),
				Rank:   child.rank,
			})
		}
		for end := start + 1; end < size; end++ {
			if _, ok := child.children[chars[end]]; !ok {
				break
			}
			child = child.children[chars[end]]
			if child.end {
				if size < end-1 && common.IsWordCell(char) && common.IsWordCell(chars[end+1]) {
					continue
				}
				if start > 0 && common.IsWordCell(char) && common.IsWordCell(chars[start-1]) {
					continue
				}
				result = append(result, &common.SearchItem{
					StartP: start,
					EndP:   end,
					Word:   string(chars[start : end+1]),
					Rank:   child.rank,
				})
			}
		}
	}

	return result
}

func (n *trieNode) Replace(content string, rank uint8) *common.FindResponse {
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
