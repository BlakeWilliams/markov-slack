package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type MarkovChain struct {
	Chain map[string][]string
}

func NewMarkov() *MarkovChain {
	mc := new(MarkovChain)
	mc.Chain = make(map[string][]string)
	return mc
}

func (mc *MarkovChain) Parse(s string) {
	words := strings.Split(s, " ")
	wordCount := len(words)

	for i := 0; i <= wordCount; i++ {
		if i+2 < wordCount {
			key := fmt.Sprintf("%s %s", words[i], words[i+1])
			value := fmt.Sprintf("%s", words[i+2])

			mc.insert(key, value)
		}
	}
}

func (mc *MarkovChain) insert(key string, value string) {
	if _, ok := mc.Chain[key]; ok {
		mc.Chain[key] = append(mc.Chain[key], value)
	} else {
		mc.Chain[key] = []string{value}
	}
}

func (mc *MarkovChain) GenerateSentence() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(mc.Chain))

	var key string

	for key = range mc.Chain {
		if i == 0 {
			break
		}
		i--
	}

	sentence := fmt.Sprintf("%s", key)

	for {
		if values, ok := mc.Chain[key]; ok {
			firstWord := strings.Split(key, " ")[1]
			secondWord := randomArrayElement(values)

			sentence = fmt.Sprintf("%s %s", sentence, secondWord)

			key = fmt.Sprintf("%s %s", firstWord, secondWord)
		} else {
			break
		}
	}

	return sentence
}

func randomArrayElement(enum []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(enum))

	return enum[i]
}
