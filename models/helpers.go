package models

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"

	"Bonalioteko/xattr"

	"github.com/charmbracelet/bubbles/list"
	"github.com/pirmd/epub"
)

func (m *Model) moveCursorUp() {
	m.highlighted--
	if m.highlighted < 0 {
		m.highlighted = 0
	}
	if m.highlighted < m.min {
		m.min--
		m.max--
	}
}

func (m *Model) moveCursorDown() {
	m.highlighted++
	if m.highlighted >= len(m.choices) {
		m.highlighted = len(m.choices) - 1
	}
	if m.highlighted > m.max {
		m.min++
		m.max++
	}
}

func (m *Model) moveTagSelectorRight() {
	m.highlightedtagpos++
	if m.highlightedtagpos >= len(m.tagnames) {
		m.highlightedtagpos = len(m.tagnames) - 1
	}
	if m.highlightedtagpos > m.maxtag {
		m.mintag++
		m.maxtag++
	}
}

func (m *Model) moveTagSelectorLeft() {
	m.highlightedtagpos--
	if m.highlightedtagpos < 0 {
		m.highlightedtagpos = 0
	}
	if m.highlightedtagpos < m.mintag {
		m.mintag--
		m.maxtag--
	}
}

func getTitlesFromPaths(paths []string) []string {
	var titles []string
	if paths == nil {
		return titles
	}

	for _, p := range paths {
		metadata, err := epub.GetMetadataFromFile(p)
		if err != nil || len(metadata.Title) == 0 {
			titles = append(titles, filepath.Base(p))
			continue
		}
		titles = append(titles, metadata.Title[0])
	}
	return titles
}

func (m *Model) selectOrDeselectTag() {
	targetTag := m.tagnames[m.highlightedtagpos]

	isSelected := slices.ContainsFunc(m.selectedTags, func(t *TagItem) bool {
		return t.Tag == targetTag.Tag
	})

	if isSelected {
		m.selectedTags = slices.DeleteFunc(m.selectedTags, func(t *TagItem) bool {
			return t.Tag == targetTag.Tag // Compare tag names instead of memory addresses
		})

		targetTag.status = false
	} else {
		m.selectedTags = append(m.selectedTags, targetTag)
		targetTag.status = true
	}

	if len(m.selectedTags) == 0 {
		m.choices = m.initialChoices
		m.ebookPaths = find(m.rootdir, ".epub")
		m.highlighted = 0
		NewTagsToPath := SetTagToPathMap(m.ebookPaths)
		uniqueTags := xattr.GetUniqueTags(NewTagsToPath)

		var newItems []list.Item
		newItems, m.tagnames = GetFilterListItems(uniqueTags, m.choices)
		m.filterModel.SetItems(newItems)
		m.highlighted = 0
		m.highlightedtagpos = 0

	} else {
		selectedTagString := GetTagStrings(m.selectedTags)
		m.ebookPaths = xattr.MultipleTagsFilter(selectedTagString, m.tags)
		m.choices = getTitlesFromPaths(m.ebookPaths)
		NewTagsToPath := SetTagToPathMap(m.ebookPaths)
		uniqueTags := xattr.GetUniqueTags(NewTagsToPath)

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

		var newItems []list.Item
		for _, tagPtr := range m.tagnames {
			newItems = append(newItems, tagPtr)
		}
		for _, choice := range m.choices {
			newItems = append(newItems, &TitleItem{Book: choice})
		}
		m.filterModel.SetItems(newItems)
		m.highlighted = 0
		m.highlightedtagpos = 0
	}
}

func SetTagToPathMap(paths []string) map[string][]string {
	result := make(map[string][]string)
	for _, path := range paths {
		tags, _ := xattr.GetTagsFromPath(path)
		xattr.AddTagAndFile(path, tags, result)
	}
	return result
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
	for _, path := range find(directory, ".epub") {
		metadata, err := epub.GetMetadataFromFile(path)
		if err != nil || len(metadata.Title) == 0 {
			log.Printf("Warning: could not read metadata for %s: %v", path, err)
			titlesSlice = append(titlesSlice, filepath.Base(path))
			continue
		}
		titlesSlice = append(titlesSlice, metadata.Title[0])
	}
	return titlesSlice
}

func GetFilterListItems(tagStrings []string, choicesinit []string) (listItems []list.Item, sharedTags []*TagItem) {
	for _, t := range tagStrings {
		sharedTags = append(sharedTags, &TagItem{Tag: t, status: false})
	}

	combinedList := initItems(choicesinit, tagStrings)
	for _, tagPtr := range sharedTags {
		listItems = append(listItems, tagPtr)
	}

	for _, tagPtr := range combinedList {
		if v, ok := tagPtr.(*TitleItem); ok {
			listItems = append(listItems, v)
		}
	}
	return listItems, sharedTags
}

// OpenFile opens a file using the default application associated with its file type
func OpenFile(path string) error {
	cleaned := filepath.Clean(path)

	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return fmt.Errorf("abs %q: %w", cleaned, err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("stat: %q %w", abs, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %q", abs)
	}

	_, inFlatpak := os.LookupEnv("FLATPAK_ID")

	if inFlatpak {
		// Reach out to the host system to run xdg-open
		if _, err := exec.LookPath("flatpak-spawn"); err != nil {
			return fmt.Errorf("flatpak-spawn not found: %w", err)
		}
		if _, err := exec.LookPath("xdg-open"); err != nil {
			return fmt.Errorf("xdg-open not found: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "flatpak-spawn", "--host", "xdg-open", abs)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("starting xdg-open for %q: %w", abs, err)
		}
	} else {
		if _, err := exec.LookPath("xdg-open"); err != nil {
			return fmt.Errorf("xdg-open not found: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "xdg-open", abs)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("starting xdg-open for %q: %w", abs, err)
		}
	}

	return nil
}
