package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URLStruct struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"originalurl"`
	ShortCutURL string    `json:"shortcuturl"`
	CreatedAt   time.Time `json:"createdat"`
}

var urlDB = make(map[string]URLStruct)

var totalUrlShort int = 0

func generateShortURL(originalURL string) string {
	hasher := md5.New()

	hasher.Write([]byte(originalURL))

	data := hasher.Sum(nil)

	// fmt.Println("hasher ;", hasher, "\n\ndata ;", data)

	hexData := hex.EncodeToString(data)

	// fmt.Println("hexData ;", hexData, hexData[:8])

	return hexData[:8]
}

func createURL(originalURL string) string {
	shortUrl := generateShortURL(originalURL)

	totalUrlShort++

	fmt.Println("totalUrlShort --->>", totalUrlShort)

	urlDB[shortUrl] = URLStruct{
		ID:          totalUrlShort,
		OriginalURL: originalURL,
		ShortCutURL: shortUrl,
		CreatedAt:   time.Now(),
	}

	return shortUrl
}

func getURL(shortUrl string) (URLStruct, error) {
	url, of := urlDB[shortUrl]

	if !of {
		return URLStruct{}, errors.New("url not found")
	}

	return url, nil
}

// w is the response writer
// r is the request object
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Method")

	fmt.Fprint(w, "Hello World")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortedUrl := createURL(data.URL)

	fmt.Println("shortUrl ;", shortedUrl)

	response := struct {
		ShortUrl string `json:"short_url"`
	}{
		ShortUrl: shortedUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "redirectURLHandler")

	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println("url ;", url)
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("Hello World")
	// originalURL := "www.google.com"
	// generateShortURL(originalURL)

	// or

	/* Register Handlers Server*/

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect", redirectURLHandler)

	PORT := 8000
	fmt.Println(`Starting Server on PORT  `, PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		fmt.Println("Server Error", err)
		return
	}

}
