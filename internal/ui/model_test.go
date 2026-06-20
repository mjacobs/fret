package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
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
