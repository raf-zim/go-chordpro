package html

import (
	"strings"
	"testing"

	"github.com/chris-skud/go-chordpro/render"
)

func TestHebrewSongIsRightToLeft(t *testing.T) {
	out := renderString(t, "{title: הבה נגילה}\n{artist: עממי}\n{meta: language Hebrew}\n[Dm]הבה נגי[Gm]לה\n[A7]ונשמחה\n", render.Options{NoStyle: true})
	for _, want := range []string{
		`<html lang="he" dir="rtl">`,
		`<header dir="rtl">`,
		`<main dir="rtl">`,
		`<p class="line" dir="auto">`,
		// Chord names stay left-to-right inside Hebrew text.
		`<span class="chord" dir="ltr" data-chord="Dm">Dm</span>`,
		// Header values are isolated from the English labels.
		`Artist: <bdi>עממי</bdi>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
}

func TestMostlyEnglishSongStaysLeftToRight(t *testing.T) {
	out := renderString(t, "{title: Mixed}\n[C]One line\n[G]Two lines\n[D]שורה אחת\n", render.Options{NoStyle: true})
	if strings.Contains(out, `dir="rtl"`) {
		t.Errorf("a mostly English song became right-to-left:\n%s", out)
	}
	if !strings.Contains(out, `<html lang="en">`) {
		t.Error("lang changed for an English song")
	}
	// Each line still finds its own direction, so the Hebrew line reads right-to-left.
	if strings.Count(out, `<p class="line" dir="auto">`) != 3 {
		t.Errorf("lines without dir=auto:\n%s", out)
	}
}

func TestDirectionComesFromLetters(t *testing.T) {
	for _, tc := range []struct{ src, lang string }{
		{"[Am]مرحبا بالعالم\n[G]أهلا\n", `<html lang="ar" dir="rtl">`},
		{"{meta: language Yiddish}\n[Am]שלום עליכם\n", `<html lang="yi" dir="rtl">`},
		{"{title: Chords only}\n[Am] [G] [C]\n", `<html lang="en">`},
	} {
		if out := renderString(t, tc.src, render.Options{NoStyle: true}); !strings.Contains(out, tc.lang) {
			t.Errorf("%q: missing %s", tc.src, tc.lang)
		}
	}
}

func TestCommentsAndLabelsFindTheirDirection(t *testing.T) {
	out := renderString(t, "{c: פעמיים}\n{soc: פזמון}\n[C]שלום\n{eoc}\n", render.Options{NoStyle: true})
	for _, want := range []string{`<p class="comment" dir="auto">פעמיים</p>`, `<p class="section-label" dir="auto">פזמון</p>`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
}
