package internal

import (
	"log"
	"net/http"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_2_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15"
)

func NewHttpGet(client *http.Client, galleryURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, galleryURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", userAgent)

	return client.Do(req)
}
