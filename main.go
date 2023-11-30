package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	words, err := readWords("/app/Noun.csv")
	if err != nil {
		return err
	}

	c := buildChain(words)
	currentWord := words[rand.Intn(len(words))]

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		newWord, ok := c.Next(currentWord)
		if !ok {
			return nil
		}
		currentWord = newWord
		fmt.Printf("%s(%s)\n", currentWord.Text, currentWord.Reading)
	}

	return nil
}

func readWords(file string) ([]Word, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	words := []Word{}
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, ",")
		words = append(words, Word{Text: items[0], Reading: items[1]})
	}
	return words, nil
}

type Word struct {
	Reading string
	Text    string
}

type Chain struct {
	initials map[string][]Word
}

func buildChain(words []Word) *Chain {
	c := &Chain{
		initials: make(map[string][]Word),
	}

	for _, word := range words {
		initial := string([]rune(word.Reading)[0])
		if _, ok := c.initials[initial]; !ok {
			c.initials[initial] = []Word{}
		}
		c.initials[initial] = append(c.initials[initial], word)
	}
	return c
}

func (c *Chain) Next(current Word) (Word, bool) {
	rs := []rune(current.Reading)
	tail := normalizeTail(string(rs[len(rs)-1]))
	if len(c.initials[tail]) == 0 {
		return Word{}, false
	}
	i := rand.Intn(len(c.initials[tail]))
	return c.initials[tail][i], true
}

func normalizeTail(s string) string {
	switch s {
	case "ャ":
		return "ヤ"
	case "ュ":
		return "ユ"
	case "ョ":
		return "ヨ"
	}
	return s
}
