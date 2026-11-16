package chordlib

import (
	"reflect"
	"testing"
)

func TestLookupCommonShapes(t *testing.T) {
	for name, frets := range map[string][]int{
		"C":  {-1, 3, 2, 0, 1, 0},
		"Am": {-1, 0, 2, 2, 1, 0},
		"G":  {3, 2, 0, 0, 0, 3},
		"E":  {0, 2, 2, 1, 0, 0},
		"D7": {-1, -1, 0, 2, 1, 2},
	} {
		d, ok := Lookup(name)
		if !ok || !reflect.DeepEqual(d.Frets, frets) || d.Name != name || d.BaseFret != 1 {
			t.Errorf("%s: got %+v %v, want frets %v", name, d, ok, frets)
		}
	}
}

func TestLookupSpellings(t *testing.T) {
	same := [][]string{
		{"C#", "Db"}, {"F#m", "Gbm", "F#min", "F#-"}, {"Bb7", "A#7"},
		{"Cmaj7", "CM7", "CMaj7", "CΔ7"}, {"Am7", "Amin7", "A-7"},
		{"Bm7b5", "Bm7-5", "Bø", "Bø7"}, {"Cdim", "Co", "C°"}, {"Caug", "C+"},
		{"Emmaj7", "Em(maj7)", "EmM7"}, {"C69", "C6/9"},
	}
	for _, group := range same {
		first, ok := Lookup(group[0])
		if !ok {
			t.Errorf("%s not found", group[0])
			continue
		}
		for _, n := range group[1:] {
			d, ok := Lookup(n)
			if !ok || !reflect.DeepEqual(d.Frets, first.Frets) || d.BaseFret != first.BaseFret {
				t.Errorf("%s differs from %s: %+v", n, group[0], d)
			}
			if d.Name != n {
				t.Errorf("%s: name %q", n, d.Name)
			}
		}
	}
}

func TestLookupSlashChords(t *testing.T) {
	if d, ok := Lookup("G/B"); !ok || d.Frets[0] == 3 {
		t.Errorf("G/B should have its own shape (B in the bass): %+v %v", d, ok)
	}
	if d, ok := Lookup("Am/G"); !ok || d.Name != "Am/G" {
		t.Errorf("Am/G: %+v %v", d, ok)
	}
	// An unusual bass falls back to the chord without it.
	plain, _ := Lookup("Dsus4")
	if d, ok := Lookup("Dsus4/Bb"); !ok || !reflect.DeepEqual(d.Frets, plain.Frets) || d.Name != "Dsus4/Bb" {
		t.Errorf("Dsus4/Bb fallback: %+v %v", d, ok)
	}
}

func TestLookupRejectsNonChords(t *testing.T) {
	for _, n := range []string{"", "N.C.", "H7", "x", "C#m7b5b9#13", "Hello"} {
		if d, ok := Lookup(n); ok {
			t.Errorf("%q gave %+v", n, d)
		}
	}
}

func TestVoicingsAndHigherPositions(t *testing.T) {
	v := Voicings("A")
	if len(v) < 2 {
		t.Fatalf("A has %d voicings", len(v))
	}
	high := false
	for _, d := range v {
		if d.BaseFret > 1 {
			high = true
		}
	}
	if !high {
		t.Error("no voicing up the neck for A")
	}
}
