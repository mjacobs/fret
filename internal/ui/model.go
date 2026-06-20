package ui

import (
	tea "charm.land/bubbletea/v2"
	"fret/internal/music"
)

type displayMode int

const (
	displayDegrees displayMode = iota
	displayNotes
	displayBoth
)

func (d displayMode) String() string {
	switch d {
	case displayNotes:
		return "notes"
	case displayBoth:
		return "both"
	default:
		return "degrees"
	}
}

type spellMode int

const (
	spellAuto spellMode = iota
	spellSharps
	spellFlats
)

func (s spellMode) String() string {
	switch s {
	case spellSharps:
		return "♯ sharps"
	case spellFlats:
		return "♭ flats"
	default:
		return "♮ auto"
	}
}

type overlayKind int

const (
	overlayNone overlayKind = iota
	overlayRoot
	overlayScale
	overlayChord
	overlayTuning
	overlayLens
	overlayHelp
)

type Model struct {
	width  int
	height int

	rootIdx   int
	kind      music.Kind
	scaleIdx  int
	arpIdx    int
	tuningIdx int
	lensIdx   int
	frets     int
	display   displayMode
	spell     spellMode

	curStr  int
	curFret int
	pins    map[[2]int]bool

	overlay   overlayKind
	pickerIdx int
}

// NewModel starts in A minor pentatonic — the people's scale — on a
// standard-tuned neck.
func NewModel() Model {
	return Model{
		rootIdx:  int(music.A),
		kind:     music.KindScale,
		scaleIdx: 5, // minor pentatonic
		frets:    12,
		curStr:   5, // low E string
		pins:     make(map[[2]int]bool),
	}
}

func (m Model) root() music.Note     { return music.Roots()[m.rootIdx] }
func (m Model) tuning() music.Tuning { return music.Tunings[m.tuningIdx] }
func (m Model) lens() music.Lens     { return music.Lenses[m.lensIdx] }

func (m Model) set() music.NoteSet {
	if m.kind == music.KindArpeggio {
		return music.Arpeggios[m.arpIdx]
	}
	return music.Scales[m.scaleIdx]
}

func (m Model) useFlats() bool {
	switch m.spell {
	case spellSharps:
		return false
	case spellFlats:
		return true
	default:
		return m.root().PrefersFlats()
	}
}

func (m Model) cursorNote() music.Note {
	return m.tuning().Open[m.curStr].Add(m.curFret)
}

// cellWidth is the column width of one fret cell, sized for the label mode.
func (m Model) cellWidth() int {
	if m.display == displayBoth {
		return 7
	}
	return 5
}

// visibleFrets caps the rendered neck at what fits the terminal width.
func (m Model) visibleFrets() int {
	if m.width == 0 {
		return m.frets
	}
	// app padding + board border/padding + gutter + open column + nut + hint.
	avail := m.width - 19 - m.openWidth()
	n := avail / (m.cellWidth() + 1)
	if n < 4 {
		n = 4
	}
	if n > m.frets {
		n = m.frets
	}
	return n
}

func (m Model) clampCursor() Model {
	if m.curStr < 0 {
		m.curStr = 0
	}
	if max := len(m.tuning().Open) - 1; m.curStr > max {
		m.curStr = max
	}
	if m.curFret < 0 {
		m.curFret = 0
	}
	if max := m.visibleFrets(); m.curFret > max {
		m.curFret = max
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.clampCursor()
	case tea.KeyPressMsg:
		if m.overlay != overlayNone {
			return m.updateOverlay(msg)
		}
		return m.updateMain(msg)
	}
	return m, nil
}

func (m Model) updateMain(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		m.curStr--
	case "down", "j":
		m.curStr++
	case "left", "h":
		m.curFret--
	case "right", "l":
		m.curFret++
	case "g", "home":
		m.curFret = 0
	case "G", "end":
		m.curFret = m.visibleFrets()
	case "enter", " ", "space":
		key := [2]int{m.curStr, m.curFret}
		if m.pins[key] {
			delete(m.pins, key)
		} else {
			m.pins[key] = true
		}
	case "x":
		m.pins = make(map[[2]int]bool)
	case ",", "<":
		m.rootIdx = (m.rootIdx + 11) % 12
	case ".", ">":
		m.rootIdx = (m.rootIdx + 1) % 12
	case "r":
		m = m.openOverlay(overlayRoot, m.rootIdx)
	case "s":
		m = m.openOverlay(overlayScale, m.scaleIdx)
	case "c", "a":
		m = m.openOverlay(overlayChord, m.arpIdx)
	case "t":
		m = m.openOverlay(overlayTuning, m.tuningIdx)
	case "p":
		m = m.openOverlay(overlayLens, m.lensIdx)
	case "?":
		m = m.openOverlay(overlayHelp, 0)
	case "tab", "n":
		m.display = (m.display + 1) % 3
	case "b":
		m.spell = (m.spell + 1) % 3
	case "[", "-":
		if m.frets > 4 {
			m.frets--
		}
	case "]", "=":
		if m.frets < 24 {
			m.frets++
		}
	case "0":
		m.frets = 12
	}
	m = m.clampCursor()
	return m, nil
}

func (m Model) openOverlay(kind overlayKind, current int) Model {
	m.overlay = kind
	m.pickerIdx = current
	return m
}

func (m Model) updateOverlay(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.overlay == overlayHelp {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.overlay = overlayNone
		}
		return m, nil
	}

	count := len(m.pickerItems())
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q":
		m.overlay = overlayNone
	case "up", "k":
		m.pickerIdx = (m.pickerIdx + count - 1) % count
	case "down", "j":
		m.pickerIdx = (m.pickerIdx + 1) % count
	case "g", "home":
		m.pickerIdx = 0
	case "G", "end":
		m.pickerIdx = count - 1
	case "enter", " ", "space":
		m = m.applyPick()
		m.overlay = overlayNone
	}
	m = m.clampCursor()
	return m, nil
}

func (m Model) applyPick() Model {
	switch m.overlay {
	case overlayRoot:
		m.rootIdx = m.pickerIdx
	case overlayScale:
		m.scaleIdx = m.pickerIdx
		m.kind = music.KindScale
	case overlayChord:
		m.arpIdx = m.pickerIdx
		m.kind = music.KindArpeggio
	case overlayTuning:
		if m.tuningIdx != m.pickerIdx {
			m.tuningIdx = m.pickerIdx
			m.pins = make(map[[2]int]bool)
		}
	case overlayLens:
		m.lensIdx = m.pickerIdx
	}
	return m
}
