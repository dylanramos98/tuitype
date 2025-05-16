package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dylanramos/tuitype/internal/text"
	"github.com/dylanramos/tuitype/internal/ui"
)

func main() {
	m := ui.NewModel(text.GetWords())
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
