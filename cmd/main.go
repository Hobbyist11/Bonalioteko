package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"Bonalioteko/config"
	"Bonalioteko/xattr"

	"github.com/charmbracelet/lipgloss"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"

	tea "github.com/charmbracelet/bubbletea"
)

// Sets the style
type Styles struct {
	cursor      lipgloss.Style
	choices     lipgloss.Style
	highlighted lipgloss.Style

	tags           lipgloss.Style
	selectedtag    lipgloss.Style
	highlightedtag lipgloss.Style
}

// MAIN MODEL
type Model struct {
	// epub title to be displayed
	title []string

	// directory of the file
	choices     []string
	cursor      string // Which item is pointed out
	highlighted int

	min int
	max int

	Height     int
	AutoHeight bool

	tags        []string
	selectedtag int

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

	return Model{
		title:       find(ebookDir, ".epub"),
		choices:     ListEpubs(ebookDir),
		cursor:      ">",
		Height:      0,
		highlighted: 0,

		Styles: DefaultStyles(),
		min:    0,
		max:    0,

		tags:        xattr.GetXattr(),
		selectedtag: 0,
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

		tags:        r.NewStyle(),
		selectedtag: r.NewStyle().Italic(true).Foreground(lipgloss.Color("21")),
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

func ListEpubs(directory string) []string {
	var sr []string
	for _, sr2 := range find(directory, ".epub") {
		sr2, err := epub.GetMetadataFromFile(sr2)
		if err != nil {
			errors.Cause(err)
		}
		sr = append(sr, sr2.Title...)
	}
	return sr
}

// Runs on start up
func (m Model) Init() tea.Cmd {
	return nil
}


// UPDATE=handle incoming events
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
			m.highlighted--
			if m.highlighted < 0 {
				m.highlighted = 0
			}
			if m.highlighted < m.min {
				m.min--
				m.max--
			}

		case "down", "j":
			m.highlighted++
			if m.highlighted >= len(m.choices) {
				m.highlighted = len(m.choices) - 1
			}
			if m.highlighted > m.max {
				m.min++
				m.max++
			}


}

	}
	return m, nil
}

// view
func (m Model) View() string {
	var s strings.Builder

	for i, tagname := range m.tags {

		if m.selectedtag == i {
			s.WriteString(m.Styles.selectedtag.Render(tagname))
			continue
		}
		s.WriteString(m.Styles.tags.Render(tagname) + " ")
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

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error %v", err)
		os.Exit(1)
	}
}
