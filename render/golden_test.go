package render_test

// Golden-file tests: each testdata/golden/*.cho is rendered to HTML and to
// text, and compared with the stored .html / .txt next to it. After a
// deliberate change, rewrite the stored output with:
//
//	go test ./render -run Golden -update
//
// and review the diff. Only made-up songs live here: the repository is public.

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chris-skud/go-chordpro/parser"
	"github.com/chris-skud/go-chordpro/render"
	"github.com/chris-skud/go-chordpro/render/html"
	"github.com/chris-skud/go-chordpro/render/text"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestGolden(t *testing.T) {
	songs, err := filepath.Glob(filepath.Join("testdata", "golden", "*.cho"))
	if err != nil || len(songs) == 0 {
		t.Fatalf("no golden songs: %v", err)
	}
	renderers := []struct {
		ext string
		r   render.Renderer
	}{
		{".html", html.New(render.Options{})},
		{".nostyle.html", html.New(render.Options{NoStyle: true})},
		{".txt", text.New(render.Options{})},
		{".transposed.txt", text.New(render.Options{Transpose: 2})},
	}
	for _, path := range songs {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		song, err := parser.ParseString(string(src))
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		for _, rd := range renderers {
			name := strings.TrimSuffix(path, ".cho") + rd.ext
			t.Run(filepath.Base(name), func(t *testing.T) {
				var buf bytes.Buffer
				if err := rd.r.Render(&buf, song); err != nil {
					t.Fatal(err)
				}
				if *update {
					if err := os.WriteFile(name, buf.Bytes(), 0o644); err != nil {
						t.Fatal(err)
					}
					return
				}
				want, err := os.ReadFile(name)
				if err != nil {
					t.Fatalf("missing golden file (run with -update): %v", err)
				}
				if !bytes.Equal(buf.Bytes(), want) {
					t.Errorf("%s changed. Review, then run with -update if intended.\n--- got ---\n%s", name, buf.String())
				}
			})
		}
	}
}
