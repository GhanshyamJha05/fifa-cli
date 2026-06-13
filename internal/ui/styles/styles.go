package styles

import "github.com/charmbracelet/lipgloss"

// Theme holds Lip Gloss styles for the dashboard.
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Danger     lipgloss.Color
	Muted      lipgloss.Color
	Background lipgloss.Color
	Surface    lipgloss.Color
	Text       lipgloss.Color
	Border     lipgloss.Color

	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Panel       lipgloss.Style
	Header      lipgloss.Style
	Selected    lipgloss.Style
	Normal      lipgloss.Style
	Live        lipgloss.Style
	Win         lipgloss.Style
	Draw        lipgloss.Style
	Loss        lipgloss.Style
	Help        lipgloss.Style
	Tab         lipgloss.Style
	TabActive   lipgloss.Style
	Banner      lipgloss.Style
	TableHeader lipgloss.Style
	StatusBar   lipgloss.Style

	// TUI-specific
	AppTitle   lipgloss.Style
	TabBar     lipgloss.Style
	Content    lipgloss.Style
	Footer     lipgloss.Style
	Chip       lipgloss.Style
	ChipActive lipgloss.Style
	Dim        lipgloss.Style
}

var themes = map[string]Theme{}

func init() {
	themes["dark"] = buildTheme("dark",
		lipgloss.Color("#00D4AA"), lipgloss.Color("#FFD700"), lipgloss.Color("#FF6B6B"),
		lipgloss.Color("#4ECDC4"), lipgloss.Color("#FFE66D"), lipgloss.Color("#FF6B6B"),
		lipgloss.Color("#6B7280"), lipgloss.Color("#0B0F19"), lipgloss.Color("#131A2B"),
		lipgloss.Color("#E5E7EB"), lipgloss.Color("#00D4AA"),
	)
	themes["light"] = buildTheme("light",
		lipgloss.Color("#00695C"), lipgloss.Color("#F57F17"), lipgloss.Color("#C62828"),
		lipgloss.Color("#00897B"), lipgloss.Color("#F9A825"), lipgloss.Color("#D32F2F"),
		lipgloss.Color("#757575"), lipgloss.Color("#F5F5F5"), lipgloss.Color("#FFFFFF"),
		lipgloss.Color("#212121"), lipgloss.Color("#00695C"),
	)
	themes["fifa"] = buildTheme("fifa",
		lipgloss.Color("#4DA3FF"), lipgloss.Color("#F5C518"), lipgloss.Color("#E31837"),
		lipgloss.Color("#00A651"), lipgloss.Color("#FFB300"), lipgloss.Color("#E31837"),
		lipgloss.Color("#8B9CB3"), lipgloss.Color("#0A1628"), lipgloss.Color("#112240"),
		lipgloss.Color("#E8EDF5"), lipgloss.Color("#4DA3FF"),
	)
}

func buildTheme(name string, primary, secondary, accent, success, warning, danger, muted, bg, surface, text, border lipgloss.Color) Theme {
	t := Theme{
		Name: name, Primary: primary, Secondary: secondary, Accent: accent,
		Success: success, Warning: warning, Danger: danger, Muted: muted,
		Background: bg, Surface: surface, Text: text, Border: border,
	}
	t.Title = lipgloss.NewStyle().Bold(true).Foreground(primary)
	t.Subtitle = lipgloss.NewStyle().Foreground(muted).Italic(true)
	t.Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(1, 2).
		Background(surface).
		Foreground(text)
	t.Header = lipgloss.NewStyle().Bold(true).Foreground(secondary)
	t.Selected = lipgloss.NewStyle().Bold(true).Foreground(bg).Background(primary).Padding(0, 1)
	t.Normal = lipgloss.NewStyle().Foreground(text)
	t.Live = lipgloss.NewStyle().Bold(true).Foreground(danger)
	t.Win = lipgloss.NewStyle().Foreground(success).Bold(true)
	t.Draw = lipgloss.NewStyle().Foreground(warning)
	t.Loss = lipgloss.NewStyle().Foreground(danger)
	t.Help = lipgloss.NewStyle().Foreground(muted)
	t.Tab = lipgloss.NewStyle().Foreground(muted).Padding(0, 1)
	t.TabActive = lipgloss.NewStyle().Bold(true).Foreground(bg).Background(primary).Padding(0, 1)
	t.Banner = lipgloss.NewStyle().Bold(true).Foreground(secondary)
	t.TableHeader = lipgloss.NewStyle().Bold(true).Foreground(primary)
	t.StatusBar = lipgloss.NewStyle().Foreground(muted).Background(surface).Padding(0, 1)
	t.AppTitle = lipgloss.NewStyle().Bold(true).Foreground(primary).Background(bg).Padding(0, 2)
	t.TabBar = lipgloss.NewStyle().Background(surface).Padding(0, 1)
	t.Content = lipgloss.NewStyle().Background(bg).Foreground(text).Padding(1, 2)
	t.Footer = lipgloss.NewStyle().Foreground(muted).Background(surface).Padding(0, 2)
	t.Chip = lipgloss.NewStyle().Foreground(muted).Padding(0, 1)
	t.ChipActive = lipgloss.NewStyle().Bold(true).Foreground(bg).Background(accent).Padding(0, 1)
	t.Dim = lipgloss.NewStyle().Foreground(muted)
	return t
}

// Get returns a theme by name, defaulting to fifa.
func Get(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["fifa"]
}

// BannerASCII returns the FIFA CLI banner for CLI output.
func BannerASCII(t Theme) string {
	lines := []string{
		"⚽ FIFA WORLD CUP 2026 ⚽",
		"━━━━━━━━━━━━━━━━━━━━━━━━━",
		"  USA · CANADA · MEXICO",
	}
	var out string
	for _, l := range lines {
		out += t.Banner.Render(l) + "\n"
	}
	return out
}
