package parser

import (
	"reflect"
	"testing"

	"github.com/chris-skud/go-chordpro/ast"
)

func TestDefineForms(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want ast.ChordDefinition
	}{
		// The current form, with keywords.
		{"{define: C#m base-fret 4 frets 1 1 3 3 2 1}", ast.ChordDefinition{Name: "C#m", BaseFret: 4, Frets: []int{1, 1, 3, 3, 2, 1}}},
		{"{define: A/D base-fret 1 frets x x 0 2 2 0}", ast.ChordDefinition{Name: "A/D", BaseFret: 1, Frets: []int{-1, -1, 0, 2, 2, 0}}},
		// Whitespace instead of a colon.
		{"{define G base-fret 1 frets 3 2 0 0 3 3}", ast.ChordDefinition{Name: "G", BaseFret: 1, Frets: []int{3, 2, 0, 0, 3, 3}}},
		// base-fret 0 means the first fret.
		{"{define: C/E base-fret 0 frets 0 3 2 0 1 0}", ast.ChordDefinition{Name: "C/E", BaseFret: 1, Frets: []int{0, 3, 2, 0, 1, 0}}},
		// The older form: base fret, then one number per string; - and -1 mute.
		{"{define: D9 4 2 2 2 1 2 -1}", ast.ChordDefinition{Name: "D9", BaseFret: 4, Frets: []int{2, 2, 2, 1, 2, -1}}},
		{"{define: G#m7-5 1 2 1 1 0 - -}", ast.ChordDefinition{Name: "G#m7-5", BaseFret: 1, Frets: []int{2, 1, 1, 0, -1, -1}}},
		// Fingers are kept when given.
		{"{define: Am base-fret 1 frets x 0 2 2 1 0 fingers 0 0 2 3 1 0}", ast.ChordDefinition{Name: "Am", BaseFret: 1, Frets: []int{-1, 0, 2, 2, 1, 0}, Fingers: []int{0, 0, 2, 3, 1, 0}}},
		// A ukulele shape (four strings).
		{"{define: C base-fret 1 frets 0 0 0 3}", ast.ChordDefinition{Name: "C", BaseFret: 1, Frets: []int{0, 0, 0, 3}}},
	} {
		song, err := ParseString(tc.src + "\n[C]la")
		if err != nil {
			t.Fatal(err)
		}
		if len(song.Definitions) != 1 || !reflect.DeepEqual(song.Definitions[0], tc.want) {
			t.Errorf("%s:\n got %+v\nwant %+v", tc.src, song.Definitions, tc.want)
		}
	}
}

func TestDefineKeepsOrderAndLastWins(t *testing.T) {
	song, _ := ParseString("{define: G 1 3 2 0 0 0 3}\n{define: D 1 x x 0 2 3 2}\n{define: G 3 1 3 3 2 1 1}\n[G]la")
	if len(song.Definitions) != 2 || song.Definitions[0].Name != "G" || song.Definitions[0].BaseFret != 3 || song.Definitions[1].Name != "D" {
		t.Fatalf("got %+v", song.Definitions)
	}
}

func TestBadDefineIsKeptAsADirective(t *testing.T) {
	song, _ := ParseString("{define: Am}\n{define: X base-fret one frets 1}\n[C]la")
	if len(song.Definitions) != 0 {
		t.Fatalf("bad defines were accepted: %+v", song.Definitions)
	}
	n := 0
	for _, it := range song.Items {
		if d, ok := it.(ast.Directive); ok && d.Name == "define" {
			n++
		}
	}
	if n != 2 {
		t.Fatalf("bad defines should stay as directives, got %d", n)
	}
}

func TestDirectiveValueAfterWhitespace(t *testing.T) {
	song, _ := ParseString("{title My Song}\n[C]la")
	if got := song.MetaFirst("title"); got != "My Song" {
		t.Fatalf("title = %q", got)
	}
}
