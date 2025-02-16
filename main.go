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

// не экспортируется из пакета, поэтому можно lowercase
const defaultFileName = "hello.txt"

type wordFreq struct {
	Word  string
	Count int
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileName := flag.String("file", defaultFileName, "The name for a file.")
	flag.Parse()

	// используем Fields, все ошибки пишем в lowercase
	logger.Info("reading file", slog.String("file", *fileName))

	wordMap := map[string]int{}
	re := regexp.MustCompile(`[\p{L}\d]+`)
	err := getWordsFromFile(*fileName, re, logger, wordMap)
	if err != nil {
		// не было  записи ошибки в лог
		logger.Error("unable to get words from file", slog.String("err", err.Error()))
		return
	}

	wordres := mapToStr(wordMap)
	sort.SliceStable(wordres, func(i, j int) bool {
		return wordres[i].Count > wordres[j].Count
	})

	logger.Info("result", "data", wordres)
}

func fillMap(gm map[string]int, str []string) {
	// небольшая оптимизация логики
	// btw, не очень понятно что такое gm
	// давай именовать более понятными именами
	for _, w := range str {
		gm[strings.ToLower(w)]++
	}
}

func getWordsFromFile(filename string, re *regexp.Regexp, logger *slog.Logger, wm map[string]int) error {
	file, err := os.Open(filename)
	if err != nil {
		// была некорректная запись в лог
		logger.Error(
			"unable to read file",
			slog.String("filename", filename),
			slog.String("err", err.Error()),
		)
		// оборачиваем ошибку, чтобы знать где именно она возникла по сообщению ошибки
		return fmt.Errorf("unable to read file: %w", err)
	}
	// потеряли запись ошибки Close в лог
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("error while closing file", slog.String("err", err.Error()))
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// в целом норм, ты тут читаешь до \n
		// но что если строка будет очень большой и файл будет в одну строку и весить 5GB?
		// в таком случае надо читать буффером
		// создавать буффер на несколько KB (например 32) и вычитывать кусками
		// https://stackoverflow.com/a/36111861
		fillMap(wm, re.FindAllString(scanner.Text(), -1))
	}

	return nil
}

// метод называется map to string, но возвращается слайс структур, а не string :)
func mapToStr(m map[string]int) []wordFreq {
	// можно добавить capacity, чтобы сразу аллоцировать нужный размер
	// и не сталкиваться постоянно с реаллокацией и эвакуацией из мапы
	words := make([]wordFreq, 0, len(m))
	for k, v := range m {
		words = append(words, wordFreq{k, v})
	}
	return words
}
