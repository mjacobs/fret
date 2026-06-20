package music

import "testing"

func TestNoteSpelling(t *testing.T) {
	cases := []struct {
		note  Note
		flats bool
		want  string
	}{
		{C, false, "C"},
		{Db, false, "C#"},
		{Db, true, "Db"},
		{Bb, true, "Bb"},
		{Bb, false, "A#"},
		{Note(13).Normalize(), true, "Db"},
		{Note(-1).Normalize(), false, "B"},
	}
	for _, c := range cases {
		if got := c.note.Name(c.flats); got != c.want {
			t.Errorf("Name(flats=%v) of %d = %q, want %q", c.flats, c.note, got, c.want)
		}
	}
}

func TestBothNames(t *testing.T) {
	if got := C.BothNames(); got != "C" {
		t.Errorf("C.BothNames() = %q, want C", got)
	}
	if got := Gb.BothNames(); got != "F#/Gb" {
		t.Errorf("Gb.BothNames() = %q, want F#/Gb", got)
	}
}

func TestPrefersFlats(t *testing.T) {
	for _, n := range []Note{F, Bb, Eb, Ab, Db} {
		if !n.PrefersFlats() {
			t.Errorf("%v should prefer flats", n)
		}
	}
	for _, n := range []Note{C, G, D, A, E, B, Gb} {
		if n.PrefersFlats() {
			t.Errorf("%v should not prefer flats", n)
		}
	}
}

func TestDegreeOf(t *testing.T) {
	major := Scales[0]
	if label, iv, ok := major.DegreeOf(C, E); !ok || label != "3" || iv != 4 {
		t.Errorf("E in C major = (%q, %d, %v), want (3, 4, true)", label, iv, ok)
	}
	if _, _, ok := major.DegreeOf(C, Bb); ok {
		t.Error("Bb should not be in C major")
	}

	var minPent NoteSet
	for _, s := range Scales {
		if s.Name == "Minor pentatonic" {
			minPent = s
		}
	}
	if label, _, ok := minPent.DegreeOf(A, C); !ok || label != "b3" {
		t.Errorf("C in A minor pentatonic = (%q, %v), want (b3, true)", label, ok)
	}
}

func TestFormulasWellFormed(t *testing.T) {
	all := append(append([]NoteSet{}, Scales...), Arpeggios...)
	for _, s := range all {
		if len(s.Formula) != len(s.Degrees) {
			t.Errorf("%s: formula and degrees length mismatch", s.Name)
		}
		if len(s.Formula) == 0 || s.Formula[0] != 0 {
			t.Errorf("%s: formula must start at 0", s.Name)
		}
		if s.Degrees[0] != "R" {
			t.Errorf("%s: first degree should be R", s.Name)
		}
		for i := 1; i < len(s.Formula); i++ {
			if s.Formula[i] <= s.Formula[i-1] || s.Formula[i] > 11 {
				t.Errorf("%s: formula not strictly ascending within an octave", s.Name)
			}
		}
	}
}

func TestFretboard(t *testing.T) {
	std := Tunings[0]
	board := Fretboard(A, Scales[5], std, 12)
	if len(board) != 6 {
		t.Fatalf("want 6 strings, got %d", len(board))
	}
	for s, row := range board {
		if len(row) != 13 {
			t.Fatalf("string %d: want 13 frets, got %d", s, len(row))
		}
	}
	// High e string, fret 5 is A: the root.
	pos := board[0][5]
	if pos.Note != A || !pos.IsRoot || !pos.InSet {
		t.Errorf("high e fret 5 = %+v, want root A", pos)
	}
	// Open low E is not in A minor pentatonic? E is the 5th — it is in set.
	if pos := board[5][0]; !pos.InSet || pos.Degree != "5" {
		t.Errorf("open low E in A min pent = %+v, want degree 5", pos)
	}
}

func TestTuningLabels(t *testing.T) {
	std := Tunings[0]
	if got := std.Label(0); got != "e" {
		t.Errorf("standard high string label = %q, want e", got)
	}
	if got := std.Label(5); got != "E" {
		t.Errorf("standard low string label = %q, want E", got)
	}
	var ebStd Tuning
	for _, tu := range Tunings {
		if tu.Name == "Eb Standard" {
			ebStd = tu
		}
	}
	if got := ebStd.Label(0); got != "eb" {
		t.Errorf("Eb standard high string label = %q, want eb", got)
	}
	if got := std.LowToHigh(false); got != "E A D G B E" {
		t.Errorf("standard LowToHigh = %q", got)
	}
}

func TestIntervalName(t *testing.T) {
	if got := IntervalName(7); got != "perfect 5th" {
		t.Errorf("IntervalName(7) = %q", got)
	}
	if got := IntervalName(0); got != "root" {
		t.Errorf("IntervalName(0) = %q", got)
	}
	if got := IntervalName(13); got != "minor 2nd" {
		t.Errorf("IntervalName(13) = %q", got)
	}
}

func TestLensIncludes(t *testing.T) {
	open := Lenses[1]
	if !open.Includes(0) || !open.Includes(4) || open.Includes(5) {
		t.Errorf("open position lens bounds wrong: %+v", open)
	}
}
