package main

import (
	"fmt"
	"strings"
)

type pageDay struct {
	NextPage     *pageDay
	PreviousPage *pageDay

	Day Day

	Bible  Bible
	Bibles []Bible
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

func (p pageDay) PathSuffixForBible(bible Bible) string {
	if p.IsBibleDefault(bible) {
		return ""
	} else {
		return fmt.Sprintf(".%s", strings.ToLower(bible.NameShort()))
	}
}

func (p pageDay) PathForBible(bible Bible) string {
	return fmt.Sprintf("%s%s", p.PathBase(), p.PathSuffixForBible(bible))
}

func (p pageDay) IsBibleDefault(bible Bible) bool {
	return bible == p.Bibles[0]
}

func (p *pageDay) LoadPassages() error {
	err := p.Day.LoadYoutubeDetails()
	if err != nil {
		return err
	}
	return p.Day.LoadPassages(p.Bible)
}
