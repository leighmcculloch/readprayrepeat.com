package main

import (
	"context"
	"time"

	"4d63.com/biblepassageapi"
	"4d63.com/youtubelength"
)

type Day struct {
	DayNumber        int
	WatchYoutubeId   string
	ReadingReference string
	PrayerReference  string

	WatchYoutubeDurationMinutes int
	ReadingBiblePassage         *biblepassageapi.Passage
	PrayerBiblePassage          *biblepassageapi.Passage
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
	var err error

	d.ReadingBiblePassage, err = bible.GetPassage(d.ReadingReference)
	if err != nil {
		return err
	}

	d.PrayerBiblePassage, err = bible.GetPassage(d.PrayerReference)
	if err != nil {
		return err
	}

	return nil
}

func (d Day) TimeToCompleteInMinutes() int {
	return d.ReadingBiblePassage.TimeToReadInMinutes() +
		d.PrayerBiblePassage.TimeToReadInMinutes() +
		d.WatchYoutubeDurationMinutes
}
