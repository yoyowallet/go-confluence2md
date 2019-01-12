package confluence2md

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConvert(t *testing.T) {
	var err error

	files, err := filepath.Glob("testdata/*.html")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file
		t.Run(file, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			output, err := ioutil.ReadFile(file[:len(file)-5] + ".md")
			if err != nil {
				t.Fatal(err)
			}

			var b bytes.Buffer
			err = Convert(f, &b)
			if err != nil {
				t.Fatal(err)
			}

			if string(output) != b.String() {
				t.Errorf("Not equal: \nexpected: %s\nactual  : %s", string(output), b.String())
			}
		})
	}
}
