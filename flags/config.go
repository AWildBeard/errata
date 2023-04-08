package flags

import (
	"gopkg.in/yaml.v3"
	"io"
)

type Config struct {
	Author              string
	DescriptionElements []DescriptionElement
	Examples            []Example
}

type Example struct {
	Example     string
	Description string
}

type DescriptionElementKind string

const (
	Paragraph   DescriptionElementKind = "Paragraph"
	BulletPoint                        = "BulletPoint"
)

type DescriptionElement struct {
	Kind    DescriptionElementKind
	Content string
}

// WriteExampleConfig writes an example YAML config to the provided writer
func WriteExampleConfig(out io.Writer) (err error) {
	exampleConfig := &Config{
		Author: "Your Name",
		DescriptionElements: []DescriptionElement{
			{
				Kind:    Paragraph,
				Content: "This is example content that get's formatted as a paragraph. This can be as short or as long as you want it",
			},
			{
				Kind:    BulletPoint,
				Content: "This get's formatted as a bullet point. Note that you don't need to supply any newlines yourself, though you can if you so choose.",
			},
		},
		Examples: []Example{
			{
				Example:     "command -x -y -n -p",
				Description: "This does x y n p",
			},
			{
				Example:     "command -a -b -c -d",
				Description: "This does a b c d",
			},
		},
	}

	enc := yaml.NewEncoder(out)
	if err = enc.Encode(exampleConfig); err != nil {
		return
	}

	if err = enc.Close(); err != nil {
		return
	}

	return nil
}
