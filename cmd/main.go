package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"Bonalioteko/config"
	"Bonalioteko/xattr"

	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	regularView sessionState = iota
	searchView
)

type Item interface {
	FilterValue() string
}

type ItemDelegate interface {
	// Render renders the item's view.
	Render(w io.Writer, m Model, index int, item Item)

	// Height is the height of the list item.
	Height() int

	// Spacing is the size of the horizontal gap between list items in cells.
	Spacing() int

	// Update is the update loop for items. All messages in the list's update
	// loop will pass through here except when the user is setting a filter.
	// Use this method to perform item-level updates appropriate to this
	// delegate.
	Update(msg tea.Msg, m *Model) tea.Cmd
}

type filteredItem struct {
	index   int
	item    Item
	matches []int
}

type filteredItems []filteredItem

func (f filteredItems) items() []Item {
	agg := make([]Item, len(f))
	for i, v := range f {
		agg[i] = v.item
	}
	return agg
}

// FilterMatchesMsg contains data about items matched during filtering. The
// message should be routed to Update for processing.
type FilterMatchesMsg []filteredItem

// FilterFunc takes a term and a list of strings to search through
// (defined by Item#FilterValue).
// It should return a sorted list of ranks.
type FilterFunc func(string, []string) []Rank

// Rank defines a rank for a given item.
type Rank struct {
	// The index of the item in the original input.
	Index int
	// Indices of the actual word that were matched against the filter term.
	MatchedIndexes []int
}

// DefaultFilter uses the sahilm/fuzzy to filter through the list.
// This is set by default.
func DefaultFilter(term string, targets []string) []Rank {
	ranks := fuzzy.Find(term, targets)
	sort.Stable(ranks)
	result := make([]Rank, len(ranks))
	for i, r := range ranks {
		result[i] = Rank{
			Index:          r.Index,
			MatchedIndexes: r.MatchedIndexes,
		}
	}
	return result
}

type statusMessageTimeoutMsg struct{}

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

// String returns a human-readable string of the current filter state.
func (f FilterState) String() string {
	return [...]string{
		"unfiltered",
		"filtering",
		"filter applied",
	}[f]
}

type Styles struct {
	cursor      lipgloss.Style
	choices     lipgloss.Style
	highlighted lipgloss.Style

	tagnames       lipgloss.Style
	highlightedtag lipgloss.Style
	selectedtag    lipgloss.Style
}

type ModalScreenUI struct {
	state   sessionState
	regular tea.Model
	search  tea.Model
}

func (m ModalScreenUI) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
// func (m ModalScreenUI) Update(msg tea.Msg) (tea.Model, tea.Cmd){
// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd
// 	switch msg := msg.(type){
//
// 	case
// 	}
// }

type Model struct {
	title []string

	choices        []string
	initialChoices []string
	cursor         string
	highlighted    int

	min int
	max int

	Height     int
	AutoHeight bool

	tags map[string][]string

	tagnames          []string
	highlightedtag    string
	highlightedtagpos int
	selectedtags      []string
	selectedtagNum    int

	mintag int
	maxtag int

	Styles Styles
}

func initialModel() Model {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	const (
		marginBottom  = 5
		fileSizewidth = 7
		paddingLeft   = 2
	)

	ebookDir := cfg.Settings.EbookDir
	tagsMap := xattr.GetXattrMapTagToFilePath()

	tagnames := xattr.GetUniqueTags(tagsMap)

	choicesinit := GetEpubTitles(ebookDir)
	return Model{
		title:          find(ebookDir, ".epub"),
		choices:        choicesinit,
		initialChoices: choicesinit,
		cursor:         ">",
		Height:         0,
		highlighted:    0,

		Styles: DefaultStyles(),
		min:    0,
		max:    0,

		tags: tagsMap,

		tagnames:          tagnames,
		highlightedtagpos: 0,
		mintag:            0,
		maxtag:            0,
		selectedtags:      nil,
		selectedtagNum:    0,
	}
}

// DefaultStyles defines the default styling for the file picker.
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

func find(root, ext string) []string {
	var filename []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			filename = append(filename, s)
		}
		return nil
	})
	return filename
}

func GetEpubTitles(directory string) []string {
	var titlesSlice []string
	for _, titles := range find(directory, ".epub") {
		titles, err := epub.GetMetadataFromFile(titles)
		if err != nil {
			errors.Cause(err)
		}
		titlesSlice = append(titlesSlice, titles.Title[0])
	}
	return titlesSlice
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = 10
		m.max = m.Height - 1

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

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

		}
	}
	return m, nil
}

// View model
func (m Model) View() string {
	var s strings.Builder

	for i, tagchoice := range m.tagnames {
		if slices.Contains(m.selectedtags, tagchoice) {
			s.WriteString(m.Styles.selectedtag.Render(tagchoice) + " ")
			continue
		}
		if m.highlightedtagpos == i {
			s.WriteString(m.Styles.highlightedtag.Render(tagchoice) + " ")
			continue
		}
		s.WriteString(m.Styles.tagnames.Render(tagchoice) + " ")
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

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()

	m := initialModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error %v", err)
		os.Exit(1)
	}
}
