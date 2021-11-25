package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Cat struct {
	Url string
}

func GetRandomCatUrl() (string, error) {
	req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search", nil)
	req.Header.Set("x-api-key", os.Getenv("API_KEY"))
	client := new(http.Client)
	resp, err := client.Do(req)
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
