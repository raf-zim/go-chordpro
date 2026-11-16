package html

import (
	"strings"
	"unicode"

	"github.com/chris-skud/go-chordpro/ast"
)

// isRTLLetter reports whether r is a letter of a right-to-left script
// (Hebrew, Arabic, Syriac, Thaana, N'Ko, … and their presentation forms).
func isRTLLetter(r rune) bool {
	if !unicode.IsLetter(r) {
		return false
	}
	return (r >= 0x0590 && r <= 0x08FF) || (r >= 0xFB1D && r <= 0xFDFF) || (r >= 0xFE70 && r <= 0xFEFF) ||
		(r >= 0x10800 && r <= 0x10FFF) || (r >= 0x1E800 && r <= 0x1EFFF)
}

// strongDir is the direction of the first letter in s ("rtl" or "ltr"),
// or "" when s has no letters (digits, punctuation, and vowel points don't
// decide direction, as in the Unicode bidirectional algorithm).
func strongDir(s string) string {
	for _, r := range s {
		if isRTLLetter(r) {
			return "rtl"
		}
		if unicode.IsLetter(r) {
			return "ltr"
		}
	}
	return ""
}

// lineText is a line's lyrics, without chords.
func lineText(l ast.Line) string {
	var b strings.Builder
	for _, t := range l.Tokens {
		if ly, ok := t.(ast.Lyric); ok {
			b.WriteString(ly.Text)
		}
	}
	return b.String()
}

// songDir is "rtl" when more of the song's lyric lines start with a
// right-to-left letter than with a left-to-right one (a tie goes to the
// title's direction), else "ltr".
func songDir(song *ast.Song) string {
	var rtl, ltr int
	count := func(l ast.Line) {
		switch strongDir(lineText(l)) {
		case "rtl":
			rtl++
		case "ltr":
			ltr++
		}
	}
	for _, it := range song.Items {
		switch v := it.(type) {
		case ast.Line:
			count(v)
		case ast.Section:
			if v.Kind == ast.SectionTab || v.Kind == ast.SectionGrid {
				continue
			}
			for _, l := range v.Lines {
				count(l)
			}
		}
	}
	if rtl > ltr || (rtl == ltr && strongDir(song.MetaFirst("title")) == "rtl") {
		return "rtl"
	}
	return "ltr"
}

// rtlLanguages maps a {meta: language …} name to its language code.
var rtlLanguages = map[string]string{
	"hebrew": "he", "arabic": "ar", "yiddish": "yi", "persian": "fa", "farsi": "fa",
	"urdu": "ur", "ladino": "lad", "aramaic": "arc", "pashto": "ps", "kurdish": "ckb",
}

// songLang is the page language for a right-to-left song: from
// {meta: language …} when it names a known language, else from its script
// (Hebrew letters → "he", otherwise "ar").
func songLang(song *ast.Song) string {
	for _, raw := range song.Meta["meta"] {
		key, value, _ := strings.Cut(strings.TrimSpace(raw), " ")
		if strings.EqualFold(key, "language") {
			if code, ok := rtlLanguages[strings.ToLower(strings.TrimSpace(value))]; ok {
				return code
			}
		}
	}
	for _, it := range song.Items {
		var lines []ast.Line
		switch v := it.(type) {
		case ast.Line:
			lines = []ast.Line{v}
		case ast.Section:
			lines = v.Lines
		}
		for _, l := range lines {
			for _, r := range lineText(l) {
				if isRTLLetter(r) {
					if (r >= 0x0590 && r <= 0x05FF) || (r >= 0xFB1D && r <= 0xFB4F) {
						return "he"
					}
					return "ar"
				}
			}
		}
	}
	return "he"
}
