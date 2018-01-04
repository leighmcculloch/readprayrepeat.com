package main

import (
	"fmt"
	"strings"

	"4d63.com/biblepassageapi"
)

type pageDay struct {
	NextPage     *pageDay
	PreviousPage *pageDay

	Day Day

	Bible  biblepassageapi.Bible
	Bibles []biblepassageapi.Bible
}

func (p pageDay) PathBase() string {
	return fmt.Sprintf("/%d", p.Day.DayNumber)
}

func (p pageDay) PathSuffix() string {
	return p.PathSuffixForBible(p.Bible)
}

func (p pageDay) Path() string {
	return p.PathForBible(p.Bible)
}

func (p pageDay) PathSuffixForBible(bible biblepassageapi.Bible) string {
	if p.IsBibleDefault(bible) {
		return ""
	} else {
		return fmt.Sprintf(".%s", strings.ToLower(bible.NameShort()))
	}
}

func (p pageDay) PathForBible(bible biblepassageapi.Bible) string {
	return fmt.Sprintf("%s%s", p.PathBase(), p.PathSuffixForBible(bible))
}

func (p pageDay) IsBibleDefault(bible biblepassageapi.Bible) bool {
	return bible == p.Bibles[0]
}

func (p *pageDay) LoadPassages() error {
	err := p.Day.LoadYoutubeDetails()
	if err != nil {
		return err
	}
	return p.Day.LoadPassages(p.Bible)
}
