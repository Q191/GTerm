package metadata

var keywordTrie *Trie

func init() {
	keywordTrie = NewTrie()
	for vendor, keywords := range networkVendorKeywords {
		for _, keyword := range keywords {
			keywordTrie.Insert(keyword, vendor, "network")
		}
	}
	for vendor, keywords := range serverVendorKeywords {
		for _, keyword := range keywords {
			keywordTrie.Insert(keyword, vendor, "server")
		}
	}
}

type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
	vendor      string
	deviceType  string
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(keyword, vendor, deviceType string) {
	node := t.root
	for _, char := range keyword {
		if _, exists := node.children[char]; !exists {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[char]
	}
	node.isEndOfWord = true
	node.vendor = vendor
	node.deviceType = deviceType
}
