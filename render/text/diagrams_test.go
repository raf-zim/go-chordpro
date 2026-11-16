package text

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chris-skud/go-chordpro/parser"
	"github.com/chris-skud/go-chordpro/render"
)

func TestTextListsDefinedChords(t *testing.T) {
	song, _ := parser.ParseString("{title: T}\n{define: Am base-fret 1 frets x 0 2 2 1 0}\n{define: C#m base-fret 4 frets 1 1 3 3 2 1}\n[Am]la")
	var buf bytes.Buffer
	New(render.Options{}).Render(&buf, song)
	out := buf.String()
	for _, want := range []string{"Chords:\n", "  Am   x 0 2 2 1 0\n", "  C#m  base-fret 4: 1 1 3 3 2 1\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	buf.Reset()
	New(render.Options{Transpose: 1}).Render(&buf, song)
	if strings.Contains(buf.String(), "Chords:") {
		t.Error("chord list shown for a transposed song")
	}
}
