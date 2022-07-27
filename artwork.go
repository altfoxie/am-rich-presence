package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type iTunesResponse struct {
	Results []struct {
		ArtworkURL100 string `json:"artworkUrl100"`
	} `json:"results"`
}

func artworkSearchITunes(query string) string {
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

const musixmatchSignatureSecret = "8d2899b2aebb97a69a4a85cc991c0b6713a1d9e2"

type musixmatchResponse struct {
	Message struct {
		Body struct {
			MacroResultList struct {
				TrackList []struct {
					Track struct {
						AlbumCoverart500 string `json:"album_coverart_500x500"`
					} `json:"track"`
				} `json:"track_list"`
			} `json:"macro_result_list"`
		} `json:"body"`
	} `json:"message"`
}

func artworkSearchMusixmatch(query string) string {
	baseURL := "https://www.musixmatch.com/ws/1.1/macro.search?" + url.Values{
		"app_id":    {"community-app-v1.0"},
		"guid":      {uuid.NewString()},
		"format":    {"json"},
		"q":         {query},
		"part":      {"artist_image"},
		"page_size": {"1"},
	}.Encode()

	mac := hmac.New(sha1.New, []byte(musixmatchSignatureSecret))
	data := []byte(baseURL + time.Now().Format("20060102"))
	if _, err := mac.Write(data); err != nil {
		return ""
	}

	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"&signature="+url.QueryEscape(sign)+"&signature_protocol=sha1",
		nil,
	)
	if err != nil {
		return ""
	}
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var respData musixmatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil ||
		len(respData.Message.Body.MacroResultList.TrackList) == 0 {
		return ""
	}

	return respData.Message.Body.MacroResultList.TrackList[0].Track.AlbumCoverart500
}
