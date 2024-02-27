package double_array_trie

import (
	"sort"
	"strings"
)

const (
	rootIndex int = 0
	failState int = -1
	extraSize int = 10
)

type DoubleArrayTrie struct {
	base  []int
	check []int
}

type AcDoubleArrayTrie struct {
	DoubleArrayTrie
	fail    []int              // fail array
	lengths []map[int]struct{} // lengths of words
}

func (dat *AcDoubleArrayTrie) GetBase() []int {
	return dat.base
}

func (dat *AcDoubleArrayTrie) GetCheck() []int {
	return dat.check
}

func (dat *AcDoubleArrayTrie) GetFail() []int {
	return dat.fail
}

func (dat *AcDoubleArrayTrie) GetLengths() []map[int]struct{} {
	return dat.lengths
}

func NewAcDoubleArrayTrie() *AcDoubleArrayTrie {
	return &AcDoubleArrayTrie{
		DoubleArrayTrie: DoubleArrayTrie{
			base:  make([]int, 0),
			check: make([]int, 0),
		},
		fail:    make([]int, 0),
		lengths: make([]map[int]struct{}, 0),
	}
}

// node
type Node struct {
	code  rune
	index int // index of base and check
	depth int
	left  int // fist child
	right int // last child
	end   bool

	children []*Node
}

func NewNode() *Node {
	return &Node{
		code:  0,
		index: 0,
		depth: 0,
		left:  0,
		right: 0,
		end:   false,

		children: make([]*Node, 0),
	}
}

// start of trie
type Trie struct {
	root *Node
}

// //////////////////////////////////////
func (ac *AcDoubleArrayTrie) resize(newSize int) {
	oldSize := len(ac.base)
	if newSize < oldSize {
		return
	}

	expand := func(src []int, size int) []int {
		dst := make([]int, size, size)
		copy(dst, src)
		return dst
	}

	newSize += extraSize
	ac.base = expand(ac.base, newSize)
	ac.check = expand(ac.check, newSize)
	ac.fail = expand(ac.fail, newSize)

	dst := make([]map[int]struct{}, newSize, newSize)
	copy(dst, ac.lengths)
	ac.lengths = dst
}

func (ac *AcDoubleArrayTrie) BuildTrie(keyWords []string) *Node {
	words := make([][]rune, 0)

	var maxCode rune = 0
	for _, keyWord := range keyWords {
		keyWord = strings.TrimSpace(keyWord)

		if len(keyWord) == 0 {
			continue
		}

		for _, code := range keyWord {
			if code > maxCode {
				maxCode = code
			}
		}

		words = append(words, []rune(keyWord))
	}

	// sort []rune, not string
	sort.Slice(words, func(i, j int) bool {
		minLen := len(words[i])
		if minLen > len(words[j]) {
			minLen = len(words[j])
		}

		for k := 0; k < minLen; k++ {
			if words[i][k] < words[j][k] {
				return true
			} else if words[i][k] > words[j][k] {
				return false
			} else {
				continue
			}
		}

		if len(words[i]) < len(words[j]) {
			return true
		}

		return false
	})

	ac.resize(int(maxCode))

	trie := &Trie{
		root: NewNode(), // root' index is zero
	}
	trie.root.right = len(words)

	queue := make([]*Node, 0)
	queue = append(queue, trie.root)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		children := ac.getChildren(node, words)

		if len(children) > 0 {
			queue = append(queue, children...)
		}

		ac.setBaseAndCheck(node)

		if node.end {
			length := len(words[node.left])
			lengths := make(map[int]struct{})
			ac.resize(length)
			lengths[length] = struct{}{}
			ac.addLens(node.index, lengths)
		}
	}

	return trie.root
}

func (ac *AcDoubleArrayTrie) addLens(index int, lengths map[int]struct{}) {
	for length := range lengths {
		if len(ac.lengths[index]) == 0 {
			ac.lengths[index] = make(map[int]struct{})
		}

		ac.lengths[index][length] = struct{}{}
	}
}

func (ac *AcDoubleArrayTrie) getChildren(parent *Node, words [][]rune) []*Node {
	var preCode rune = 0
	for i := parent.left; i < parent.right; i++ {
		var curCode rune = 0
		if len(words[i]) != parent.depth {
			curCode = words[i][parent.depth]
		} else {
			parent.end = true
		}

		if preCode != curCode {
			node := &Node{
				code:  curCode,
				depth: parent.depth + 1,
				left:  i,
			}

			if len(parent.children) > 0 {
				parent.children[len(parent.children)-1].right = i
			}
			parent.children = append(parent.children, node)
		}

		preCode = curCode
	}

	if len(parent.children) > 0 {
		parent.children[len(parent.children)-1].right = parent.right
	}

	return parent.children
}

func (ac *AcDoubleArrayTrie) setBaseAndCheck(parent *Node) {
	base := ac.getBase(parent)
	ac.resize(parent.index)
	ac.base[parent.index] = base // base[s]+c = t, reset parent's base value

	index := 0
	for _, child := range parent.children {
		index = base + int(child.code)
		child.index = index

		ac.resize(index)
		ac.base[index] = base          // placeholder, may not be the corrent base value
		ac.check[index] = parent.index // check[t] = s
	}
}

func (ac *AcDoubleArrayTrie) getBase(parent *Node) int {
	if len(parent.children) == 0 {
		// different from the paper, negative values are not set.
		return 0 // base=0, has no children
	}

	firstChildCode := int(parent.children[0].code)
	pos := firstChildCode
	base := 0
	for {
	next:
		pos++

		// index begin from 0, use base/check/lengths array to check whether the position is occupied
		if len(ac.base) > pos &&
			(ac.base[pos] != 0 || ac.check[pos] != 0 || len(ac.lengths[pos]) > 0) {
			continue
		}

		base = pos - firstChildCode
		for _, child := range parent.children {
			index := base + int(child.code)
			if len(ac.base) > index &&
				(ac.base[index] != 0 || ac.check[index] != 0 || len(ac.lengths[index]) > 0) {
				goto next
			}
		}

		break
	}
	return base
}

func (ac *AcDoubleArrayTrie) findSameCode(parentIndex int, code rune) int {
	base := ac.base[parentIndex] // base >= 0
	childIndex := base + int(code)

	if len(ac.base) <= childIndex {
		if parentIndex == rootIndex {
			return rootIndex
		}
		return failState
	}

	if (ac.base[childIndex] != 0 || ac.check[childIndex] != 0 || len(ac.lengths[childIndex]) > 0) &&
		parentIndex == ac.check[childIndex] {
		return childIndex
	}

	if parentIndex == rootIndex {
		return rootIndex
	}

	return failState
}

func (ac *AcDoubleArrayTrie) BuildFailPointer(root *Node) {
	queue := make([]*Node, 0)

	// depth=0 and depth=1, fail pointer point to rootIndex
	ac.fail[root.index] = rootIndex
	for _, child := range root.children {
		ac.fail[child.index] = rootIndex
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		parent := queue[0]
		queue = queue[1:]

		for _, child := range parent.children {
			// parent.index's fail pointer point to failIndex
			parentFailIndex := ac.fail[parent.index]
			failPointer := ac.findSameCode(parentFailIndex, child.code)
			for failPointer == failState {
				parentFailIndex = ac.fail[parentFailIndex]
				failPointer = ac.findSameCode(parentFailIndex, child.code)
			}

			ac.fail[child.index] = failPointer
			if len(ac.lengths[failPointer]) > 0 {
				ac.addLens(child.index, ac.lengths[failPointer])
			}

			queue = append(queue, child)
		}
	}
}

type Hit struct {
	Begin int
	End   int
	Value []rune
}

func (ac *AcDoubleArrayTrie) Search(content []rune) []*Hit {
	hits := []*Hit{}
	currentIndex := rootIndex
	for index, code := range content {
		nextIndex := ac.findSameCode(currentIndex, code)
		for nextIndex == failState {
			currentIndex = ac.fail[currentIndex]
			nextIndex = ac.findSameCode(currentIndex, code)
		}

		currentIndex = nextIndex

		for length := range ac.lengths[currentIndex] {
			if length == 0 {
				continue
			}

			hit := &Hit{
				Begin: index + 1 - length,
				End:   index,
				Value: make([]rune, length),
			}
			copy(hit.Value, content[index+1-length:index+1])
			hits = append(hits, hit)
		}
	}

	return hits
}

// ignore whitelist
func IgnoreWithWhiteList(acHits []*Hit, whitelistHits []*Hit) []*Hit {
	hits := make([]*Hit, 0)

	ignore := func(hit *Hit, whitelists []*Hit) bool {
		if hit.Begin > hit.End {
			hit.Begin, hit.End = hit.End, hit.Begin
		}

		for _, whitelist := range whitelists {
			if whitelist.Begin > whitelist.End {
				whitelist.Begin, whitelist.End = whitelist.End, whitelist.Begin
			}

			if hit.Begin >= whitelist.Begin &&
				hit.End <= whitelist.End {
				return true
			}
		}

		return false
	}

	for _, acHit := range acHits {
		if ignore(acHit, whitelistHits) {
			continue
		}

		hit := &Hit{
			Begin: acHit.Begin,
			End:   acHit.End,
			Value: make([]rune, len(acHit.Value)),
		}
		copy(hit.Value, acHit.Value)
		hits = append(hits, hit)
	}

	return hits
}
