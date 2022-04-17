package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/shuntagami/go-dl/helper"
	"golang.org/x/sync/errgroup"
)

// Config holds value from user input
type Config struct {
	URL            string
	Concurrency    int
	OutFilename    string
	CopyBufferSize int
	Resume         bool
}

// downloader holds necessary value for downloading
type downloader struct {
	acceptRanges bool
	Paused       bool
	config       *Config
	contentSize  int
	bar          *progressbar.ProgressBar
}

// Initialize initializes package and check if server support range download by making http head request
func Initialize(config *Config) (*downloader, error) {
	if config.URL == "" {
		return nil, errors.New("URL is empty")
	}
	if config.Concurrency < 1 {
		config.Concurrency = 1
		log.Print("Concurrency level: 1")
	}
	if config.OutFilename == "" {
		config.OutFilename = helper.FileNameFromURL(config.URL)
	}
	if config.CopyBufferSize == 0 {
		config.CopyBufferSize = 1024
	}

	d := &downloader{config: config}

	// rename file if such file already exist
	if !d.config.Resume {
		d.renameFilenameIfNecessary()
	}
	res, err := http.Head(d.config.URL)
	if err != nil {
		return nil, err
	}

	if res.Header.Get("Accept-Ranges") == "bytes" {
		d.acceptRanges = true
		var err error
		if d.contentSize, err = strconv.Atoi(res.Header.Get("Content-Length")); err != nil {
			return nil, err
		}
	}

	log.Printf("Output file: %s", filepath.Base(config.OutFilename))
	return d, nil
}

func (d *downloader) Download(ctx context.Context) error {
	if d.acceptRanges && d.config.Concurrency != 1 {
		if err := d.multiDownload(ctx); err != nil {
			return err
		}
	} else {
		if err := d.singleDownload(); err != nil {
			return err
		}
	}
	return nil
}

// Add a number to the filename if file already exist
// For example, if filename `hello.pdf` already exist
// it returns hello(1).pdf
func (d *downloader) renameFilenameIfNecessary() {
	if _, err := os.Stat(d.config.OutFilename); err == nil {
		counter := 1
		filename, ext := helper.FileNameAndExt(d.config.OutFilename)
		outDir := filepath.Dir(d.config.OutFilename)

		for err == nil {
			log.Printf("File %s%s already exist", filename, ext)
			newFilename := fmt.Sprintf("%s(%d)%s", filename, counter, ext)
			d.config.OutFilename = filepath.Join(outDir, newFilename)
			_, err = os.Stat(d.config.OutFilename)
			counter += 1
		}
	}
}

func (d *downloader) partFileName(partNum int) string {
	return d.config.OutFilename + ".part" + strconv.Itoa(partNum)
}

// Server does not support partial download for this file
func (d *downloader) singleDownload() error {
	if d.config.Resume {
		return errors.New("cannot resume. Must be downloaded again")
	}

	res, err := http.Get(d.config.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// create the output file
	f, err := os.OpenFile(d.config.OutFilename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	d.bar = progressbar.DefaultBytes(int64(res.ContentLength), "downloading")

	// copy to output file
	buffer := make([]byte, d.config.CopyBufferSize)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), res.Body, buffer)
	if err != nil {
		return err
	}
	return nil
}

// download concurrently
func (d *downloader) multiDownload(ctx context.Context) error {
	d.bar = progressbar.DefaultBytes(int64(d.contentSize), "downloading")

	eg := errgroup.Group{}
	for i := 1; i <= d.config.Concurrency; i++ {
		i := i

		// handle resume
		downloaded := 0
		if d.config.Resume {
			filePath := d.partFileName(i)
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			fileInfo, err := f.Stat()
			if err != nil {
				return err
			}
			downloaded = int(fileInfo.Size())
			// update progress bar
			d.bar.Add64(int64(downloaded))
		}

		eg.Go(func() error {
			return d.downloadPartial(ctx, downloaded, i)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if !d.Paused {
		if err := d.merge(); err != nil {
			return err
		}
	}
	return nil
}

func (d *downloader) merge() error {
	dest, err := os.OpenFile(d.config.OutFilename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer dest.Close()

	for i := 1; i <= d.config.Concurrency; i++ {
		filename := d.partFileName(i)
		source, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dest, source); err != nil {
			return err
		}
		if err := source.Close(); err != nil {
			return err
		}
		if err := os.Remove(filename); err != nil {
			return err
		}
	}
	return nil
}

func (d *downloader) downloadPartial(ctx context.Context, downloaded int, partialNum int) error {
	var from, to int
	partSize := d.contentSize / d.config.Concurrency

	if partialNum == 1 {
		from = downloaded
		to = from + partSize - downloaded
	} else if partialNum == d.config.Concurrency {
		from = downloaded + (partialNum-1)*(partSize+1)
		to = d.contentSize
	} else {
		from = downloaded + (partialNum-1)*(partSize+1)
		to = from + partSize - downloaded
	}

	req, err := http.NewRequest("GET", d.config.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", from, to))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// create the output file
	outputPath := d.partFileName(partialNum)
	flags := os.O_CREATE | os.O_WRONLY
	if d.config.Resume {
		flags = flags | os.O_APPEND
	}
	f, err := os.OpenFile(outputPath, flags, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// copy to output file
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err = io.CopyN(io.MultiWriter(f, d.bar), res.Body, int64(d.config.CopyBufferSize))
			if err != nil {
				if err == io.EOF {
					return nil
				} else {
					return err
				}
			}
		}
	}
}
