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

type URL struct {
	Id           string    `json:"id"`
	OriginalUrl  string    `json:"uri"`
	ShortUrl     string    `json:"url"`
	CreationDate time.Time `json:"created_at"`
}

var urlDB = make(map[string]URL)

func GenerateShortUrl(OriginalUrl string) string {
	Harsher := md5.New()
	Harsher.Write([]byte(OriginalUrl)) //it converts OriginalUrl String to a byte slices
	data := Harsher.Sum(nil)
	fmt.Println("data:", data)
	hash := hex.EncodeToString(data) //hash is a string of hexadecimal digits
	fmt.Println("Encoded string:", hash)
	fmt.Println("Final Hash:", hash[:8])

	return hash[:8]
}

func CreateURL(OriginalUrl string) string {

	ShortURL := GenerateShortUrl(OriginalUrl)
	id := ShortURL
	urlDB[id] = URL{
		Id:           id,
		OriginalUrl:  OriginalUrl,
		ShortUrl:     ShortURL,
		CreationDate: time.Now(),
	}
	return ShortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")

	}
	return url, nil
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")

}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var DATA struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&DATA)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	ShortURL_ := CreateURL(DATA.URL)
	//fmt.Fprintf(w, ShortURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: ShortURL_}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/r/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "invalid", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}
func main() {
	//fmt.Println("URL Shortener")
	OriginalUrl := "https:www.youtube.com"
	GenerateShortUrl(OriginalUrl)
	//fmt.Println("short Url is:", ShortUrl)
	//register the handler function to handle all the request to the root URL("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/r/", RedirectURLHandler)
	http.HandleFunc("/shorten", ShortURLHandler)
	//start the server on port 8080
	fmt.Println("starting server on port 8080...")
	err := http.ListenAndServe(":5500", nil)
	if err != nil {
		fmt.Println("error found at port 8080:", err)
	}

}
