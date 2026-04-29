package keymaps

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu.
type KeyMap struct {
	// Keybindings used when browsing the list.
	CursorRight key.Binding
	CursorLeft  key.Binding
	CursorUp    key.Binding
	CursorDown  key.Binding
	Filter      key.Binding
	ClearFilter key.Binding
	Edit key.Binding
	Enter key.Binding
	SpaceBar key.Binding

	// Keybindings used when setting a filter.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Browsing.
		CursorRight: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("->", "l"),
		),
		CursorLeft: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("<-", "h"),
		),
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
			),
SpaceBar: key.NewBinding(
			key.WithKeys(" "),
		key.WithHelp("SpaceBar", "selectTag"),
			),
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}
