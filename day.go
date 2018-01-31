package main

import (
	"context"
	"time"

	"4d63.com/biblepassageapi"
	"4d63.com/youtubelength"
)

type Day struct {
	DayNumber      int
	WatchYoutubeId string
	References     []string

	WatchYoutubeDurationMinutes int
	BiblePassages               []*biblepassageapi.Passage
}

func (d *Day) LoadYoutubeDetails() error {
	switch d.WatchYoutubeId {
	case "", "...":
		return nil
	}
	duration, err := youtubelength.Get(context.TODO(), d.WatchYoutubeId)
	if err != nil {
		return err
	}
	d.WatchYoutubeDurationMinutes = int(duration / time.Minute)
	return nil
}

func (d *Day) LoadPassages(bible biblepassageapi.Bible) error {
	d.BiblePassages = []*biblepassageapi.Passage{}
	for _, r := range d.References {
		bp, err := bible.GetPassage(r)
		if err != nil {
			return err
		}
		d.BiblePassages = append(d.BiblePassages, bp)
	}

	return nil
}

func (d Day) TimeToCompleteInMinutes() int {
	mins := d.WatchYoutubeDurationMinutes
	for _, bp := range d.BiblePassages {
		mins += bp.TimeToReadInMinutes()
	}
	return mins
}
