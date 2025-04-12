package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var buttonStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFF7DB")).
	Background(lipgloss.Color("#888B7E")).
	Padding(0, 3).
	MarginRight(2)

var activeButtonStyle = buttonStyle.
	Foreground(lipgloss.Color("#FFF7DB")).
	Background(lipgloss.Color("#F25D94")).
	Underline(true)

type model struct {
	tabs      []string
	activeTab int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isActive := i == m.activeTab
		if isActive {
			style = activeButtonStyle
		} else {
			style = buttonStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	row := lipgloss.JoinHorizontal(lipgloss.Center, renderedTabs...)
	return row
}

func main() {
	getGitRemotes()
	tabs := []string{"Local"}
	tabs = append(tabs, getGitRemotes()...)
	m := model{tabs: tabs, activeTab: 0}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// p := tea.NewProgram(tui.NewModel())
	//
	// err := p.Start()
	// if err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
}

func getGitRemotes() []string {
	cmd := exec.Command("git", "remote")

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	remotes := strings.Split(strings.TrimSpace(string(out)), "\n")
	return remotes
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
