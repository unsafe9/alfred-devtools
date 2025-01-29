package main

import "testing"

func TestLoremIpsumGenerator(t *testing.T) {
	t.Log(LoremIpsum.Words(10))
	t.Log(LoremIpsum.Sentences(2))
	t.Log(LoremIpsum.Paragraphs(1))
}
