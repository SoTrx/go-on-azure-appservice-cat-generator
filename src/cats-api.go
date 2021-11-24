package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Cat struct {
	Url string
}

func GetRandomCatUrl() (string, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		return "", err
	}
	// Free body when function returns
	defer resp.Body.Close()
	var cats []Cat
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if json.Unmarshal(bytes, &cats); err != nil {
		return "", err
	}
	return cats[0].Url, nil
}
