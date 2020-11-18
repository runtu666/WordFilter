package dfa

type (
	trieNode struct {
		children map[rune]*trieNode
		rank     uint8
		end      bool
	}
)

func NewTire() *trieNode {
	n := new(trieNode)
	return n
}

func (n *trieNode) LoadWords(words []string) {
	for _, word := range words {
		n.add(word)
	}
}

func (n *trieNode) add(word string, rank uint8) {
	chars := []rune(word)
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

func (n *trieNode) findKeywordScopes(chars []rune) []string {
	var words []string
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
			words = append(words, string(chars[start:i+1]))
		}

		for j := i + 1; j < size; j++ {
			grandchild, ok := child.children[chars[j]]
			if !ok {
				break
			}

			child = grandchild
			if child.end {
				words = append(words, string(chars[start:j+1]))
			}
		}

		start = -1
	}

	return words
}
