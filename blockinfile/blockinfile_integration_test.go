package blockinfile_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/czerasz/go-blockinfile/blockinfile"
)

func TestUpdate(t *testing.T) {
	variants := [...]struct {
		name            string
		inputContent    string
		expectedContent string
	}{
		{
			"default marker present",
			`1
# BEGIN MANAGED BLOCK
2
# END MANAGED BLOCK
3`,
			`1
# BEGIN MANAGED BLOCK
block content
# END MANAGED BLOCK
3`,
		},
	}

	for _, tc := range variants {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "prefix")
			if err != nil {
				t.Errorf("could NOT create temporary file: %s", err)
			}
			defer os.Remove(file.Name())
			_, err = file.Write([]byte(tc.inputContent))
			if err != nil {
				t.Errorf("could NOT write initial content to temporary file: %s", err)
			}

			err = blockinfile.Update(file.Name(), "", []byte("block content"))
			if err != nil {
				t.Errorf("Update() with %s: expected no error, got: %s", tc.name, err)
			}

			_, err = file.Seek(0, 0)
			if err != nil {
				t.Errorf("could NOT jump to the beginning of the file: %s", err)
			}

			out, err := ioutil.ReadAll(file)
			if err != nil {
				t.Errorf("could NOT read updated content from temporary file: %s", err)
			}

			if !bytes.Equal([]byte(tc.expectedContent), out) {
				t.Errorf("updater.update() with %s: expected file content:\n%s\ngot:\n%s", tc.name, tc.expectedContent, string(out))
			}
		})
	}
}
