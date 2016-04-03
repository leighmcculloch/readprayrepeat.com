package main

import "strings"

type Bible interface {
	Source() string
	NameShort() string
	NameCommon() string
	Name() string
	GetPassage(reference string) (*BiblePassage, error)
}

type BiblePassage struct {
	Html         string
	TrackingCode string
	Copyright    string
}

func (p *BiblePassage) TimeToReadInMinutes() int {
	const READING_WORDS_PER_MINUTE = 220
	text := p.Html
	wordCount := strings.Count(text, " ")
	wordsPerMinute := wordCount / READING_WORDS_PER_MINUTE
	if wordsPerMinute == 0 {
		wordsPerMinute = 1
	}
	return wordsPerMinute
}
