package blockinfile

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// ErrorMarkerMultiline represents the error when marker has multiple lines.
var ErrorMarkerMultiline = errors.New("marker must be single line")

// DefaultMarkerTemplate represents the default marker template.
const DefaultMarkerTemplate = "# {{.Mark}} MANAGED BLOCK"

// updater is a updateruration struct.
type updater struct {
	Marker string
	Path   string
	Block  []byte
}

func (c *updater) validate() error {
	r := regexp.MustCompile("\n")

	m := r.FindIndex([]byte(c.Marker))
	if m != nil {
		return ErrorMarkerMultiline
	}

	return nil
}

func (c *updater) markerLine(mark string) ([]byte, error) {
	data := struct {
		Mark string
	}{
		Mark: mark,
	}
	tpl := DefaultMarkerTemplate

	if strings.TrimSpace(c.Marker) != "" {
		tpl = c.Marker
	}

	t, err := template.New("line").Parse(tpl)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, data)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (c *updater) update() error {
	beginMarker, err := c.markerLine("BEGIN")
	if err != nil {
		return err
	}

	endMarker, err := c.markerLine("END")
	if err != nil {
		return err
	}

	fileRead, err := os.Open(c.Path)
	if err != nil {
		return err
	}
	defer fileRead.Close()

	in, err := ioutil.ReadAll(fileRead)
	if err != nil {
		return err
	}

	lines := bytes.Split(in, []byte("\n"))
	newLines := [][]byte{}
	ignore := false
	found := false

	for _, line := range lines {
		if !found && bytes.Equal(line, beginMarker) {
			found = true
			ignore = true

			newLines = append(newLines, beginMarker)
			newLines = append(newLines, c.Block)

			continue
		} else if ignore && bytes.Equal(line, endMarker) {
			ignore = false

			newLines = append(newLines, endMarker)

			continue
		}

		if !ignore {
			newLines = append(newLines, line)
		}
	}

	if found && ignore {
		newLines = lines
		newLines = append(newLines, beginMarker)
		newLines = append(newLines, c.Block)
		newLines = append(newLines, endMarker)
	} else if !found {
		newLines = append(newLines, beginMarker)
		newLines = append(newLines, c.Block)
		newLines = append(newLines, endMarker)
	}

	fileWrite, err := os.Create(c.Path)
	if err != nil {
		return err
	}

	_, err = fileWrite.Write(bytes.Join(newLines, []byte("\n")))
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	return nil
}

// Update inserts the block in the file.
func Update(path, marker string, block []byte) error {
	c := updater{
		Marker: marker,
		Path:   path,
		Block:  block,
	}

	err := c.validate()
	if err != nil {
		return err
	}

	err = c.update()

	return err
}
