package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/TohoEnjoyer2000/hc-dl/pkg/utils"
	"github.com/schollz/progressbar/v3"
)

var (
	wg                       sync.WaitGroup
	resolutionOverrideRegexp = regexp.MustCompile(`p=[0-9]{3,}`)
)

func Run(
	urls []string,
	dirname string,
	concurrency int,
	bar *progressbar.ProgressBar,
) {
	err := os.Mkdir(dirname, os.ModePerm)
	pipeline := make(chan struct{}, concurrency)

	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(urls))

	for _, galleryURL := range urls {
		pipeline <- struct{}{}

		go func(url string) {
			download(resolutionOverrideRegexp.ReplaceAllString(url, ""), dirname)
			bar.Add(1)
			wg.Done()
			<-pipeline
		}(galleryURL)
	}

	wg.Wait()
}

func download(galleryURL, output string) error {
	res, err := utils.PerformHttpGet(http.DefaultClient, galleryURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	tokens := strings.Split(galleryURL, "/")
	fd, err := os.Create(filepath.Join(output, tokens[len(tokens)-1]))
	if err != nil {
		return err
	}

	defer fd.Close()

	_, err = io.Copy(fd, res.Body)

	return err
}
