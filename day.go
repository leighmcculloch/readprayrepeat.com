package main

import (
	"4d63.com/biblepassageapi"
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
	durationMinutes, err := getYoutubeVideoDurationMinutes(d.WatchYoutubeId)
	if err != nil {
		return err
	}
	d.WatchYoutubeDurationMinutes = durationMinutes
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
