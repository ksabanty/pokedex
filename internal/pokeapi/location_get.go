package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

// Explore
func (c *Client) GetLocation(loc string) (Location, error) {
	url := baseURL + "location-area/" + loc

	if val, ok := c.cache.Get(url); ok {
		locationResp := Location{}
		err := json.Unmarshal(val, &locationResp)
		if err != nil {
			return Location{}, err
		}
		return locationResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}

	locationsRes := Location{}
	err = json.Unmarshal(dat, &locationsRes)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, dat)

	return locationsRes, nil
}
