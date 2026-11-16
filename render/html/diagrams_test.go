package html

import (
	"strings"
	"testing"

	"github.com/chris-skud/go-chordpro/render"
)

func TestDiagramsForDefinedChords(t *testing.T) {
	out := renderString(t, "{title: T}\n{define: Am base-fret 1 frets x 0 2 2 1 0}\n{define: C#m base-fret 4 frets 1 1 3 3 2 1}\n[Am]la [C#m]lee", render.Options{NoStyle: true})
	for _, want := range []string{
		`<section class="chord-diagrams">`,
		`<p class="section-label">Chords</p>`,
		`<figure class="diagram">`,
		`<figcaption>Am</figcaption>`,
		`<figcaption>C#m</figcaption>`,
		`aria-label="Chord diagram: Am"`,
		`class="nut"`,      // Am starts at the nut
		`class="base-fret"`, // C#m starts at fret 4
		`>4</text>`,
		`class="muted"`, // Am's low E
		`class="open"`,  // Am's A string
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
	// Diagrams come after the header and before the song.
	if h, d, m := strings.Index(out, "</header>"), strings.Index(out, "chord-diagrams"), strings.Index(out, "<main>"); !(h < d && d < m) {
		t.Errorf("diagrams are not between the header and the song")
	}
	// No inline styles: sites with a strict Content Security Policy block them.
	if strings.Contains(out, "style=") {
		t.Error("diagram uses a style attribute")
	}
	// Each fretted string gets a dot: Am has 3, C#m 6.
	if n := strings.Count(out, `class="dot"`); n != 9 {
		t.Errorf("dots = %d, want 9", n)
	}
}

func TestNoDiagramsWithoutDefinitionsOrWhenTransposed(t *testing.T) {
	if strings.Contains(renderString(t, "[Am]la", render.Options{}), `<section class="chord-diagrams">`) {
		t.Error("diagrams without any {define}")
	}
	// A transposed song would need different shapes, so none are shown.
	if strings.Contains(renderString(t, "{define: Am 1 x 0 2 2 1 0}\n[Am]la", render.Options{Transpose: 2}), `<section class="chord-diagrams">`) {
		t.Error("diagrams shown for a transposed song")
	}
}

func TestDiagramNamesAreEscaped(t *testing.T) {
	out := renderString(t, "{define: <b>&x 1 0 0 0 0 0 0}\n[C]la", render.Options{})
	if strings.Contains(out, "<b>&x") || !strings.Contains(out, "&lt;b&gt;&amp;x") {
		t.Errorf("name not escaped:\n%s", out)
	}
}
