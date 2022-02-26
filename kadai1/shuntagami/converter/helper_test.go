// http://craigwickesser.com/2015/01/capture-stdout-in-go/

package converter_test

import (
	"bytes"
	"io"
	"os"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
