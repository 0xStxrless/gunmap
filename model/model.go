package model

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/StxrlessLabs/gunmap/choices"
	"github.com/StxrlessLabs/gunmap/styles"
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
		m.selected = make(map[int]struct{})
		if already {
			return false
		}
		m.selected[m.cursor] = struct{}{}
		if takesInput {
			m.activeFlag = key
			m.inputBuf = m.inputValues[key]
			m.inputActive = true
			return true
		}
		return false
	}

	if _, ok := m.selected[m.cursor]; ok {
		delete(m.selected, m.cursor)
		return false
	}

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

func (m *PageModel) SelectedCount() int {
	return len(m.selected)
}

func (m *PageModel) SelectedCommands() []string {
	out := make([]string, 0, len(m.selected))
	for idx := range m.selected {
		raw := m.keys[idx]
		flag, _ := styles.SplitFlagDescription(raw)
		if val, ok := m.inputValues[raw]; ok && val != "" {
			out = append(out, flag+" "+val)
		} else {
			out = append(out, flag)
		}
	}
	sort.Strings(out)
	return out
}

func (m *PageModel) flagColWidth() int {
	width := 0
	for _, k := range m.keys {
		flag, _ := styles.SplitFlagDescription(k)
		if w := lipgloss.Width(flag); w > width {
			width = w
		}
	}
	return width
}

func (m *PageModel) View() string {
	return m.viewWidth(0)
}

func (m *PageModel) ViewWidth(width int) string {
	return m.viewWidth(width)
}

func (m *PageModel) viewWidth(width int) string {
	var s strings.Builder

	flagColWidth := m.flagColWidth()

	for i, key := range m.keys {
		takesInput := m.Page.Commands[key]
		isCursor := m.cursor == i
		_, isSelected := m.selected[i]

		flag, desc := styles.SplitFlagDescription(key)

		if takesInput {
			if val, ok := m.inputValues[key]; ok && val != "" {
				desc = fmt.Sprintf("%s = %s", desc, val)
			}
		}

		s.WriteString(styles.RenderOption(flag, desc, isCursor, isSelected, flagColWidth, width))
		s.WriteString("\n")
	}

	if m.inputActive {
		flag, _ := styles.SplitFlagDescription(m.activeFlag)
		prompt := fmt.Sprintf("  %s › %s", flag, m.inputBuf)
		cursor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00C8E8")).Render("█")
		s.WriteString("\n")
		s.WriteString(styles.Accent(prompt))
		s.WriteString(cursor)
		s.WriteString("\n")
	}

	return s.String()
}
