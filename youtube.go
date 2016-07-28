package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var youtubeApiKey = os.Getenv("GOOGLE_YOUTUBE_API_KEY")

func GetYoutubeVideoDurationMinutes(youtubeId string) (int, error) {
	q := url.Values{}
	q.Add("key", youtubeApiKey)
	q.Add("part", "contentDetails")
	q.Add("id", youtubeId)

	u := url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/videos",
		RawQuery: q.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	var response struct {
		Items []struct {
			ContentDetails struct {
				Duration string
			}
		}
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&response)
	if err != nil {
		return 0, err
	}

	if len(response.Items) != 1 {
		return 0, errors.New(fmt.Sprintf("Response items for Youtube ID %v was expected to be a single item, but was: %v", youtubeId, response.Items))
	}
	durationString := response.Items[0].ContentDetails.Duration
	r := regexp.MustCompile(`^PT([0-9]+)M`)
	m := r.FindAllStringSubmatch(durationString, -1)
	durationStringMinutes := m[0][1]
	durationMinutes, err := strconv.Atoi(durationStringMinutes)
	return durationMinutes, err
}
