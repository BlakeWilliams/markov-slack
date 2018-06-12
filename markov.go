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
	word_count := len(words)

	for i := 0; i <= word_count; i++ {
		if i+2 < word_count {
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
			first_word := strings.Split(key, " ")[1]
			second_word := randomArrayElement(values)

			sentence = fmt.Sprintf("%s %s", sentence, second_word)

			key = fmt.Sprintf("%s %s", first_word, second_word)
		} else {
			break
		}
	}

	return sentence
}

func randomArrayElement(enum []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(enum))
	var key string

	for _, key = range enum {
		if i == 0 {
			break
		}
		i--
	}

	return key
}
