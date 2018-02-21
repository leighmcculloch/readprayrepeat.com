package main

import (
	"fmt"
	"strings"

	"4d63.com/biblepassageapi"
)

type pagePassage struct {
	Reference     string
	AbbrReference string

	PreviousPage *pagePassage
	NextPage     *pagePassage

	Bible        biblepassageapi.Bible
	BiblePassage *biblepassageapi.Passage
}

func (p pagePassage) Path() string {
	const chars = " :"
	slug := strings.ToLower(p.AbbrReference)
	for _, c := range chars {
		slug = strings.Replace(slug, string(c), ".", -1)
	}
	return fmt.Sprintf("/%s", slug)
}

func (p *pagePassage) LoadPassages() error {
	bp, err := p.Bible.GetPassage(p.Reference)
	if err != nil {
		return err
	}
	p.BiblePassage = bp
	return nil
}
