package main

import (
	"fmt"
	"time"

	"4d63.com/biblepassageapi"
)

type pageDay struct {
	NextPage     *pageDay
	PreviousPage *pageDay

	Day Day

	Bible biblepassageapi.Bible
}

func (p pageDay) Date() time.Time {
	return time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, p.Day.DayNumber-1)
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
