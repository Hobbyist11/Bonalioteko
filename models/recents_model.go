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
)

type RecentModel struct {
	err error

	width   int
	choices []string
	titles  []string

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
		m.height = msg.Height
		m.width = msg.Width
		m.max = m.height - 1
		m.titles = getTitlesFromPaths(m.choices)

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
			if len(m.choices) == 0 || m.highlighted < 0 || m.highlighted >= len(m.choices) {
				break
			}
			if err := OpenFile(m.choices[m.highlighted]); err != nil {
				m.err = err
				break
			}
			if err := config.AddToRecentFileList(m.choices[m.highlighted]); err != nil {
				m.err = err
				break
			}
			choices, err := config.GetRecentsSlice()
			if err != nil {
				m.err = err
				break
			}
			m.choices = choices

		case key.Matches(msg, m.KeyMap.Quit):
			cmd = func() tea.Msg { return ExitTagViewMsg{"Exit"} }

		case key.Matches(msg, m.KeyMap.Tab):
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

	for i, items := range m.titles {

		if m.highlighted == i {
			highlighted := fmt.Sprint(m.Styles.highlighted.Render(items))
			s.WriteString(m.Styles.cursor.Render(m.cursor) + highlighted)
			s.WriteRune('\n')
			continue
		}

		s.WriteString(m.Styles.choices.Render(items))

		s.WriteRune('\n')

	}

	return lipgloss.Place(50, 50, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, "Recent eBooks:\n", s.String(), m.helpView()))
}

func NewRecentsModel() RecentModel {
	choices, err := config.GetRecentsSlice()

	titles := getTitlesFromPaths(choices)
	return RecentModel{
		err: err,

		width:   10,
		choices: choices,
		titles:  titles,

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
