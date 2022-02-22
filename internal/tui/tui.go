package tui

import (
	"fmt"
	"io"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/internal/git"
	"github.com/joaom00/gh-b/internal/tui/keys"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct {
	Name          string
	AuthorName    string
	CommitterDate string
	Track         string
	RemoteName    string
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(string(i)))
}

type state int

const (
	browsing state = iota
	creating
	deleting
)

type Model struct {
	items       []item
	createInput textinput.Model
	create      *createModel
	delete      *deleteModel
	keyMap      *keys.KeyMap
	list        list.Model
	state       state
}

func NewModel() Model {
	branches, err := git.GetAllBranches()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}
	for _, b := range branches {
		items = append(items, item(b.Name))
	}

	const defaultWidth = 20
	const listHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Your Branches"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return Model{
		state:  browsing,
		list:   l,
		create: newCreateModel(),
		delete: newDeleteModel(),
		keyMap: keys.NewKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		}
	}
	switch m.state {
	case browsing:
		return listUpdate(msg, m)

	case creating:
		return createUpdate(msg, m)

	case deleting:
		return deleteUpdate(msg, m)

	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.state {
	case browsing:
		return "\n" + m.list.View()

	case creating:
		return m.createView()

	case deleting:
		return m.deleteView()

	default:
		return ""
	}
}

func listUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Create):
			m.state = creating
			m.keyMap.State = "creating"
			m.create.textinput.Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Delete):
			m.state = deleting
			m.keyMap.State = "deleting"
			m.delete.confirmInput.Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Enter):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				out := git.CheckoutBranch(string(i))

				fmt.Println("\n", out)
			}
			return m, tea.Quit
		}
	}

	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateKeybindins() {
	switch m.state {
	case creating, deleting:
		m.keyMap.Quit.SetEnabled(false)
		// m.keyMap.CheckOut.SetEnabled(true)
		// m.keyMap.ClearAction.SetEnabled(true)
		// m.keyMap.ClearAction.Keys()
		// m.keyMap.CursorUp.SetEnabled(false)
		// m.keyMap.CursorDown.SetEnabled(false)
		// m.keyMap.Delete.SetEnabled(false)
		// m.keyMap.Track.SetEnabled(false)
		// m.keyMap.Merge.SetEnabled(false)
		// m.keyMap.Rebase.SetEnabled(false)
		// m.keyMap.Help.SetEnabled(true)
	case browsing:
		m.keyMap.Cancel.SetEnabled(false)
		// default:
		// 	hasItems := len(m.branches) != 0

		// 	m.keyMap.CursorUp.SetEnabled(hasItems)
		// 	m.keyMap.CursorDown.SetEnabled(hasItems)
		// 	m.keyMap.CheckOut.SetEnabled(hasItems)
		// 	m.keyMap.Create.SetEnabled(hasItems)
		// 	m.keyMap.Delete.SetEnabled(hasItems)
		// 	m.keyMap.Help.SetEnabled(true)
	}
}
