package main

import (
	"math/rand/v2"
	"strconv"
	"strings"
)

type LoremIpsumGenerator []string

var loremIpsum = LoremIpsumGenerator{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	"adipiscing", "elit", "sed", "do", "eiusmod", "tempor",
	"incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua",
	"enim", "ad", "minim", "veniam", "quis", "nostrud", "exercitation",
	"ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
	"in", "voluptate", "velit", "esse", "cillum", "eu", "fugiat", "nulla",
	"pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non",
	"proident", "sunt", "in", "culpa", "qui", "officia", "deserunt",
	"mollit", "anim", "id", "est", "laborum",
}

func (r LoremIpsumGenerator) word() string {
	return r[rand.IntN(len(r))]
}

func (r LoremIpsumGenerator) Words(input string) string {
	n := atoi(input, 10)
	if n <= 0 {
		return ""
	}
	buf := strings.Builder{}
	buf.Grow(n * 10)
	for i := 0; i < n; i++ {
		buf.WriteString(r.word())
		if i < n-1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func (r LoremIpsumGenerator) Sentences(input string) string {
	n := atoi(input, 2)
	if n <= 0 {
		return ""
	}
	buf := strings.Builder{}
	buf.Grow(n * 100)
	for i := 0; i < n; i++ {
		words := rand.IntN(10) + 5
		commaAt := rand.IntN(3) + 2
		for j := 0; j < words; j++ {
			word := r.word()
			if j == 0 {
				buf.WriteByte(word[0] - ('a' - 'A'))
				buf.WriteString(word[1:])
			} else {
				buf.WriteString(word)
			}
			if j == commaAt {
				buf.WriteByte(',')
			}
			if j < words-1 {
				buf.WriteByte(' ')
			} else {
				buf.WriteByte('.')
			}
		}
		if i < n-1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String()
}

func (r LoremIpsumGenerator) Paragraphs(input string) string {
	n := atoi(input, 1)
	if n <= 0 {
		return ""
	}
	buf := strings.Builder{}
	buf.Grow(n * 500)
	for i := 0; i < n; i++ {
		sentences := r.Sentences(strconv.Itoa(rand.IntN(5) + 3))
		buf.WriteString(sentences)
		if i < n-1 {
			buf.WriteString("\n\n")
		}
	}
	return buf.String()
}
