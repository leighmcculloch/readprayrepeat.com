package main

type Day struct {
	DayNumber        int
	WatchYoutubeId   string
	ReadingReference string
	PrayerReference  string

	WatchYoutubeDurationMinutes int
	ReadingBiblePassage         *BiblePassage
	PrayerBiblePassage          *BiblePassage
}

func (d *Day) LoadYoutubeDetails() error {
	switch d.WatchYoutubeId {
	case "", "...":
		return nil
	}
	durationMinutes, err := GetYoutubeVideoDurationMinutes(d.WatchYoutubeId)
	if err != nil {
		return err
	}
	d.WatchYoutubeDurationMinutes = durationMinutes
	return nil
}

func (d *Day) LoadPassages(bible Bible) error {
	var err error

	d.ReadingBiblePassage, err = GetBiblePassageWithCache(bible, d.ReadingReference)
	if err != nil {
		return err
	}

	d.PrayerBiblePassage, err = GetBiblePassageWithCache(bible, d.PrayerReference)
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
