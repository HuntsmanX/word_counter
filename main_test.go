package main

import (
	"log/slog"
	"os"
	"regexp"
	"testing"
)

const TestFile = "tmp.txt"

func TestFillMap(t *testing.T) {
	wordMap := map[string]int{}
	words := []string{"Test", "test", "TEST", "Hello", "HELLO"}
	fillMap(wordMap, words)
	if wordMap["test"] != 3 {
		t.Errorf("Expected 'test' to appear 3 times, got %d", wordMap["test"])
	}
	if wordMap["hello"] != 2 {
		t.Errorf("Expected 'test' to appear 2 times, got %d", wordMap["hello"])
	}
}

func TestGetWordsFromFile(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wordMap := map[string]int{}
	re := regexp.MustCompile(`[\p{L}\d]+`)
	ts := []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa.")
	c := 13
	_, err := os.Stat(TestFile)
	if err != nil {
		err := os.WriteFile(TestFile, ts, 0644)
		if err != nil {
			t.Errorf("Error creating file %s", TestFile)
		} else {
			getWordsFromFile(TestFile, re, logger, wordMap)
			if len(wordMap) != c {
				t.Errorf("Expected 13 words got %d", len(wordMap))
			}
			err := os.Remove(TestFile)
			if err != nil {
				t.Errorf("Failed to remove test file %s", TestFile)
			}
		}
	} else {
		err := os.Remove(TestFile)
		if err != nil {
			t.Errorf("Failed to remove test file %s", TestFile)
		}

	}

}

func TestMapToStr(t *testing.T) {
	testm := map[string]int{"test": 1, "hello": 2}
	var resSl = mapToStr(testm)
	if len(resSl) != len(testm) {
		t.Errorf("Expected %d words, got %d", len(testm), len(resSl))
	}
}

func TestEndToEnd(t *testing.T) {
	wordMap := map[string]int{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	re := regexp.MustCompile(`[\p{L}\d]+`)
	ts := []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa.")
	_, err := os.Stat(TestFile)
	if err != nil {
		err := os.WriteFile(TestFile, ts, 0644)
		if err != nil {
			t.Fatalf("Error creating file %s", TestFile)
		} else {
			getWordsFromFile(TestFile, re, logger, wordMap)
			if wordMap["adipiscing"] != 1 {
				t.Errorf("Expected 1 word got %d", wordMap["adipiscing"])
			}
			if wordMap["dolor"] != 2 {
				t.Errorf("Expected 2 word got %d", wordMap["dolor"])
			}
			if wordMap["aenean"] != 2 {
				t.Errorf("Expected 2 word got %d", wordMap["aenean"])
			}
			err := os.Remove(TestFile)
			if err != nil {
				t.Fatalf("Failed to remove test file %s", TestFile)
			}
		}
	} else {
		err := os.Remove(TestFile)
		if err != nil {
			t.Fatalf("Failed to remove test file %s", TestFile)
		}

	}
}
