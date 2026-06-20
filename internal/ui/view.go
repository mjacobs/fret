package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fret/internal/music"
)

const gutterWidth = 3 // string-name column, e.g. " e "

// openWidth is the open-string cell before the nut, sized so the widest
// label of the current display mode fits.
func (m Model) openWidth() int {
	return m.cellWidth() - 2
}

func (m Model) View() tea.View {
	var content string
	if m.overlay != overlayNone {
		content = m.renderOverlay()
	} else {
		content = m.render()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) render() string {
	m = m.clampCursor()
	board := m.renderBoard()
	boardWidth := lipgloss.Width(board)

	sections := []string{
		m.renderHeader(),
		"",
		board,
		m.renderInspector(boardWidth),
		m.renderLegend(boardWidth),
		"",
		m.renderFooter(),
	}
	return appStyle.Render(strings.Join(sections, "\n"))
}

// --- header ---------------------------------------------------------------

var logoColors = []color.Color{cPink, cMauve, cBlue, cTeal}

func (m Model) renderHeader() string {
	var logo strings.Builder
	logo.WriteString(lipgloss.NewStyle().Foreground(cYellow).Render("♪ "))
	for i, r := range "fret" {
		logo.WriteString(lipgloss.NewStyle().Foreground(logoColors[i%len(logoColors)]).Bold(true).Render(string(r)))
	}
	title := logo.String() + "  " + taglineStyle.Render("fretboard cartography for curious fingers")

	chips := []string{
		keyChipStyle.Render(m.setTitle()),
		chipStyle.Foreground(cBlue).Render("✻ " + m.tuning().Name),
		chipStyle.Foreground(cTeal).Render("⌖ " + m.lens().Name),
		chipStyle.Foreground(cYellow).Render("⊙ " + m.display.String()),
		chipStyle.Foreground(cGreen).Render(fmt.Sprintf("%d frets", m.frets)),
		chipStyle.Foreground(cPink).Render(m.spell.String()),
	}
	return title + "\n\n" + strings.Join(chips, " ")
}

// setTitle is the headline state, e.g. "A minor pentatonic" or "C Major 7 arpeggio".
func (m Model) setTitle() string {
	set := m.set()
	name := m.root().Name(m.useFlats()) + " " + strings.ToLower(set.Name)
	if set.Kind == music.KindArpeggio {
		name += " arpeggio"
	}
	return name
}

// --- fretboard ------------------------------------------------------------

func inlayMarker(fret int) string {
	switch fret {
	case 12, 24:
		return "∙∙"
	case 3, 5, 7, 9, 15, 17, 19, 21:
		return "∙"
	default:
		return ""
	}
}

func (m Model) renderBoard() string {
	root, set, tun, lens := m.root(), m.set(), m.tuning(), m.lens()
	flats := m.useFlats()
	visible := m.visibleFrets()
	cw := m.cellWidth()
	ow := m.openWidth()
	grid := music.Fretboard(root, set, tun, visible)

	var b strings.Builder

	// Fret numbers.
	b.WriteString(strings.Repeat(" ", gutterWidth))
	b.WriteString(centerText("0", ow, fretNumStyle))
	b.WriteString(" ")
	for f := 1; f <= visible; f++ {
		st := fretNumStyle
		if inlayMarker(f) != "" {
			st = fretNumDotStyle
		}
		b.WriteString(centerText(strconv.Itoa(f), cw, st))
		b.WriteString(" ")
	}
	if visible < m.frets {
		b.WriteString(moreStyle.Render(fmt.Sprintf("⋯+%d", m.frets-visible)))
	}
	b.WriteString("\n")

	// Strings.
	for s := range grid {
		b.WriteString(stringLabelStyle.Render(fmt.Sprintf("%2s ", tun.Label(s))))
		b.WriteString(m.renderCell(grid[s][0], lens, flats, ow, " "))
		b.WriteString(nutStyle.Render("┃"))
		fill := "─"
		if tun.Wound(s) {
			fill = "═"
		}
		for f := 1; f <= visible; f++ {
			b.WriteString(m.renderCell(grid[s][f], lens, flats, cw, fill))
			b.WriteString(wireStyle.Render("│"))
		}
		b.WriteString("\n")
	}

	// Inlay markers.
	b.WriteString(strings.Repeat(" ", gutterWidth+ow+1))
	for f := 1; f <= visible; f++ {
		b.WriteString(centerText(inlayMarker(f), cw, markerStyle))
		b.WriteString(" ")
	}

	return boardStyle.Render(b.String())
}

func (m Model) renderCell(pos music.Position, lens music.Lens, flats bool, w int, fill string) string {
	isCursor := pos.String == m.curStr && pos.Fret == m.curFret
	isPinned := m.pins[[2]int{pos.String, pos.Fret}]

	var label string
	switch {
	case pos.InSet:
		switch m.display {
		case displayNotes:
			label = pos.Note.Name(flats)
		case displayBoth:
			label = pos.Note.Name(flats) + "·" + pos.Degree
		default:
			label = pos.Degree
		}
	case isCursor || isPinned:
		label = pos.Note.Name(flats)
	}

	if label == "" {
		return stringStyle.Render(strings.Repeat(fill, w))
	}

	inLens := lens.Includes(pos.Fret)
	var st lipgloss.Style
	pad := false
	switch {
	case isCursor && isPinned:
		st = cursorPinnedStyle
		pad = true
	case isCursor:
		st = cursorStyle
		pad = true
	case isPinned:
		st = pinnedStyle
		pad = true
	case !pos.InSet:
		st = dimNoteStyle
	case pos.IsRoot && inLens:
		st = rootStyle
		pad = true
	case pos.IsRoot:
		st = lipgloss.NewStyle().Foreground(cPeach).Bold(true)
	case inLens:
		st = lipgloss.NewStyle().Foreground(intervalColor(pos.Interval)).Bold(true)
	default:
		st = dimNoteStyle
	}

	// Echo: every cell sharing the cursor's pitch class glows a little, so
	// one hovered note reveals all its octave twins across the neck.
	if !isCursor && pos.InSet && pos.Note == m.cursorNote() {
		st = st.Underline(true)
	}

	return centerLabel(label, w, fill, st, stringStyle, pad)
}

// centerLabel centers a styled label inside a cell of width w, filling the
// remainder with the string-line rune.
func centerLabel(label string, w int, fill string, st, lineSt lipgloss.Style, pad bool) string {
	r := []rune(label)
	if pad && len(r)+2 <= w {
		label = " " + label + " "
		r = []rune(label)
	}
	if len(r) > w {
		r = r[:w]
		label = string(r)
	}
	left := (w - len(r)) / 2
	right := w - len(r) - left
	return lineSt.Render(strings.Repeat(fill, left)) + st.Render(label) + lineSt.Render(strings.Repeat(fill, right))
}

func centerText(s string, w int, st lipgloss.Style) string {
	r := []rune(s)
	if len(r) > w {
		r = r[:w]
		s = string(r)
	}
	left := (w - len(r)) / 2
	right := w - len(r) - left
	return strings.Repeat(" ", left) + st.Render(s) + strings.Repeat(" ", right)
}

// --- inspector ------------------------------------------------------------

func (m Model) renderInspector(width int) string {
	root, set, tun, lens := m.root(), m.set(), m.tuning(), m.lens()
	flats := m.useFlats()
	note := m.cursorNote()
	degree, interval, in := set.DegreeOf(root, note)

	where := fmt.Sprintf("fret %d", m.curFret)
	if m.curFret == 0 {
		where = "open"
	}

	var sound string
	noteName := lipgloss.NewStyle().Foreground(intervalColor(interval)).Bold(true).Render(note.Name(flats))
	switch {
	case in && interval == 0:
		sound = noteName + dimStyle.Render(" — R, the root: home base")
	case in:
		sound = noteName + dimStyle.Render(fmt.Sprintf(" — %s (%s)", degree, music.IntervalName(interval)))
	default:
		sound = dimNoteStyle.Render(note.Name(flats)) +
			dimStyle.Render(fmt.Sprintf(" — visitor, outside this %s", set.Kind))
	}

	var focus string
	if lens.Includes(m.curFret) {
		focus = subtleStyle.Render("in " + lens.Name)
	} else {
		focus = dimStyle.Render("outside " + lens.Name)
	}

	parts := []string{
		subtleStyle.Render(fmt.Sprintf("%s string", strings.TrimSpace(tun.Label(m.curStr)))),
		subtleStyle.Render(where),
		sound,
		focus,
	}
	if m.pins[[2]int{m.curStr, m.curFret}] {
		parts = append(parts, pinnedStyle.Render(" ◉ pinned "))
	}

	sep := dimStyle.Render("  ·  ")
	return inspectorStyle.Width(width).Render("▸ " + strings.Join(parts, sep))
}

// --- legend ---------------------------------------------------------------

func (m Model) renderLegend(width int) string {
	root, set, lens := m.root(), m.set(), m.lens()
	flats := m.useFlats()

	chips := make([]string, len(set.Formula))
	for i, iv := range set.Formula {
		st := lipgloss.NewStyle().Foreground(intervalColor(iv)).Bold(true)
		chips[i] = st.Render(set.Degrees[i] + "·" + root.Add(iv).Name(flats))
	}

	lines := []string{
		subtleStyle.Bold(true).Render(m.setTitle()) + dimStyle.Render("  —  ") + vibeStyle.Render(set.Vibe),
		strings.Join(chips, dimStyle.Render("   ")),
		dimStyle.Render(fmt.Sprintf("lens: %s (frets %d–%d) — %s", lens.Name, lens.Min, lens.Max, lens.Blurb)),
	}
	return legendStyle.Width(width).Render(strings.Join(lines, "\n"))
}

// --- footer ---------------------------------------------------------------

func key(k, desc string) string {
	return helpKeyStyle.Render(k) + " " + helpDescStyle.Render(desc)
}

func (m Model) renderFooter() string {
	sep := helpDescStyle.Render(" · ")
	line1 := strings.Join([]string{
		key("←↓↑→", "move"), key("g/G", "ends"), key("space", "pin"), key("x", "unpin all"),
		key("tab", "labels"), key("[ ]", "frets"),
	}, sep)
	line2 := strings.Join([]string{
		key("r", "root"), key("s", "scales"), key("c", "chords"), key("t", "tunings"),
		key("p", "lens"), key("b", "spelling"), key("?", "help"), key("q", "quit"),
	}, sep)
	return line1 + "\n" + line2
}

// --- overlays ---------------------------------------------------------------

type pickerItem struct {
	name string
	desc string
}

func (m Model) pickerTitle() string {
	switch m.overlay {
	case overlayRoot:
		return "♪ Pick a key"
	case overlayScale:
		return "♪ Pick a scale"
	case overlayChord:
		return "♪ Pick a chord to arpeggiate"
	case overlayTuning:
		return "♪ Pick a tuning"
	case overlayLens:
		return "♪ Pick a practice lens"
	default:
		return ""
	}
}

func (m Model) pickerItems() []pickerItem {
	flats := m.useFlats()
	switch m.overlay {
	case overlayRoot:
		items := make([]pickerItem, 0, 12)
		set := m.set()
		for _, n := range music.Roots() {
			spellFlat := n.PrefersFlats()
			switch m.spell {
			case spellSharps:
				spellFlat = false
			case spellFlats:
				spellFlat = true
			}
			notes := make([]string, len(set.Formula))
			for i, iv := range set.Formula {
				notes[i] = n.Add(iv).Name(spellFlat)
			}
			items = append(items, pickerItem{n.BothNames(), strings.Join(notes, " ")})
		}
		return items
	case overlayScale:
		items := make([]pickerItem, len(music.Scales))
		for i, s := range music.Scales {
			items[i] = pickerItem{s.Name, s.Vibe}
		}
		return items
	case overlayChord:
		items := make([]pickerItem, len(music.Arpeggios))
		for i, a := range music.Arpeggios {
			items[i] = pickerItem{a.Name, a.Vibe}
		}
		return items
	case overlayTuning:
		items := make([]pickerItem, len(music.Tunings))
		for i, t := range music.Tunings {
			items[i] = pickerItem{t.Name + "  (" + t.LowToHigh(flats) + ")", t.Desc}
		}
		return items
	case overlayLens:
		items := make([]pickerItem, len(music.Lenses))
		for i, l := range music.Lenses {
			items[i] = pickerItem{fmt.Sprintf("%s  (frets %d–%d)", l.Name, l.Min, l.Max), l.Blurb}
		}
		return items
	default:
		return nil
	}
}

func (m Model) renderOverlay() string {
	var panel string
	if m.overlay == overlayHelp {
		panel = m.renderHelpPanel()
	} else {
		panel = m.renderPickerPanel()
	}
	w, h := m.width, m.height
	if w == 0 {
		w = 100
	}
	if h == 0 {
		h = 32
	}
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, panel)
}

func (m Model) renderPickerPanel() string {
	items := m.pickerItems()
	nameW := 0
	for _, it := range items {
		if n := len([]rune(it.name)); n > nameW {
			nameW = n
		}
	}

	var b strings.Builder
	b.WriteString(overlayTitleStyle.Render(m.pickerTitle()))
	b.WriteString("\n\n")
	for i, it := range items {
		name := it.name + strings.Repeat(" ", nameW-len([]rune(it.name)))
		if i == m.pickerIdx {
			b.WriteString(pickerSelStyle.Render(" ▸ " + name + "  "))
			b.WriteString("  " + pickerDescStyle.Italic(true).Render(it.desc))
		} else {
			b.WriteString(pickerItemStyle.Render("   " + name + "  "))
			b.WriteString("  " + pickerDescStyle.Render(it.desc))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(overlayHintStyle.Render("↑↓ choose · enter select · esc cancel"))
	return overlayStyle.Render(b.String())
}

func (m Model) renderHelpPanel() string {
	rows := [][2]string{
		{"←↓↑→ / hjkl", "move the cursor around the neck"},
		{"g / G", "jump to the nut / the last fret"},
		{"space / enter", "pin the hovered position"},
		{"x", "clear every pin"},
		{"r", "pick a key (root note)"},
		{", / .", "transpose down / up a half step"},
		{"s", "pick a scale"},
		{"c (or a)", "pick a chord to see as an arpeggio"},
		{"t", "pick a tuning"},
		{"p", "pick a practice lens"},
		{"tab / n", "cycle labels: degrees → notes → both"},
		{"b", "cycle spelling: auto → sharps → flats"},
		{"[ / ]", "shrink / grow the neck"},
		{"0", "reset to 12 frets"},
		{"q / ctrl+c", "quit"},
	}

	var b strings.Builder
	b.WriteString(overlayTitleStyle.Render("♪ fret — every key"))
	b.WriteString("\n\n")
	for _, r := range rows {
		b.WriteString(helpKeyStyle.Render(fmt.Sprintf("  %-14s", r[0])))
		b.WriteString(helpDescStyle.Render(r[1]))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(overlayHintStyle.Render("press any key to close"))
	return overlayStyle.Render(b.String())
}
