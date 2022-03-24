package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joaom00/gh-b/tui/keys"
	"github.com/joaom00/gh-b/tui/styles"
)

type itemDelegate struct {
	keys   *keys.KeyMap
	styles *styles.Styles
}

func newItemDelegate(keys *keys.KeyMap, styles *styles.Styles) *itemDelegate {
	return &itemDelegate{
		keys:   keys,
		styles: styles,
	}
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title := d.styles.NormalTitle.Render
	desc := d.styles.NormalDesc.Render

	if index == m.Index() {
		title = func(s string) string {
			return d.styles.SelectedTitle.Render("> " + s)
		}
		desc = func(s string) string {
			return d.styles.SelectedDesc.Render(s)
		}
	}

	branch := title(i.Name)
	author := desc(i.AuthorName)
	committerDate := desc(fmt.Sprintf("(%s)", i.CommitterDate))

	itemListStyle := fmt.Sprintf("%s %s %s", branch, author, committerDate)

	fmt.Fprint(w, itemListStyle)
}

func (d itemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (d itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{d.keys.Track, d.keys.Create, d.keys.Delete, d.keys.Merge, d.keys.Rebase, d.keys.Rename},
	}
}
