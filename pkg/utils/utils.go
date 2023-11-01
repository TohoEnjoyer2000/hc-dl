package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_2_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15"
)

func DetectName(galleryURL string) string {
	clean := strings.TrimSuffix(galleryURL, "/")
	parts := strings.Split(clean, "/")
	return parts[len(parts)-1]
}

func PerformHttpGet(client *http.Client, galleryURL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, galleryURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", userAgent)

	return client.Do(req)
}

func ExtractImagesFromRoot(root *soup.Root) *[]string {
	images := root.Find("div", "id", "display_image_detail").FindAll("img")

	retval := make([]string, len(images))

	for i := range images {
		retval[i] = images[i].Attrs()["src"]
	}

	return &retval
}

func ExtractImagesDiv(body []byte) *[]string {
	doc := soup.HTMLParse(string(body))
	images := doc.Find("div", "id", "display_image_detail").FindAll("img")

	retval := make([]string, len(images))

	for i := range images {
		retval[i] = images[i].Attrs()["src"]
	}

	return &retval
}

func ExtractPaginatorData(galleryURL string) (*[]string, *soup.Root, error) {
	res, err := PerformHttpGet(http.DefaultClient, galleryURL)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	parsedUrl, err := url.Parse(galleryURL)
	if err != nil {
		return nil, nil, err
	}

	var (
		doc     = soup.HTMLParse(string(page))
		host    = fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
		pages   = doc.Find("div", "id", "paginator").FindAll("a")
		content []string
	)

	for i, a := range pages {
		t := strings.ToLower(a.Text())
		if t != "next>" && t != "last>>" && t != "<<first" && t != "<prev" {
			u := (a.Attrs()["href"])
			content = append(content, host+u)
			fmt.Println("Discovered page", i+1, "...")
		}
	}

	return &content, &doc, nil
}

func ExtractImagesFromPage(page string) (*[]string, error) {
	res, err := PerformHttpGet(http.DefaultClient, page)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	buff, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	images := ExtractImagesDiv(buff)

	return images, nil
}
