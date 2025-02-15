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
