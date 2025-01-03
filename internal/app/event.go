package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type SquarespaceEvent struct {
	Item struct {
		UrlId string `json:"urlId"`
		Title string `json:"title"`
	} `json:"item"`
}

func getEventFromSquarespace(slug string) (SquarespaceEvent, error) {
	var event SquarespaceEvent
	if req, err := http.NewRequest(http.MethodGet, "https://outclimb.gay/events/"+url.PathEscape(slug)+"?format=json", nil); err != nil {
		return SquarespaceEvent{}, err
	} else if res, err := http.DefaultClient.Do(req); err != nil {
		return SquarespaceEvent{}, err
	} else if res.StatusCode != http.StatusOK {
		return SquarespaceEvent{}, errors.New("Squarespace returned non-200 status code")
	} else if resBody, err := io.ReadAll(res.Body); err != nil {
		return SquarespaceEvent{}, err
	} else if err = json.Unmarshal(resBody, &event); err != nil {
		return SquarespaceEvent{}, err
	} else {
		return event, nil
	}
}

func (a *appLayer) CheckEventExists(slug string) bool {
	if _, err := a.store.GetEvent(slug); err == nil {
		return true
	}

	if squarespaceEvent, err := getEventFromSquarespace(slug); err != nil {
		return false
	} else {
		a.store.CreateEvent(squarespaceEvent.Item.Title, slug)
		return true
	}
}
