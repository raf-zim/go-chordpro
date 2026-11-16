package html

import (
	"bufio"
	"fmt"
	"html"

	"github.com/chris-skud/go-chordpro/ast"
)

// Diagram geometry, in SVG user units.
const (
	diagramPad    = 12 // left and right, room for the base-fret number
	diagramTop    = 12 // room for the open/muted marks above the nut
	stringGap     = 10
	fretGap       = 12
	minFretRows   = 4
	dotRadius     = 3.5
	markRadius    = 2.5
	diagramBottom = 4
)

// writeDiagrams draws one chord diagram per {define}, after the header.
// They are skipped for a transposed song, whose chords need other shapes.
func (r *Renderer) writeDiagrams(w *bufio.Writer, song *ast.Song) {
	if len(song.Definitions) == 0 || r.Transpose != 0 {
		return
	}
	fmt.Fprintln(w, `<section class="chord-diagrams">`)
	fmt.Fprintln(w, `  <p class="section-label">Chords</p>`)
	fmt.Fprintln(w, `  <div class="diagrams">`)
	for _, d := range song.Definitions {
		fmt.Fprintf(w, "    <figure class=\"diagram\">%s<figcaption>%s</figcaption></figure>\n", diagramSVG(d), html.EscapeString(d.Name))
	}
	fmt.Fprintln(w, `  </div>`)
	fmt.Fprintln(w, `</section>`)
}

// diagramSVG draws a chord shape: strings left (lowest) to right, frets top
// to bottom, a thick nut at fret 1 or the base-fret number beside the first
// row, ○ for open and × for muted strings, and a dot per fretted string.
// It uses presentation attributes only (no style attributes) and
// currentColor, so a page's CSS decides the colour.
func diagramSVG(d ast.ChordDefinition) string {
	strings := len(d.Frets)
	rows := minFretRows
	for _, f := range d.Frets {
		rows = max(rows, f)
	}
	width := 2*diagramPad + (strings-1)*stringGap
	height := diagramTop + rows*fretGap + diagramBottom
	x := func(s int) int { return diagramPad + s*stringGap }
	y := func(row int) int { return diagramTop + row*fretGap }
	right := x(strings - 1)

	out := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d" role="img" aria-label="Chord diagram: %s" fill="none" stroke="currentColor" stroke-width="1">`,
		width, height, width, height, html.EscapeString(d.Name))
	for row := 0; row <= rows; row++ {
		out += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x(0), y(row), right, y(row))
	}
	for s := 0; s < strings; s++ {
		out += fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d"/>`, x(s), y(0), x(s), y(rows))
	}
	if d.BaseFret <= 1 {
		out += fmt.Sprintf(`<line class="nut" x1="%d" y1="%d" x2="%d" y2="%d" stroke-width="3"/>`, x(0), y(0), right, y(0))
	} else {
		out += fmt.Sprintf(`<text class="base-fret" x="%d" y="%d" text-anchor="end" font-size="8" font-family="sans-serif" fill="currentColor" stroke="none">%d</text>`,
			x(0)-4, y(0)+fretGap/2+3, d.BaseFret)
	}
	for s, f := range d.Frets {
		switch {
		case f < 0:
			c, t := x(s), diagramTop/2
			out += fmt.Sprintf(`<path class="muted" d="M%d %d L%d %d M%d %d L%d %d"/>`, c-2, t-2, c+2, t+2, c-2, t+2, c+2, t-2)
		case f == 0:
			out += fmt.Sprintf(`<circle class="open" cx="%d" cy="%d" r="%g"/>`, x(s), diagramTop/2, markRadius)
		default:
			out += fmt.Sprintf(`<circle class="dot" cx="%d" cy="%g" r="%g" fill="currentColor" stroke="none"/>`, x(s), float64(y(f-1))+fretGap/2.0, dotRadius)
		}
	}
	return out + `</svg>`
}
