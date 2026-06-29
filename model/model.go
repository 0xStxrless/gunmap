package model

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/StxrlessLabs/webber/choices"
)

type PageModel struct {
	Page        *choices.SubCategory
	keys        []string
	cursor      int
	selected    map[int]struct{}
	inputValues map[string]string // flag -> value
	activeFlag  string            // which flag is being typed for
	inputBuf    string            // current typing buffer
	inputActive bool
}

func (m *PageModel) IsInputActive() bool {
	return m.inputActive
}

func NewPageModel(page *choices.SubCategory) *PageModel {
	keys := make([]string, 0, len(page.Commands))
	for k := range page.Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &PageModel{
		Page:        page,
		keys:        keys,
		selected:    make(map[int]struct{}),
		inputValues: make(map[string]string),
	}
}

func (m *PageModel) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *PageModel) CursorDown() {
	if m.cursor < len(m.keys)-1 {
		m.cursor++
	}
}

func (m *PageModel) ToggleSelected() (needsInput bool) {
	if len(m.keys) == 0 {
		return false
	}

	key := m.keys[m.cursor]
	takesInput := m.Page.Commands[key]

	if !m.Page.IsMultiChoice {
		_, already := m.selected[m.cursor]
		// always clear first for single-choice
		m.selected = make(map[int]struct{})
		if already {
			// was selected, now deselected — done
			return false
		}
		// select it
		m.selected[m.cursor] = struct{}{}
		if takesInput {
			m.activeFlag = key
			m.inputBuf = m.inputValues[key]
			m.inputActive = true
			return true
		}
		return false
	}

	// multi-choice
	if _, ok := m.selected[m.cursor]; ok {
		// deselect
		delete(m.selected, m.cursor)
		return false
	}

	// select
	m.selected[m.cursor] = struct{}{}
	if takesInput {
		m.activeFlag = key
		m.inputBuf = m.inputValues[key]
		m.inputActive = true
		return true
	}
	return false
}

func (m *PageModel) UpdateInput(msg tea.KeyPressMsg) {
	switch msg.String() {
	case "backspace":
		if len(m.inputBuf) > 0 {
			m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.inputBuf += msg.String()
		}
	}
}

func (m *PageModel) ConfirmInput() {
	if m.activeFlag != "" {
		if m.inputBuf == "" {
			// Empty input = deselect the flag
			for i, k := range m.keys {
				if k == m.activeFlag {
					delete(m.selected, i)
					break
				}
			}
		} else {
			m.inputValues[m.activeFlag] = m.inputBuf
		}
	}
	m.inputBuf = ""
	m.activeFlag = ""
	m.inputActive = false
}

// SelectedCommands returns flags ready to be assembled into an nmap command.
// Flags with values are returned as "flag value", plain flags as "flag".
func (m *PageModel) SelectedCommands() []string {
	out := make([]string, 0, len(m.selected))
	for idx := range m.selected {
		flag := m.keys[idx]
		if val, ok := m.inputValues[flag]; ok && val != "" {
			out = append(out, flag+" "+val)
		} else {
			out = append(out, flag)
		}
	}
	sort.Strings(out)
	return out
}

func (m *PageModel) View() string {
	var s strings.Builder
	s.WriteString(m.Page.Title)
	s.WriteString("\n")
	if m.Page.IsMultiChoice {
		s.WriteString("Pick more than one\n")
	} else {
		s.WriteString("Pick one\n")
	}

	for i, key := range m.keys {
		takesInput := m.Page.Commands[key]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		extra := ""
		if takesInput {
			if val, ok := m.inputValues[key]; ok && val != "" {
				extra = fmt.Sprintf(" = %q", val)
			} else {
				extra = " (needs input)"
			}
		}
		fmt.Fprintf(&s, "%s [%s] %s%s\n", cursor, checked, key, extra)
	}

	if m.inputActive {
		fmt.Fprintf(&s, "\nValue for %s: %s█\n", m.activeFlag, m.inputBuf)
	}

	return s.String()
}
