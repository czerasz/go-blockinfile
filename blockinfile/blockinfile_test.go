package blockinfile

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	u := updater{
		Marker: "# {{.Mark}} MY TAG\nMULTILINE",
	}
	if err := u.validate(); !errors.Is(err, ErrorMarkerMultiline) {
		t.Errorf("Expected error: %s, got: %s", ErrorMarkerMultiline, err)
	}
}

func TestMarkerLine(t *testing.T) {
	variants := [...]struct {
		name           string
		updater        updater
		expectedMarker string
	}{
		{"default marker", updater{}, "# BEGIN MANAGED BLOCK"},
		{"custom marker", updater{Marker: "# {{.Mark}} MY TAG"}, "# BEGIN MY TAG"},
	}

	for _, tc := range variants {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.updater.markerLine("BEGIN")
			if err != nil {
				t.Errorf("updater.markerLine(\"BEGIN\") with %s: expected no error, got: %s", tc.name, err)
			}
			if !bytes.Equal([]byte(tc.expectedMarker), b) {
				t.Errorf("updater.markerLine(\"BEGIN\") with %s: expected default marker \"%s\", got: \"%s\"", tc.name, tc.expectedMarker, string(b))
			}
		})
	}
}

var testUpdateVariants = [...]struct {
	name            string
	updater         updater
	inputContent    string
	expectedContent string
}{
	{
		"default marker present",
		updater{},
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
	{
		"default marker missing",
		updater{},
		`1
2
3`,
		`1
2
3
# BEGIN MANAGED BLOCK
block content
# END MANAGED BLOCK`,
	},
	{
		"default marker only open present",
		updater{},
		`1
# BEGIN MANAGED BLOCK
2
3`,
		`1
# BEGIN MANAGED BLOCK
2
3
# BEGIN MANAGED BLOCK
block content
# END MANAGED BLOCK`,
	},
	{
		"custom marker present",
		updater{Marker: "# {{.Mark}} MY TAG"},
		`1
# BEGIN MY TAG
2
# END MY TAG
3`,
		`1
# BEGIN MY TAG
block content
# END MY TAG
3`,
	},
	{
		"custom marker missing",
		updater{Marker: "# {{.Mark}} MY TAG"},
		`1
2
3`,
		`1
2
3
# BEGIN MY TAG
block content
# END MY TAG`,
	},
	{
		"custom marker only open present",
		updater{Marker: "# {{.Mark}} MY TAG"},
		`1
# BEGIN MY TAG
2
3`,
		`1
# BEGIN MY TAG
2
3
# BEGIN MY TAG
block content
# END MY TAG`,
	},
	{
		"custom marker open and close the same",
		updater{Marker: "# MY TAG"},
		`1
# MY TAG
2
# MY TAG
3`,
		`1
# MY TAG
block content
# MY TAG
3`,
	},
}

func TestUpdate(t *testing.T) {
	for _, tc := range testUpdateVariants {
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

			tc.updater.Path = file.Name()
			tc.updater.Block = []byte("block content")

			err = tc.updater.update()
			if err != nil {
				t.Errorf("updater.update() with %s: expected no error, got: %s", tc.name, err)
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
