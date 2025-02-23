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

const (
	defaultFileName = "hello.txt"
	bufferSize      = 4096
)

type wordFreq struct {
	Word  string
	Count int
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileName := flag.String("file", defaultFileName, "The name for a file.")
	flag.Parse()
	logger.Info("reading file", slog.String("file", *fileName))
	wordMap := map[string]int{}
	expression := regexp.MustCompile(`[\p{L}\d]+`)
	err := getWordsFromFile(*fileName, expression, logger, wordMap)
	if err != nil {
		logger.Error("unable to get words from file", slog.String("err", err.Error()))
		os.Exit(-1)
	}
	wordResult := mapToSlice(wordMap)
	sort.SliceStable(wordResult, func(i, j int) bool {
		return wordResult[i].Count > wordResult[j].Count
	})
	logger.Info("result", "data", wordResult)
	output(wordResult)
}

func fillMap(wordMap map[string]int, str []string) {
	for _, word := range str {
		wordMap[strings.ToLower(word)]++
	}
}

func getWordsFromFile(filename string, regexp *regexp.Regexp, logger *slog.Logger, wordMap map[string]int) error {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error(
			"unable to read file",
			slog.String("filename", filename),
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("unable to read file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("error while closing file", slog.String("err", err.Error()))
		}
	}()
	reader := bufio.NewReaderSize(file, bufferSize)
	buffer := make([]byte, bufferSize)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			text := string(buffer[:n])
			words := regexp.FindAllString(text, -1)
			fillMap(wordMap, words)
		}
		if err != nil {
			if err.Error() != "EOF" {
				logger.Error("error while reading file", slog.String("err", err.Error()))
			}
			break
		}
	}

	return nil
}

func mapToSlice(m map[string]int) []wordFreq {
	words := make([]wordFreq, 0, len(m))
	for k, v := range m {
		words = append(words, wordFreq{k, v})
	}
	return words
}

func output(wordsFre []wordFreq) {
	for _, kv := range wordsFre {
		fmt.Printf("%s %d\n", kv.Word, kv.Count)
	}
}
