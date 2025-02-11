package main

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

const FileName = "hello.txt"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileName := flag.String("file", FileName, "The name for a file.")
	flag.Parse()
	logger.Info("Get File", "File", *fileName)
	wordMap := map[string]int{}
	re := regexp.MustCompile(`[\p{L}\d]+`)
	lines, err := getWordsFromFile(*fileName, re, logger)
	if err != nil {
		os.Exit(-1)
	}
	fillMap(wordMap, lines)
	logger.Info("Result", "Data", wordMap)
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

func getWordsFromFile(filename string, re *regexp.Regexp, logger *slog.Logger) ([]string, error) {
	lines := []string{}
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("Parse File", "Error", err)
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		lines = append(lines, re.FindAllString(str, -1)...)

	}
	return lines, nil
}
