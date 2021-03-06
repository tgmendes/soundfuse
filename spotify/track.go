package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Track struct {
	ID          string      `json:"id"`
	Duration    int         `json:"duration_ms"`
	Genres      []string    `json:"genres"`
	Images      []Image     `json:"images"`
	Name        string      `json:"name"`
	Popularity  int         `json:"popularity"`
	Type        string      `json:"type"`
	HRef        string      `json:"href"`
	Album       Album       `json:"album"`
	Artists     []Artist    `json:"artists"`
	ExternalIDs ExternalIDs `json:"external_ids"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type ExternalIDs struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

type TracksResponse struct {
	Href     string  `json:"href"`
	Items    []Track `json:"items"`
	Previous string  `json:"previous"`
	Next     string  `json:"next"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
}

func (c *Client) UserTopTracks(ctx context.Context, limit, offset int, timeRange TopTimeRange) (TracksResponse, error) {
	url := fmt.Sprintf("%s/me/top/tracks", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TracksResponse{}, err
	}
	q := req.URL.Query()
	q.Add("time_range", string(timeRange))
	if limit != 0 {
		q.Add("limit", strconv.Itoa(limit))
	}

	if offset != 0 {
		q.Add("offset", strconv.Itoa(offset))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TracksResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return TracksResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tracks TracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return TracksResponse{}, errors.New("unable to unmarshal response")
	}
	return tracks, nil
}
