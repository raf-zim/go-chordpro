// Package chordlib looks up guitar chord shapes by ChordPro chord name, from
// the chords-db guitar data (MIT License; see LICENSE-chords-db.txt).
package chordlib

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"

	"github.com/chris-skud/go-chordpro/ast"
)

//go:embed guitar.json
var guitarJSON []byte

type position struct {
	Frets    []int `json:"frets"`
	Fingers  []int `json:"fingers"`
	BaseFret int   `json:"baseFret"`
}

type entry struct {
	Key       string     `json:"key"`
	Suffix    string     `json:"suffix"`
	Positions []position `json:"positions"`
}

var (
	loadOnce sync.Once
	shapes   map[string][]position // "C#|m7" → positions
)

func load() {
	var db struct {
		Chords map[string][]entry `json:"chords"`
	}
	if err := json.Unmarshal(guitarJSON, &db); err != nil {
		panic("chordlib: bad embedded guitar.json: " + err.Error())
	}
	shapes = map[string][]position{}
	for _, list := range db.Chords {
		for _, e := range list {
			shapes[e.Key+"|"+e.Suffix] = e.Positions
		}
	}
}

// The keys chords-db uses, by any spelling of the note.
var keys = map[string]string{
	"C": "C", "B#": "C", "C#": "C#", "Db": "C#", "D": "D", "D#": "Eb", "Eb": "Eb",
	"E": "E", "Fb": "E", "E#": "F", "F": "F", "F#": "F#", "Gb": "F#", "G": "G",
	"G#": "Ab", "Ab": "Ab", "A": "A", "A#": "Bb", "Bb": "Bb", "B": "B", "Cb": "B",
}

// Other ways of writing a chord type, mapped to chords-db's suffix.
var suffixes = map[string]string{
	"": "major", "maj": "major", "M": "major",
	"m": "minor", "min": "minor", "-": "minor",
	"o": "dim", "°": "dim", "o7": "dim7", "°7": "dim7",
	"+": "aug", "#5": "aug", "+7": "aug7", "7#5": "aug7", "7+5": "aug7",
	"M7": "maj7", "Maj7": "maj7", "Δ7": "maj7", "Δ": "maj7", "ma7": "maj7",
	"min7": "m7", "-7": "m7", "m7-5": "m7b5", "ø": "m7b5", "ø7": "m7b5", "-7b5": "m7b5",
	"min6": "m6", "-6": "m6", "min9": "m9", "-9": "m9", "m(maj7)": "mmaj7", "mM7": "mmaj7", "m(M7)": "mmaj7",
	"6/9": "69", "2": "add9", "sus2/4": "sus2sus4", "7sus": "7sus4", "M9": "maj9", "Maj9": "maj9",
}

// split reads a chord name into its root, chord type, and bass note.
func split(name string) (root, quality, bass string, ok bool) {
	name = strings.TrimSpace(name)
	if name == "" || name[0] < 'A' || name[0] > 'G' {
		return "", "", "", false
	}
	n := 1
	if len(name) > 1 && (name[1] == '#' || name[1] == 'b') {
		n = 2
	}
	root, rest := name[:n], name[n:]
	if i := strings.LastIndex(rest, "/"); i >= 0 {
		b := rest[i+1:]
		if b != "" && b[0] >= 'A' && b[0] <= 'G' && (len(b) == 1 || (len(b) == 2 && (b[1] == '#' || b[1] == 'b'))) {
			rest, bass = rest[:i], b
		}
	}
	return root, rest, bass, true
}

// suffixFor maps a chord type to chords-db's suffix name, if it has one.
func suffixFor(quality string) (string, bool) {
	if s, ok := suffixes[quality]; ok {
		return s, true
	}
	return quality, quality != ""
}

// find returns the positions for a root, suffix, and bass (bass may be "").
func find(root, quality, bass string) ([]position, bool) {
	loadOnce.Do(load)
	key, ok := keys[root]
	if !ok {
		return nil, false
	}
	suffix, ok := suffixFor(quality)
	if !ok {
		return nil, false
	}
	if bass != "" {
		base := ""
		switch suffix {
		case "major":
		case "minor":
			base = "m"
		default:
			base = suffix
		}
		for _, b := range bassSpellings(bass) {
			if p, ok := shapes[key+"|"+base+"/"+b]; ok {
				return p, true
			}
		}
	}
	p, ok := shapes[key+"|"+suffix]
	return p, ok
}

// bassSpellings lists the ways chords-db might name a bass note.
func bassSpellings(b string) []string {
	out := []string{b}
	for spelling, k := range keys {
		if spelling != b && k == keys[b] {
			out = append(out, spelling)
		}
	}
	return out
}

func toDefinition(name string, p position) ast.ChordDefinition {
	d := ast.ChordDefinition{Name: name, BaseFret: max(p.BaseFret, 1), Frets: append([]int(nil), p.Frets...)}
	if len(p.Fingers) == len(p.Frets) {
		d.Fingers = append([]int(nil), p.Fingers...)
	}
	return d
}

// Lookup returns the most common guitar shape for a chord name such as "Am",
// "C#m7", "Bbmaj7", or "G/B" (a slash chord falls back to the chord without
// its bass when there is no shape for it). Frets count from BaseFret, as in
// {define}; -1 is a muted string.
func Lookup(name string) (ast.ChordDefinition, bool) {
	v := Voicings(name)
	if len(v) == 0 {
		return ast.ChordDefinition{}, false
	}
	return v[0], true
}

// Voicings returns every guitar shape for a chord name, most common first.
func Voicings(name string) []ast.ChordDefinition {
	root, quality, bass, ok := split(name)
	if !ok {
		return nil
	}
	ps, ok := find(root, quality, bass)
	if !ok {
		return nil
	}
	out := make([]ast.ChordDefinition, len(ps))
	for i, p := range ps {
		out[i] = toDefinition(strings.TrimSpace(name), p)
	}
	return out
}
