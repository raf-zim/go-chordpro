package text

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/chris-skud/go-chordpro/ast"
)

// writeDefinitions lists the {define} shapes after the header, one per line:
//
//	Chords:
//	  Am   x 0 2 2 1 0
//	  C#m  base-fret 4: 1 1 3 3 2 1
//
// x is a muted string. Skipped for a transposed song.
func (r *Renderer) writeDefinitions(w *bufio.Writer, song *ast.Song) {
	if len(song.Definitions) == 0 || r.Transpose != 0 {
		return
	}
	width := 0
	for _, d := range song.Definitions {
		width = max(width, len([]rune(d.Name)))
	}
	fmt.Fprintln(w, "Chords:")
	for _, d := range song.Definitions {
		frets := make([]string, len(d.Frets))
		for i, f := range d.Frets {
			if f < 0 {
				frets[i] = "x"
			} else {
				frets[i] = strconv.Itoa(f)
			}
		}
		pad := strings.Repeat(" ", width-len([]rune(d.Name)))
		if d.BaseFret > 1 {
			fmt.Fprintf(w, "  %s%s  base-fret %d: %s\n", d.Name, pad, d.BaseFret, strings.Join(frets, " "))
		} else {
			fmt.Fprintf(w, "  %s%s  %s\n", d.Name, pad, strings.Join(frets, " "))
		}
	}
	fmt.Fprintln(w)
}
