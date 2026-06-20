package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Catppuccin Mocha palette.
var (
	cBase     = lipgloss.Color("#1e1e2e")
	cSurface0 = lipgloss.Color("#313244")
	cSurface1 = lipgloss.Color("#45475a")
	cSurface2 = lipgloss.Color("#585b70")
	cOverlay0 = lipgloss.Color("#6c7086")
	cOverlay1 = lipgloss.Color("#7f849c")
	cSubtext  = lipgloss.Color("#a6adc8")
	cText     = lipgloss.Color("#cdd6f4")

	cRosewater = lipgloss.Color("#f5e0dc")
	cPink      = lipgloss.Color("#f5c2e7")
	cMauve     = lipgloss.Color("#cba6f7")
	cRed       = lipgloss.Color("#f38ba8")
	cMaroon    = lipgloss.Color("#eba0ac")
	cPeach     = lipgloss.Color("#fab387")
	cYellow    = lipgloss.Color("#f9e2af")
	cGreen     = lipgloss.Color("#a6e3a1")
	cTeal      = lipgloss.Color("#94e2d5")
	cSky       = lipgloss.Color("#89dceb")
	cBlue      = lipgloss.Color("#89b4fa")
	cLavender  = lipgloss.Color("#b4befe")
)

// intervalColor gives each scale function its own hue so the eye can learn
// "pink means a third" before the brain catches up.
func intervalColor(semitones int) color.Color {
	switch ((semitones % 12) + 12) % 12 {
	case 0:
		return cPeach // root
	case 3, 4:
		return cPink // thirds
	case 7:
		return cBlue // fifth
	case 10, 11:
		return cMauve // sevenths
	case 2:
		return cTeal // second / ninth
	case 5:
		return cSky // fourth / eleventh
	case 9:
		return cGreen // sixth / thirteenth
	case 6:
		return cRed // the blue note
	default:
		return cMaroon // b2, b6
	}
}

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	taglineStyle = lipgloss.NewStyle().Foreground(cOverlay1).Italic(true)

	// Header chips.
	keyChipStyle = lipgloss.NewStyle().
			Foreground(cBase).
			Background(cPeach).
			Bold(true).
			Padding(0, 1)

	chipStyle = lipgloss.NewStyle().
			Background(cSurface0).
			Padding(0, 1)

	// Fretboard furniture.
	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cSurface1).
			Padding(1, 2)

	stringLabelStyle = lipgloss.NewStyle().Foreground(cSubtext).Bold(true)
	stringStyle      = lipgloss.NewStyle().Foreground(cSurface2)
	wireStyle        = lipgloss.NewStyle().Foreground(cSurface1)
	nutStyle         = lipgloss.NewStyle().Foreground(cRosewater).Bold(true)
	markerStyle      = lipgloss.NewStyle().Foreground(cSurface2)
	fretNumStyle     = lipgloss.NewStyle().Foreground(cOverlay0)
	fretNumDotStyle  = lipgloss.NewStyle().Foreground(cOverlay1).Bold(true)
	moreStyle        = lipgloss.NewStyle().Foreground(cOverlay0).Italic(true)

	// Note cells.
	rootStyle = lipgloss.NewStyle().
			Foreground(cBase).
			Background(cPeach).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(cBase).
			Background(cBlue).
			Bold(true)

	pinnedStyle = lipgloss.NewStyle().
			Foreground(cBase).
			Background(cPink).
			Bold(true)

	cursorPinnedStyle = lipgloss.NewStyle().
				Foreground(cBase).
				Background(cLavender).
				Bold(true)

	dimNoteStyle = lipgloss.NewStyle().Foreground(cOverlay0)

	// Panels under the board.
	inspectorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cSurface1).
			Padding(0, 2)

	legendStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cSurface1).
			Padding(0, 2)

	dimStyle    = lipgloss.NewStyle().Foreground(cOverlay0)
	subtleStyle = lipgloss.NewStyle().Foreground(cSubtext)
	vibeStyle   = lipgloss.NewStyle().Foreground(cOverlay1).Italic(true)

	// Footer help.
	helpKeyStyle  = lipgloss.NewStyle().Foreground(cLavender)
	helpDescStyle = lipgloss.NewStyle().Foreground(cOverlay0)

	// Overlays.
	overlayStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(cMauve).
			Padding(1, 3)

	overlayTitleStyle = lipgloss.NewStyle().Foreground(cMauve).Bold(true)

	pickerSelStyle = lipgloss.NewStyle().
			Foreground(cText).
			Background(cSurface0).
			Bold(true)

	pickerItemStyle  = lipgloss.NewStyle().Foreground(cSubtext)
	pickerDescStyle  = lipgloss.NewStyle().Foreground(cOverlay0)
	overlayHintStyle = lipgloss.NewStyle().Foreground(cOverlay0).Italic(true)
)
