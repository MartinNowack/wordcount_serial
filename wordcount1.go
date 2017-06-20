// Reads a file and counts how often each word exists
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
)

func main() {
	// Initialize profiling support (Remove before deploying in production)
	//	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	//	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	inFilePtr := flag.String("i", "", "Input file")
	flag.Parse()

	inputF, err := os.Open(*inFilePtr)
	defer inputF.Close()
	check(err)

	lines := make([]string, 0)

	lineitem := make([]byte, 0)
	for {
		buffer := make([]byte, 1)
		n, err := inputF.Read(buffer)
		if n == 0 && err == io.EOF {
			break
		}
		check(err)

		// check if we read an EOL \n
		if buffer[0] == '\r' {
			continue
		}
		if buffer[0] == '\n' {
			lines = append(lines, string(lineitem))
			lineitem = []byte{}
			continue
		}
		lineitem = append(lineitem, buffer[0])
	}

	fmt.Print("## Read text:\n")
	for _, v := range lines {
		fmt.Print(string(v))
	}
	fmt.Print("\n\n")

	fmt.Print("## Split lines:\n")
	split_lines := make([][]string, 0)
	for _, v := range lines {
		split_lines = append(split_lines, splitString2(string(v)))
	}

	var wordCounters map[string]int = make(map[string]int)

	for _, line := range split_lines {
		for _, word := range line {
			wordCounters[word]++
		}
	}

	for word, counter := range wordCounters {
		fmt.Printf("'%s': %v\n", word, counter)
	}

	fmt.Printf("\nTop 10\n")

	for k, v := range getTopN(wordCounters, 10) {
		fmt.Printf("'%s': %v\n", k, v)
	}
}

func getTopN(wordMap map[string]int, n int) map[string]int {
	word_pairs := make(PairList, len(wordMap))
	i := 0
	for k, v := range wordMap {
		word_pairs[i] = StrCntPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(word_pairs))

	result := make(map[string]int, n)
	for _, v := range word_pairs[:n] {
		result[v.Key] = v.Value
	}

	return result
}

func splitString2(s string) []string {
	re := regexp.MustCompile(`[^a-zA-Z]`)
	fields := re.Split(s, -1)

	results := make([]string, 0)
	for _, f := range fields {
		if f != "" {
			results = append(results, f)
		}
	}
	return results
}

type StrCntPair struct {
	Key   string
	Value int
}
type PairList []StrCntPair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Check for errors and quit if needed
func check(e error) {
	if e != nil {
		panic(e)
	}
}
