package mdbook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type RenderContext struct {
	Version     string  `json:"version"`
	Root        string  `json:"root"`
	Book        *Book   `json:"book"`
	Config      *Config `json:"config"`
	Destination string  `json:"destination"`
}

type Book struct {
	Sections []*BookItem `json:"sections"`
}

type BookItem struct {
	Chapter   *Chapter   `json:"Chapter,omitempty"`
	Separator *Separator `json:"Separator,omitempty"`
	PartTitle *PartTitle `json:"PartTitle,omitempty"`
}

type Chapter struct {
	Name        string      `json:"name"`
	Content     string      `json:"content"`
	Number      []int       `json:"number"`
	SubItems    []*BookItem `json:"sub_items"`
	Path        string      `json:"path"`
	SourcePath  string      `json:"source_path"`
	ParentNames []string    `json:"parent_names"`
}

type Separator struct {
}

type PartTitle string

type Config struct {
	Book   *BookConfig            `json:"book,omitempty"`
	Build  *BuildConfig           `json:"build,omitempty"`
	Rust   *RustConfig            `json:"rust,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`
}

type BookConfig struct {
	Title         string   `json:"title,omitempty"`
	Authors       []string `json:"authors,omitempty"`
	Description   string   `json:"description,omitempty"`
	Src           string   `json:"src"`
	Multilingual  bool     `json:"multilingual"`
	Language      string   `json:"language,omitempty"`
	TextDirection string   `json:"text_direction,omitempty"`
}

type BuildConfig struct {
	BuildDir                string   `json:"build_dir,omitempty"`
	CreateMissing           bool     `json:"create_missing"`
	UseDefaultPreprocessors bool     `json:"use_default_preprocessors"`
	ExtraWatchDirs          []string `json:"extra_watch_dirs,omitempty"`
}

type RustConfig struct {
	Edition string `json:"edition"`
}

type chapterHolder struct {
	Chapter Chapter `json:"Chapter"`
}

type partTitleHolder struct {
	PartTitle PartTitle `json:"PartTitle"`
}

func (item *BookItem) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	switch v := token.(type) {
	case json.Delim:
		if v == '{' {
			token, err = decoder.Token()
			if err != nil {
				return err
			}
			switch v := token.(type) {
			case string:
				switch v {
				case "Chapter":
					var holder chapterHolder
					err = json.Unmarshal(b, &holder)
					if err != nil {
						return err
					}
					item.Chapter = &holder.Chapter
					return nil
				case "PartTitle":
					var holder partTitleHolder
					err = json.Unmarshal(b, &holder)
					if err != nil {
						return err
					}
					item.PartTitle = &holder.PartTitle
					return nil
				}
			}
		}
	case string:
		if v == "Separator" {
			item.Separator = &Separator{}
			return nil
		}
	}
	return fmt.Errorf("could not decode BookItem")
}

func (item *BookItem) MarshalJSON() ([]byte, error) {
	var output any = nil
	if item.Chapter != nil {
		output = chapterHolder{
			Chapter: *item.Chapter,
		}
	} else if item.Separator != nil {
		output = "Separator"
	} else if item.PartTitle != nil {
		output = partTitleHolder{
			PartTitle: *item.PartTitle,
		}
	}
	return json.Marshal(output)
}

func ParseRenderContextFromReader(r io.Reader) (*RenderContext, error) {
	jsonDecoder := json.NewDecoder(r)
	var renderContext RenderContext
	err := jsonDecoder.Decode(&renderContext)
	if err != nil {
		return nil, err
	}
	return &renderContext, nil
}
