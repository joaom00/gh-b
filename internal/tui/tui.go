package tui

import (
	"fmt"
	"io"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joaom00/gh-b/internal/git"
	"github.com/joaom00/gh-b/internal/tui/keys"
	"github.com/joaom00/gh-b/internal/tui/styles"
)

type item struct {
	Name          string
	AuthorName    string
	CommitterDate string
	Track         string
	RemoteName    string
}

func (i item) FilterValue() string { return i.Name }

type itemDelegate struct {
	style *styles.Styles
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title := d.style.NormalTitle.Render
	desc := d.style.NormalDesc.Render

	if index == m.Index() {
		title = func(s string) string {
			return d.style.SelectedTitle.Render("> " + s)
		}
		desc = func(s string) string {
			return d.style.SelectedDesc.Render(s)
		}
	}

	branch := title(i.Name)
	author := desc(i.AuthorName)
	committerDate := desc(fmt.Sprintf("(%s)", i.CommitterDate))

	itemListStyle := lipgloss.NewStyle().
		Render(fmt.Sprintf("%s %s %s", branch, author, committerDate))

	fmt.Fprint(w, itemListStyle)
}

type state int

const (
	browsing state = iota
	creating
	deleting
)

type Model struct {
	items  []item
	create *createModel
	delete *deleteModel
	keyMap *keys.KeyMap
	list   list.Model
	style  styles.Styles
	state  state
}

func NewModel() Model {
	branches, err := git.GetAllBranches()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}
	for _, b := range branches {
		items = append(items, item{
			Name:          b.Name,
			AuthorName:    b.AuthorName,
			CommitterDate: b.CommitterDate,
			Track:         b.Track,
			RemoteName:    b.RemoteName,
		})
	}

	const defaultWidth = 20
	const listHeight = 20

	s := styles.DefaultStyles()

	l := list.New(items, itemDelegate{style: &s}, defaultWidth, listHeight)
	l.Title = "Your Branches"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = s.Pagination
	l.Styles.HelpStyle = s.Help

	return Model{
		create: newCreateModel(),
		delete: newDeleteModel(),
		keyMap: keys.NewKeyMap(),
		list:   l,
		style:  s,
		state:  browsing,
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

		case key.Matches(msg, m.keyMap.Track):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				out := git.TrackBranch(i.Name)

				fmt.Println("\n", out)
			}
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Enter):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				out := git.CheckoutBranch(i.Name)

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
