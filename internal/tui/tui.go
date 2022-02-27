package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joaom00/gh-b/internal/git"
	"github.com/joaom00/gh-b/internal/tui/keys"
	"github.com/joaom00/gh-b/internal/tui/styles"
)

const (
	defaultWidth = 20
	listHeight   = 15
)

type item git.Branch

func (i item) FilterValue() string { return i.Name }

type state int

const (
	browsing state = iota
	creating
	deleting
	merge
	rebasing
)

type Model struct {
	items  []item
	create *createModel
	delete *deleteModel
	merge  *mergeModel
	rebase *rebaseModel
	keyMap *keys.KeyMap
	list   list.Model
	styles styles.Styles
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
		})
	}

	styles := styles.DefaultStyles()
	keys := keys.NewKeyMap()

	l := list.New(items, newItemDelegate(keys, &styles), defaultWidth, listHeight)
	l.Title = "Your Branches"
	l.SetShowStatusBar(false)
	l.Styles.PaginationStyle = styles.Pagination
	l.Styles.HelpStyle = styles.Help

	return Model{
		create: newCreateModel(),
		delete: newDeleteModel(),
		merge:  newMergeModel(),
		rebase: newRebaseModel(),
		keyMap: keys,
		list:   l,
		styles: styles,
		state:  browsing,
	}
}

func (m *Model) updateListItem() {
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
		})
	}

	m.list.SetItems(items)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.list.SettingFilter() {
		m.keyMap.Enter.SetEnabled(false)
	}

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

	case merge:
		return mergeUpdate(msg, m)

	case rebasing:
		return rebaseUpdate(msg, m)

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

	case merge:
		return m.mergeView()

	case rebasing:
		return m.rebaseView()

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
		case key.Matches(msg, m.list.KeyMap.AcceptWhileFiltering):
			m.state = browsing
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.CursorUp):
			m.list.CursorUp()

		case key.Matches(msg, m.keyMap.CursorDown):
			m.list.CursorDown()

		case key.Matches(msg, m.keyMap.Create):
			m.state = creating
			m.keyMap.State = "creating"
			m.create.inputs[0].Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Delete):
			m.state = deleting
			m.keyMap.State = "deleting"
			m.delete.confirmInput.Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Track):
			if i, ok := m.list.SelectedItem().(item); ok {
				i.Name = strings.TrimSuffix(i.Name, "*")
				out := git.TrackBranch(i.Name)

				fmt.Println(m.styles.NormalTitle.Render(out))

				return m, tea.Quit
			}

		case key.Matches(msg, m.keyMap.Merge):
			m.state = merge
			m.keyMap.State = "merge"
			m.merge.confirmInput.Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Rebase):
			m.state = rebasing
			m.keyMap.State = "rebasing"
			m.rebase.confirmInput.Focus()
			m.updateKeybindins()

		case key.Matches(msg, m.keyMap.Enter):
			if i, ok := m.list.SelectedItem().(item); ok {
				i.Name = strings.TrimSuffix(i.Name, "*")
				out := git.CheckoutBranch(i.Name)

				fmt.Println(m.styles.NormalTitle.Copy().MarginTop(1).Render(out))

				return m, tea.Quit
			}

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
	if m.list.SettingFilter() {
		m.keyMap.Enter.SetEnabled(false)
	}

	switch m.state {
	case creating, deleting, merge, rebasing:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Cancel.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.Delete.SetEnabled(false)
		m.keyMap.Track.SetEnabled(false)
		m.keyMap.Merge.SetEnabled(false)
		m.keyMap.Rebase.SetEnabled(false)

		m.list.KeyMap.AcceptWhileFiltering.SetEnabled(false)
		m.list.KeyMap.CancelWhileFiltering.SetEnabled(false)
	case browsing:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Create.SetEnabled(true)
		m.keyMap.Delete.SetEnabled(true)
		m.keyMap.Merge.SetEnabled(true)
		m.keyMap.Rebase.SetEnabled(true)
		m.keyMap.Track.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		m.keyMap.Cancel.SetEnabled(false)

	default:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Create.SetEnabled(true)
		m.keyMap.Delete.SetEnabled(true)
		m.keyMap.Merge.SetEnabled(true)
		m.keyMap.Rebase.SetEnabled(true)
		m.keyMap.Track.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		m.keyMap.Cancel.SetEnabled(false)
	}
}
