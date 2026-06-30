package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/StxrlessLabs/gunmap/choices"
	"github.com/StxrlessLabs/gunmap/model"
	"github.com/StxrlessLabs/gunmap/styles"
	"github.com/charmbracelet/colorprofile"
)

const sidebarWidth = 22

// pageSubtitle gives each page a short one-line description for the
// section-title row since choices.SubCategory has no subtitle field.
var pageSubtitle = map[string]string{
	"HostDiscovery":                 "host discovery & ping options",
	"ScanTechniques":                "choose a scan technique",
	"PortSelection":                 "which ports to scan",
	"ServiceDetection":              "probe services for version info",
	"OS Detection":                  "fingerprint the target OS",
	"Timing And Performance":        "controls scan aggressiveness",
	"Firewall And Evasion Spoofing": "evade firewalls & IDS",
	"Output":                        "output formats & verbosity",
	"Misc":                          "everything else",
}

type appModel struct {
	pages   []*choices.SubCategory
	models  []*model.PageModel
	pageIdx int
	height  int
	width   int

	// Target-input mode: entered by pressing "r". Collects the scan target
	// (IP/hostname/CIDR) before handing off to the real nmap process.
	targetMode bool
	targetBuf  string
	runErr     error
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
	if sz, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = sz.Width
		m.height = sz.Height
		return m, nil
	}

	// nmap has finished (or failed to start). Quit immediately rather than
	// resuming the TUI, so the real terminal with nmap's own output still
	// on screen is what the user is left looking at.
	if doneMsg, ok := msg.(nmapDoneMsg); ok {
		m.runErr = doneMsg.err
		return m, tea.Quit
	}

	if m.targetMode {
		if key, ok := msg.(tea.KeyPressMsg); ok {
			switch key.String() {
			case "esc":
				m.targetMode = false
				m.targetBuf = ""
			case "enter":
				target := strings.TrimSpace(m.targetBuf)
				if target == "" {
					return m, nil
				}
				cmd := buildNmapCmd(m.assembledCommand(), target)
				m.targetMode = false
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return nmapDoneMsg{err: err}
				})
			case "backspace":
				if len(m.targetBuf) > 0 {
					m.targetBuf = m.targetBuf[:len(m.targetBuf)-1]
				}
			default:
				if len(key.String()) == 1 {
					m.targetBuf += key.String()
				}
			}
		}
		return m, nil
	}

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

		case "e":
			current.ToggleSelected()
			current.ToggleSelected()
			// ^ this is the worst hack possible but it does the trick lol
		case "r":
			m.targetMode = true
			m.targetBuf = ""
		}
	}
	return m, nil
}

type nmapDoneMsg struct{ err error }

// buildNmapCmd assembles the real *exec.Cmd for nmap, wired to inherit the
// program's stdin/stdout/stderr so its live scan output appears directly in
// the terminal exactly as if it had been typed by hand.
// this is very haacky but present me is happy
func buildNmapCmd(flagsLine, target string) *exec.Cmd {
	args := strings.Fields(flagsLine)
	args = append(args, target)
	cmd := exec.Command("nmap", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func sidebarLabel(title string) string {
	return strings.ToLower(title)
}

func (m appModel) assembledCommand() string {
	var parts []string
	for _, pm := range m.models {
		parts = append(parts, pm.SelectedCommands()...)
	}
	return strings.Join(parts, " ")
}

func (m appModel) View() tea.View {
	if m.width == 0 {
		return tea.NewView("loading...")
	}

	current := m.models[m.pageIdx]
	currentPage := m.pages[m.pageIdx]

	crumbs := []string{"scan", sidebarLabel(currentPage.Title)}
	header := styles.RenderHeader(crumbs, m.pageIdx+1, len(m.pages), m.width)
	headerH := lipgloss.Height(header)

	footer := styles.RenderFooter(m.width)
	footerH := lipgloss.Height(footer)

	var cmdPreview string
	if m.targetMode {
		cmdPreview = styles.RenderTargetPrompt(m.assembledCommand(), m.targetBuf, m.width)
	} else {
		cmdPreview = styles.RenderCommandPreview(m.assembledCommand(), m.width)
	}
	cmdPreviewH := lipgloss.Height(cmdPreview)

	bodyH := max(m.height-headerH-footerH-cmdPreviewH, 1)

	items := make([]styles.SidebarItem, len(m.pages))
	for i, p := range m.pages {
		items[i] = styles.SidebarItem{
			Label: sidebarLabel(p.Title),
			Count: m.models[i].SelectedCount(),
		}
	}
	sidebar := styles.RenderSidebar(items, m.pageIdx, sidebarWidth, bodyH)
	rule := styles.RenderVerticalRule(bodyH)

	contentWidth := max(m.width-sidebarWidth-1, 1)

	subtitle := pageSubtitle[currentPage.Title]
	sectionTitle := styles.RenderSectionTitle(currentPage.Title, subtitle, contentWidth)
	sectionTitleH := lipgloss.Height(sectionTitle)

	listH := max(bodyH-sectionTitleH, 1)
	optionList := current.ViewWidth(contentWidth)
	content := styles.RenderContent(optionList, contentWidth, listH)

	contentColumn := lipgloss.JoinVertical(lipgloss.Top, sectionTitle, content)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, rule, contentColumn)

	v := tea.NewView(lipgloss.JoinVertical(lipgloss.Top, header, body, cmdPreview, footer))
	v.AltScreen = true
	v.WindowTitle = "gunmap · nmap builder"
	return v
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Printf("%s you must run as sudo!\n", choices.RandomErrFace())
		return
	}

	// this must be here cuz otherwise you couldn't run this as sudo and
	// get em pretty colors
	profile := colorprofile.TrueColor
	p := tea.NewProgram(initialModel(), tea.WithColorProfile(profile))
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
