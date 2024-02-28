	ac := NewAcDoubleArrayTrie()
	
	words := []string{"i", "he", "his", "she", "hers"}
	root := ac.BuildTrie(words)
	ac.BuildFailPointer(root)
	
	content := "ifindhehishehersall"
	result := ac.Search([]rune(content))
	for _, res := range result {
		fmt.Printf("%-3d %-3d %s\n", res.Begin, res.End, string(res.Value))
	}