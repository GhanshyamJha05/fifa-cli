package apifootball

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
)

// Client talks to API-Football (api-sports.io).
type Client struct {
	client   *resty.Client
	leagueID int
	season   int
	logger   *slog.Logger
}

// New creates an API-Football client.
func New(baseURL, apiKey string, leagueID, season int, logger *slog.Logger) *Client {
	c := resty.New().
		SetBaseURL(baseURL).
		SetHeader("x-apisports-key", apiKey).
		SetTimeout(15 * time.Second)
	return &Client{client: c, leagueID: leagueID, season: season, logger: logger}
}

type apiResponse[T any] struct {
	Response T `json:"response"`
	Errors   map[string]string `json:"errors"`
}

type apiTeam struct {
	Team struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"team"`
}

type apiFixture struct {
	Fixture struct {
		ID       int    `json:"id"`
		Date     string `json:"date"`
		Status   struct {
			Short string `json:"short"`
			Elapsed int  `json:"elapsed"`
		} `json:"status"`
		Venue struct {
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"venue"`
		Referee string `json:"referee"`
	} `json:"fixture"`
	League struct {
		Round string `json:"round"`
	} `json:"league"`
	Teams struct {
		Home struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	} `json:"goals"`
}

type apiStanding struct {
	League struct {
		Standings [][]struct {
			Rank int `json:"rank"`
			Team struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"team"`
			All struct {
				Played int `json:"played"`
				Win    int `json:"win"`
				Draw   int `json:"draw"`
				Lose   int `json:"lose"`
				Goals  struct {
					For     int `json:"for"`
					Against int `json:"against"`
				} `json:"goals"`
			} `json:"all"`
			GoalsDiff int    `json:"goalsDiff"`
			Points    int    `json:"points"`
			Form      string `json:"form"`
			Group     string `json:"group"`
		} `json:"standings"`
	} `json:"league"`
}

type apiPlayer struct {
	Player struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	} `json:"player"`
	Statistics []struct {
		Games struct {
			Number  int    `json:"number"`
			Minutes int    `json:"minutes"`
			Position string `json:"position"`
			Captain bool   `json:"captain"`
		} `json:"games"`
		Team struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"team"`
		Goals struct {
			Total    *int `json:"total"`
			Assists  *int `json:"assists"`
		} `json:"goals"`
		Cards struct {
			Yellow int `json:"yellow"`
			Red    int `json:"red"`
		} `json:"cards"`
	} `json:"statistics"`
}

func (c *Client) checkErr(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("api request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("api error: status %d", resp.StatusCode())
	}
	return nil
}

func parseStatus(short string) domain.MatchStatus {
	switch short {
	case "NS", "TBD":
		return domain.StatusScheduled
	case "1H", "2H", "HT", "ET", "P", "LIVE":
		return domain.StatusLive
	case "FT", "AET", "PEN":
		return domain.StatusFinished
	case "PST", "CANC":
		return domain.StatusPostponed
	default:
		return domain.StatusScheduled
	}
}

func (c *Client) fetchFixtures(ctx context.Context, params map[string]string) ([]domain.Match, error) {
	p := map[string]string{
		"league": fmt.Sprintf("%d", c.leagueID),
		"season": fmt.Sprintf("%d", c.season),
	}
	for k, v := range params {
		p[k] = v
	}

	var result apiResponse[[]apiFixture]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(p).
		SetResult(&result).
		Get("/fixtures")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	var matches []domain.Match
	for _, f := range result.Response {
		dt, _ := time.Parse(time.RFC3339, f.Fixture.Date)
		m := domain.Match{
			ID: f.Fixture.ID,
			HomeTeam: domain.Team{ID: f.Teams.Home.ID, Name: f.Teams.Home.Name},
			AwayTeam: domain.Team{ID: f.Teams.Away.ID, Name: f.Teams.Away.Name},
			HomeScore: f.Goals.Home, AwayScore: f.Goals.Away,
			Date: dt, Stage: f.League.Round,
			Stadium: f.Fixture.Venue.Name, City: f.Fixture.Venue.City,
			Status: parseStatus(f.Fixture.Status.Short),
			Minute: f.Fixture.Status.Elapsed, Referee: f.Fixture.Referee,
		}
		matches = append(matches, m)
	}
	return matches, nil
}

func (c *Client) GetTeams(ctx context.Context) ([]domain.Team, error) {
	var result apiResponse[[]apiTeam]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"league": fmt.Sprintf("%d", c.leagueID),
			"season": fmt.Sprintf("%d", c.season),
		}).
		SetResult(&result).
		Get("/teams")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	var teams []domain.Team
	for _, t := range result.Response {
		teams = append(teams, domain.Team{
			ID: t.Team.ID, Name: t.Team.Name, Code: t.Team.Code,
		})
	}
	return teams, nil
}

func (c *Client) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	teams, err := c.GetTeams(ctx)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(name)
	for i := range teams {
		if strings.Contains(strings.ToLower(teams[i].Name), q) || strings.EqualFold(teams[i].Code, name) {
			return &teams[i], nil
		}
	}
	return nil, fmt.Errorf("team not found: %s", name)
}

func (c *Client) GetTeamByID(ctx context.Context, id int) (*domain.Team, error) {
	teams, err := c.GetTeams(ctx)
	if err != nil {
		return nil, err
	}
	for i := range teams {
		if teams[i].ID == id {
			return &teams[i], nil
		}
	}
	return nil, fmt.Errorf("team not found: id %d", id)
}

func (c *Client) GetPlayerByID(ctx context.Context, id int) (*domain.Player, error) {
	var result apiResponse[[]apiPlayer]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"id": fmt.Sprintf("%d", id),
		}).
		SetResult(&result).
		Get("/players")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}
	if len(result.Response) == 0 {
		return nil, fmt.Errorf("player not found: id %d", id)
	}
	return c.mapPlayer(result.Response[0])
}

func (c *Client) mapPlayer(p apiPlayer) (*domain.Player, error) {
	if len(p.Statistics) == 0 {
		return &domain.Player{ID: p.Player.ID, Name: p.Player.Name, Age: p.Player.Age}, nil
	}
	s := p.Statistics[0]
	goals, assists := 0, 0
	if s.Goals.Total != nil {
		goals = *s.Goals.Total
	}
	if s.Goals.Assists != nil {
		assists = *s.Goals.Assists
	}
	return &domain.Player{
		ID: p.Player.ID, Name: p.Player.Name, Age: p.Player.Age,
		TeamID: s.Team.ID, TeamName: s.Team.Name,
		Position: s.Games.Position, Number: s.Games.Number,
		Captain: s.Games.Captain, Goals: goals, Assists: assists,
		Yellow: s.Cards.Yellow, Red: s.Cards.Red, Minutes: s.Games.Minutes,
	}, nil
}

func (c *Client) GetSquad(ctx context.Context, teamName string) ([]domain.Player, error) {
	team, err := c.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	var result apiResponse[[]apiPlayer]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"team":   fmt.Sprintf("%d", team.ID),
			"season": fmt.Sprintf("%d", c.season),
		}).
		SetResult(&result).
		Get("/players/squads")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	var players []domain.Player
	for _, entry := range result.Response {
		for _, s := range entry.Statistics {
			goals, assists := 0, 0
			if s.Goals.Total != nil {
				goals = *s.Goals.Total
			}
			if s.Goals.Assists != nil {
				assists = *s.Goals.Assists
			}
			players = append(players, domain.Player{
				ID: entry.Player.ID, Name: entry.Player.Name,
				TeamID: team.ID, TeamName: team.Name,
				Position: s.Games.Position, Number: s.Games.Number,
				Age: entry.Player.Age, Captain: s.Games.Captain,
				Goals: goals, Assists: assists,
				Yellow: s.Cards.Yellow, Red: s.Cards.Red,
				Minutes: s.Games.Minutes,
			})
		}
	}
	return players, nil
}

func (c *Client) GetPlayer(ctx context.Context, name string) (*domain.Player, error) {
	var result apiResponse[[]apiPlayer]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"league": fmt.Sprintf("%d", c.leagueID),
			"season": fmt.Sprintf("%d", c.season),
			"search": name,
		}).
		SetResult(&result).
		Get("/players")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}
	if len(result.Response) == 0 {
		return nil, fmt.Errorf("player not found: %s", name)
	}
	p := result.Response[0]
	if len(p.Statistics) == 0 {
		return &domain.Player{ID: p.Player.ID, Name: p.Player.Name, Age: p.Player.Age}, nil
	}
	s := p.Statistics[0]
	goals, assists := 0, 0
	if s.Goals.Total != nil {
		goals = *s.Goals.Total
	}
	if s.Goals.Assists != nil {
		assists = *s.Goals.Assists
	}
	return &domain.Player{
		ID: p.Player.ID, Name: p.Player.Name, Age: p.Player.Age,
		TeamID: s.Team.ID, TeamName: s.Team.Name,
		Position: s.Games.Position, Number: s.Games.Number,
		Captain: s.Games.Captain, Goals: goals, Assists: assists,
		Yellow: s.Cards.Yellow, Red: s.Cards.Red, Minutes: s.Games.Minutes,
	}, nil
}

func (c *Client) GetMatches(ctx context.Context) ([]domain.Match, error) {
	return c.fetchFixtures(ctx, nil)
}

func (c *Client) GetMatchesToday(ctx context.Context) ([]domain.Match, error) {
	return c.fetchFixtures(ctx, map[string]string{
		"date": time.Now().Format("2006-01-02"),
	})
}

func (c *Client) GetUpcoming(ctx context.Context, limit int) ([]domain.Match, error) {
	matches, err := c.GetMatches(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	var upcoming []domain.Match
	for _, m := range matches {
		if m.Status == domain.StatusScheduled && m.Date.After(now) {
			upcoming = append(upcoming, m)
		}
	}
	if limit > 0 && len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}
	return upcoming, nil
}

func (c *Client) GetResults(ctx context.Context) ([]domain.Match, error) {
	matches, err := c.GetMatches(ctx)
	if err != nil {
		return nil, err
	}
	var out []domain.Match
	for _, m := range matches {
		if m.Status == domain.StatusFinished {
			out = append(out, m)
		}
	}
	return out, nil
}

func (c *Client) GetStandings(ctx context.Context) ([]domain.GroupStanding, error) {
	var result apiResponse[[]apiStanding]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"league": fmt.Sprintf("%d", c.leagueID),
			"season": fmt.Sprintf("%d", c.season),
		}).
		SetResult(&result).
		Get("/standings")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	var standings []domain.GroupStanding
	for _, block := range result.Response {
		for _, group := range block.League.Standings {
			if len(group) == 0 {
				continue
			}
			gs := domain.GroupStanding{Group: strings.TrimPrefix(group[0].Group, "Group ")}
			for _, row := range group {
				gs.Rows = append(gs.Rows, domain.StandingRow{
					Rank: row.Rank,
					Team: domain.Team{ID: row.Team.ID, Name: row.Team.Name},
					Played: row.All.Played, Won: row.All.Win, Drawn: row.All.Draw,
					Lost: row.All.Lose, GoalsFor: row.All.Goals.For,
					GoalsAgainst: row.All.Goals.Against, GoalDiff: row.GoalsDiff,
					Points: row.Points, Form: row.Form,
				})
			}
			standings = append(standings, gs)
		}
	}
	return standings, nil
}

func (c *Client) GetBracket(ctx context.Context) (*domain.Bracket, error) {
	matches, err := c.GetMatches(ctx)
	if err != nil {
		return nil, err
	}
	b := &domain.Bracket{}
	for _, m := range matches {
		stage := strings.ToLower(m.Stage)
		bm := domain.BracketMatch{
			ID: m.ID, Stage: m.Stage,
			HomeTeam: m.HomeTeam.Name, AwayTeam: m.AwayTeam.Name,
			HomeScore: m.HomeScore, AwayScore: m.AwayScore, Date: m.Date,
		}
		switch {
		case strings.Contains(stage, "round of 16"):
			b.RoundOf16 = append(b.RoundOf16, bm)
		case strings.Contains(stage, "quarter"):
			b.Quarter = append(b.Quarter, bm)
		case strings.Contains(stage, "semi"):
			b.Semi = append(b.Semi, bm)
		case strings.Contains(stage, "3rd") || strings.Contains(stage, "third"):
			b.ThirdPlace = bm
		case strings.Contains(stage, "final"):
			b.Final = bm
		}
	}
	return b, nil
}

func (c *Client) GetStats(ctx context.Context) (*domain.TournamentStats, error) {
	scorers, err := c.fetchTop(ctx, "goals")
	if err != nil {
		return nil, err
	}
	assists, err := c.fetchTop(ctx, "assists")
	if err != nil {
		return nil, err
	}
	return &domain.TournamentStats{TopScorers: scorers, TopAssists: assists}, nil
}

func (c *Client) fetchTop(ctx context.Context, stat string) ([]domain.PlayerStat, error) {
	var result apiResponse[[]struct {
		Player struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"player"`
		Statistics []struct {
			Goals struct {
				Total   *int `json:"total"`
				Assists *int `json:"assists"`
			} `json:"goals"`
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
		} `json:"statistics"`
	}]

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"league": fmt.Sprintf("%d", c.leagueID),
			"season": fmt.Sprintf("%d", c.season),
		}).
		SetResult(&result).
		Get("/players/top" + stat)
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	var stats []domain.PlayerStat
	for i, entry := range result.Response {
		if len(entry.Statistics) == 0 {
			continue
		}
		s := entry.Statistics[0]
		val := 0
		if stat == "goals" && s.Goals.Total != nil {
			val = *s.Goals.Total
		}
		if stat == "assists" && s.Goals.Assists != nil {
			val = *s.Goals.Assists
		}
		stats = append(stats, domain.PlayerStat{
			Rank: i + 1,
			Player: domain.Player{ID: entry.Player.ID, Name: entry.Player.Name, TeamName: s.Team.Name},
			Value: val,
		})
	}
	return stats, nil
}

func (c *Client) GetTournamentInfo(ctx context.Context) (*domain.TournamentInfo, error) {
	matches, err := c.GetMatches(ctx)
	if err != nil {
		return nil, err
	}
	completed := 0
	for _, m := range matches {
		if m.Status == domain.StatusFinished {
			completed++
		}
	}
	return &domain.TournamentInfo{
		Name: "FIFA World Cup", Season: c.season,
		StartDate: time.Date(c.season, 6, 11, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(c.season, 7, 19, 0, 0, 0, 0, time.UTC),
		Host: "USA, Canada & Mexico", TotalMatches: len(matches),
		Completed: completed, CurrentStage: "Group Stage",
	}, nil
}

func (c *Client) Search(ctx context.Context, query string) ([]domain.SearchResult, error) {
	var results []domain.SearchResult
	teams, _ := c.GetTeams(ctx)
	q := strings.ToLower(query)
	for _, t := range teams {
		if strings.Contains(strings.ToLower(t.Name), q) {
			results = append(results, domain.SearchResult{Type: "team", Title: t.Name, ID: t.ID})
		}
	}
	matches, _ := c.GetMatches(ctx)
	for _, m := range matches {
		if strings.Contains(strings.ToLower(m.HomeTeam.Name), q) || strings.Contains(strings.ToLower(m.AwayTeam.Name), q) {
			results = append(results, domain.SearchResult{
				Type: "match", Title: fmt.Sprintf("%s vs %s", m.HomeTeam.Name, m.AwayTeam.Name),
				Subtitle: m.Date.Format("Jan 2"), ID: m.ID,
			})
		}
	}
	return results, nil
}

func (c *Client) GetHeadToHead(ctx context.Context, teamA, teamB string) (*domain.HeadToHead, error) {
	a, err := c.GetTeam(ctx, teamA)
	if err != nil {
		return nil, err
	}
	b, err := c.GetTeam(ctx, teamB)
	if err != nil {
		return nil, err
	}

	var result apiResponse[[]apiFixture]
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"h2h": fmt.Sprintf("%d-%d", a.ID, b.ID),
		}).
		SetResult(&result).
		Get("/fixtures/headtohead")
	if err := c.checkErr(resp, err); err != nil {
		return nil, err
	}

	h2h := &domain.HeadToHead{TeamA: *a, TeamB: *b}
	for _, f := range result.Response {
		dt, _ := time.Parse(time.RFC3339, f.Fixture.Date)
		m := domain.Match{
			ID: f.Fixture.ID,
			HomeTeam: domain.Team{ID: f.Teams.Home.ID, Name: f.Teams.Home.Name},
			AwayTeam: domain.Team{ID: f.Teams.Away.ID, Name: f.Teams.Away.Name},
			HomeScore: f.Goals.Home, AwayScore: f.Goals.Away, Date: dt,
			Status: parseStatus(f.Fixture.Status.Short),
		}
		h2h.Matches = append(h2h.Matches, m)
	}
	return h2h, nil
}

func (c *Client) GetTeamForm(ctx context.Context, teamName string) (*domain.TeamForm, error) {
	team, err := c.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}
	matches, err := c.fetchFixtures(ctx, map[string]string{
		"team": fmt.Sprintf("%d", team.ID),
		"last": "5",
	})
	if err != nil {
		return nil, err
	}
	form := &domain.TeamForm{Team: *team}
	for _, m := range matches {
		if m.Status != domain.StatusFinished || m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		hs, as := *m.HomeScore, *m.AwayScore
		if m.HomeTeam.ID == team.ID {
			form.Scores = append(form.Scores, fmt.Sprintf("%d-%d", hs, as))
			form.Results = append(form.Results, resultChar(hs, as))
		} else {
			form.Scores = append(form.Scores, fmt.Sprintf("%d-%d", as, hs))
			form.Results = append(form.Results, resultChar(as, hs))
		}
	}
	return form, nil
}

func resultChar(forTeam, against int) string {
	switch {
	case forTeam > against:
		return "W"
	case forTeam < against:
		return "L"
	default:
		return "D"
	}
}
