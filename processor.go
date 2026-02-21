package mdbook

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type Handler interface {
	HandleChapter(chapter *Chapter, contentHandler func(walker ast.Walker) error) error
	HandleSeparator(separator *Separator) error
	HandlePartTitle(partTitle *PartTitle) error
}

type Processor struct {
	renderContext *RenderContext
	handler       Handler
	md            goldmark.Markdown
}

func NewProcessor(renderContext *RenderContext, handler Handler) *Processor {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
		goldmark.WithParserOptions(),
	)
	return &Processor{
		renderContext: renderContext,
		handler:       handler,
		md:            md,
	}
}

func Process(renderContext *RenderContext, handler Handler) error {
	p := NewProcessor(renderContext, handler)
	return p.Process()
}

func (p *Processor) Process() error {
	for _, item := range p.renderContext.Book.GetItems() {
		if err := p.handleBookItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) handleBookItem(bookItem *BookItem) error {
	if bookItem.Chapter != nil {
		return p.handleChapter(bookItem.Chapter)
	} else if bookItem.Separator != nil {
		return p.handler.HandleSeparator(bookItem.Separator)
	} else if bookItem.PartTitle != nil {
		return p.handler.HandlePartTitle(bookItem.PartTitle)
	} else {
		return fmt.Errorf("invalid book item")
	}
}

func (p *Processor) handleChapter(chapter *Chapter) error {
	err := p.handler.HandleChapter(chapter, func(walker ast.Walker) error {
		sourceBytes := []byte(chapter.Content)
		doc := p.md.Parser().Parse(text.NewReader(sourceBytes))
		return ast.Walk(doc, walker)
	})
	if err != nil {
		return err
	}

	for _, subItem := range chapter.SubItems {
		err := p.handleBookItem(subItem)
		if err != nil {
			return err
		}
	}

	return nil
}
