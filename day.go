package main

type Day struct {
	DayNumber   int
	DayPrevious *int
	DayNext     *int

	WatchYoutubeId   string
	ReadingReference string
	PrayerReference  string

	ReadingBiblePassage *BiblePassage
	PrayerBiblePassage  *BiblePassage
}

func NewDay(dayNumber, dayCount int, readingReference, prayerReference, watchYoutubeId string) *Day {
	dayPreviousValue := dayNumber - 1
	dayPrevious := &dayPreviousValue
	if *dayPrevious == 0 {
		dayPrevious = nil
	}

	dayNextValue := dayNumber + 1
	dayNext := &dayNextValue
	if *dayNext > dayCount {
		dayNext = nil
	}

	day := Day{
		DayNumber:        dayNumber,
		DayPrevious:      dayPrevious,
		DayNext:          dayNext,
		ReadingReference: readingReference,
		PrayerReference:  prayerReference,
		WatchYoutubeId:   watchYoutubeId,
	}

	return &day
}

func (d *Day) LoadPassages() error {
	var err error

	d.ReadingBiblePassage, err = GetBiblePassageWithCache(d.ReadingReference)
	if err != nil {
		return err
	}

	d.PrayerBiblePassage, err = GetBiblePassageWithCache(d.PrayerReference)
	if err != nil {
		return err
	}

	return nil
}
