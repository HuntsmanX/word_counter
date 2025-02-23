package main

import (
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestFile = "tmp.txt"

func TestFillMap(t *testing.T) {
	wordMap := map[string]int{}
	words := []string{"Test", "test", "TEST", "Hello", "HELLO"}
	fillMap(wordMap, words)

	assert.Equal(t, 3, wordMap["test"], "Expected 'test' to appear 3 times")
	assert.Equal(t, 2, wordMap["hello"], "Expected 'hello' to appear 2 times")
}

func TestMapToSlice(t *testing.T) {
	testMap := map[string]int{"test": 1, "hello": 2}
	result := mapToSlice(testMap)

	assert.Equal(t, len(testMap), len(result), "Mismatch in slice length and map keys count")

	expected := map[string]int{"test": 1, "hello": 2}
	for _, entry := range result {
		assert.Equal(t, expected[entry.Word], entry.Count, "Mismatch in word count")
	}

}

func TestGetWordsFromFile(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wordMap := map[string]int{}
	expression := regexp.MustCompile(`[\p{L}\d]+`)
	testContent := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."

	err := os.WriteFile(TestFile, []byte(testContent), 0644)
	assert.NoError(t, err, "Error creating test file")
	defer func() {
		assert.NoError(t, os.Remove(TestFile), "Failed to remove test file")
	}()

	err = getWordsFromFile(TestFile, expression, logger, wordMap)
	assert.NoError(t, err, "Error reading test file")
	assert.Equal(t, 13, len(wordMap), "Expected 13 unique words in file")
}

func TestEndToEnd(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wordMap := map[string]int{}
	expression := regexp.MustCompile(`[\p{L}\d]+`)
	testContent := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa."

	err := os.WriteFile(TestFile, []byte(testContent), 0644)
	assert.NoError(t, err, "Error creating test file")
	defer func() {
		assert.NoError(t, os.Remove(TestFile), "Failed to remove test file")
	}()

	err = getWordsFromFile(TestFile, expression, logger, wordMap)
	assert.NoError(t, err, "Error processing test file")

	assert.Equal(t, 1, wordMap["adipiscing"], "Expected 1 occurrence of 'adipiscing'")
	assert.Equal(t, 2, wordMap["dolor"], "Expected 2 occurrences of 'dolor'")
	assert.Equal(t, 2, wordMap["aenean"], "Expected 2 occurrences of 'aenean'")
}
