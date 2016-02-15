package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type BibleNET struct{}

func NewBibleNET() BibleNET {
	return BibleNET{}
}

func (b BibleNET) Source() string {
	return "netbibleapi"
}

func (b BibleNET) NameShort() string {
	return "NET"
}

func (b BibleNET) NameCommon() string {
	return "NET Bible"
}

func (b BibleNET) Name() string {
	return "New English Translation Bible"
}

func (b BibleNET) GetPassage(reference string) (*BiblePassage, error) {
	q := url.Values{}
	q.Add("formatting", "full")
	q.Add("passage", reference)

	u := url.URL{
		Scheme:   "http",
		Host:     "labs.bible.org",
		Path:     "/api/",
		RawQuery: q.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// .gsub(%r{<n id="\d+" />}, "")

	biblePassage := BiblePassage{
		Html:         string(body),
		TrackingCode: "",
		Copyright:    "Scripture quoted by permission. All scripture quotations, unless otherwise indicated, are taken from the NET Bible&reg; copyright &copy;1996-2006 by Biblical Studies Press, L.L.C. http://netbible.com All rights reserved.",
	}

	return &biblePassage, nil
}
