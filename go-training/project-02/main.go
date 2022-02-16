package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string, value int) {
	c.mu.Lock()
	c.v[key] += value
	c.mu.Unlock()
}

func (c *SafeCounter) Value() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

func readFile(fname string, c chan map[string]int) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordCountFile := map[string]int{}
	for scanner.Scan() {
		wordCountFile[processWord(scanner.Text())] += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	c <- wordCountFile
}

func processWord(w string) string {
	var word string
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	word = reg.ReplaceAllString(strings.ToLower(w), "")
	return word
}

func main() {
	totalCount := SafeCounter{v: make(map[string]int)}

	// Initialize channel for file level word count
	ch := make(chan map[string]int)
	// Run go routines for each file to produce file level word count
	for _, fname := range os.Args[1:] {
		go readFile(fname, ch)
	}

	// Initialize channel for total level word count
	fileWordCount_ch := <-ch
	// Run go routines for each file level word count to produce total word count
	for key, value := range fileWordCount_ch {
		go totalCount.Inc(key, value)
	}

	// Prepare to print total word count list from all files
	totalCountList := totalCount.Value()

	// Create list of keys to sort
	keys := make([]string, 0, len(totalCountList))
	for k := range totalCountList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Print all word and corresponding frequency in alphabetical order
	for _, k := range keys {
		fmt.Printf("%v %v\n", k, totalCountList[k])
	}

}

// References used:
// https://golangbyexample.com/read-large-file-word-by-word-go/
