package converter_test

import (
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
)

func TestConvert(t *testing.T) {
	tests := []struct {
		from, to, targetDIRName string
		expectedSTDOuts         []string
		expectedFilespath       []string
	}{
		{
			from:          ".png",
			to:            ".jpg",
			targetDIRName: "sample",
			expectedSTDOuts: []string{
				"Successfully converted /workspace/sample/sample2/sample3/dojo4.png, to /workspace/result/sample/sample2/sample3/dojo4.jpg",
				"Successfully converted /workspace/sample/sample2/dojo3.png, to /workspace/result/sample/sample2/dojo3.jpg",
				"Successfully converted /workspace/sample/sample4/sample5/dojo6.PNG, to /workspace/result/sample/sample4/sample5/dojo6.jpg",
				"Successfully converted /workspace/sample/sample4/dojo2.png, to /workspace/result/sample/sample4/dojo2.jpg",
				"Successfully converted /workspace/sample/sample4/dojo5.png, to /workspace/result/sample/sample4/dojo5.jpg",
				"Successfully converted /workspace/sample/dojo1.png, to /workspace/result/sample/dojo1.jpg",
				"Successfully converted /workspace/sample/dojo2.png, to /workspace/result/sample/dojo2.jpg",
				"",
			},
			expectedFilespath: []string{
				filepath.Join(rootDir, saveLocationName),
				filepath.Join(rootDir, saveLocationName, "sample"),
				filepath.Join(rootDir, saveLocationName, "sample", "dojo1.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "dojo2.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2", "dojo3.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2", "sample3"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2", "sample3", "dojo4.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "dojo2.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "dojo5.jpg"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "sample5"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "sample5", "dojo6.jpg"),
			},
		},
		{
			from:          ".png",
			to:            ".jpg",
			targetDIRName: "sample5",
			expectedSTDOuts: []string{
				"Successfully converted /workspace/sample/sample4/sample5/dojo6.PNG, to /workspace/result/sample/sample4/sample5/dojo6.jpg",
				"",
			},
			expectedFilespath: []string{
				filepath.Join(rootDir, saveLocationName),
				filepath.Join(rootDir, saveLocationName, "sample"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "sample5"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "sample5", "dojo6.jpg"),
			},
		},
		{
			from:          ".jpg",
			to:            ".png",
			targetDIRName: "sample",
			expectedSTDOuts: []string{
				"Successfully converted /workspace/sample/sample2/dojo5.jpg, to /workspace/result/sample/sample2/dojo5.png",
				"Successfully converted /workspace/sample/sample4/dojo3.jpg, to /workspace/result/sample/sample4/dojo3.png",
				"",
			},
			expectedFilespath: []string{
				filepath.Join(rootDir, saveLocationName),
				filepath.Join(rootDir, saveLocationName, "sample"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample2", "dojo5.png"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4"),
				filepath.Join(rootDir, saveLocationName, "sample", "sample4", "dojo3.png"),
			},
		},
		{
			from:          ".jpg",
			to:            ".png",
			targetDIRName: "sample5",
			expectedSTDOuts: []string{
				"",
			},
			expectedFilespath: []string{
				filepath.Join(rootDir, saveLocationName),
			},
		},
	}

	for _, test := range tests {
		func() {
			defer os.RemoveAll(filepath.Join(rootDir, saveLocationName, "sample"))

			err = converter.Initialize(test.from, test.to, rootDir)
			out := captureStdout(func() {
				converter.Client.Convert(test.targetDIRName)
			})
			outToSlice := strings.Split(out, "\n")

			// 変換処理を並行処理で行なっているためソートする
			sort.Strings(outToSlice)
			sort.Strings(test.expectedSTDOuts)

			// stdout 確認
			if !reflect.DeepEqual(outToSlice, test.expectedSTDOuts) {
				t.Errorf("actual: %v, expected: %v", outToSlice, test.expectedSTDOuts)
			}

			var filesPath []string
			filesPath, _ = filePathWalkDir(filepath.Join(rootDir, saveLocationName))

			// 生成されたファイルのパスを確認
			if !reflect.DeepEqual(filesPath, test.expectedFilespath) {
				t.Errorf("actual: %v, expected: %v", filesPath, test.expectedFilespath)
			}
		}()
	}
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	return files, err
}
