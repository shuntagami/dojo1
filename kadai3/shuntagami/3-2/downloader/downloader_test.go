package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSingleDownload(t *testing.T) {
	files := http.Dir("../testdata/")
	portCh := make(chan int, 1)

	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}
		// notify the port to others
		portCh <- listener.Addr().(*net.TCPAddr).Port
		log.Fatal(http.Serve(listener, http.FileServer(files)))
	}()

	port := <-portCh

	// wait for fileserver to initialize
	time.Sleep(2 * time.Second)

	outFile, err := ioutil.TempFile("", "go_dl_temp_file")
	if err != nil {
		t.Fatal("Coudn't create the output file")
	}
	outFile.Close()
	os.Remove(outFile.Name())

	d, err := Initialize(&Config{
		URL:         fmt.Sprintf("http://localhost:%d/1KB.png", port),
		Concurrency: 1,
		OutFilename: outFile.Name(),
	})
	if err != nil {
		t.Fatal("Coudn't initialize downloader", err)
	}
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := d.Download(cancelCtx); err != nil {
		t.Fatal("Download failed", err)
	}

	original, err := ioutil.ReadFile(filepath.Join("..", "testdata", "1KB.png"))
	if err != nil {
		t.Fatal("Cannot read ../testdata/1KB.png")
	}

	downloaded, err := ioutil.ReadFile(outFile.Name())
	if err != nil {
		t.Fatalf("Cannot read %s", outFile.Name())
	}

	equal := bytes.Equal(original, downloaded)
	if !equal {
		t.Error("Downloaded file is not the same as original file")
	}

	os.Remove(outFile.Name())
}

func TestParallelDownload(t *testing.T) {
	files := http.Dir("../testdata/")
	portCh := make(chan int, 1)

	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}
		// notify the port to others
		portCh <- listener.Addr().(*net.TCPAddr).Port
		log.Fatal(http.Serve(listener, http.FileServer(files)))
	}()

	port := <-portCh

	// wait for fileserver to initialize
	time.Sleep(2 * time.Second)

	outFile, err := ioutil.TempFile("", "go_dl_temp_file")
	if err != nil {
		t.Fatal("Coudn't create the output file")
	}
	outFile.Close()
	// We just want to use this temp filename, so we delete the file,
	// otherwise downloader creates a new file
	os.Remove(outFile.Name())

	d, err := Initialize(&Config{
		URL:         fmt.Sprintf("http://localhost:%d/1KB.png", port),
		Concurrency: 10,
		OutFilename: outFile.Name(),
	})
	if err != nil {
		t.Fatal("Coudn't initialize downloader")
	}

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := d.Download(cancelCtx); err != nil {
		t.Fatal("Download failed", err)
	}

	original, err := ioutil.ReadFile("../testdata/1KB.png")
	if err != nil {
		t.Fatal("Cannot read ../testdata/1KB.png")
	}

	downloaded, err := ioutil.ReadFile(outFile.Name())
	if err != nil {
		t.Fatalf("Cannot read %s", outFile.Name())
	}

	equal := bytes.Equal(original, downloaded)
	if !equal {
		fmt.Println(original)
		fmt.Println(downloaded)
		t.Error("Downloaded file is not the same as original file")
	}

	os.Remove(outFile.Name())
}
