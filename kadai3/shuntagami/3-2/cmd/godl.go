package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/shuntagami/go-dl/downloader"
)

func main() {
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	url := flag.String("u", "", "* Download url")
	concurrency := flag.Int("n", 1, "Concurrency level")
	filename := flag.String("f", "", "Output file name")
	bufferSize := flag.Int("buffer-size", 32*1024, "The buffer size to copy from http response body")
	resume := flag.Bool("resume", false, "Resume the download")

	flag.Parse()
	if *url == "" {
		log.Fatal("Please specify the url using -u parameter")
	}

	d, err := downloader.Initialize(&downloader.Config{
		URL:            *url,
		Concurrency:    *concurrency,
		OutFilename:    *filename,
		CopyBufferSize: *bufferSize,
		Resume:         *resume,
	})
	if err != nil {
		log.Fatal(err)
	}

	termCh := make(chan os.Signal)
	signal.Notify(termCh, os.Interrupt)
	go func() {
		<-termCh
		fmt.Println("\nExiting ...")
		d.Paused = true
		cancel()
		fmt.Println("\nDownload has paused. Resume it again with -resume=true parameter.")
	}()

	if err := d.Download(cancelCtx); err != nil {
		cancel()
		log.Fatal(err, "\nDownload failed. Try again.")
	}

	if !d.Paused {
		fmt.Println("\nDownload complete")
	}
}
