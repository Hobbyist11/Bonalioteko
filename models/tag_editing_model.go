package models

import (
	"fmt"
	"strings"

	"Bonalioteko/xattr"

	keymaps "Bonalioteko/Keymaps"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap interface {
	// ShortHelp returns a slice of bindings to be displayed in the short
	// version of the help. The help bubble will render help in the order in
	// which the help items are returned here.
	ShortHelp() []key.Binding

	// FullHelp returns an extended group of help items, grouped by columns.
	// The help bubble will render the help in the order in which the help
	// items are returned here.
	FullHelp() [][]key.Binding
}

const (
	defaultView modelState = iota
	editTagView
)

type TagEditModel struct {
	modelState modelState
	fileName   string
	Tags       []string
	cursor     int
	Styles     Styles
	Width      int
	max        int
	Help       help.Model
	KeyMap     keymaps.KeyMap
	height     int

	textInput textinput.Model
	err       error
}

type ExitTagViewMsg struct {
	Message string
}
type TagsUpdatedMsg struct {
	NewTags  []string
	filename string
}

func initialTextInputModel() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "New tag:"
	ti.Focus()
	ti.CharLimit = 50
	return ti
}

func NewTagEditModel(fileName string, Tags []string) TagEditModel {
	return TagEditModel{
		fileName:  fileName,
		Tags:      Tags,
		cursor:    0,
		Styles:    DefaultStyles(),
		Width:     30,
		KeyMap:    keymaps.DefaultKeyMap(),
		height:    10,
		textInput: initialTextInputModel(),
		Help:      help.New(),
	}
}

func (m TagEditModel) Init() tea.Cmd {
	return nil
}

func (m TagEditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = 30
		m.max = m.Width - 1

	case tea.KeyMsg:
		if m.err != nil {
			m.err = nil
			return m, nil
		}
		if m.modelState == editTagView {
			m.textInput, cmd = m.textInput.Update(msg)

			switch msg.String() {
			case "enter":
				if err := xattr.Addtag(m.fileName, []byte(m.textInput.Value())); err != nil {
					m.err = err
					return m, nil
				}
				m.textInput.Reset()
				var err error
				m.Tags, err = xattr.GetTagsFromPath(m.fileName)
				if err != nil {
					m.err = err
					return m, nil
				}
				cmd = func() tea.Msg { return TagsUpdatedMsg{NewTags: m.Tags, filename: m.fileName} }

			case "esc":
				m.textInput.Blur()
				m.modelState = defaultView

			}

			return m, cmd

		} else {
			switch msg.String() {
			case "right", "l":
				if m.cursor < len(m.Tags)-1 {
					m.cursor++
				}

			case "left", "h":
				if m.cursor > 0 {
					m.cursor--
				}

			case "d":
				if err := xattr.RemoveTag(m.fileName, m.Tags[m.cursor]); err != nil {
					m.err = err
					return m, nil
				}

				var err error
				m.Tags, err = xattr.GetTagsFromPath(m.fileName)
				if err != nil {
					m.err = err
					return m, func() tea.Msg { return TagsUpdatedMsg{NewTags: m.Tags, filename: m.fileName} }
				}
				m.cursor = 0
				return m, func() tea.Msg { return TagsUpdatedMsg{NewTags: m.Tags, filename: m.fileName} }

			case "a":
				m.modelState = editTagView

			case "esc":
				cmd = func() tea.Msg { return ExitTagViewMsg{"Exit"} }
			}
		}

	}
	return m, cmd
}

func (m TagEditModel) View() string {
	var s string
	if m.err != nil {
		s += fmt.Sprintf("error: %v\n\nPress any key to continue", m.err)
	}
	switch m.modelState {
	case defaultView:
		s = lipgloss.Place(50, 50, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.helpView()))

	case editTagView:
		s = lipgloss.Place(50, 50, lipgloss.Top, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.helpView()))
	}
	return s
}

func (m TagEditModel) headerView() string {
	var (
		sections    []string
		availHeight = m.height
	)

	var s strings.Builder

	for i, tagPtr := range m.Tags {
		if m.cursor == i {
			s.WriteString(m.Styles.highlightedtag.Render(tagPtr) + " ")
		} else {
			s.WriteString(m.Styles.tagnames.Render(tagPtr) + " ")
		}
	}

	content := lipgloss.NewStyle().Height(availHeight).Render(s.String())
	sections = append(sections, m.fileName)
	sections = append(sections, content)

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m TagEditModel) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}

func (m TagEditModel) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.CursorRight,
		m.KeyMap.CursorLeft,
	}}

	listLevelBindings := []key.Binding{
		// m.KeyMap.Filter,
		// m.KeyMap.ClearFilter,
		// m.KeyMap.AcceptWhileFiltering,
		// m.KeyMap.CancelWhileFiltering,
	}

	return append(kb,
		listLevelBindings,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m TagEditModel) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.CursorLeft,
		m.KeyMap.CursorRight,
	}

	return append(kb,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	)
}
