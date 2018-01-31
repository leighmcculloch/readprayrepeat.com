package main

import (
	"fmt"

	"4d63.com/biblepassageapi"
)

type pageDay struct {
	NextPage     *pageDay
	PreviousPage *pageDay

	Day Day

	Bible biblepassageapi.Bible
}

func (p pageDay) Path() string {
	return fmt.Sprintf("/%d", p.Day.DayNumber)
}

func (p *pageDay) LoadPassages() error {
	err := p.Day.LoadYoutubeDetails()
	if err != nil {
		return err
	}
	return p.Day.LoadPassages(p.Bible)
}
