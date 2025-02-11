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
	re := regexp.MustCompile(`\p{L}+`)
	ts := []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa.")
	c := 15
	_, err := os.Stat(TestFile)
	if err != nil {
		err := os.WriteFile(TestFile, ts, 0644)
		if err != nil {
			t.Errorf("Error creating file %s", TestFile)
		} else {
			lines, _ := getWordsFromFile(TestFile, re, logger)
			if len(lines) != c {
				t.Errorf("Expected 15 words got %d", len(lines))
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

func TestEndToEnd(t *testing.T) {
	wordMap := map[string]int{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	re := regexp.MustCompile(`\p{L}+`)
	ts := []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa.")
	_, err := os.Stat(TestFile)
	if err != nil {
		err := os.WriteFile(TestFile, ts, 0644)
		if err != nil {
			t.Fatalf("Error creating file %s", TestFile)
		} else {
			lines, _ := getWordsFromFile(TestFile, re, logger)
			fillMap(wordMap, lines)
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
