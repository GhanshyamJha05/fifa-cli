package mock

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
)

// Provider serves rich demo data for offline use.
type Provider struct {
	teams    []domain.Team
	players  []domain.Player
	matches  []domain.Match
	standings []domain.GroupStanding
}

// New creates a mock provider with 2026 World Cup demo data.
func New() *Provider {
	p := &Provider{}
	p.initTeams()
	p.initPlayers()
	p.initMatches()
	p.initStandings()
	return p
}

func (p *Provider) initTeams() {
	groups := map[string][]domain.Team{
		"A": {
			{ID: 1, Name: "United States", Code: "USA", Flag: "🇺🇸", Group: "A", Coach: "Gregg Berhalter", FIFAWorldRank: 11},
			{ID: 2, Name: "Mexico", Code: "MEX", Flag: "🇲🇽", Group: "A", Coach: "Jaime Lozano", FIFAWorldRank: 14},
			{ID: 3, Name: "Canada", Code: "CAN", Flag: "🇨🇦", Group: "A", Coach: "Jesse Marsch", FIFAWorldRank: 48},
			{ID: 4, Name: "Ecuador", Code: "ECU", Flag: "🇪🇨", Group: "A", Coach: "Sebastián Beccacece", FIFAWorldRank: 31},
		},
		"B": {
			{ID: 5, Name: "Brazil", Code: "BRA", Flag: "🇧🇷", Group: "B", Coach: "Dorival Júnior", FIFAWorldRank: 5},
			{ID: 6, Name: "Argentina", Code: "ARG", Flag: "🇦🇷", Group: "B", Coach: "Lionel Scaloni", FIFAWorldRank: 1},
			{ID: 7, Name: "Uruguay", Code: "URU", Flag: "🇺🇾", Group: "B", Coach: "Marcelo Bielsa", FIFAWorldRank: 9},
			{ID: 8, Name: "Colombia", Code: "COL", Flag: "🇨🇴", Group: "B", Coach: "Néstor Lorenzo", FIFAWorldRank: 12},
		},
		"C": {
			{ID: 9, Name: "France", Code: "FRA", Flag: "🇫🇷", Group: "C", Coach: "Didier Deschamps", FIFAWorldRank: 2},
			{ID: 10, Name: "England", Code: "ENG", Flag: "🏴󠁧󠁢󠁥󠁮󠁧󠁿", Group: "C", Coach: "Gareth Southgate", FIFAWorldRank: 4},
			{ID: 11, Name: "Germany", Code: "GER", Flag: "🇩🇪", Group: "C", Coach: "Julian Nagelsmann", FIFAWorldRank: 16},
			{ID: 12, Name: "Spain", Code: "ESP", Flag: "🇪🇸", Group: "C", Coach: "Luis de la Fuente", FIFAWorldRank: 8},
		},
		"D": {
			{ID: 13, Name: "Portugal", Code: "POR", Flag: "🇵🇹", Group: "D", Coach: "Roberto Martínez", FIFAWorldRank: 6},
			{ID: 14, Name: "Netherlands", Code: "NED", Flag: "🇳🇱", Group: "D", Coach: "Ronald Koeman", FIFAWorldRank: 7},
			{ID: 15, Name: "Belgium", Code: "BEL", Flag: "🇧🇪", Group: "D", Coach: "Domenico Tedesco", FIFAWorldRank: 3},
			{ID: 16, Name: "Croatia", Code: "CRO", Flag: "🇭🇷", Group: "D", Coach: "Zlatko Dalić", FIFAWorldRank: 10},
		},
		"E": {
			{ID: 17, Name: "Japan", Code: "JPN", Flag: "🇯🇵", Group: "E", Coach: "Hajime Moriyasu", FIFAWorldRank: 18},
			{ID: 18, Name: "South Korea", Code: "KOR", Flag: "🇰🇷", Group: "E", Coach: "Hong Myung-bo", FIFAWorldRank: 23},
			{ID: 19, Name: "Australia", Code: "AUS", Flag: "🇦🇺", Group: "E", Coach: "Graham Arnold", FIFAWorldRank: 24},
			{ID: 20, Name: "Morocco", Code: "MAR", Flag: "🇲🇦", Group: "E", Coach: "Walid Regragui", FIFAWorldRank: 13},
		},
		"F": {
			{ID: 21, Name: "Senegal", Code: "SEN", Flag: "🇸🇳", Group: "F", Coach: "Aliou Cissé", FIFAWorldRank: 17},
			{ID: 22, Name: "Nigeria", Code: "NGA", Flag: "🇳🇬", Group: "F", Coach: "Finidi George", FIFAWorldRank: 38},
			{ID: 23, Name: "Ghana", Code: "GHA", Flag: "🇬🇭", Group: "F", Coach: "Otto Addo", FIFAWorldRank: 60},
			{ID: 24, Name: "Cameroon", Code: "CMR", Flag: "🇨🇲", Group: "F", Coach: "Rigobert Song", FIFAWorldRank: 43},
		},
	}

	for _, groupTeams := range groups {
		p.teams = append(p.teams, groupTeams...)
	}
	sort.Slice(p.teams, func(i, j int) bool {
		return p.teams[i].Name < p.teams[j].Name
	})
}

func (p *Provider) initPlayers() {
	squads := map[string][]struct {
		name, pos, club string
		num, age, goals, assists int
		captain bool
	}{
		"Argentina": {
			{"Lionel Messi", "FW", "Inter Miami", 10, 38, 4, 3, true},
			{"Lautaro Martínez", "FW", "Inter Milan", 22, 28, 3, 1, false},
			{"Emiliano Martínez", "GK", "Aston Villa", 23, 32, 0, 0, false},
			{"Rodrigo De Paul", "MF", "Atlético Madrid", 7, 30, 0, 2, false},
			{"Cristian Romero", "DF", "Tottenham", 13, 26, 0, 0, false},
			{"Enzo Fernández", "MF", "Chelsea", 24, 23, 1, 2, false},
		},
		"Brazil": {
			{"Vinícius Júnior", "FW", "Real Madrid", 7, 25, 3, 2, false},
			{"Neymar", "FW", "Al-Hilal", 10, 34, 2, 4, true},
			{"Alisson", "GK", "Liverpool", 1, 32, 0, 0, false},
			{"Casemiro", "MF", "Manchester United", 5, 32, 1, 0, false},
			{"Marquinhos", "DF", "PSG", 4, 30, 0, 0, false},
			{"Rodrygo", "FW", "Real Madrid", 20, 23, 2, 1, false},
		},
		"France": {
			{"Kylian Mbappé", "FW", "Real Madrid", 10, 27, 5, 2, true},
			{"Antoine Griezmann", "FW", "Atlético Madrid", 7, 34, 2, 3, false},
			{"Hugo Lloris", "GK", "LAFC", 1, 38, 0, 0, false},
			{"Ousmane Dembélé", "FW", "PSG", 11, 27, 1, 4, false},
			{"Aurélien Tchouaméni", "MF", "Real Madrid", 8, 25, 0, 1, false},
			{"William Saliba", "DF", "Arsenal", 17, 24, 0, 0, false},
		},
		"England": {
			{"Harry Kane", "FW", "Bayern Munich", 9, 31, 4, 1, true},
			{"Jude Bellingham", "MF", "Real Madrid", 10, 22, 2, 2, false},
			{"Bukayo Saka", "FW", "Arsenal", 7, 23, 1, 3, false},
			{"Jordan Pickford", "GK", "Everton", 1, 30, 0, 0, false},
			{"Declan Rice", "MF", "Arsenal", 4, 25, 0, 1, false},
			{"John Stones", "DF", "Manchester City", 5, 30, 0, 0, false},
		},
		"United States": {
			{"Christian Pulisic", "FW", "AC Milan", 10, 26, 2, 1, true},
			{"Tyler Adams", "MF", "Bournemouth", 4, 25, 0, 1, false},
			{"Matt Turner", "GK", "Nottingham Forest", 1, 30, 0, 0, false},
			{"Weston McKennie", "MF", "Juventus", 8, 26, 1, 0, false},
			{"Antonee Robinson", "DF", "Fulham", 5, 27, 0, 2, false},
			{"Gio Reyna", "MF", "Borussia Dortmund", 7, 22, 1, 1, false},
		},
		"Germany": {
			{"Jamal Musiala", "MF", "Bayern Munich", 10, 22, 3, 2, false},
			{"Florian Wirtz", "MF", "Bayer Leverkusen", 17, 22, 2, 3, false},
			{"Manuel Neuer", "GK", "Bayern Munich", 1, 38, 0, 0, true},
			{"Kai Havertz", "FW", "Arsenal", 7, 25, 2, 0, false},
			{"Joshua Kimmich", "MF", "Bayern Munich", 6, 29, 0, 2, false},
			{"Antonio Rüdiger", "DF", "Real Madrid", 2, 31, 0, 0, false},
		},
	}

	pid := 1
	for _, team := range p.teams {
		squad, ok := squads[team.Name]
		if !ok {
			squad = []struct {
				name, pos, club string
				num, age, goals, assists int
				captain bool
			}{
				{fmt.Sprintf("%s Captain", team.Name), "MF", "Club FC", 10, 28, 1, 1, true},
				{fmt.Sprintf("%s Striker", team.Name), "FW", "Club FC", 9, 26, 2, 0, false},
				{fmt.Sprintf("%s Keeper", team.Name), "GK", "Club FC", 1, 30, 0, 0, false},
			}
		}
		for _, pl := range squad {
			p.players = append(p.players, domain.Player{
				ID: pid, Name: pl.name, TeamID: team.ID, TeamName: team.Name,
				Position: pl.pos, Number: pl.num, Age: pl.age, Club: pl.club,
				Captain: pl.captain, Goals: pl.goals, Assists: pl.assists,
				Yellow: pl.goals % 2, Red: 0, Minutes: pl.goals * 90 + 180,
			})
			pid++
		}
	}
}

func (p *Provider) initMatches() {
	now := time.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	findTeam := func(name string) domain.Team {
		for _, t := range p.teams {
			if t.Name == name {
				return t
			}
		}
		return domain.Team{Name: name}
	}

	fixtures := []struct {
		home, away, stage, group, stadium, city string
		dayOffset, hour int
		homeScore, awayScore *int
		status domain.MatchStatus
		minute int
	}{
		{"United States", "Ecuador", "Group Stage", "A", "MetLife Stadium", "East Rutherford", -3, 18, intPtr(2), intPtr(1), domain.StatusFinished, 0},
		{"Mexico", "Canada", "Group Stage", "A", "Azteca Stadium", "Mexico City", -2, 20, intPtr(1), intPtr(1), domain.StatusFinished, 0},
		{"Argentina", "Colombia", "Group Stage", "B", "Hard Rock Stadium", "Miami", -1, 16, intPtr(3), intPtr(0), domain.StatusFinished, 0},
		{"Brazil", "Uruguay", "Group Stage", "B", "NRG Stadium", "Houston", 0, 14, intPtr(1), intPtr(1), domain.StatusLive, 67},
		{"France", "England", "Group Stage", "C", "SoFi Stadium", "Los Angeles", 0, 17, nil, nil, domain.StatusScheduled, 0},
		{"Germany", "Spain", "Group Stage", "C", "Levi's Stadium", "Santa Clara", 0, 20, nil, nil, domain.StatusScheduled, 0},
		{"Portugal", "Netherlands", "Group Stage", "D", "Lincoln Financial Field", "Philadelphia", 1, 15, nil, nil, domain.StatusScheduled, 0},
		{"Belgium", "Croatia", "Group Stage", "D", "Mercedes-Benz Stadium", "Atlanta", 1, 18, nil, nil, domain.StatusScheduled, 0},
		{"Japan", "Morocco", "Group Stage", "E", "Lumen Field", "Seattle", 2, 16, nil, nil, domain.StatusScheduled, 0},
		{"South Korea", "Australia", "Group Stage", "E", "BC Place", "Vancouver", 2, 19, nil, nil, domain.StatusScheduled, 0},
		{"Senegal", "Nigeria", "Group Stage", "F", "AT&T Stadium", "Arlington", -4, 18, intPtr(2), intPtr(0), domain.StatusFinished, 0},
		{"Ghana", "Cameroon", "Group Stage", "F", "Arrowhead Stadium", "Kansas City", -3, 15, intPtr(0), intPtr(2), domain.StatusFinished, 0},
	}

	mid := 1
	for _, f := range fixtures {
		m := domain.Match{
			ID: mid, HomeTeam: findTeam(f.home), AwayTeam: findTeam(f.away),
			HomeScore: f.homeScore, AwayScore: f.awayScore,
			Date: base.AddDate(0, 0, f.dayOffset).Add(time.Duration(f.hour) * time.Hour),
			Stage: f.stage, Group: f.group, Stadium: f.stadium, City: f.city,
			Country: "USA/Canada/Mexico", Status: f.status, Minute: f.minute,
		}
		p.matches = append(p.matches, m)
		mid++
	}
}

func (p *Provider) initStandings() {
	groups := map[string][]struct {
		team string
		p, w, d, l, gf, ga, pts int
		form string
	}{
		"A": {
			{"United States", 2, 1, 1, 0, 3, 2, 4, "WD"},
			{"Mexico", 2, 1, 1, 0, 2, 1, 4, "DW"},
			{"Ecuador", 2, 1, 0, 1, 2, 3, 3, "LW"},
			{"Canada", 2, 0, 0, 2, 1, 2, 0, "LL"},
		},
		"B": {
			{"Argentina", 2, 2, 0, 0, 5, 1, 6, "WW"},
			{"Brazil", 1, 0, 1, 0, 1, 1, 1, "D"},
			{"Uruguay", 1, 0, 1, 0, 1, 1, 1, "D"},
			{"Colombia", 2, 0, 0, 2, 1, 5, 0, "LL"},
		},
		"C": {
			{"France", 0, 0, 0, 0, 0, 0, 0, ""},
			{"England", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Germany", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Spain", 0, 0, 0, 0, 0, 0, 0, ""},
		},
		"D": {
			{"Portugal", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Netherlands", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Belgium", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Croatia", 0, 0, 0, 0, 0, 0, 0, ""},
		},
		"E": {
			{"Morocco", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Japan", 0, 0, 0, 0, 0, 0, 0, ""},
			{"South Korea", 0, 0, 0, 0, 0, 0, 0, ""},
			{"Australia", 0, 0, 0, 0, 0, 0, 0, ""},
		},
		"F": {
			{"Senegal", 1, 1, 0, 0, 2, 0, 3, "W"},
			{"Cameroon", 1, 1, 0, 0, 2, 0, 3, "W"},
			{"Nigeria", 1, 0, 0, 1, 0, 2, 0, "L"},
			{"Ghana", 1, 0, 0, 1, 0, 2, 0, "L"},
		},
	}

	findTeam := func(name string) domain.Team {
		for _, t := range p.teams {
			if t.Name == name {
				return t
			}
		}
		return domain.Team{Name: name}
	}

	for g := 'A'; g <= 'F'; g++ {
		group := string(g)
		rows := groups[group]
		var standing domain.GroupStanding
		standing.Group = group
		for rank, r := range rows {
			standing.Rows = append(standing.Rows, domain.StandingRow{
				Rank: rank + 1, Team: findTeam(r.team),
				Played: r.p, Won: r.w, Drawn: r.d, Lost: r.l,
				GoalsFor: r.gf, GoalsAgainst: r.ga, GoalDiff: r.gf - r.ga,
				Points: r.pts, Form: r.form,
			})
		}
		p.standings = append(p.standings, standing)
	}
}

func intPtr(v int) *int { return &v }

func matchQuery(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (p *Provider) findTeam(name string) (*domain.Team, error) {
	q := matchQuery(name)
	for i := range p.teams {
		t := &p.teams[i]
		if strings.Contains(matchQuery(t.Name), q) || strings.EqualFold(t.Code, name) {
			return t, nil
		}
	}
	return nil, fmt.Errorf("team not found: %s", name)
}

func (p *Provider) GetTeams(_ context.Context) ([]domain.Team, error) {
	return p.teams, nil
}

func (p *Provider) GetTeam(_ context.Context, name string) (*domain.Team, error) {
	return p.findTeam(name)
}

func (p *Provider) GetSquad(_ context.Context, teamName string) ([]domain.Player, error) {
	team, err := p.findTeam(teamName)
	if err != nil {
		return nil, err
	}
	var squad []domain.Player
	for _, pl := range p.players {
		if pl.TeamID == team.ID {
			squad = append(squad, pl)
		}
	}
	sort.Slice(squad, func(i, j int) bool {
		if squad[i].Number != squad[j].Number {
			return squad[i].Number < squad[j].Number
		}
		return squad[i].Name < squad[j].Name
	})
	return squad, nil
}

func (p *Provider) GetPlayer(_ context.Context, name string) (*domain.Player, error) {
	q := matchQuery(name)
	for i := range p.players {
		if strings.Contains(matchQuery(p.players[i].Name), q) {
			return &p.players[i], nil
		}
	}
	return nil, fmt.Errorf("player not found: %s", name)
}

func (p *Provider) GetMatches(_ context.Context) ([]domain.Match, error) {
	return p.matches, nil
}

func (p *Provider) GetMatchesToday(_ context.Context) ([]domain.Match, error) {
	today := time.Now()
	var out []domain.Match
	for _, m := range p.matches {
		if sameDay(m.Date, today) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (p *Provider) GetUpcoming(_ context.Context, limit int) ([]domain.Match, error) {
	now := time.Now()
	var upcoming []domain.Match
	for _, m := range p.matches {
		if m.Status == domain.StatusScheduled && m.Date.After(now) {
			upcoming = append(upcoming, m)
		}
	}
	sort.Slice(upcoming, func(i, j int) bool { return upcoming[i].Date.Before(upcoming[j].Date) })
	if limit > 0 && len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}
	return upcoming, nil
}

func (p *Provider) GetResults(_ context.Context) ([]domain.Match, error) {
	var out []domain.Match
	for _, m := range p.matches {
		if m.Status == domain.StatusFinished {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date.After(out[j].Date) })
	return out, nil
}

func (p *Provider) GetStandings(_ context.Context) ([]domain.GroupStanding, error) {
	return p.standings, nil
}

func (p *Provider) GetBracket(_ context.Context) (*domain.Bracket, error) {
	now := time.Now()
	return &domain.Bracket{
		Quarter: []domain.BracketMatch{
			{ID: 101, Stage: "Quarter-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 14)},
			{ID: 102, Stage: "Quarter-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 15)},
			{ID: 103, Stage: "Quarter-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 16)},
			{ID: 104, Stage: "Quarter-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 17)},
		},
		Semi: []domain.BracketMatch{
			{ID: 105, Stage: "Semi-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 21)},
			{ID: 106, Stage: "Semi-Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 22)},
		},
		ThirdPlace: domain.BracketMatch{ID: 107, Stage: "Third Place", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 25)},
		Final:      domain.BracketMatch{ID: 108, Stage: "Final", HomeTeam: "TBD", AwayTeam: "TBD", Date: now.AddDate(0, 0, 26)},
	}, nil
}

func (p *Provider) GetStats(_ context.Context) (*domain.TournamentStats, error) {
	var scorers []domain.PlayerStat
	for _, pl := range p.players {
		if pl.Goals > 0 {
			scorers = append(scorers, domain.PlayerStat{Player: pl, Value: pl.Goals})
		}
	}
	sort.Slice(scorers, func(i, j int) bool { return scorers[i].Value > scorers[j].Value })
	for i := range scorers {
		scorers[i].Rank = i + 1
	}
	if len(scorers) > 10 {
		scorers = scorers[:10]
	}

	var assists []domain.PlayerStat
	for _, pl := range p.players {
		if pl.Assists > 0 {
			assists = append(assists, domain.PlayerStat{Player: pl, Value: pl.Assists})
		}
	}
	sort.Slice(assists, func(i, j int) bool { return assists[i].Value > assists[j].Value })
	for i := range assists {
		assists[i].Rank = i + 1
	}
	if len(assists) > 10 {
		assists = assists[:10]
	}

	return &domain.TournamentStats{TopScorers: scorers, TopAssists: assists}, nil
}

func (p *Provider) GetTournamentInfo(_ context.Context) (*domain.TournamentInfo, error) {
	finished := 0
	for _, m := range p.matches {
		if m.Status == domain.StatusFinished {
			finished++
		}
	}
	return &domain.TournamentInfo{
		Name: "FIFA World Cup", Season: 2026,
		StartDate: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
		Host: "USA, Canada & Mexico", TotalTeams: len(p.teams),
		TotalMatches: 104, Completed: finished, CurrentStage: "Group Stage",
	}, nil
}

func (p *Provider) Search(_ context.Context, query string) ([]domain.SearchResult, error) {
	q := matchQuery(query)
	if q == "" {
		return nil, nil
	}
	var results []domain.SearchResult
	for _, t := range p.teams {
		if strings.Contains(matchQuery(t.Name), q) || strings.Contains(strings.ToLower(t.Code), q) {
			results = append(results, domain.SearchResult{Type: "team", Title: t.Name, Subtitle: "Group " + t.Group, ID: t.ID})
		}
	}
	for _, pl := range p.players {
		if strings.Contains(matchQuery(pl.Name), q) {
			results = append(results, domain.SearchResult{Type: "player", Title: pl.Name, Subtitle: pl.TeamName + " · " + pl.Position, ID: pl.ID})
		}
	}
	for _, m := range p.matches {
		if strings.Contains(matchQuery(m.HomeTeam.Name), q) || strings.Contains(matchQuery(m.AwayTeam.Name), q) ||
			strings.Contains(m.Date.Format("2006-01-02"), q) {
			results = append(results, domain.SearchResult{
				Type: "match", Title: fmt.Sprintf("%s vs %s", m.HomeTeam.Name, m.AwayTeam.Name),
				Subtitle: m.Date.Format("Jan 2, 15:04") + " · " + m.Stage, ID: m.ID,
			})
		}
	}
	return results, nil
}

func (p *Provider) GetHeadToHead(_ context.Context, teamA, teamB string) (*domain.HeadToHead, error) {
	a, err := p.findTeam(teamA)
	if err != nil {
		return nil, err
	}
	b, err := p.findTeam(teamB)
	if err != nil {
		return nil, err
	}
	var matches []domain.Match
	winsA, winsB, draws := 0, 0, 0
	for _, m := range p.matches {
		if (m.HomeTeam.ID == a.ID && m.AwayTeam.ID == b.ID) || (m.HomeTeam.ID == b.ID && m.AwayTeam.ID == a.ID) {
			matches = append(matches, m)
			if m.Status != domain.StatusFinished || m.HomeScore == nil || m.AwayScore == nil {
				continue
			}
			hs, as := *m.HomeScore, *m.AwayScore
			if hs == as {
				draws++
			} else if (m.HomeTeam.ID == a.ID && hs > as) || (m.AwayTeam.ID == a.ID && as > hs) {
				winsA++
			} else {
				winsB++
			}
		}
	}
	return &domain.HeadToHead{TeamA: *a, TeamB: *b, Matches: matches, WinsA: winsA, WinsB: winsB, Draws: draws}, nil
}

func (p *Provider) GetTeamForm(_ context.Context, teamName string) (*domain.TeamForm, error) {
	team, err := p.findTeam(teamName)
	if err != nil {
		return nil, err
	}
	var results, scores []string
	for _, m := range p.matches {
		if m.Status != domain.StatusFinished || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		hs, as := *m.HomeScore, *m.AwayScore
		if m.HomeTeam.ID == team.ID {
			scores = append(scores, fmt.Sprintf("%d-%d", hs, as))
			switch {
			case hs > as:
				results = append(results, "W")
			case hs < as:
				results = append(results, "L")
			default:
				results = append(results, "D")
			}
		} else if m.AwayTeam.ID == team.ID {
			scores = append(scores, fmt.Sprintf("%d-%d", as, hs))
			switch {
			case as > hs:
				results = append(results, "W")
			case as < hs:
				results = append(results, "L")
			default:
				results = append(results, "D")
			}
		}
	}
	if len(results) > 5 {
		results = results[len(results)-5:]
		scores = scores[len(scores)-5:]
	}
	return &domain.TeamForm{Team: *team, Results: results, Scores: scores}, nil
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
