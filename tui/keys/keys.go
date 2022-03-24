package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Enter      key.Binding
	Create     key.Binding
	Delete     key.Binding
	Track      key.Binding
	Merge      key.Binding
	Rebase     key.Binding
	Rename     key.Binding
	Cancel     key.Binding
	Quit       key.Binding
	ForceQuit  key.Binding

	State string
}

func (k KeyMap) ShortHelp() []key.Binding {
	var kb []key.Binding

	if k.State != "browsing" {
		kb = append(kb, k.Cancel, k.ForceQuit)
	}

	return kb
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "move up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Check out the currently selected branch"),
		),
		Create: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp(
				"ctrl+a",
				"Create a new branch, with confirmation",
			),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp(
				"ctrl+d",
				"Delete the currently selected branch, with confirmation",
			),
		),
		Track: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "Track the currently selected branch"),
		),
		Merge: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp(
				"ctrl+y",
				"Merge the currently selected branch, with confirmation",
			),
		),
		Rebase: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp(
				"ctrl+u",
				"Rebase the currently selected branch, with confirmation",
			),
		),
		Rename: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "Rename the currently selected branch"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Force quit"),
		),
	}
}
