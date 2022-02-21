package converter

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/shuntagami/dojo1/kadai1/shuntagami/helper"
	"github.com/shuntagami/dojo1/kadai1/shuntagami/validator"
)

var Client Converter

type Converter interface {
	Convert(targetDir string) error
}

type ConverterClient struct {
	From    string
	To      string
	RootDir string
}

// Initialize initializes Converter Client
func Initialize(from, to, rootDir string) error {
	if err := validator.ValidateInput(from, to); err != nil {
		return err
	}

	Client = &ConverterClient{
		From:    from,
		To:      to,
		RootDir: rootDir,
	}

	return nil
}

// Convert converts files under targetDir recursively
func (c *ConverterClient) Convert(targetDirName string) error {
	pathToTargetDir, err := helper.Fullpath(targetDirName)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(pathToTargetDir)
	if err != nil {
		return err
	}

	// ターゲットディレクトリ配下の画像ファイルをループして形式変換する
	ch := make(chan string)
	var wg sync.WaitGroup

	for _, file := range files {
		if !file.IsDir() {
			extension := strings.ToLower(filepath.Ext(file.Name()))
			name := file.Name()[0 : len(file.Name())-len(extension)]
			if extension == c.From {
				wg.Add(1)

				pathToTargetFile := filepath.Join(pathToTargetDir, file.Name())
				targetFile, err := os.Open(pathToTargetFile)
				if err != nil {
					return err
				}
				defer targetFile.Close()

				// ./result配下にPROJECT_ROOT_DIR名を含める必要ないためカットする
				destDir := filepath.Join(c.RootDir, "result", pathToTargetDir[len(c.RootDir):])
				pathToDestFile := filepath.Join(destDir, name+c.To)

				// ディレクトリ作成
				if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
					return err
				}

				// 変換後のファイルを作成
				destFile, err := os.Create(pathToDestFile)
				if err != nil {
					return err
				}

				// 画像ファイルの変換を実行
				go c.convertSingleIMGFile(pathToTargetFile, pathToDestFile, destFile, targetFile, ch, &wg)
			}
		} else {
			// ターゲットディレクトリ配下にさらにディレクトリがあった場合は再起的に変換処理を実行する
			c.Convert(file.Name())
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		fmt.Println(msg)
	}
	return nil
}

// convertSingleIMGFile converts a single img file to jpg or png
func (c *ConverterClient) convertSingleIMGFile(target, dest string, w io.Writer, r io.Reader, ch chan string, wg *sync.WaitGroup) {
	defer (*wg).Done()

	var img image.Image
	var err error

	switch c.From {
	case helper.JPG:
		img, err = jpeg.Decode(r)
		if err == nil {
			err = png.Encode(w, img)
		}
	case helper.PNG:
		img, err = png.Decode(r)
		if err == nil {
			err = jpeg.Encode(w, img, nil)
		}
	}

	if err == nil {
		ch <- fmt.Sprintf("Successfully converted %s, to %s", target, dest)
	} else {
		ch <- fmt.Sprintf("Failed to convert %s, to %s: %s", target, dest, err.Error())
	}
}
