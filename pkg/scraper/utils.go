package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TohoEnjoyer2000/hc-dl/internal"
	"github.com/schollz/progressbar/v3"
)

func DetectName(galleryURL string) string {
	clean := strings.TrimSuffix(galleryURL, "/")
	parts := strings.Split(clean, "/")
	return parts[len(parts)-1]
}

func ExtractImagesFromRoot(root *goquery.Document) *[]string {
	var images []string

	root.
		Find("div#display_image_detail").
		Find("img").
		Each(func(i int, s *goquery.Selection) {
			src, ok := s.Attr("src")
			if !ok {
				return
			}

			images = append(images, src)
		})

	return &images
}

func ExtractImagesDiv(r io.Reader) (*[]string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var images []string

	doc.
		Find("div#display_image_detail").
		Find("img").
		Each(func(i int, s *goquery.Selection) {
			src, ok := s.Attr("src")
			if !ok {
				return
			}

			images = append(images, src)
		})

	return &images, nil
}

func ExtractPaginatorData(galleryURL string) (*[]string, *goquery.Document, error) {
	res, err := internal.NewHttpGet(http.DefaultClient, galleryURL)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	parsedUrl, err := url.Parse(galleryURL)
	if err != nil {
		return nil, nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var (
		host    = fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
		content []string
	)

	s := doc.
		Find("div#paginator").
		Find("a").
		FilterFunction(func(i int, s *goquery.Selection) bool {
			title := s.Text()
			i, err := strconv.Atoi(title)
			if err != nil {
				return false
			}
			return i > 0
		})

	bar := progressbar.NewOptions(
		s.Length()+1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(
			fmt.Sprintf(
				"[bold][cyan][Analyzer][reset] %s",
				"discovering pages",
			),
		),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	bar.Add(1)

	s.Each(func(i int, s *goquery.Selection) {
		u, ok := (s.Attr("href"))
		if !ok {
			return
		}

		content = append(content, host+u)
		bar.Add(1)
	})

	return &content, doc, nil
}

func ExtractImagesFromPage(page string) (*[]string, error) {
	res, err := internal.NewHttpGet(http.DefaultClient, page)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ExtractImagesDiv(res.Body)
}
