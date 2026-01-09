package main

import (
	"slices"

	"Bonalioteko/xattr"
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

func (m *Model) selectOrDeselectTag() {
	m.selectedtagNum = m.highlightedtagpos

	if len(m.selectedtags) > 0 && slices.Contains(m.selectedtags, m.tagnames[m.selectedtagNum]) {
		m.selectedtags = slices.DeleteFunc(m.selectedtags, func(s string) bool {
			return m.tagnames[m.selectedtagNum] == s
		})

		m.choices = m.initialChoices
	} else {
		m.selectedtags = append(m.selectedtags, m.tagnames[m.selectedtagNum])

		m.choices = xattr.MultipleTagsFilter(m.selectedtags)
		m.choices = getTitlesFromPaths(m.choices)
	}
}
