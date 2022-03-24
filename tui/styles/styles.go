package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title lipgloss.Style

	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	Pagination lipgloss.Style
	Help       lipgloss.Style
	QuitText   lipgloss.Style
}

func DefaultStyles() (s Styles) {
	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	s.NormalTitle = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	s.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	s.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	s.Pagination = list.DefaultStyles().
		PaginationStyle.
		PaddingLeft(4)

	s.Help = list.DefaultStyles().
		HelpStyle.
		PaddingLeft(4).
		PaddingBottom(1)

	s.QuitText = lipgloss.NewStyle().
		Margin(1, 0, 2, 4)

	return s
}
