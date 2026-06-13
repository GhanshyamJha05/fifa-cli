package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/styles"
)

// Teams renders a team list panel.
func Teams(teams []domain.Team, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("🏆 World Cup Teams") + "\n\n")
	for _, team := range teams {
		line := fmt.Sprintf("%s  %-22s  %s  Group %s  #%d",
			team.Flag, team.Name, team.Code, team.Group, team.FIFAWorldRank)
		b.WriteString(t.Normal.Render(line) + "\n")
	}
	return t.Panel.Render(b.String())
}

// TeamDetail renders detailed team info.
func TeamDetail(team *domain.Team, form *domain.TeamForm, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("%s %s", team.Flag, team.Name)) + "\n\n")
	b.WriteString(fmt.Sprintf("  Code:     %s\n", team.Code))
	b.WriteString(fmt.Sprintf("  Group:    %s\n", team.Group))
	b.WriteString(fmt.Sprintf("  Coach:    %s\n", team.Coach))
	b.WriteString(fmt.Sprintf("  FIFA Rank: #%d\n", team.FIFAWorldRank))
	if form != nil && len(form.Results) > 0 {
		b.WriteString("\n  Form (last 5): ")
		for i, r := range form.Results {
			switch r {
			case "W":
				b.WriteString(t.Win.Render(r))
			case "D":
				b.WriteString(t.Draw.Render(r))
			case "L":
				b.WriteString(t.Loss.Render(r))
			default:
				b.WriteString(r)
			}
			if i < len(form.Results)-1 {
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}
	return t.Panel.Render(b.String())
}

// Squad renders a team squad table.
func Squad(players []domain.Player, teamName string, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("👥 %s Squad", teamName)) + "\n\n")
	header := fmt.Sprintf("  %-4s %-24s %-6s %-4s %-20s",
		"#", "Player", "Pos", "Age", "Club")
	b.WriteString(t.TableHeader.Render(header) + "\n")
	b.WriteString(strings.Repeat("─", 60) + "\n")
	for _, p := range players {
		cap := ""
		if p.Captain {
			cap = "©"
		}
		line := fmt.Sprintf("  %-4d %-24s %-6s %-4d %-20s %s",
			p.Number, p.Name, p.Position, p.Age, p.Club, cap)
		if p.Captain {
			b.WriteString(t.Selected.Render(line) + "\n")
		} else {
			b.WriteString(t.Normal.Render(line) + "\n")
		}
	}
	return t.Panel.Render(b.String())
}

// PlayerProfile renders player details.
func PlayerProfile(p *domain.Player, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("⚽ %s", p.Name)) + "\n\n")
	b.WriteString(fmt.Sprintf("  Team:      %s\n", p.TeamName))
	b.WriteString(fmt.Sprintf("  Position:  %s\n", p.Position))
	b.WriteString(fmt.Sprintf("  Number:    %d\n", p.Number))
	b.WriteString(fmt.Sprintf("  Age:       %d\n", p.Age))
	b.WriteString(fmt.Sprintf("  Club:      %s\n", p.Club))
	b.WriteString(fmt.Sprintf("  Goals:     %d\n", p.Goals))
	b.WriteString(fmt.Sprintf("  Assists:   %d\n", p.Assists))
	b.WriteString(fmt.Sprintf("  Yellow:    %d  Red: %d\n", p.Yellow, p.Red))
	b.WriteString(fmt.Sprintf("  Minutes:   %d\n", p.Minutes))
	if p.Captain {
		b.WriteString(lipgloss.NewStyle().Foreground(t.Secondary).Render("\n  © Team Captain\n"))
	}
	return t.Panel.Render(b.String())
}

// MatchLine renders a single match row.
func MatchLine(m domain.Match, t styles.Theme) string {
	score := "vs"
	if m.HomeScore != nil && m.AwayScore != nil {
		score = fmt.Sprintf("%d - %d", *m.HomeScore, *m.AwayScore)
	}
	status := m.Date.Format("Jan 2 15:04")
	switch m.Status {
	case domain.StatusLive:
		status = t.Live.Render(fmt.Sprintf("🔴 LIVE %d'", m.Minute))
	case domain.StatusFinished:
		status = lipgloss.NewStyle().Foreground(t.Muted).Render("FT")
	case domain.StatusScheduled:
		status = m.Date.Format("Jan 2 15:04")
	}
	return fmt.Sprintf("  %-3s  %-18s  %s  %-18s  %s  %s",
		m.Group, m.HomeTeam.Name, score, m.AwayTeam.Name, status, m.Stadium)
}

// Matches renders a match list.
func Matches(matches []domain.Match, title string, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(title) + "\n\n")
	if len(matches) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(t.Muted).Render("  No matches found.") + "\n")
	} else {
		for _, m := range matches {
			b.WriteString(MatchLine(m, t) + "\n")
		}
	}
	return t.Panel.Render(b.String())
}

// Standings renders group tables.
func Standings(groups []domain.GroupStanding, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("📊 Group Standings") + "\n\n")
	for _, g := range groups {
		b.WriteString(t.Header.Render(fmt.Sprintf("Group %s", g.Group)) + "\n")
		hdr := fmt.Sprintf("  %-3s %-18s %3s %3s %3s %3s %4s %4s %3s %5s",
			"#", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts")
		b.WriteString(t.TableHeader.Render(hdr) + "\n")
		for _, row := range g.Rows {
			line := fmt.Sprintf("  %-3d %-18s %3d %3d %3d %3d %4d %4d %+3d %5d  %s",
				row.Rank, row.Team.Name, row.Played, row.Won, row.Drawn, row.Lost,
				row.GoalsFor, row.GoalsAgainst, row.GoalDiff, row.Points, row.Form)
			if row.Rank <= 2 {
				b.WriteString(t.Selected.Render(line) + "\n")
			} else {
				b.WriteString(t.Normal.Render(line) + "\n")
			}
		}
		b.WriteString("\n")
	}
	return t.Panel.Render(b.String())
}

// Bracket renders knockout stage.
func Bracket(b *domain.Bracket, t styles.Theme) string {
	var out strings.Builder
	out.WriteString(t.Title.Render("🏟️  Knockout Bracket") + "\n\n")
	renderRound := func(name string, matches []domain.BracketMatch) {
		if len(matches) == 0 {
			return
		}
		out.WriteString(t.Header.Render(name) + "\n")
		for _, m := range matches {
			score := " vs "
			if m.HomeScore != nil && m.AwayScore != nil {
				score = fmt.Sprintf(" %d-%d ", *m.HomeScore, *m.AwayScore)
			}
			line := fmt.Sprintf("  %s%s%s  (%s)", m.HomeTeam, score, m.AwayTeam, m.Date.Format("Jan 2"))
			out.WriteString(t.Normal.Render(line) + "\n")
		}
		out.WriteString("\n")
	}
	renderRound("Quarter-Finals", b.Quarter)
	renderRound("Semi-Finals", b.Semi)
	if b.ThirdPlace.HomeTeam != "" {
		renderRound("Third Place", []domain.BracketMatch{b.ThirdPlace})
	}
	if b.Final.HomeTeam != "" {
		renderRound("Final 🏆", []domain.BracketMatch{b.Final})
	}
	return t.Panel.Render(out.String())
}

// Stats renders tournament statistics.
func Stats(s *domain.TournamentStats, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("📈 Tournament Statistics") + "\n\n")

	b.WriteString(t.Header.Render("Top Scorers") + "\n")
	for _, ps := range s.TopScorers {
		b.WriteString(fmt.Sprintf("  %2d. %-24s %s  ⚽ %d\n",
			ps.Rank, ps.Player.Name, ps.Player.TeamName, ps.Value))
	}
	b.WriteString("\n")

	b.WriteString(t.Header.Render("Top Assists") + "\n")
	for _, ps := range s.TopAssists {
		b.WriteString(fmt.Sprintf("  %2d. %-24s %s  🅰️  %d\n",
			ps.Rank, ps.Player.Name, ps.Player.TeamName, ps.Value))
	}
	return t.Panel.Render(b.String())
}

// HeadToHead renders a team comparison.
func HeadToHead(h *domain.HeadToHead, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("⚔️  %s vs %s", h.TeamA.Name, h.TeamB.Name)) + "\n\n")
	b.WriteString(fmt.Sprintf("  Record: %s %dW - %dD - %dW %s\n\n",
		h.TeamA.Name, h.WinsA, h.Draws, h.WinsB, h.TeamB.Name))
	if len(h.Matches) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(t.Muted).Render("  No meetings in this tournament.") + "\n")
	} else {
		for _, m := range h.Matches {
			b.WriteString(MatchLine(m, t) + "\n")
		}
	}
	return t.Panel.Render(b.String())
}

// SearchResults renders search hits.
func SearchResults(results []domain.SearchResult, query string, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("🔍 Search: %q", query)) + "\n\n")
	if len(results) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(t.Muted).Render("  No results found.") + "\n")
	} else {
		for _, r := range results {
			icon := "•"
			switch r.Type {
			case "team":
				icon = "🏳️"
			case "player":
				icon = "👤"
			case "match":
				icon = "⚽"
			}
			line := fmt.Sprintf("  %s [%s] %s", icon, r.Type, r.Title)
			if r.Subtitle != "" {
				line += " — " + r.Subtitle
			}
			b.WriteString(t.Normal.Render(line) + "\n")
		}
	}
	return t.Panel.Render(b.String())
}

// ProgressBar renders tournament progress.
func ProgressBar(info *domain.TournamentInfo, t styles.Theme, width int) string {
	pct := service.TournamentProgress(info)
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	bar := lipgloss.NewStyle().Foreground(t.Primary).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(t.Muted).Render(strings.Repeat("░", width-filled))
	return fmt.Sprintf("  %s  %.0f%% (%d/%d matches)", bar, pct, info.Completed, info.TotalMatches)
}

// TournamentOverview renders tournament header info.
func TournamentOverview(info *domain.TournamentInfo, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(styles.BannerASCII(t))
	b.WriteString(fmt.Sprintf("\n  Host: %s\n", info.Host))
	b.WriteString(fmt.Sprintf("  Stage: %s\n", info.CurrentStage))
	b.WriteString(ProgressBar(info, t, 30) + "\n")
	return b.String()
}
