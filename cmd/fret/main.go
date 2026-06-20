package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"fret/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fret: %v\n", err)
		os.Exit(1)
	}
}
