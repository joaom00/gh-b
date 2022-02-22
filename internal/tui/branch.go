package tui

// import (
// 	"fmt"
// 	"strings"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/joaom00/gh-b/internal/git"
// )

// func (m Model) checkoutBranch() (tea.Model, tea.Cmd) {
// 	branch := strings.TrimSuffix(m.branches[m.cursor].Branch, "*")
// 	out := git.CheckoutBranch(branch)

// 	fmt.Println("\n\n", out)
// 	return m, tea.Quit
// }

// func (m *Model) createBranch() (tea.Model, tea.Cmd) {
// 	switch m.confirmInput.Value() {

// 	case "y":
// 		out := git.CreateBranch(m.createInput.Value())

// 		fmt.Println("\n\n", out)
// 		return m, tea.Quit
// 	case "n":
// 		m.createInput.SetValue("")
// 		m.confirmInput.SetValue("")
// 		m.createInput.Blur()
// 		m.confirmInput.Blur()
// 		m.showConfirmInput = false
// 		m.actionState = idle

// 		return m, nil
// 	default:
// 		m.NewBranch = m.createInput.Value()
// 		m.createInput.Blur()
// 		m.showConfirmInput = true
// 		m.confirmInput.Focus()

// 		return m, nil
// 	}
// }

// func (m *Model) deleteBranch() (tea.Model, tea.Cmd) {
// 	switch m.confirmInput.Value() {

// 	case "n":
// 		m.confirmInput.SetValue("")
// 		m.confirmInput.Blur()
// 		m.showConfirmInput = false
// 		m.actionState = idle

// 		return m, nil

// 	default:
// 		branch := strings.TrimSuffix(m.branches[m.cursor].Branch, "*")
// 		out := git.DeleteBranch(branch)

// 		fmt.Println("\n\n", out)
// 		return m, tea.Quit
// 	}
// }

// // var (
// // 	normalTitle = lipgloss.NewStyle().
// // 			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
// // 			Padding(0, 0, 0, 2)
// // 	normalDesc = normalTitle.Copy().
// // 			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

// // 	selectedTitle = lipgloss.NewStyle().
// // 			Border(lipgloss.NormalBorder(), false, false, false, true).
// // 			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
// // 			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
// // 			Padding(0, 0, 0, 1)

// // 	selectedDesc = selectedTitle.Copy().
// // 			Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})
// // )

// // func (b *Branch) render(isSelected bool) string {
// // 	var title, desc string
// // 	// branch := lipgloss.NewStyle().
// // 	// 	Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
// // 	// 	// Padding(0, 0, 0, 2).
// // 	// 	Render(b.branch)
// // 	// authorname := lipgloss.NewStyle().
// // 	// 	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
// // 	// 	Render(fmt.Sprintf("%s (%s)", b.authorName, b.committerDate))

// // 	// remoteName := lipgloss.NewStyle().
// // 	// 	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
// // 	// 	// Padding(0, 0, 0, 2).
// // 	// 	Render(b.remoteName)

// // 	// track := lipgloss.NewStyle().
// // 	// 	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
// // 	// 	// Padding(0, 0, 0, 2).
// // 	// 	Render(b.track)

// // 	if isSelected {
// // 		title = selectedTitle.Bold(true).Render(b.Branch)
// // 		desc = selectedDesc.Render(fmt.Sprintf("%s (%s)", b.AuthorName, b.CommitterDate))
// // 	} else {
// // 		title = normalTitle.Render(b.Branch)
// // 		desc = normalDesc.Render(fmt.Sprintf("%s (%s)", b.AuthorName, b.CommitterDate))
// // 	}

// // 	return lipgloss.NewStyle().
// // 		MarginBottom(1).
// // 		Render(lipgloss.JoinVertical(lipgloss.Left, title, desc))
// // }
