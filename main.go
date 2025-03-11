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
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

var urlDB = make(map[string]URLStruct)

func generateShortURL(_originalURL string) string {
	hasher := md5.New()

	hasher.Write([]byte(_originalURL))

	data := hasher.Sum(nil)

	// fmt.Println("hasher ;", hasher, "\n\ndata ;", data)

	hexData := hex.EncodeToString(data)

	// fmt.Println("hexData ;", hexData, hexData[:8])

	return hexData[:8]
}

func createURL(_originalURL string) string {
	shortUrl := generateShortURL(_originalURL)

	urlDB[shortUrl] = URLStruct{
		ID:          shortUrl,
		OriginalURL: _originalURL,
		ShortURL:    shortUrl,
		CreatedAt:   time.Now(),
	}

	return shortUrl
}

func getURL(_shortUrl string) (URLStruct, error) {
	url, status := urlDB[_shortUrl]

	if !status {
		return URLStruct{}, errors.New("url not found")
	}

	return url, nil
}

// w is the response writer
// r is the request object
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello User")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URLStruct string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortedUrl := createURL(data.URLStruct)

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

	fmt.Println("redirect URL Handler")
	fmt.Fprint(w, "redirect URL Handler")

	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Println("url Data;", url.OriginalURL)
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("Hello User")
	// originalURL := "www.google.com"
	// generateShortURL(originalURL)

	// or

	/* Register Handlers Server*/

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	PORT := 8000
	fmt.Println(`Starting Server on PORT  `, PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		fmt.Println("Server Error", err)
		return
	}

}
