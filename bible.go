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

const bibleVersion = "eng-CEVD" // 2006 Contemporary English Version, Second Edition (US Version)

var bibleApiUrl = url.URL{
	Scheme: "https",
	Host:   "bibles.org",
	Path:   fmt.Sprintf("/v2/%s/passages.js", bibleVersion),
}

type BiblePassage struct {
	Html         string
	TrackingCode string
	Copyright    string
}

func GetBiblePassage(reference string) (*BiblePassage, error) {
	u := bibleApiUrl
	q := url.Values{}
	q.Add("q[]", reference)
	u.RawQuery = q.Encode()

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
