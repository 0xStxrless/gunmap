package styles

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	bgBase  = lipgloss.Color("#0A0B0F")
	bgPanel = lipgloss.Color("#0F1017")
	bgCard  = lipgloss.Color("#13141C")

	borderSubtle = lipgloss.Color("#1C1C2E")
	borderActive = lipgloss.Color("#2A2A48")

	cyan      = lipgloss.Color("#00C8E8") // primary accent
	cyanDim   = lipgloss.Color("#005F72") // borders/rules
	cyanGhost = lipgloss.Color("#002830") // bg accents

	green = lipgloss.Color("#00E87A") // selected
	amber = lipgloss.Color("#E8A200") // warning
	red   = lipgloss.Color("#E84040") // error

	textPrimary = lipgloss.Color("#D8D8F0") // main text
	textMuted   = lipgloss.Color("#6666A0") // secondary / labels
	textGhost   = lipgloss.Color("#2E2E50") // disabled / placeholder
	textInverse = lipgloss.Color("#0A0B0F") // text on bright backgrounds

	crumbActive = lipgloss.Color("#E8C468") // current crumb segment
	crumbDim    = lipgloss.Color("#5A5A78") // separator / trailing crumb

	sidebarBgActive = lipgloss.Color("#15151F")
	sidebarRule     = lipgloss.Color("#00C8E8") // left accent bar on active page row
)

func Primary(s string) string {
	return lipgloss.NewStyle().Foreground(textPrimary).Render(s)
}

func Muted(s string) string {
	return lipgloss.NewStyle().Foreground(textMuted).Render(s)
}

func Ghost(s string) string {
	return lipgloss.NewStyle().Foreground(textGhost).Render(s)
}

func Accent(s string) string {
	return lipgloss.NewStyle().Foreground(cyan).Bold(true).Render(s)
}

func Active(s string) string {
	return lipgloss.NewStyle().Foreground(green).Bold(true).Render(s)
}

func Warn(s string) string {
	return lipgloss.NewStyle().Foreground(amber).Render(s)
}

func Danger(s string) string {
	return lipgloss.NewStyle().Foreground(red).Render(s)
}

// Single-line breadcrumb: "gunmap / scan › current-page"  ...  "page/total"
// A thin borderSubtle rule separates it from content the only decoration.
// RenderHeader renders the top breadcrumb bar.
// crumbs is the trail of segments after the root wordmark, e.g. []string{"scan", "timing"}.
// The last segment is rendered as the active (highlighted) crumb.
func RenderHeader(crumbs []string, page, total, width int) string {
	root := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(bgBase).
		Bold(true).
		Render("gunmap")

	sep := lipgloss.NewStyle().Foreground(crumbDim).Background(bgBase).Render(" / ")

	var b strings.Builder
	b.WriteString(root)

	for i, c := range crumbs {
		if i == 0 {
			b.WriteString(sep)
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(crumbDim).Background(bgBase).Render(" › "))
		}
		if i == len(crumbs)-1 {
			b.WriteString(lipgloss.NewStyle().Foreground(crumbActive).Background(bgBase).Render(c))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(textMuted).Background(bgBase).Render(c))
		}
	}

	left := b.String()

	progress := lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(bgBase).
		Bold(true).
		Render(fmt.Sprintf("%d", page))

	slash := lipgloss.NewStyle().Foreground(textGhost).Background(bgBase).Render("/")

	tot := lipgloss.NewStyle().Foreground(textMuted).Background(bgBase).Render(fmt.Sprintf("%d", total))

	right := progress + slash + tot

	usedWidth := lipgloss.Width(left) + lipgloss.Width(right) + 4 // 4 = side padding×2
	gapWidth := max(width-usedWidth, 0)
	gap := lipgloss.NewStyle().Background(bgBase).Render(strings.Repeat(" ", gapWidth))

	row := left + gap + right

	content := lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Background(bgBase).
		Render(row)

	rule := lipgloss.NewStyle().Foreground(borderSubtle).Render(strings.Repeat("─", width))

	return content + "\n" + rule
}

func RenderFooter(width int) string {
	bindings := []struct{ key, action string }{
		{"↑↓", "move"},
		{"space", "toggle"},
		{"↔", "page"},
		{"e", "edit value"},
		{"r", "run"},
		{"q", "quit"},
	}

	sep := lipgloss.NewStyle().Background(bgBase).Render("    ")
	gapStyle := lipgloss.NewStyle().Background(bgBase)

	parts := make([]string, len(bindings))
	for i, b := range bindings {
		action := lipgloss.NewStyle().Foreground(textMuted).Background(bgBase).Render(b.action)
		parts[i] = renderKey(b.key) + gapStyle.Render(" ") + action
	}

	inner := gapStyle.Render("  ") + strings.Join(parts, sep)

	rowWidth := lipgloss.Width(inner) + 2 // +2 for Padding(0, 1)
	if pad := width - rowWidth; pad > 0 {
		inner += gapStyle.Render(strings.Repeat(" ", pad))
	}

	row := lipgloss.NewStyle().
		Padding(0, 1).
		Background(bgBase).
		Render(inner)

	return row
}

func renderKey(s string) string {
	return lipgloss.NewStyle().
		Foreground(textPrimary).
		Background(bgCard).
		Padding(0, 1).
		Bold(true).
		Render(s)
}

func RenderContent(inner string, width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 2).
		Background(bgBase).
		Render(inner)
}

type SidebarItem struct {
	Label string // page name, lowercased, e.g. "timing"
	Count int    // number of selected options on that page; 0 = no badge
}

func RenderSidebar(items []SidebarItem, activeIdx, width, height int) string {
	var rows []string

	eyebrow := lipgloss.NewStyle().
		Foreground(textGhost).
		Bold(true).
		Padding(0, 2).
		Render("PAGES")
	rows = append(rows, eyebrow, "")

	for i, item := range items {
		active := i == activeIdx

		rowBg := bgBase
		if active {
			rowBg = sidebarBgActive
		}

		labelFg := textMuted
		labelBold := false
		if active {
			labelFg, labelBold = textPrimary, true
		}
		label := lipgloss.NewStyle().Foreground(labelFg).Bold(labelBold).Background(rowBg).Render(item.Label)

		badge := ""
		if item.Count > 0 {
			badgeFg := textMuted
			if active {
				badgeFg = textPrimary
			}
			badge = lipgloss.NewStyle().Foreground(badgeFg).Background(rowBg).Render(fmt.Sprintf("%d", item.Count))
		}

		innerWidth := width - 4 // minus left bar + outer padding
		gapWidth := max(innerWidth-lipgloss.Width(label)-lipgloss.Width(badge), 1)
		gap := lipgloss.NewStyle().Background(rowBg).Render(strings.Repeat(" ", gapWidth))
		line := label + gap + badge

		borderColor := bgBase // invisible border on inactive rows, keeps alignment consistent
		if active {
			borderColor = sidebarRule
		}

		rowStyle := lipgloss.NewStyle().
			Background(rowBg).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(borderColor).
			BorderBackground(rowBg).
			Padding(0, 1).
			Width(width - 1)

		rows = append(rows, rowStyle.Render(line))
	}

	body := strings.Join(rows, "\n")

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(bgBase).
		Render(body)
}

func RenderVerticalRule(height int) string {
	line := lipgloss.NewStyle().Foreground(borderSubtle).Render("│")
	lines := make([]string, height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// Three states: resting (ghosted), focused (cursor here), selected (checked).
// Flag and description sit in aligned columns; focused row gets a full-width
// highlight background, resting rows are dimmed.
// SplitFlagDescription splits a raw command-map key like
// "-T<0-5> (Paranoid/.../Insane) timing template" into a short flag token
// ("-T<0-5>") and the remaining description. It splits on the first space.
func SplitFlagDescription(raw string) (flag, description string) {
	raw = strings.TrimSpace(raw)
	before, after, ok := strings.Cut(raw, " ")
	if !ok {
		return raw, ""
	}
	return before, strings.TrimSpace(after)
}

// RenderOption renders a single selectable flag item in the list.
// focused = cursor is on this item; selected = the flag is toggled on.
// flagColWidth controls the alignment of the description column.
// IMPORTANT: every fragment's style carries its own explicit Background.
// Wrapping already-rendered (ANSI-containing) text in an outer
// Width+Background style does not reliably repaint the background behind
// existing escape codes — it can leave a patchy, partial highlight instead
// of a clean full-width bar. Setting the background on each fragment up
// front, then padding with plain spaces in that same background, avoids
// the issue entirely.
func RenderOption(flag, description string, focused, selected bool, flagColWidth, width int) string {
	rowBg := bgBase
	if focused {
		rowBg = sidebarBgActive
	}

	var indicatorFg color.Color
	switch {
	case selected:
		indicatorFg = green
	case focused:
		indicatorFg = cyan
	default:
		indicatorFg = textGhost
	}
	indicatorChar := " "
	switch {
	case selected:
		indicatorChar = "✓"
	case focused:
		indicatorChar = "-"
	}
	indicator := lipgloss.NewStyle().Foreground(indicatorFg).Background(rowBg).Render(indicatorChar)

	var flagFg color.Color
	flagBold := false
	switch {
	case selected:
		flagFg, flagBold = green, true
	case focused:
		flagFg, flagBold = textPrimary, true
	default:
		flagFg = textMuted
	}
	flagRendered := lipgloss.NewStyle().
		Foreground(flagFg).
		Bold(flagBold).
		Background(rowBg).
		Render(fmt.Sprintf("%-*s", flagColWidth, flag))

	var descColor color.Color
	switch {
	case focused, selected:
		descColor = textPrimary
	default:
		descColor = textGhost
	}
	desc := lipgloss.NewStyle().Foreground(descColor).Background(rowBg).Render(description)

	spacer := lipgloss.NewStyle().Background(rowBg).Render("  ")
	leadIn := lipgloss.NewStyle().Background(rowBg).Render("  ")

	row := leadIn + indicator + spacer + flagRendered + spacer + desc

	// pad the remainder of the row width with plain, same-background spaces
	// rather than re-wrapping the whole (already-styled) row in another
	// Width()+Background() style.
	rowWidth := lipgloss.Width(row)
	if pad := width - rowWidth; pad > 0 {
		row += lipgloss.NewStyle().Background(rowBg).Render(strings.Repeat(" ", pad))
	}

	return row
}

func RenderSectionTitle(title, subtitle string, width int) string {
	label := lipgloss.NewStyle().
		Foreground(textPrimary).
		Bold(true).
		Render(strings.ToUpper(title))

	sub := lipgloss.NewStyle().
		Foreground(textGhost).
		Render("  " + subtitle)

	row := label + sub

	heading := lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Background(bgBase).
		Render(row)

	rule := lipgloss.NewStyle().Foreground(borderSubtle).Render(strings.Repeat("─", width))

	return heading + "\n" + rule
}

func RenderCommandPreview(cmd string, width int) string {
	prompt := lipgloss.NewStyle().
		Foreground(crumbActive).
		Render("$ nmap ")

	command := lipgloss.NewStyle().
		Foreground(textPrimary).
		Render(cmd)

	inner := prompt + command

	rule := lipgloss.NewStyle().Foreground(borderSubtle).Render(strings.Repeat("─", width))

	row := lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Background(bgBase).
		Render(inner)

	return rule + "\n" + row
}

func RenderTargetPrompt(flagsLine, targetBuf string, width int) string {
	prompt := lipgloss.NewStyle().
		Foreground(textGhost).
		Render("$ nmap ")

	flags := lipgloss.NewStyle().
		Foreground(textMuted).
		Render(flagsLine + " ")

	var typed string
	if targetBuf == "" {
		typed = lipgloss.NewStyle().Foreground(textGhost).Italic(true).Render("target")
	} else {
		typed = lipgloss.NewStyle().Foreground(green).Bold(true).Render(targetBuf)
	}

	cursor := lipgloss.NewStyle().Foreground(cyan).Render("█")

	inner := prompt + flags + typed + cursor

	hint := lipgloss.NewStyle().
		Foreground(textGhost).
		Render("enter to run · esc to cancel")

	gap := max(width-lipgloss.Width(inner)-lipgloss.Width(hint)-4, 1)

	row := inner + strings.Repeat(" ", gap) + hint

	rule := lipgloss.NewStyle().Foreground(cyanDim).Render(strings.Repeat("─", width))

	rendered := lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Background(bgBase).
		Render(row)

	return rule + "\n" + rendered
}

type StatusKind int

const (
	StatusInfo StatusKind = iota
	StatusSuccess
	StatusWarn
	StatusError
)

func RenderStatus(msg string, kind StatusKind, width int) string {
	var col color.Color
	var icon string
	switch kind {
	case StatusSuccess:
		col, icon = green, "✔"
	case StatusWarn:
		col, icon = amber, "!"
	case StatusError:
		col, icon = red, "✖"
	default:
		col, icon = cyan, "·"
	}

	badge := lipgloss.NewStyle().Foreground(textInverse).Background(col).
		Padding(0, 1).Bold(true).Render(icon)

	text := lipgloss.NewStyle().Foreground(textPrimary).Render("  " + msg)

	return lipgloss.NewStyle().
		Width(width).
		Padding(0, 2).
		Background(bgPanel).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, badge, text))
}
