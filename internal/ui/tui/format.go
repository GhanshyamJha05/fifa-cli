package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"github.com/GhanshyamJha05/fifa-cli/internal/ui/styles"
)

func formatHome(info *domain.TournamentInfo, matches []domain.Match, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("FIFA World Cup 2026") + "\n")
	b.WriteString(t.Dim.Render("USA · Canada · Mexico") + "\n\n")

	if info != nil {
		pct := service.TournamentProgress(info)
		filled := int(pct / 100 * 30)
		bar := lipgloss.NewStyle().Foreground(t.Primary).Render(strings.Repeat("█", filled)) +
			lipgloss.NewStyle().Foreground(t.Muted).Render(strings.Repeat("░", 30-filled))
		b.WriteString(fmt.Sprintf("Stage: %s\n", info.CurrentStage))
		b.WriteString(fmt.Sprintf("Progress: %s %.0f%% (%d/%d)\n\n", bar, pct, info.Completed, info.TotalMatches))
	}

	b.WriteString(t.Header.Render("Today's Matches") + "\n")
	b.WriteString(strings.Repeat("─", 50) + "\n")
	if len(matches) == 0 {
		b.WriteString(t.Dim.Render("  No matches scheduled for today.") + "\n")
	} else {
		for _, m := range matches {
			b.WriteString(formatMatchRow(m, t) + "\n")
		}
	}
	return b.String()
}

func formatMatchRow(m domain.Match, t styles.Theme) string {
	score := "  -  "
	if m.HomeScore != nil && m.AwayScore != nil {
		score = fmt.Sprintf(" %d-%d ", *m.HomeScore, *m.AwayScore)
	}

	status := m.Date.Format("Jan 2  15:04")
	switch m.Status {
	case domain.StatusLive:
		status = t.Live.Render(fmt.Sprintf("LIVE %d'", m.Minute))
	case domain.StatusFinished:
		status = t.Dim.Render("FT")
	}

	home := truncate(m.HomeTeam.Name, 16)
	away := truncate(m.AwayTeam.Name, 16)
	return fmt.Sprintf("  %-16s%s%-16s  %s  %s", home, score, away, status, truncate(m.Stadium, 18))
}

func formatMatches(matches []domain.Match, filter string, t styles.Theme) string {
	var b strings.Builder
	title := "All Matches"
	switch filter {
	case "today":
		title = "Today's Matches"
	case "live":
		title = "Live Matches"
	case "upcoming":
		title = "Upcoming Fixtures"
	case "results":
		title = "Results"
	}
	b.WriteString(t.Title.Render(title) + "\n")
	b.WriteString(t.Dim.Render(fmt.Sprintf("%d matches", len(matches))) + "\n\n")
	b.WriteString(fmt.Sprintf("  %-16s %-7s %-16s %-12s %s\n",
		"HOME", "SCORE", "AWAY", "STATUS", "VENUE"))
	b.WriteString(strings.Repeat("─", 72) + "\n")

	if len(matches) == 0 {
		b.WriteString(t.Dim.Render("\n  No matches in this view.") + "\n")
		return b.String()
	}

	for _, m := range matches {
		b.WriteString(formatMatchRow(m, t) + "\n")
	}
	return b.String()
}

func formatStandings(groups []domain.GroupStanding, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("Group Standings") + "\n\n")

	for _, g := range groups {
		b.WriteString(t.Header.Render("Group "+g.Group) + "\n")
		b.WriteString(fmt.Sprintf("  %-3s %-18s %2s %2s %2s %2s %3s %3s %3s %4s %s\n",
			"#", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts", "Form"))
		b.WriteString(t.Dim.Render("  "+strings.Repeat("─", 58)) + "\n")
		for _, row := range g.Rows {
			line := fmt.Sprintf("  %-3d %-18s %2d %2d %2d %2d %3d %3d %+3d %4d %s",
				row.Rank, truncate(row.Team.Name, 18),
				row.Played, row.Won, row.Drawn, row.Lost,
				row.GoalsFor, row.GoalsAgainst, row.GoalDiff, row.Points, row.Form)
			if row.Rank <= 2 {
				b.WriteString(t.Selected.Render(line) + "\n")
			} else {
				b.WriteString("  " + line[2:] + "\n")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func formatStats(s *domain.TournamentStats, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render("Tournament Stats") + "\n\n")

	b.WriteString(t.Header.Render("Top Scorers") + "\n")
	if len(s.TopScorers) == 0 {
		b.WriteString(t.Dim.Render("  No data yet.") + "\n")
	} else {
		for _, ps := range s.TopScorers {
			b.WriteString(fmt.Sprintf("  %2d  %-22s %-16s  %d goals\n",
				ps.Rank, truncate(ps.Player.Name, 22), truncate(ps.Player.TeamName, 16), ps.Value))
		}
	}

	b.WriteString("\n" + t.Header.Render("Top Assists") + "\n")
	if len(s.TopAssists) == 0 {
		b.WriteString(t.Dim.Render("  No data yet.") + "\n")
	} else {
		for _, ps := range s.TopAssists {
			b.WriteString(fmt.Sprintf("  %2d  %-22s %-16s  %d assists\n",
				ps.Rank, truncate(ps.Player.Name, 22), truncate(ps.Player.TeamName, 16), ps.Value))
		}
	}
	return b.String()
}

func formatTeamDetail(team *domain.Team, squad []domain.Player, form *domain.TeamForm, t styles.Theme) string {
	var b strings.Builder
	flag := team.Flag
	if flag == "" {
		flag = "🏳️"
	}
	b.WriteString(t.Title.Render(fmt.Sprintf("%s %s", flag, team.Name)) + "\n")
	b.WriteString(strings.Repeat("─", 50) + "\n\n")

	b.WriteString(fmt.Sprintf("  Code:       %s\n", team.Code))
	if team.Group != "" {
		b.WriteString(fmt.Sprintf("  Group:      %s\n", team.Group))
	}
	if team.Coach != "" {
		b.WriteString(fmt.Sprintf("  Coach:      %s\n", team.Coach))
	}
	if team.FIFAWorldRank > 0 {
		b.WriteString(fmt.Sprintf("  FIFA Rank:  #%d\n", team.FIFAWorldRank))
	}

	if form != nil && len(form.Results) > 0 {
		b.WriteString("\n  Recent form: ")
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

	if len(squad) > 0 {
		b.WriteString("\n" + t.Header.Render("Squad") + "\n")
		b.WriteString(fmt.Sprintf("  %-4s %-22s %-5s %-4s %s\n", "#", "Player", "Pos", "Age", "Club"))
		b.WriteString(t.Dim.Render("  "+strings.Repeat("─", 50)) + "\n")
		for _, p := range squad {
			cap := ""
			if p.Captain {
				cap = " ©"
			}
			b.WriteString(fmt.Sprintf("  %-4d %-22s %-5s %-4d %s%s\n",
				p.Number, truncate(p.Name, 22), p.Position, p.Age, truncate(p.Club, 20), cap))
		}
	}
	return b.String()
}

func formatPlayer(p *domain.Player, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(p.Name) + "\n")
	b.WriteString(strings.Repeat("─", 40) + "\n\n")
	b.WriteString(fmt.Sprintf("  Team:      %s\n", p.TeamName))
	b.WriteString(fmt.Sprintf("  Position:  %s\n", p.Position))
	b.WriteString(fmt.Sprintf("  Number:    %d\n", p.Number))
	b.WriteString(fmt.Sprintf("  Age:       %d\n", p.Age))
	if p.Club != "" {
		b.WriteString(fmt.Sprintf("  Club:      %s\n", p.Club))
	}
	b.WriteString(fmt.Sprintf("  Goals:     %d\n", p.Goals))
	b.WriteString(fmt.Sprintf("  Assists:   %d\n", p.Assists))
	b.WriteString(fmt.Sprintf("  Cards:     %dY / %dR\n", p.Yellow, p.Red))
	if p.Minutes > 0 {
		b.WriteString(fmt.Sprintf("  Minutes:   %d\n", p.Minutes))
	}
	if p.Captain {
		b.WriteString(lipgloss.NewStyle().Foreground(t.Secondary).Render("\n  Team Captain\n"))
	}
	return b.String()
}

func formatSearchResults(results []domain.SearchResult, query string, t styles.Theme) string {
	var b strings.Builder
	b.WriteString(t.Title.Render(fmt.Sprintf("Results for \"%s\"", query)) + "\n\n")
	if len(results) == 0 {
		b.WriteString(t.Dim.Render("  No matches found. Try a team or player name.") + "\n")
		return b.String()
	}
	for i, r := range results {
		icon := "•"
		switch r.Type {
		case "team":
			icon = "🏳️"
		case "player":
			icon = "👤"
		case "match":
			icon = "⚽"
		}
		b.WriteString(fmt.Sprintf("  %d. %s [%s] %s", i+1, icon, r.Type, r.Title))
		if r.Subtitle != "" {
			b.WriteString(t.Dim.Render(" — " + r.Subtitle))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-1] + "…"
}
