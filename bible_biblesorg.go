package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var bibleApiKey = os.Getenv("API_KEY_BIBLESORG")

type translation struct {
	ID         string
	Name       string
	NameCommon string
	NameShort  string
}

var (
	CEV = "CEV"
	GNT = "GNT"
)

var translations = map[string]translation{
	CEV: translation{
		ID:         "eng-CEVD",
		NameShort:  "CEV",
		NameCommon: "Contemporary English Version",
		Name:       "2006 Contemporary English Version, Second Edition (US Version)",
	},
	GNT: translation{
		ID:         "eng-GNTD",
		NameShort:  "GNT",
		NameCommon: "Good News Translation",
		Name:       "1992 Good News Translation, Second Edition (US Version)",
	},
}

type BibleBiblesOrg struct {
	translation
}

func NewBiblesOrg(translation string) BibleBiblesOrg {
	return BibleBiblesOrg{translation: translations[translation]}
}

func (b BibleBiblesOrg) Source() string {
	return "biblesorg"
}

func (b BibleBiblesOrg) NameShort() string {
	return b.translation.NameShort
}

func (b BibleBiblesOrg) NameCommon() string {
	return b.translation.NameCommon
}

func (b BibleBiblesOrg) Name() string {
	return b.translation.Name
}

func (b BibleBiblesOrg) GetPassage(reference string) (*BiblePassage, error) {
	q := url.Values{}
	q.Add("q[]", reference)

	u := url.URL{
		Scheme:   "https",
		Host:     "bibles.org",
		Path:     fmt.Sprintf("/v2/%s/passages.js", b.translation.ID),
		RawQuery: q.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(bibleApiKey, "X")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var searchRes struct {
		Response struct {
			Search struct {
				Result struct {
					Passages []struct {
						Copyright string
						Text      string
					}
				}
			}
			Meta struct {
				FumsNoscript string
			}
		}
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&searchRes)
	if err != nil {
		return nil, err
	}

	passageCount := len(searchRes.Response.Search.Result.Passages)
	if passageCount != 1 {
		return nil, fmt.Errorf("Reference %s returned %d passages, expected 1.", reference, passageCount)
	}

	biblePassage := BiblePassage{
		Html:         transposePassageHtml(searchRes.Response.Search.Result.Passages[0].Text),
		TrackingCode: searchRes.Response.Meta.FumsNoscript,
		Copyright:    transposePassageCopyright(searchRes.Response.Search.Result.Passages[0].Copyright),
	}

	return &biblePassage, nil
}

func transposePassageHtml(s string) string {
	s = strings.Replace(s, "<sup", " <sup", -1)
	s = strings.Replace(s, "</sup>", "</sup> ", -1)
	s = strings.Replace(s, "[Lord]", "Lord", -1)
	return s
}

func transposePassageCopyright(s string) string {
	s = regexp.MustCompile("(&#169;)").ReplaceAllString(s, "</br>$1")
	s = regexp.MustCompile("(</p>)").ReplaceAllString(s, "</br>Used with permission.$1")
	return s
}
