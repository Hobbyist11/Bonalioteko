package main


import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"Bonalioteko/config"
	"Bonalioteko/xattr"

	"github.com/charmbracelet/lipgloss"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"

	tea "github.com/charmbracelet/bubbletea"
)

// getTitlesFromPaths converts a slice of file paths into a slice of ePub titles.
func getTitlesFromPaths(paths []string) []string {
	var titles []string
	if paths == nil {
		return titles // Return empty slice, not nil
	}

	for _, p := range paths {
		metadata, err := epub.GetMetadataFromFile(p)
		if err != nil || len(metadata.Title) == 0 {
			// Fallback to the file name if metadata fails
			titles = append(titles, filepath.Base(p))
			continue
		}
		titles = append(titles, metadata.Title[0])
	}
	return titles
}

// resetNavigation resets the cursor and viewport after the choices list changes.
func resetNavigation(m *Model) {
	m.highlighted = 0
	m.min = 0
	if m.Height > 0 {
		m.max = m.Height - 1
	} else {
		m.max = 0 // Or some default
	}
}

// Sets the style
type Styles struct {
	cursor      lipgloss.Style
	choices     lipgloss.Style
	highlighted lipgloss.Style

	tagnames       lipgloss.Style
	highlightedtag lipgloss.Style
	selectedtag    lipgloss.Style
}

// MAIN MODEL
type Model struct {
	// epub title to be displayed
	title []string

	choices     []string
	cursor      string // Which item is pointed out
	highlighted int

	min int
	max int

	Height     int
	AutoHeight bool

	// Maps string tag to files (Tag -> []FilePath)
	tags map[string][]string

	// slice of available tagnames
	tagnames []string
	// Rendering the highlightedtag
	highlightedtag string
	// Which tag the pointer is pointing to
	highlightedtagpos int
	// Which tags are selected
	selectedtags []string
	// Which tagname is selected
	selectedtagNum int

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
	// tags := xattr.GetXattrMapTagToFilePath()
	tagsMap := xattr.GetXattrMapTagToFilePath()

	tagnames := xattr.GetUniqueTags(tagsMap)

	return Model{
		title:       find(ebookDir, ".epub"),
		choices:     ListEpubs(ebookDir),
		cursor:      ">",
		Height:      0,
		highlighted: 0,

		Styles: DefaultStyles(),
		min:    0,
		max:    0,

		// Maps string tag to files
		tags: tagsMap,
		// slice of available tagnames
		tagnames:          tagnames,
		highlightedtagpos: 0,
		mintag:            0,
		maxtag:            0,
		// Which tagname is selected
		selectedtags:   nil,
		selectedtagNum: 0,
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

func ListEpubs(directory string) []string {
	var sr []string
	for _, sr2 := range find(directory, ".epub") {
		sr2, err := epub.GetMetadataFromFile(sr2)
		if err != nil {
			errors.Cause(err)
		}
		sr = append(sr, sr2.Title[0])
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

		case "l":
			m.highlightedtagpos++
			if m.highlightedtagpos >= len(m.tagnames) {
				m.highlightedtagpos = len(m.tagnames) - 1
			}
			if m.highlightedtagpos > m.maxtag {
				m.mintag++
				m.maxtag++
			}

		case "h":
			m.highlightedtagpos--
			if m.highlightedtagpos < 0 {
				m.highlightedtagpos = 0
			}
			if m.highlightedtagpos < m.mintag {
				m.mintag--
				m.maxtag--
			}

		// TODO: Unselect a tag if " " is pressed again, simply check if tagname is selected if yes, then unselect (remove from slice)
		// If not then select it
		case " ":
			m.selectedtagNum = m.highlightedtagpos

			// It deletes it even if I'm not pointing to it
			// go to this tagname, see if it's part of selectedtags, if yes, delete it
			// even if this is true, appen still happens?
			if len(m.selectedtags) > 0 && m.selectedtags[0] == m.tagnames[m.selectedtagNum] {
				m.selectedtags = slices.DeleteFunc(m.selectedtags, func(s string) bool {
					return m.tagnames[m.selectedtagNum] == s
				})
			} else {
				m.selectedtags = append(m.selectedtags, m.tagnames[m.selectedtagNum])

				// TODO: render the title and not the filepath
				m.choices = xattr.MultipleTagsFilter(m.selectedtags)

			}

		}
	}
	return m, nil
}

// View model
func (m Model) View() string {
	var s strings.Builder

	// Renders the UI for moving between tags
	for i, tagchoice := range m.tagnames {
		// Checks if selectedtags slice contains this tag, if yes render as selected
		if slices.Contains(m.selectedtags, tagchoice) {
			s.WriteString(m.Styles.selectedtag.Render(tagchoice) + " ")
			continue
		}
		if m.highlightedtagpos == i {
			// highlighted := fmt.Sprint(m.Styles.highlightedtag.Render(tagchoice))
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

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there has been an error %v", err)
		os.Exit(1)
	}
}
