package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func sized(m Model) Model {
	next, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return next.(Model)
}

func press(t *testing.T, m Model, keys ...string) Model {
	t.Helper()
	for _, k := range keys {
		next, _ := m.Update(tea.KeyPressMsg{Code: rune(k[0]), Text: k})
		m = next.(Model)
	}
	return m
}

func TestRenderSmoke(t *testing.T) {
	m := sized(NewModel())
	out := m.render()
	if out == "" {
		t.Fatal("render produced nothing")
	}
	if !strings.Contains(out, "fret") {
		t.Error("render missing app title")
	}
	if !strings.Contains(out, "minor pentatonic") {
		t.Error("render missing default scale title")
	}
}

func TestCursorStaysOnBoard(t *testing.T) {
	m := sized(NewModel())
	for range 30 {
		m = press(t, m, "l")
	}
	if m.curFret > m.visibleFrets() {
		t.Errorf("cursor fret %d beyond visible %d", m.curFret, m.visibleFrets())
	}
	for range 10 {
		m = press(t, m, "k")
	}
	if m.curStr != 0 {
		t.Errorf("cursor string = %d, want 0", m.curStr)
	}
}

func TestPinToggle(t *testing.T) {
	m := sized(NewModel())
	m = press(t, m, " ")
	if len(m.pins) != 1 {
		t.Fatalf("want 1 pin, got %d", len(m.pins))
	}
	m = press(t, m, " ")
	if len(m.pins) != 0 {
		t.Fatalf("pin did not toggle off, got %d", len(m.pins))
	}
}

func TestTranspose(t *testing.T) {
	m := sized(NewModel())
	start := m.rootIdx
	m = press(t, m, ".")
	if m.rootIdx != (start+1)%12 {
		t.Errorf("transpose up: rootIdx = %d", m.rootIdx)
	}
	m = press(t, m, ",", ",")
	if m.rootIdx != (start+11)%12 {
		t.Errorf("transpose down: rootIdx = %d", m.rootIdx)
	}
}

func TestPickerFlow(t *testing.T) {
	m := sized(NewModel())
	m = press(t, m, "s") // open scale picker
	if m.overlay != overlayScale {
		t.Fatalf("overlay = %v, want scale picker", m.overlay)
	}
	if m.renderOverlay() == "" {
		t.Fatal("overlay rendered nothing")
	}
	m = press(t, m, "j") // move selection
	idx := m.pickerIdx
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = next.(Model)
	if m.overlay != overlayNone {
		t.Fatal("picker did not close on enter")
	}
	if m.scaleIdx != idx {
		t.Errorf("scaleIdx = %d, want %d", m.scaleIdx, idx)
	}
}

func TestChordPickerSwitchesKind(t *testing.T) {
	m := sized(NewModel())
	m = press(t, m, "c")
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = next.(Model)
	out := m.render()
	if !strings.Contains(out, "arpeggio") {
		t.Error("after chord pick, title should mention arpeggio")
	}
}

func TestDisplayAndSpellingCycle(t *testing.T) {
	m := sized(NewModel())
	if m.useFlats() {
		t.Error("A should spell with sharps by default")
	}
	m = press(t, m, "b") // auto -> sharps
	m = press(t, m, "b") // sharps -> flats
	if !m.useFlats() {
		t.Error("explicit flats mode should use flats")
	}
	m = press(t, m, "n")
	if m.display != displayNotes {
		t.Errorf("display = %v, want notes", m.display)
	}
}

// Cycling label modes must never reflow the board: the same number of frets
// stays visible and every rendered line keeps its width.
func TestLabelModesKeepBoardWidth(t *testing.T) {
	m := sized(NewModel())
	want := lipgloss.Width(m.renderBoard())
	wantFrets := m.visibleFrets()
	for _, mode := range []displayMode{displayNotes, displayBoth, displayDegrees} {
		m = press(t, m, "n")
		if m.display != mode {
			t.Fatalf("display = %v, want %v", m.display, mode)
		}
		if got := lipgloss.Width(m.renderBoard()); got != want {
			t.Errorf("mode %v: board width %d, want %d", mode, got, want)
		}
		if got := m.visibleFrets(); got != wantFrets {
			t.Errorf("mode %v: visible frets %d, want %d", mode, got, wantFrets)
		}
	}
}

// A highlighted cell always spans the full cell width, however long its
// label, so cursor and root highlights are the same size everywhere.
func TestHighlightFillsCell(t *testing.T) {
	st := lipgloss.NewStyle()
	for _, label := range []string{"1", "b3", "C·1", "F#·b3"} {
		out := placeLabel(label, 1, cellWidth, "─", st, st, true)
		if strings.Contains(out, "─") {
			t.Errorf("padded label %q left string fill in cell: %q", label, out)
		}
		if w := lipgloss.Width(out); w != cellWidth {
			t.Errorf("padded label %q width %d, want %d", label, w, cellWidth)
		}
	}
}

// Labels are anchored to fixed columns so neighbouring cells line up: the
// note letter, the degree digit, and the "both" separator never move as the
// label length changes.
func TestLabelSlotsAlign(t *testing.T) {
	st := lipgloss.NewStyle()
	cell := func(label string, at int) string {
		return placeLabel(label, at, cellWidth, "─", st, st, false)
	}
	cases := []struct {
		label string
		at    int
		want  string
	}{
		{"F", noteCol, "──F──"},
		{"F#", noteCol, "──F#─"},
		{"5", degreeCol, "──5──"},
		{"b3", degreeCol - 1, "─b3──"},
		{"F·5", 1, "─F·5─"},
		{"Eb·4", 0, "Eb·4─"},
		{"F#·b3", 0, "F#·b3"},
		{"A·R", 1, "─A·R─"},
	}
	for _, c := range cases {
		if got := cell(c.label, c.at); got != c.want {
			t.Errorf("place(%q, %d) = %q, want %q", c.label, c.at, got, c.want)
		}
	}
	for _, c := range []struct {
		name, deg string
		at        int
	}{
		{"F", "5", 1}, {"Eb", "4", 0}, {"F#", "b3", 0}, {"A", "R", 1},
	} {
		if _, at := bothSlot(c.name, c.deg); at != c.at {
			t.Errorf("bothSlot(%q,%q) at %d, want %d", c.name, c.deg, at, c.at)
		}
	}
	for _, c := range []struct {
		deg string
		at  int
	}{{"5", 2}, {"b3", 1}, {"R", 2}} {
		if _, at := degreeSlot(c.deg); at != c.at {
			t.Errorf("degreeSlot(%q) at %d, want %d", c.deg, at, c.at)
		}
	}
}
