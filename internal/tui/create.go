package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/internal/git"
)

type createModel struct {
	help             help.Model
	inputs           []textinput.Model
	focusIndex       int
	showConfirmInput bool
}

func newCreateModel() *createModel {
	ti := textinput.New()
	ti.Placeholder = "Try feature..."

	ci := textinput.New()
	ci.CharLimit = 1

	return &createModel{
		help:             help.New(),
		inputs:           []textinput.Model{ti, ci},
		showConfirmInput: false,
	}
}

func (m *createModel) prevFocus() tea.Cmd {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex--

	if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs) - 1
	}

	return m.inputs[m.focusIndex].Focus()
}

func (m *createModel) nextFocus() tea.Cmd {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex++

	if m.focusIndex > len(m.inputs)-1 {
		m.focusIndex = 0
	}

	return m.inputs[m.focusIndex].Focus()
}

func createUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab", "up"))):
			cmd := m.create.prevFocus()

			return m, cmd

		case key.Matches(msg, key.NewBinding(key.WithKeys("tab", "down"))):
			cmd := m.create.nextFocus()

			return m, cmd

		case key.Matches(msg, m.keyMap.Enter):
			if !m.create.showConfirmInput {
				m.create.inputs[0].Blur()
				m.create.showConfirmInput = true
				m.create.inputs[1].Focus()

				return m, nil
			}

			switch m.create.inputs[1].Value() {
			case "y", "Y", "":
				out := git.CreateBranch(m.create.inputs[0].Value())

				fmt.Println(m.styles.NormalTitle.Render(out))

				return m, tea.Quit

			case "n", "N":
				m.create.inputs[0].Reset()
				m.create.inputs[1].Reset()
				m.create.showConfirmInput = false
				m.state = browsing
				m.updateKeybindins()

				return m, nil

			default:
				m.create.inputs[1].SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.create.inputs[0].Reset()
			m.create.inputs[1].Reset()
			m.create.showConfirmInput = false
			m.state = browsing
			m.updateKeybindins()
		}

	case tea.WindowSizeMsg:
		m.create.help.Width = msg.Width
	}

	cmds := make([]tea.Cmd, len(m.create.inputs))

	for i := range m.create.inputs {
		m.create.inputs[i], cmds[i] = m.create.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) createView() string {
	title := m.styles.Title.MarginLeft(2).Render("Type name of the new branch")
	textInput := lipgloss.NewStyle().MarginLeft(4).Render(m.create.inputs[0].View())
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.create.help.View(m.keyMap))

	if m.create.showConfirmInput {
		branch := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(m.create.inputs[0].Value())

		confirmInput := lipgloss.NewStyle().
			MarginLeft(4).
			Render(lipgloss.JoinHorizontal(
				lipgloss.Left,
				fmt.Sprintf("Create new branch \"%s\"? [Y/n]", branch),
				m.create.inputs[1].View(),
			))

		return lipgloss.NewStyle().
			MarginTop(1).
			Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", confirmInput, "\n", help))
	}

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", textInput, "\n", help))
}
