package models

import (
	"fmt"
	"io"
	"os"
	"strings"

	keymaps "Bonalioteko/Keymaps"
	"Bonalioteko/xattr"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	filterView modelState = iota
	normalView
	tagView
)

type modelState int

type Model struct {
	dump    io.Writer
	err     error
	rootdir string

	state modelState

	filterModel list.Model
	tagModel    tea.Model

	ebookPaths []string

	choices        []string
	initialChoices []string
	cursor         string
	highlighted    int

	min int
	max int

	Height     int
	AutoHeight bool

	tags map[string][]string

	highlightedtagpos int

	tagnames       []*TagItem
	selectedTags   []*TagItem
	selectedtagNum int

	mintag int
	maxtag int

	Styles Styles

	pathTags map[string][]string

	KeyMap keymaps.KeyMap
	Help   help.Model
}

type Styles struct {
	cursor      lipgloss.Style
	choices     lipgloss.Style
	highlighted lipgloss.Style

	tagnames       lipgloss.Style
	highlightedtag lipgloss.Style
	selectedtag    lipgloss.Style
	HelpStyle      lipgloss.Style
}

type delegateStyles struct {
	cursor    lipgloss.Style
	greyed    lipgloss.Style
	item      lipgloss.Style
	selected  lipgloss.Style
	tag       lipgloss.Style
	HelpStyle lipgloss.Style
}

func NewStyles() delegateStyles {
	var s delegateStyles
	s.greyed = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	s.cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	s.item = lipgloss.NewStyle().Foreground(lipgloss.Color("02"))
	s.selected = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	s.tag = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	return s
}

func initItems(choices []string, tags []string) []ResultItem {
	var items []ResultItem
	for _, title := range choices {
		items = append(items, &TitleItem{Book: title})
	}
	for _, tag := range tags {
		items = append(items, &TagItem{Tag: tag})
	}
	return items
}

func InitialModel(dump *os.File, rootdir string) Model {
	tagsMap := xattr.GetXattrMapTagToFilePath(rootdir)

	tagStrings := xattr.GetUniqueTags(tagsMap)
	choicesinit := GetEpubTitles(rootdir)

	listItems, sharedTags := GetFilterListItems(tagStrings, choicesinit)

	return Model{
		dump:        dump,
		state:       normalView,
		rootdir:     rootdir,
		filterModel: list.New(listItems, Bonadelegate{styles: NewStyles()}, 80, 40),

		ebookPaths:     find(rootdir, ".epub"),
		choices:        choicesinit,
		initialChoices: choicesinit,
		cursor:         ">",
		Height:         0,
		highlighted:    0,

		Styles: DefaultStyles(),
		min:    0,
		max:    0,

		tags: tagsMap,

		tagnames: sharedTags,

		highlightedtagpos: 0,
		mintag:            0,
		maxtag:            0,
		selectedTags:      nil,
		selectedtagNum:    0,

		pathTags: xattr.GetXattrMapFilePathToTag(rootdir),
		KeyMap:   keymaps.DefaultKeyMap(),
		Help:     help.New(),
	}
}

func DefaultStyles() Styles {
	return DefaultStylesWithRenderer(lipgloss.DefaultRenderer())
}

func DefaultStylesWithRenderer(r *lipgloss.Renderer) Styles {
	return Styles{
		cursor:      r.NewStyle().Foreground(lipgloss.Color("212")),
		choices:     r.NewStyle(),
		highlighted: r.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),

		tagnames:       r.NewStyle().Foreground(lipgloss.Color("5")),
		selectedtag:    r.NewStyle().Italic(true).Foreground(lipgloss.Color("2")),
		highlightedtag: r.NewStyle().Foreground(lipgloss.Color("12")),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type SpecialString string

func (s SpecialString) FilterValue() string {
	return ""
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Height = 100
		m.max = m.Height - 1
		m.filterModel.SetSize(30, 30)

	case TagFilterMsg:
		m.selectedTags = nil
		for _, tag := range m.tagnames {
			if tag.status {
				m.selectedTags = append(m.selectedTags, tag)
			}
		}
		if len(m.selectedTags) == 0 {
			m.choices = m.initialChoices
			m.ebookPaths = find(m.rootdir, ".epub")
		} else {
			tagStrings := GetTagStrings(m.selectedTags)
			m.ebookPaths = xattr.MultipleTagsFilter(tagStrings, m.tags)
			m.choices = getTitlesFromPaths(m.ebookPaths)
		}

	case ExitTagViewMsg:
		m.state = normalView

	case TagsUpdatedMsg:
		m.pathTags[msg.filename] = msg.NewTags

		allTagsMap := xattr.GetXattrMapTagToFilePath(m.rootdir)
		m.tags = allTagsMap
		uniqueTags := xattr.GetUniqueTags(allTagsMap)

		var newTagItems []*TagItem
		for _, tagName := range uniqueTags {
			status := false

			for _, oldTag := range m.tagnames {
				if oldTag.Tag == tagName {
					status = oldTag.status
					break
				}
			}

			newTagItems = append(newTagItems, &TagItem{
				Tag:    tagName,
				status: status,
			})
		}

		m.tagnames = newTagItems
		if m.highlightedtagpos >= len(m.tagnames) {
			m.highlightedtagpos = max(0, len(m.tagnames)-1)
		}

		m.choices = getTitlesFromPaths(m.ebookPaths)
		m.selectedTags = nil
		updatedFilterItems, updatedSharedTags := GetFilterListItems(xattr.GetUniqueTags(m.tags), m.choices)
		m.filterModel.SetItems(updatedFilterItems)
		m.tagnames = updatedSharedTags
		return m, func() tea.Msg { return TagFilterMsg{} }

	case tea.KeyMsg:
		if m.err != nil {
			m.err = nil
			return m, nil
		}
		switch state := m.state; state {
		case filterView:
			m.filterModel, cmd = m.filterModel.Update(msg)
			cmds = append(cmds, cmd)

			if m.filterModel.FilterState() == list.Unfiltered {
				m.state = normalView
			}

			return m, tea.Batch(cmds...)

		case tagView:
			m.tagModel, cmd = m.tagModel.Update(msg)
			cmds = append(cmds, cmd)

		default:
			switch {

			case key.Matches(msg, m.KeyMap.CursorUp):
				m.moveCursorUp()

			case key.Matches(msg, m.KeyMap.CursorDown):
				m.moveCursorDown()

			case key.Matches(msg, m.KeyMap.CursorRight):
				m.moveTagSelectorRight()

			case key.Matches(msg, m.KeyMap.CursorLeft):
				m.moveTagSelectorLeft()

			case key.Matches(msg, m.KeyMap.SpaceBar):
				m.selectOrDeselectTag()

			case key.Matches(msg, m.KeyMap.Filter):
				m.state = filterView
				m.filterModel, cmd = m.filterModel.Update(msg)

			case key.Matches(msg, m.KeyMap.Edit):
				if len(m.ebookPaths) == 0 || m.highlighted < 0 || m.highlighted >= len(m.ebookPaths) {
					break
				}
				m.tagModel = NewTagEditModel(m.ebookPaths[m.highlighted], m.pathTags[m.ebookPaths[m.highlighted]])

				m.state = tagView

			case key.Matches(msg, m.KeyMap.Enter):
				if len(m.ebookPaths) == 0 || m.highlighted < 0 || m.highlighted >= len(m.ebookPaths){
					break
				}
				err := OpenFile(m.ebookPaths[m.highlighted])
				if err != nil {
					m.err = err
				}
			case key.Matches(msg, m.KeyMap.Quit):
				return m, tea.Quit
			}
		}

	default:
		m.filterModel, cmd = m.filterModel.Update(msg)
		cmds = append(cmds, cmd)

	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View model
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n\nPress any key to continue", m.err)
	}
	switch m.state {

	case filterView:
		return m.filterModel.View()

	case tagView:
		return m.tagModel.View()

	default:
		var s strings.Builder

		for i, tagPtr := range m.tagnames {
			if m.highlightedtagpos == i && !tagPtr.status {
				s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.highlightedtag.Render(tagPtr.Tag) + " ")
			} else if tagPtr.status && m.highlightedtagpos != i {
				s.WriteString(m.Styles.selectedtag.Render(tagPtr.Tag) + " ")
			} else if m.highlightedtagpos == i && tagPtr.status {
				s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.selectedtag.Render(tagPtr.Tag) + " ")
			} else {
				s.WriteString(m.Styles.tagnames.Render(tagPtr.Tag) + " ")
			}
		}
		s.WriteString("\n")

		for i, items := range m.choices {
			if i < m.min || i > m.max {
				continue
			}

			if m.highlighted == i {
				highlighted := fmt.Sprint(m.Styles.highlighted.Render(items))
				s.WriteString(m.Styles.cursor.Render(m.cursor) + m.Styles.highlighted.Render(highlighted))
				s.WriteRune('\n')
				continue
			}

			s.WriteString(m.Styles.choices.Render(items))
			s.WriteRune('\n')

		}
		return lipgloss.Place(50, 50, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, s.String(), m.helpView()))
	}
}

func (m Model) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}

func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.CursorLeft,
		m.KeyMap.CursorRight,
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.SpaceBar,
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
func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.CursorLeft,
		m.KeyMap.CursorRight,
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.Filter,
	}

	return append(kb,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	)
}
