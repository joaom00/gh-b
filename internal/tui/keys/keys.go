package keys

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Enter     key.Binding
	Create    key.Binding
	Delete    key.Binding
	Track     key.Binding
	Merge     key.Binding
	Rebase    key.Binding
	Cancel    key.Binding
	Help      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding

	State string
}

func (k KeyMap) ShortHelp() []key.Binding {
	var kb []key.Binding

	if k.State == "creating" || k.State == "deleting" || k.State == "merge" {
		kb = append(kb, k.Cancel, k.ForceQuit)
	}

	return kb
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Create, k.Delete},
		{k.Track, k.Merge, k.Rebase},
		{k.Cancel, k.Help, k.Quit},
	}
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Check out the currently selected branch"),
		),
		Create: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp(
				"ctrl+a",
				"Create a new branch, with confirmation prompt before creation",
			),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp(
				"ctrl+d",
				"Delete the currently selected branch, with confirmation prompt before deletion",
			),
		),
		Track: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "Track currently selected branch"),
		),
		Merge: key.NewBinding(
			key.WithKeys("ctrl+y"),
			key.WithHelp(
				"ctrl+y",
				"Merge the currently selected branch, with confirmation prompt before merge",
			),
		),
		Rebase: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "Rebase currently selected branch"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "Toggle help"),
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

func NewHelpModel() help.Model {
	return help.New()
}
