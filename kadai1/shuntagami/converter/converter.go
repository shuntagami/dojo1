package converter

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	for _, file := range files {
		if !file.IsDir() {
			extension := strings.ToLower(filepath.Ext(file.Name()))
			name := file.Name()[0 : len(file.Name())-len(extension)]
			if extension == c.From {
				target, err := os.Open(filepath.Join(pathToTargetDir, file.Name()))
				if err != nil {
					return err
				}
				defer target.Close()

				// converted配下にPROJECT_ROOT_DIR名を含める必要ないためカットする
				destpath := pathToTargetDir[len(c.RootDir):]

				// ディレクトリ作成
				if err := os.MkdirAll(filepath.Join("./converted/", destpath), os.ModePerm); err != nil {
					return err
				}

				// 変換後のファイルを作成
				dest, err := os.Create(filepath.Join("./converted/", destpath, name+c.To))
				if err != nil {
					return err
				}

				// 画像ファイルの変換を実行
				if err := convertSingleIMGFile(c.From, dest, target); err != nil {
					return err
				}
			}
		} else {
			// ターゲットディレクトリ配下にさらにディレクトリがあった場合は再起的に変換処理を実行する
			c.Convert(file.Name())
		}
	}
	return nil
}

// convertSingleIMGFile converts a single img file to jpg or png
func convertSingleIMGFile(from string, w io.Writer, r io.Reader) error {
	var img image.Image
	var err error

	switch from {
	case helper.JPG:
		img, err = jpeg.Decode(r)
		if err != nil {
			return err
		}
		err = png.Encode(w, img)
	case helper.PNG:
		img, err = png.Decode(r)
		if err != nil {
			return err
		}
		err = jpeg.Encode(w, img, nil)
	}
	return err
}
