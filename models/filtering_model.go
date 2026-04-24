package models

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

type ResultItem interface {
	list.Item
	isTag() bool
}

type (
	TitleItem struct{ Book string }
	TagItem   struct {
		Tag    string
		status bool
	}
)

func (i TitleItem) Title() string       { return "" }
func (i TitleItem) Description() string { return "" }
func (t TitleItem) FilterValue() string { return t.Book }
func (t TitleItem) isTag() bool         { return false }

func (i TagItem) Title() string       { return "" }
func (i TagItem) Description() string { return "" }
func (t TagItem) FilterValue() string { return t.Tag }
func (t TagItem) isTag() bool         { return true }

func GetTitles(items []*TagItem) []string {
	s := make([]string, len(items))
	for i, t := range items {
		s[i] = t.Tag
	}
	return s
}

type (
	TagFilterMsg struct {
		TagItem *TagItem
	}
	Bonadelegate struct {
		styles styles
	}
)

func (d Bonadelegate) Height() int  { return 1 }
func (d Bonadelegate) Spacing() int { return 0 }

func (d Bonadelegate) Update(msg tea.Msg, lm *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == " " {
			if tag, ok := lm.SelectedItem().(*TagItem); ok {
				tag.status = !tag.status

				return func() tea.Msg {
					return TagFilterMsg{TagItem: tag}
				}
			}
		}
	}
	return nil
}

func (d Bonadelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	switch m.FilterState() {
	case list.Filtering:
		switch v := item.(type) {
		case *TitleItem:
			fmt.Fprintf(w, "📖%s ", d.styles.item.Render(v.Book))
		case *TagItem:
			if !v.status {
				fmt.Fprintf(w, "  🏷 %s", d.styles.item.Render(v.Tag))
			} else {
				fmt.Fprintf(w, "  🏷 %s", d.styles.selected.Render(v.Tag))
			}
		}

	case list.FilterApplied, list.Unfiltered:
		if index == m.Index() {
			prefix := "> "
			switch v := item.(type) {
			case *TitleItem:
				fmt.Fprintf(w, "📖%s ", d.styles.cursor.Render(prefix+v.Book))
			case *TagItem:
				fmt.Fprintf(w, "🏷 %s", d.styles.cursor.Render(v.Tag))
			}
			return
		}

		switch v := item.(type) {
		case *TitleItem:
			fmt.Fprintf(w, " 📖 %s ", d.styles.greyed.Render(v.Book))
		case *TagItem:
			if !v.status {
				fmt.Fprintf(w, "🏷 %s", d.styles.item.Render(v.Tag))
			} else {
				fmt.Fprintf(w, "🏷 %s", d.styles.selected.Render(v.Tag))
			}
		}
	}
}
