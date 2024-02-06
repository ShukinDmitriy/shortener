package main

import (
	"io"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}
	return string(shortKey)
}

func (us *URLShortener) HandleShorten(res http.ResponseWriter, req *http.Request) {
	originalURL, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(res, "can't read body. internal error",
			http.StatusInternalServerError)
		return
	}

	if string(originalURL) == "" {
		http.Error(res, "empty url",
			http.StatusBadRequest)
		return
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	us.urls[shortKey] = string(originalURL)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusCreated)
	// TODO как получить текущий протокол из запроса?
	res.Write([]byte("http://" + req.Host + "/" + shortKey))
}

func (us *URLShortener) HandleRedirect(res http.ResponseWriter, req *http.Request) {
	shortKey := req.URL.Path[len("/"):]

	if shortKey == "" {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(res, "", http.StatusBadRequest)
	}

	http.Redirect(res, req, originalURL, http.StatusTemporaryRedirect)
}

var shortener = &URLShortener{
	urls: make(map[string]string),
}

func mainHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		shortener.HandleRedirect(res, req)
		return
	}

	if req.Method == http.MethodPost {
		shortener.HandleShorten(res, req)
		return
	}

	http.Error(res, "", http.StatusBadRequest)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
