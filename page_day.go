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

func (p pageDay) Path() string {
	if p.Bible == p.Bibles[0] {
		return fmt.Sprintf("/%d", p.Day.DayNumber)
	} else {
		return fmt.Sprintf("/%d.%s", p.Day.DayNumber, strings.ToLower(p.Bible.NameShort()))
	}
}

func (p *pageDay) LoadPassages() error {
	return p.Day.LoadPassages(p.Bible)
}
