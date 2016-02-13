package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var esvAPIKey = os.Getenv("API_KEY_ESVAPI")

type BibleESVAPI struct{}

func NewESVAPI() BibleESVAPI {
	return BibleESVAPI{}
}

func (b BibleESVAPI) Source() string {
	return "esvapi"
}

func (b BibleESVAPI) NameShort() string {
	return "ESV"
}

func (b BibleESVAPI) NameCommon() string {
	return "English Standard Version"
}

func (b BibleESVAPI) Name() string {
	return "English Standard Version"
}

func (b BibleESVAPI) GetPassage(reference string) (*BiblePassage, error) {
	q := url.Values{}
	q.Add("key", esvAPIKey)
	q.Add("include-passage-references", "true")
	q.Add("include-first-verse-numbers", "true")
	q.Add("include-footnotes", "false")
	q.Add("include-footnote-links", "false")
	q.Add("include-surrounding-chapters", "false")
	q.Add("include-headings", "false")
	q.Add("include-subheadings", "false")
	q.Add("include-short-copyright", "false")
	q.Add("include-audio-link", "false")
	q.Add("passage", reference)

	u := url.URL{
		Scheme:   "http",
		Host:     "www.esvapi.org",
		Path:     "/v2/rest/passageQuery",
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

	biblePassage := BiblePassage{
		Html:         string(body),
		TrackingCode: "",
		Copyright:    "Scripture taken from The Holy Bible, English Standard Version and Copyright &copy;2001 by <a href=\"http://www.crosswaybibles.org\">Crossway Bibles</a>, a publishing ministry of Good News Publishers. Used by permission. All rights reserved. Text provided by the <a href=\"http://www.gnpcb.org/esv/share/services/\">Crossway Bibles Web Service</a>.",
	}

	return &biblePassage, nil
}
