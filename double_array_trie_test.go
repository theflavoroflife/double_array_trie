package double_array_trie

import (
	"fmt"
	"sort"
	"testing"
)

func TestNewAcDoubleArrayTrieSort(t *testing.T) {
	a := make([][]rune, 0)
	a = append(a, []rune("i"))
	a = append(a, []rune("he"))
	a = append(a, []rune("his"))
	a = append(a, []rune("she"))
	a = append(a, []rune("hers"))
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

func PrintWords(ac *AcDoubleArrayTrie, words []string) {
	fmt.Printf("print words:\n")
	codes := make(map[rune]struct{})
	for _, word := range words {
		for _, v := range word {
			codes[v] = struct{}{}
		}
	}

	for code := range codes {
		fmt.Printf("%s:%d ", string(code), code)
	}

	fmt.Printf("\n")
}

func PrintBaseAndCheck(ac *AcDoubleArrayTrie, params ...int) {
	fmt.Println("base and check")
	fn := func(i int) bool {
		if ac.base[i] == 0 {
			if len(params) > 1 {
				if i >= params[0] && i <= params[1] {
					return false
				}
			}
			return true
		}

		return false
	}

	var s string
	s = "\nindex:   "
	fmt.Printf("%s", s)
	for i := 0; i < len(ac.base); i++ {
		if fn(i) {
			continue
		}
		fmt.Printf("%-4d ", i)
	}
	s = "\nbase:    "
	fmt.Printf("%s", s)
	for i := 0; i < len(ac.base); i++ {
		if fn(i) {
			continue
		}
		fmt.Printf("%-4d ", ac.base[i])
	}
	s = "\ncheck:   "
	fmt.Printf("%s", s)
	for i := 0; i < len(ac.base); i++ {
		if fn(i) {
			continue
		}
		fmt.Printf("%-4d ", ac.check[i])
	}
	s = "\nfail:    "
	fmt.Printf("%s", s)
	for i := 0; i < len(ac.base); i++ {
		if fn(i) {
			continue
		}
		fmt.Printf("%-4d ", ac.fail[i])
	}
	s = "\nlengths: "
	fmt.Printf("%s", s)
	for i := 0; i < len(ac.base); i++ {
		if fn(i) {
			continue
		}
		var lenStr string
		if len(ac.lengths[i]) > 1 {
			TotalLen := len(ac.lengths[i])
			var j int
			var lastNum int
			for length := range ac.lengths[i] {
				lastNum = length
				if j+1 == TotalLen {
					break
				}
				j++
				lenStr += fmt.Sprintf("%d,", length)
			}
			lenStr += fmt.Sprintf("%d", lastNum)

			fmt.Printf("%-5s", lenStr)
		} else {
			lens := 0
			for length := range ac.lengths[i] {
				lens = length
			}
			fmt.Printf("%-4d ", lens)
		}
	}
	fmt.Printf("\n")
}

func TestNewAcDoubleArrayTrieAnother(t *testing.T) {
	ac := NewAcDoubleArrayTrie()

	words := []string{"i", "he", "his", "she", "hers"}
	root := ac.BuildTrie(words)
	PrintWords(ac, words)
	ac.BuildFailPointer(root)
	PrintBaseAndCheck(ac, 100, 120)

	//content := "i he his she hers"
	content := "ifindhehishehersall"
	result := ac.Search([]rune(content))
	for _, res := range result {
		fmt.Printf("%-3d %-3d %s\n", res.Begin, res.End, string(res.Value))
	}
}
