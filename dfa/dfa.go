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
