package models

import (
	"fmt"
	"strings"

	keymaps "Bonalioteko/Keymaps"
	"Bonalioteko/config"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pirmd/epub"
)

type RecentModel struct {
	err error

	width       int
	choices     []string
	recentFiles []string

	min int
	max int

	cursor      string
	highlighted int
	height      int
	AutoHeight  bool

	Styles Styles

	KeyMap keymaps.KeyMap
	Help   help.Model
}

func (m RecentModel) Init() tea.Cmd {
	return nil
}

func (m RecentModel) Update(msg tea.Msg) (RecentModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = 50
		m.max = m.height - 1

	case tea.KeyMsg:
		if m.err != nil {
			m.err = nil
			return m, nil
		}
		switch {
		case key.Matches(msg, m.KeyMap.CursorDown):
			m.moveCursorDown()
		case key.Matches(msg, m.KeyMap.CursorUp):
			m.moveCursorUp()

		case key.Matches(msg, m.KeyMap.Enter):
			// if len(m.ebookPaths) == 0 || m.highlighted < 0 || m.highlighted >= len(m.ebookPaths) {
			// 	break
			// }
			// err := OpenFile(m.ebookPaths[m.highlighted])
			// if err != nil {
			// 	m.err = err
			// }
			// // Add opened file to Recents slice. If it already exists, append
			// config.AddToRecentFileList(m.ebookPaths[m.highlighted])

		case key.Matches(msg, m.KeyMap.Quit):
			cmd = func() tea.Msg { return ExitTagViewMsg{"Exit"} }
		}
	}
	return m, cmd
}

func (m RecentModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n\nPress any key to continue", m.err)
	}
	var s strings.Builder

	for i, items := range m.choices {
		metadata, err := epub.GetMetadataFromFile(items)
		if err != nil {
			return fmt.Sprintf("error: %v\n\nPress any key to continue", m.err)
		}
		// if i < m.min || i > m.max {
		// 	continue
		// }

		if m.highlighted == i {
			highlighted := fmt.Sprint(m.Styles.highlighted.Render(metadata.Title[0]))
			s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.highlighted.Render(highlighted))
			s.WriteRune('\n')
			continue
		}

		s.WriteString(m.Styles.choices.Render(metadata.Title[0]))

		s.WriteRune('\n')

	}

	return lipgloss.Place(50, 50, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, "Recent eBooks:\n", s.String(), m.helpView()))
}

func (m RecentModel) RecentsView() string {
	var s strings.Builder

	for i, items := range m.choices {
		if i < m.min || i > m.max {
			continue
		}

		if m.highlighted == i {
			highlighted := fmt.Sprint(m.Styles.highlighted.Render(items))
			s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.highlighted.Render(highlighted))
			continue
		}

		s.WriteString(m.Styles.choices.Render(items))

		s.WriteRune('\n')

	}
	return s.String()
}

func NewRecentsModel() RecentModel {
	choices, err := config.GetRecentsSlice()

	return RecentModel{
		err: err,

		width:   10,
		choices: choices,
		// recentFiles []string

		min: 0,
		max: 0,

		cursor:      ">",
		highlighted: 0,
		height:      10,
		// AutoHeight  bool

		Styles: DefaultStyles(),

		KeyMap: keymaps.DefaultKeyMap(),
		Help:   help.New(),
	}
}

func (m *RecentModel) moveCursorUp() {
	m.highlighted--
	if m.highlighted < 0 {
		m.highlighted = 0
	}
	if m.highlighted < m.min {
		m.min--
		m.max--
	}
}

func (m *RecentModel) moveCursorDown() {
	m.highlighted++
	if m.highlighted >= len(m.choices) {
		m.highlighted = len(m.choices) - 1
	}
	if m.highlighted > m.max {
		m.min++
		m.max++
	}
}

func (m RecentModel) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.Filter,
		m.KeyMap.Edit,
	}}

	return append(kb,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m RecentModel) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.Filter,
		m.KeyMap.Tab,
	}

	return append(kb,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	)
}

func (m RecentModel) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}
