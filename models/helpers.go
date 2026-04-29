package models

import (
	"io/fs"
	"log"
	"path/filepath"
	"slices"

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
	// m.tagnames changes len when I add or remove a tag
	targetTag := m.tagnames[m.highlightedtagpos]

	isSelected := slices.Contains(m.selectedTags, targetTag)

	if isSelected {
		m.selectedTags = slices.DeleteFunc(m.selectedTags, func(t *TagItem) bool {
			return t == targetTag // Compare memory addresses
		})

		targetTag.status = false
	} else {
		// Add to selected list
		m.selectedTags = append(m.selectedTags, targetTag)
		targetTag.status = true
	}

	if len(m.selectedTags) == 0 {
		m.choices = m.initialChoices
		m.ebookPaths = find(m.rootdir, ".epub")
		m.highlighted = 0
	} else {
		tagStrings := GetTagStrings(m.selectedTags)
		m.ebookPaths = xattr.MultipleTagsFilter(tagStrings, m.tags)
		m.choices = getTitlesFromPaths(m.ebookPaths)
		m.highlighted = 0
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
