package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/StxrlessLabs/webber/choices"
	"github.com/StxrlessLabs/webber/model"
)

type appModel struct {
	pages   []*choices.SubCategory
	models  []*model.PageModel
	pageIdx int
}

func initialModel() appModel {
	pages := choices.GetAllPages()
	models := make([]*model.PageModel, len(pages))
	for i, p := range pages {
		models[i] = model.NewPageModel(p)
	}
	return appModel{pages: pages, models: models, pageIdx: 0}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := m.models[m.pageIdx]

	if current.IsInputActive() {
		if key, ok := msg.(tea.KeyPressMsg); ok {
			switch key.String() {
			case "enter", "esc":
				current.ConfirmInput()
			default:
				current.UpdateInput(key)
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			current.CursorUp()
		case "down", "j":
			current.CursorDown()
		case "space", "enter":
			current.ToggleSelected()
		case "right", "l", "tab":
			m.pageIdx = (m.pageIdx + 1) % len(m.pages)
		case "left", "h", "shift+tab":
			m.pageIdx = (m.pageIdx - 1 + len(m.pages)) % len(m.pages)
		}
	}

	return m, nil
}

func (m appModel) View() tea.View {
	current := m.models[m.pageIdx]
	pageCounter := fmt.Sprintf("\n\t%d/%d\n", m.pageIdx+1, len(m.pages))
	footer := "\n[tab/shift+tab: switch page] [space/enter: select] [q: quit]\n"
	return tea.NewView(current.View() + pageCounter + footer)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
