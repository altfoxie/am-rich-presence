package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type iTunesResponse struct {
	Results []struct {
		ArtworkURL100 string `json:"artworkUrl100"`
	} `json:"results"`
}

func artworkSearch(query string) string {
	resp, err := http.Get("https://itunes.apple.com/search?" + url.Values{
		"term":   {query},
		"media":  {"music"},
		"entity": {"musicArtist,musicTrack,album,mix,song"},
		"limit":  {"1"},
	}.Encode())
	if err != nil {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var respData iTunesResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil ||
		len(respData.Results) == 0 {
		return ""
	}

	// https://github.com/paambaati/itunes-artwork/blob/master/api.js#L72
	return strings.Replace(respData.Results[0].ArtworkURL100, "100x100bb.", "512x512bb.", 1)
}
