package view

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

func PrintHelp(out io.Writer) {
	theme := helpTheme()

	title := lipgloss.JoinHorizontal(
		lipgloss.Top,
		theme.Title.Render("BGO"),
		theme.Subtitle.Render(" Weather CLI"),
	)
	subtitle := theme.Muted.Render("Fast, clear weather in your terminal.")

	usage := strings.Join([]string{
		theme.Section.Render("[*] USAGE"),
		"  bgo [now|forecast] [flags]",
		"",
		theme.Section.Render("[*] COMMANDS"),
		"  now       Show current weather",
		"  forecast  Show forecast for multiple days",
	}, "\n")

	flags := theme.Section.Render("[*] FLAGS") + "\n" + renderHelpTable(theme)

	notes := strings.Join([]string{
		theme.Section.Render("[*] ENV"),
		"  OPEN_WEATHER_API_KEY  API key for OpenWeather",
		"  .env                 Optional file in project root",
	}, "\n")

	examples := strings.Join([]string{
		theme.Section.Render("[*] EXAMPLES"),
		"  bgo now",
		"  bgo forecast -d 3",
		"  bgo forecast -d 2 --with-time",
		"  bgo now --high-precision",
		"  bgo now --full-info",
	}, "\n")

	content := strings.Join([]string{
		title,
		subtitle,
		"",
		usage,
		"",
		flags,
		"",
		notes,
		"",
		examples,
	}, "\n")

	panel := renderGradientPanel(content, theme)
	fmt.Fprintln(out, panel)
}

type helpStyles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Muted    lipgloss.Style
	Section  lipgloss.Style
	Panel    lipgloss.Style
}

func helpTheme() helpStyles {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58A6FF"))
	subtitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F0F6FC"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))
	section := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7EE787"))
	panel := lipgloss.NewStyle().
		Padding(1, 2)

	return helpStyles{
		Title:    title,
		Subtitle: subtitle,
		Muted:    muted,
		Section:  section,
		Panel:    panel,
	}
}

func renderHelpTable(theme helpStyles) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateRows = false
	t.Style().Options.DrawBorder = false

	t.AppendHeader(table.Row{"Option", "Default", "Description"})
	t.AppendRows([]table.Row{
		{"-d, --days <1-5>", "5", "Set forecast length in days"},
		{"--high-precision", "false", "Show values with max precision"},
		{"--full-info", "false", "Print full API response"},
		{"--with-time", "false", "Include time in forecast output"},
		{"-h, --help", "", "Show this help"},
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Colors: text.Colors{text.FgHiCyan},
		},
		{
			Number: 2,
			Colors: text.Colors{text.FgHiYellow},
		},
		{
			Number: 3,
			Colors: text.Colors{text.FgHiWhite},
		},
	})

	t.Style().Color.Header = text.Colors{text.FgHiBlue, text.Bold}
	return t.Render()
}

func renderGradientPanel(content string, theme helpStyles) string {
	width := terminalWidth()
	if width < 60 {
		width = 60
	}
	if width > 96 {
		width = 96
	}

	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	wrapped := lipgloss.NewStyle().Width(innerWidth).Render(content)
	lines := strings.Split(wrapped, "\n")

	top := gradientLine(width-2, []string{"#1B2B34", "#2D4F67", "#4F9CD4"})
	bottom := gradientLine(width-2, []string{"#4F9CD4", "#2D4F67", "#1B2B34"})

	var b strings.Builder
	b.WriteString("+")
	b.WriteString(top)
	b.WriteString("+\n")

	leftBar := lipgloss.NewStyle().Foreground(lipgloss.Color("#2D4F67")).Render("|")
	rightBar := lipgloss.NewStyle().Foreground(lipgloss.Color("#4F9CD4")).Render("|")

	for _, line := range lines {
		padded := padRightVisible(line, innerWidth)
		b.WriteString(leftBar)
		b.WriteString(" ")
		b.WriteString(theme.Panel.Render(padded))
		b.WriteString(" ")
		b.WriteString(rightBar)
		b.WriteString("\n")
	}

	b.WriteString("+")
	b.WriteString(bottom)
	b.WriteString("+")
	return b.String()
}

func gradientLine(width int, stops []string) string {
	if width <= 0 {
		return ""
	}
	if len(stops) < 2 {
		return strings.Repeat("-", width)
	}
	var b strings.Builder
	for i := 0; i < width; i++ {
		t := float64(i) / float64(max(1, width-1))
		color := interpolateStops(stops, t)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("-"))
	}
	return b.String()
}

func interpolateStops(stops []string, t float64) string {
	segment := 1.0 / float64(len(stops)-1)
	index := int(t / segment)
	if index >= len(stops)-1 {
		return stops[len(stops)-1]
	}
	localT := (t - float64(index)*segment) / segment
	return lerpHex(stops[index], stops[index+1], localT)
}

func lerpHex(a string, b string, t float64) string {
	ar, ag, ab := hexToRGB(a)
	br, bg, bb := hexToRGB(b)
	r := int(float64(ar) + (float64(br)-float64(ar))*t)
	g := int(float64(ag) + (float64(bg)-float64(ag))*t)
	bb2 := int(float64(ab) + (float64(bb)-float64(ab))*t)
	return fmt.Sprintf("#%02X%02X%02X", clamp(r), clamp(g), clamp(bb2))
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	r := parseHexByte(hex[0:2])
	g := parseHexByte(hex[2:4])
	b := parseHexByte(hex[4:6])
	return r, g, b
}

func parseHexByte(s string) int {
	var v int
	fmt.Sscanf(s, "%02x", &v)
	return v
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func padRightVisible(s string, width int) string {
	visible := lipgloss.Width(s)
	if visible >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visible)
}

func terminalWidth() int {
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		if w, _, err := term.GetSize(fd); err == nil && w > 0 {
			return w
		}
	}
	return 80
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
