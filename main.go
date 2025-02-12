package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
)

const FileName = "hello.txt"

type wordFreeq struct {
	Word  string
	Count int
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileName := flag.String("file", FileName, "The name for a file.")
	flag.Parse()
	logger.Info("Get File", "File", *fileName)
	wordMap := map[string]int{}
	re := regexp.MustCompile(`[\p{L}\d]+`)
	err := getWordsFromFile(*fileName, re, logger, wordMap)
	if err != nil {
		os.Exit(-1)
	}
	wordres := mapToStr(wordMap)
	sort.SliceStable(wordres, func(i, j int) bool {
		return wordres[i].Count > wordres[j].Count
	})
	logger.Info("Result", "Data", wordres)
	output(wordres)
}

func fillMap(gm map[string]int, str []string) {
	for _, w := range str {
		w = strings.ToLower(w)
		_, ok := gm[w]
		if ok {
			gm[w]++
		} else {
			gm[w] = 1
		}
	}
}

func getWordsFromFile(filename string, re *regexp.Regexp, logger *slog.Logger, wm map[string]int) error {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("Parse File", "Error", err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fillMap(wm, re.FindAllString(scanner.Text(), -1))
	}

	return nil
}

func mapToStr(m map[string]int) []wordFreeq {
	words := make([]wordFreeq, 0)
	for k, v := range m {
		words = append(words, wordFreeq{k, v})
	}
	return words
}

func output(wordsFre []wordFreeq) {
	for _, kv := range wordsFre {
		fmt.Printf("%s %d\n", kv.Word, kv.Count)
	}
}
