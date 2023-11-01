package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/TohoEnjoyer2000/hc-dl/pkg/downloader"
	"github.com/TohoEnjoyer2000/hc-dl/pkg/utils"
	"github.com/schollz/progressbar/v3"
)

var (
	workingSet  []string
	galleryURL  string
	concurrency int
)

func init() {
	flag.StringVar(&galleryURL, "u", "", "hentai-cosplay.com url")
	flag.IntVar(&concurrency, "c", runtime.NumCPU(), "concurrent downloads")
	flag.Parse()
}

func main() {
	if galleryURL == "" {
		fmt.Println("Please provide a valid URL")
		os.Exit(1)
	}

	splash()

	pages, doc, err := utils.ExtractPaginatorData(galleryURL)
	if err != nil {
		fmt.Println("Please provide a valid URL")
		os.Exit(1)
	}

	fmt.Println("Downloading", len(*pages), "pages.")

	images := utils.ExtractImagesFromRoot(doc)

	workingSet = append(workingSet, *images...)

	for _, page := range *pages {
		images, err := utils.ExtractImagesFromPage(page)
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
				utils.DetectName(galleryURL),
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

	dirname := filepath.Join(".", utils.DetectName(galleryURL))

	downloader.Run(workingSet, dirname, concurrency, bar)

	fmt.Println("")
}

func splash() {
	fmt.Println("HC-dl")
	fmt.Println("")
	fmt.Println("hentai-cosplay.com / hentai-img.com dowloader")
	fmt.Println("v0303-2023-1")
	fmt.Println("")
}
