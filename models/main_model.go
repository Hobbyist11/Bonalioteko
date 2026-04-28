package models

import (
	"fmt"
	"io"
	"os"
	"strings"

	"Bonalioteko/xattr"

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
	dump io.Writer

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

type styles struct {
	cursor    lipgloss.Style
	greyed    lipgloss.Style
	item      lipgloss.Style
	selected  lipgloss.Style
	tag       lipgloss.Style
	HelpStyle lipgloss.Style
}

func NewStyles() styles {
	var s styles
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

func InitialModel(dump *os.File) Model {
	tagsMap := xattr.GetXattrMapTagToFilePath()

	tagStrings := xattr.GetUniqueTags(tagsMap)
	choicesinit := GetEpubTitles(xattr.Ebookdir)

	listItems, sharedTags := GetFilterListItems(tagStrings, choicesinit)

	return Model{
		dump:        dump,
		state:       normalView,
		filterModel: list.New(listItems, Bonadelegate{styles: NewStyles()}, 80, 40),

		ebookPaths:     find(xattr.Ebookdir, ".epub"),
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

		pathTags: xattr.GetXattrMapFilePathToTag(),
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
			m.ebookPaths = find(xattr.Ebookdir, ".epub")
		} else {
			tagStrings := GetTagStrings(m.selectedTags)
			m.ebookPaths = xattr.MultipleTagsFilter(tagStrings)
			m.choices = getTitlesFromPaths(m.ebookPaths)
		}

	case ExitTagViewMsg:
		m.state = normalView

	case TagsUpdatedMsg:
		m.pathTags[msg.filename] = msg.NewTags

		allTagsMap := xattr.GetXattrMapTagToFilePath()
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
			switch msg.String() {

			case "up", "k":
				m.moveCursorUp()

			case "down", "j":
				m.moveCursorDown()

			case "l":
				m.moveTagSelectorRight()

			case "h":
				m.moveTagSelectorLeft()

			case " ":
				m.selectOrDeselectTag()

			case "/":
				m.state = filterView
				m.filterModel, cmd = m.filterModel.Update(msg)

			case "e":
				if len(m.ebookPaths) == 0 || m.highlighted < 0 || m.highlighted >= len(m.ebookPaths) {
					break
				}
				m.tagModel = NewTagEditModel(m.ebookPaths[m.highlighted], m.pathTags[m.ebookPaths[m.highlighted]])

				m.state = tagView

			case "esc", "ctrl+c":
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
	switch m.state {

	case filterView:
		return m.filterModel.View()

	case tagView:
		return m.tagModel.View()

	default:
		var s strings.Builder

		for i, tagPtr := range m.tagnames {
			if tagPtr.status {
				s.WriteString(m.Styles.selectedtag.Render(tagPtr.Tag) + " ")
			} else if m.highlightedtagpos == i {
				s.WriteString(m.Styles.highlightedtag.Render(tagPtr.Tag) + " ")
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
		return s.String()
	}
}
