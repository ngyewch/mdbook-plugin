package mdbook

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRenderContext(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		version string
	}{
		{name: "v0.4", file: "render_context_v04.json", version: "0.4.52"},
		{name: "v0.5", file: "render_context_v05.json", version: "0.5.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.file))
			if err != nil {
				t.Fatalf("failed to open testdata file: %v", err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					t.Fatalf("failed to close testdata file: %v", err)
				}
			}()

			ctx, err := ParseRenderContextFromReader(f)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			if ctx.Version != tt.version {
				t.Errorf("version = %q, want %q", ctx.Version, tt.version)
			}
			if ctx.Root != "/tmp/test-book" {
				t.Errorf("root = %q, want %q", ctx.Root, "/tmp/test-book")
			}
			if ctx.Destination != "/tmp/test-book/book" {
				t.Errorf("destination = %q, want %q", ctx.Destination, "/tmp/test-book/book")
			}

			// Book items (via GetItems for backward compat)
			if ctx.Book == nil {
				t.Fatal("book is nil")
			}
			items := ctx.Book.GetItems()
			if len(items) != 3 {
				t.Fatalf("GetItems() length = %d, want 3", len(items))
			}

			if items[0].Chapter == nil {
				t.Fatal("item 0 should be a Chapter")
			}
			if items[0].Chapter.Name != "Chapter 1" {
				t.Errorf("chapter name = %q, want %q", items[0].Chapter.Name, "Chapter 1")
			}
			if items[0].Chapter.Path != "chapter_1.md" {
				t.Errorf("chapter path = %q, want %q", items[0].Chapter.Path, "chapter_1.md")
			}
			if items[1].Separator == nil {
				t.Error("item 1 should be a Separator")
			}
			if items[2].PartTitle == nil {
				t.Fatal("item 2 should be a PartTitle")
			}
			if *items[2].PartTitle != "Part 2" {
				t.Errorf("part title = %q, want %q", *items[2].PartTitle, "Part 2")
			}

			// Config
			if ctx.Config == nil {
				t.Fatal("config is nil")
			}
			if ctx.Config.Book == nil {
				t.Fatal("config.book is nil")
			}
			if ctx.Config.Book.Title != "Test Book" {
				t.Errorf("config.book.title = %q, want %q", ctx.Config.Book.Title, "Test Book")
			}
			if ctx.Config.Book.Language != "en" {
				t.Errorf("config.book.language = %q, want %q", ctx.Config.Book.Language, "en")
			}

			if ctx.Config.Build == nil {
				t.Fatal("config.build is nil")
			}
			if ctx.Config.Build.BuildDir != "output" {
				t.Errorf("config.build.build-dir = %q, want %q", ctx.Config.Build.BuildDir, "output")
			}
			if !ctx.Config.Build.CreateMissing {
				t.Error("config.build.create-missing = false, want true")
			}
			if !ctx.Config.Build.UseDefaultPreprocessors {
				t.Error("config.build.use-default-preprocessors = false, want true")
			}

			if ctx.Config.Rust == nil {
				t.Fatal("config.rust is nil")
			}
			if ctx.Config.Rust.Edition != "2021" {
				t.Errorf("config.rust.edition = %q, want %q", ctx.Config.Rust.Edition, "2021")
			}

			if ctx.Config.Preprocessor == nil {
				t.Fatal("config.preprocessor is nil")
			}
			if _, ok := ctx.Config.Preprocessor["test"]; !ok {
				t.Error("config.preprocessor[\"test\"] missing")
			}
		})
	}
}
