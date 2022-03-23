package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type renameModel struct {
	input textinput.Model
}

func newRenameModel() *renameModel {
	ti := textinput.New()

	return &renameModel{
		input: ti,
	}
}

func renameUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			return m, nil
		}
	}

	return m, nil
}

func (m Model) renameView() string {
	return lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(4).
		Render(lipgloss.JoinVertical(lipgloss.Left, m.rename.input.View()))
}
