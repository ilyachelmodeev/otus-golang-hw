package hw03frequencyanalysis

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

type wordCountStruct struct {
	w string
	c int32
}

const TopNumber = 10

var globalRegexp = regexp.MustCompile("[a-яА-Я-]+")

func getTokenFrequency(inp string) map[string]int32 {
	tokenFrequency := make(map[string]int32)
	tokens := strings.Fields(inp)
	for _, t := range tokens {
		if globalRegexp.MatchString(t) {
			tokenFrequency[t]++
		}
	}
	return tokenFrequency
}

func getWordCounts(tokenFrequency map[string]int32) []wordCountStruct {
	var wordCounts []wordCountStruct
	for key, val := range tokenFrequency {
		wordCounts = append(wordCounts, wordCountStruct{w: key, c: val})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].c == wordCounts[j].c {
			return wordCounts[i].w < wordCounts[j].w
		}
		return wordCounts[i].c > wordCounts[j].c
	})

	return wordCounts
}

func getTopWords(wordCounts []wordCountStruct, topN int) []string {
	result := make([]string, topN)
	for i, curWordCount := range wordCounts[:topN] {
		result[i] = curWordCount.w
	}
	return result
}

func Top10(inp string) []string {
	tokenFrequency := getTokenFrequency(inp)
	wordCounts := getWordCounts(tokenFrequency)
	maxRetLen := math.Min(float64(TopNumber), float64(len(wordCounts)))
	return getTopWords(wordCounts, int(maxRetLen))
}
