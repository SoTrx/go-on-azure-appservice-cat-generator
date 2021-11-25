package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// These can be harcoded, there is a nginx proxy anyway
const SERVER_PORT = 8081
const SERVER_HOST = "0.0.0.0"

// Global logger
var logger = MakeLogger()

func serveRandomCat(w http.ResponseWriter, req *http.Request) {
	// Fetch a random Cat url
	url, urlErr := GetRandomCatUrl()
	if urlErr != nil {
		handleServerError(urlErr, w)
		return
	}
	logger.info.Println(fmt.Sprintf("Serving cat %s", url))

	// From the Url, get the image
	resp, getErr := http.Get(url)
	if getErr != nil {
		handleServerError(getErr, w)
		return
	}

	// Transfert the images bytes into the body
	defer resp.Body.Close()
	imageBytes, imageErr := ioutil.ReadAll(resp.Body)
	if imageErr != nil {
		handleServerError(imageErr, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(imageBytes)
}

// Internal function used to handle 500 Errors
func handleServerError(err error, w http.ResponseWriter) {
	logger.err.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Unexpected error\n"))
}

func main() {
	// Try to get the API Key or die
	err := godotenv.Load()
	if err != nil {
		logger.warn.Println("No .env file detected")
	}
	if apiKey := os.Getenv("API_KEY"); len(apiKey) == 0 {
		logger.err.Panicln("No API KEY detected, aborting")
	}

	logger.info.Println(fmt.Sprintf("Now starting server on %s:%d", SERVER_HOST, SERVER_PORT))
	http.HandleFunc("/", serveRandomCat)
	logger.err.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT), nil))
}
