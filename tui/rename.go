package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/git"
)

type renameModel struct {
	input            textinput.Model
	confirmInput     textinput.Model
	showConfirmInput bool
	help             help.Model
}

func newRenameModel() *renameModel {
	ti := textinput.New()
	ci := textinput.New()
	ci.CharLimit = 1

	return &renameModel{
		input:            ti,
		confirmInput:     ci,
		showConfirmInput: false,
		help:             help.New(),
	}
}

func renameUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			if i, ok := m.selectedItem(); ok {
				out := git.RenameBranch(i.Name, m.rename.input.Value())

				m.state = browsing
				m.updateListItem()
				m.list.NewStatusMessage(out)

				return m, nil
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.rename.input.Reset()
			m.state = browsing
			m.updateKeybindins()
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.rename.input, cmd = m.rename.input.Update(msg)
	cmds = append(cmds, cmd)

	m.rename.confirmInput, cmd = m.rename.confirmInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) renameView() string {
	m.rename.input.Placeholder = strings.TrimSuffix(m.list.SelectedItem().(item).Name, "*")

	title := m.styles.Title.MarginLeft(2).Render("Rename Branch")
	textInput := lipgloss.NewStyle().MarginLeft(4).Render(m.rename.input.View())
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.create.help.View(m.keyMap))

	if m.rename.showConfirmInput {
		confirmInput := lipgloss.NewStyle().
			MarginLeft(4).
			Render(lipgloss.JoinHorizontal(lipgloss.Left, "You would like rename remote branch? [y/N]", m.rename.confirmInput.View()))

		return lipgloss.NewStyle().MarginTop(1).Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", confirmInput, "\n", help))
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", help))
}
