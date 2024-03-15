package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TohoEnjoyer2000/hc-dl/pkg/downloader"
	"github.com/TohoEnjoyer2000/hc-dl/pkg/scraper"
	"github.com/schollz/progressbar/v3"
)

var (
	workingSet  []string
	galleryURL  string
	listFile    string
	concurrency int
)

func init() {
	flag.StringVar(&galleryURL, "u", "", "hentai-cosplay.com url")
	flag.StringVar(&listFile, "a", "", "consume all url from file")
	flag.IntVar(&concurrency, "c", runtime.NumCPU(), "concurrent downloads")
	flag.Parse()
}

func main() {
	info()

	if listFile == "" {
		run(galleryURL)
		return
	}

	file, err := os.ReadFile(listFile)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(file), "\n") {
		run(line)
	}
}

func run(galleryURL string) {
	if galleryURL == "" {
		fmt.Println("Please provide a valid URL")
		os.Exit(1)
	}

	pages, doc, err := scraper.ExtractPaginatorData(galleryURL)
	if err != nil {
		fmt.Println("Please provide a valid URL")
		os.Exit(1)
	}

	images := scraper.ExtractImagesFromRoot(doc)

	workingSet = append(workingSet, *images...)

	for _, page := range *pages {
		images, err := scraper.ExtractImagesFromPage(page)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		workingSet = append(workingSet, *images...)
	}

	fmt.Println()

	bar := progressbar.NewOptions(
		len(workingSet),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(
			fmt.Sprintf(
				"[bold][cyan][Download][reset] %s",
				scraper.DetectName(galleryURL),
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

	dirname := filepath.Join(".", scraper.DetectName(galleryURL))

	downloader.Run(workingSet, dirname, concurrency, bar)

	fmt.Println("")
}

func info() {
	fmt.Println("HC-dl")
	fmt.Println("")
	fmt.Println("hentai-cosplay.com / hentai-img.com dowloader")
	fmt.Println("v1503-2024")
	fmt.Println("")
}
