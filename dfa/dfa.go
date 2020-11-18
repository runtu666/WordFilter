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
	start := -1
	for i := 0; i < size; i++ {
		child, ok := n.children[chars[i]]
		if !ok {
			continue
		}

		if start < 0 {
			start = i
		}
		if child.end {
			result = append(result, &common.SearchItem{
				StartP: start,
				EndP:   i + 1,
				Words:  string(chars[start : i+1]),
				Rank:   child.rank,
			})

		}
		for j := i + 1; j < size; j++ {
			grandchild, ok := child.children[chars[j]]
			if !ok {
				break
			}
			child = grandchild
			if child.end {
				result = append(result, &common.SearchItem{
					StartP: start,
					EndP:   j + 1,
					Words:  string(chars[start : j+1]),
					Rank:   child.rank,
				})
			}
		}
		start = -1
	}

	return result
}
