package main

type Day struct {
	DayNumber        int
	WatchYoutubeId   string
	ReadingReference string
	PrayerReference  string

	ReadingBiblePassage *BiblePassage
	PrayerBiblePassage  *BiblePassage
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
