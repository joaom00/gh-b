package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joaom00/gh-b/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel())

	err := p.Start()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
