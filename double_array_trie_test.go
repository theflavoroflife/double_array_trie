package double_array_trie

import (
	"fmt"
	"sort"
	"testing"
)

func TestNewAcDoubleArrayTrieSort(t *testing.T) {
	a := make([][]rune, 0)
	a = append(a, []rune("aquaj"))
	a = append(a, []rune("jack"))
	a = append(a, []rune("ack"))
	a = append(a, []rune("card"))
	a = append(a, []rune("care"))
	sort.Slice(a, func(i, j int) bool {
		length := len(a[i])
		if length > len(a[j]) {
			length = len(a[j])
		}

		for k := 0; k < length; k++ {
			if a[i][k] < a[j][k] {
				return true
			} else if a[i][k] > a[j][k] {
				return false
			} else {
				continue
			}
		}

		if len(a[i]) < len(a[j]) {
			return true
		}

		return false
	})

	for _, v := range a {
		fmt.Printf("%v ", string(v))
	}
	fmt.Println("")
	fmt.Println(a)
}

func TestNewAcDoubleArrayTrie(t *testing.T) {
	ac := NewAcDoubleArrayTrie()

	words := []string{"aquaj", "jack", "ack", "card", "care"}
	root := ac.BuildTrie(words)
	ac.BuildFailPointer(root)

	content := "aquajackcarecardcatdogcancare"
	result := ac.Search([]rune(content))
	for _, res := range result {
		fmt.Printf("%-3d %-3d %s\n", res.Begin, res.End, string(res.Value))
	}
}

func TestNewAcDoubleArrayTrieSearchWithWhiteList(t *testing.T) {

	fn := func(words []string) *AcDoubleArrayTrie {
		ac := NewAcDoubleArrayTrie()
		root := ac.BuildTrie(words)
		ac.BuildFailPointer(root)
		return ac
	}

	// dirty world
	dirtyWords := []string{
		"日",
		"156",
		"不思蜀",
		"乐生",
		"乐不",
	}

	// whitelist dirty world
	whitelistWords := []string{
		"生日快乐",
		"1565",
		"乐不思蜀",
		"生日",
		"快乐",
	}

	content := "123456祝你生日快乐生日快乐不思蜀！1565456"

	ac := fn(dirtyWords)
	whitelistAc := fn(whitelistWords)

	acHits := ac.Search([]rune(content))
	whitelistHits := whitelistAc.Search([]rune(content))
	for _, v := range whitelistHits {
		t.Log("whitelist: ", v.Begin, v.End, string(v.Value))
	}

	result := IgnoreWithWhiteList(acHits, whitelistHits)

	for _, res := range result {
		fmt.Printf("dirty word: %-3d %-3d %s\n", res.Begin, res.End, string(res.Value))
	}
}
