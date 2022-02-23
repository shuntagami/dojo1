package converter_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/shuntagami/dojo1/kadai1/shuntagami/converter"
)

const (
	saveLocationName = "result"
)

var (
	rootDir = os.Getenv("PROJECT_ROOT_DIR")
	err     error
	files   []fs.FileInfo
)

func TestConvert_sample_from_png_to_jpg(t *testing.T) {
	defer os.RemoveAll(filepath.Join(rootDir, "result", "sample"))

	from := ".png"
	to := ".jpg"
	targetDirName := "sample"

	// 初期化
	if err = converter.Initialize(from, to, rootDir); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// 実行
	out := captureStdout(func() {
		converter.Client.Convert(targetDirName)
	})
	out = strings.TrimSuffix(out, "")

	// stdout 確認
	outToSlice := strings.Split(out, "\n")
	if len(outToSlice) != 8 {
		t.Errorf("actual: %v, expected: %v", len(outToSlice), 8)
	}
	expectedOuts := []string{
		"Successfully converted /workspace/sample/sample2/sample3/dojo4.png, to /workspace/result/sample/sample2/sample3/dojo4.jpg",
		"Successfully converted /workspace/sample/sample2/dojo3.png, to /workspace/result/sample/sample2/dojo3.jpg",
		"Successfully converted /workspace/sample/sample4/sample5/dojo6.PNG, to /workspace/result/sample/sample4/sample5/dojo6.jpg",
		"Successfully converted /workspace/sample/sample4/dojo2.png, to /workspace/result/sample/sample4/dojo2.jpg",
		"Successfully converted /workspace/sample/sample4/dojo5.png, to /workspace/result/sample/sample4/dojo5.jpg",
		"Successfully converted /workspace/sample/dojo1.png, to /workspace/result/sample/dojo1.jpg",
		"Successfully converted /workspace/sample/dojo2.png, to /workspace/result/sample/dojo2.jpg",
		"",
	}
	sort.Strings(outToSlice)
	sort.Strings(expectedOuts)
	if !reflect.DeepEqual(outToSlice, expectedOuts) {
		t.Errorf("actual: %v, expected: %v", outToSlice, expectedOuts)
	}

	// /workspace/result/sample 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 4 {
		t.Errorf("actual: %v, expected: %v", len(files), 4)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo1.jpg" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo1.jpg")
	}
	if files[1].Name() != "dojo2.jpg" {
		t.Errorf("actual: %v, expected: %v", files[1].Name(), "dojo2.jpg")
	}
	if files[2].Name() != "sample2" {
		t.Errorf("actual: %v, expected: %v", files[2].Name(), "sample2")
	}
	if files[3].Name() != "sample4" {
		t.Errorf("actual: %v, expected: %v", files[3].Name(), "sample4")
	}

	// /workspace/result/sample/sample2 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample2"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 2 {
		t.Errorf("actual: %v, expected: %v", len(files), 2)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo3.jpg" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo3.jpg")
	}
	if files[1].Name() != "sample3" {
		t.Errorf("actual: %v, expected: %v", files[1].Name(), "sample3")
	}

	// /workspace/result/sample/sample2/sample3 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample2", "sample3"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo4.jpg" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo4.jpg")
	}
}

func TestConvert_sample_from_jpg_to_png(t *testing.T) {
	defer os.RemoveAll(filepath.Join(rootDir, "result", "sample"))

	from := ".jpg"
	to := ".png"
	targetDirName := "sample"

	// 初期化
	if err = converter.Initialize(from, to, rootDir); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// 実行
	if err := converter.Client.Convert(targetDirName); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// /workspace/result/sample 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 2 {
		t.Errorf("actual: %v, expected: %v", len(files), 2)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "sample2" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "sample2")
	}
	if files[1].Name() != "sample4" {
		t.Errorf("actual: %v, expected: %v", files[1].Name(), "sample4")
	}

	// /workspace/result/sample/sample2 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample2"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo5.png" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo5.png")
	}

	// /workspace/result/sample/sample4 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample4"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo3.png" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo3.png")
	}
}

func TestConvert_sample5_from_png_to_jpg(t *testing.T) {
	defer os.RemoveAll(filepath.Join(rootDir, "result", "sample"))

	from := ".png"
	to := ".jpg"
	targetDirName := "sample5"

	// 初期化
	if err = converter.Initialize(from, to, rootDir); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// 実行
	if err := converter.Client.Convert(targetDirName); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// /workspace/result/sample 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "sample4" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "sample4")
	}

	// /workspace/result/sample/sample4 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample4"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "sample5" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "sample5")
	}

	// /workspace/result/sample/sample4/sample5 配下のファイル数を確認
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName, "sample", "sample4", "sample5"))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 1 {
		t.Errorf("actual: %v, expected: %v", len(files), 1)
	}
	//  生成されたファイル名を確認
	if files[0].Name() != "dojo6.jpg" {
		t.Errorf("actual: %v, expected: %v", files[0].Name(), "dojo6.jpg")
	}
}

func TestConvert_sample5_from_jpg_to_png(t *testing.T) {
	defer os.RemoveAll(filepath.Join(rootDir, "result", "sample"))

	from := ".jpg"
	to := ".png"
	targetDirName := "sample5"

	// 初期化
	if err = converter.Initialize(from, to, rootDir); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// 実行
	if err := converter.Client.Convert(targetDirName); err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}

	// /workspace/result/sample 配下のファイル数を確認(変換対象ファイルないのでファイル数は0)
	files, err = ioutil.ReadDir(filepath.Join(rootDir, saveLocationName))
	if err != nil {
		t.Errorf("actual: %v, expected: %v", err, nil)
	}
	if len(files) != 0 {
		t.Errorf("actual: %v, expected: %v", len(files), 0)
	}
}
