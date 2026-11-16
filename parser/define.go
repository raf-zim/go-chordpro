package parser

import (
	"strconv"
	"strings"

	"github.com/chris-skud/go-chordpro/ast"
)

// parseDefine reads a {define} value in either form:
//
//	name base-fret N frets f1 f2 … [fingers n1 n2 …] [keys …]
//	name N f1 f2 …                                    (the older form)
//
// Frets are numbers, 0 for open, or x / X / - / -1 / N for muted.
// base-fret 0 is taken as 1. A colon right after the name ("Em7/B:") is
// dropped. It reports false for anything it can't read.
func parseDefine(value string) (ast.ChordDefinition, bool) {
	f := strings.Fields(value)
	if len(f) < 2 {
		return ast.ChordDefinition{}, false
	}
	name := strings.TrimSuffix(f[0], ":")
	if name == "" {
		return ast.ChordDefinition{}, false
	}
	d := ast.ChordDefinition{Name: name, BaseFret: 1}
	rest := f[1:]
	if isKeyword(rest[0]) {
		var list *[]int
		var fingers bool
		for i := 0; i < len(rest); i++ {
			switch strings.ToLower(rest[i]) {
			case "base-fret":
				if i+1 >= len(rest) {
					return d, false
				}
				n, err := strconv.Atoi(rest[i+1])
				if err != nil || n < 0 {
					return d, false
				}
				d.BaseFret = max(n, 1)
				i++
				list = nil
			case "frets":
				list, fingers = &d.Frets, false
			case "fingers":
				list, fingers = &d.Fingers, true
			case "keys":
				list = nil // keyboard chords are not drawn
				for i+1 < len(rest) && !isKeyword(rest[i+1]) {
					i++
				}
			default:
				if list == nil {
					return d, false
				}
				n, ok := fretValue(rest[i], fingers)
				if !ok {
					return d, false
				}
				*list = append(*list, n)
			}
		}
	} else {
		n, err := strconv.Atoi(rest[0])
		if err != nil || n < 0 {
			return d, false
		}
		d.BaseFret = max(n, 1)
		for _, s := range rest[1:] {
			v, ok := fretValue(s, false)
			if !ok {
				return d, false
			}
			d.Frets = append(d.Frets, v)
		}
	}
	if len(d.Frets) == 0 || (d.Fingers != nil && len(d.Fingers) != len(d.Frets)) {
		return d, false
	}
	return d, true
}

func isKeyword(s string) bool {
	switch strings.ToLower(s) {
	case "base-fret", "frets", "fingers", "keys":
		return true
	}
	return false
}

// fretValue reads one fret (or finger) entry; -1 means muted (no finger).
func fretValue(s string, finger bool) (int, bool) {
	switch s {
	case "x", "X", "-", "N", "n":
		if finger {
			return 0, true
		}
		return -1, true
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < -1 || n > 24 {
		return 0, false
	}
	if finger && n < 0 {
		return 0, true
	}
	return n, true
}

// addDefinition records a shape, replacing an earlier one of the same name.
func addDefinition(song *ast.Song, d ast.ChordDefinition) {
	for i := range song.Definitions {
		if song.Definitions[i].Name == d.Name {
			song.Definitions[i] = d
			return
		}
	}
	song.Definitions = append(song.Definitions, d)
}
