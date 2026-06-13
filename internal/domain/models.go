package domain

import "time"

// Team represents a World Cup participating nation.
type Team struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Flag        string `json:"flag"`
	Group       string `json:"group,omitempty"`
	Coach       string `json:"coach,omitempty"`
	FIFAWorldRank int  `json:"fifa_rank,omitempty"`
	Stadium     string `json:"stadium,omitempty"`
}

// Player represents a squad member.
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`
	Position string `json:"position"`
	Number   int    `json:"number"`
	Age      int    `json:"age"`
	Club     string `json:"club"`
	Captain  bool   `json:"captain"`
	Goals    int    `json:"goals"`
	Assists  int    `json:"assists"`
	Yellow   int    `json:"yellow"`
	Red      int    `json:"red"`
	Minutes  int    `json:"minutes"`
	Photo    string `json:"photo,omitempty"`
}

// MatchStatus represents the state of a fixture.
type MatchStatus string

const (
	StatusScheduled MatchStatus = "scheduled"
	StatusLive      MatchStatus = "live"
	StatusFinished  MatchStatus = "finished"
	StatusPostponed MatchStatus = "postponed"
)

// Match represents a World Cup fixture.
type Match struct {
	ID         int         `json:"id"`
	HomeTeam   Team        `json:"home_team"`
	AwayTeam   Team        `json:"away_team"`
	HomeScore  *int        `json:"home_score,omitempty"`
	AwayScore  *int        `json:"away_score,omitempty"`
	Date       time.Time   `json:"date"`
	Stage      string      `json:"stage"`
	Group      string      `json:"group,omitempty"`
	Stadium    string      `json:"stadium"`
	City       string      `json:"city"`
	Country    string      `json:"country"`
	Status     MatchStatus `json:"status"`
	Minute     int         `json:"minute,omitempty"`
	Referee    string      `json:"referee,omitempty"`
	Attendance int         `json:"attendance,omitempty"`
}

// StandingRow is a single row in a group table.
type StandingRow struct {
	Rank         int    `json:"rank"`
	Team         Team   `json:"team"`
	Played       int    `json:"played"`
	Won          int    `json:"won"`
	Drawn        int    `json:"drawn"`
	Lost         int    `json:"lost"`
	GoalsFor     int    `json:"goals_for"`
	GoalsAgainst int    `json:"goals_against"`
	GoalDiff     int    `json:"goal_diff"`
	Points       int    `json:"points"`
	Form         string `json:"form"`
}

// GroupStanding holds standings for one group.
type GroupStanding struct {
	Group string        `json:"group"`
	Rows  []StandingRow `json:"rows"`
}

// BracketMatch is a knockout stage fixture.
type BracketMatch struct {
	ID       int    `json:"id"`
	Stage    string `json:"stage"`
	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`
	HomeScore *int  `json:"home_score,omitempty"`
	AwayScore *int  `json:"away_score,omitempty"`
	Winner   string `json:"winner,omitempty"`
	Date     time.Time `json:"date"`
}

// Bracket holds knockout stage structure.
type Bracket struct {
	RoundOf32  []BracketMatch `json:"round_of_32,omitempty"`
	RoundOf16  []BracketMatch `json:"round_of_16,omitempty"`
	Quarter    []BracketMatch `json:"quarter"`
	Semi       []BracketMatch `json:"semi"`
	ThirdPlace BracketMatch   `json:"third_place"`
	Final      BracketMatch   `json:"final"`
}

// PlayerStat holds tournament statistics for a player.
type PlayerStat struct {
	Player Player `json:"player"`
	Value  int    `json:"value"`
	Rank   int    `json:"rank"`
}

// TournamentStats aggregates top performers.
type TournamentStats struct {
	TopScorers    []PlayerStat `json:"top_scorers"`
	TopAssists    []PlayerStat `json:"top_assists"`
	CleanSheets   []PlayerStat `json:"clean_sheets"`
	YellowCards   []PlayerStat `json:"yellow_cards"`
	RedCards      []PlayerStat `json:"red_cards"`
}

// TournamentInfo holds overall tournament progress.
type TournamentInfo struct {
	Name       string    `json:"name"`
	Season     int       `json:"season"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Host       string    `json:"host"`
	TotalTeams int       `json:"total_teams"`
	TotalMatches int     `json:"total_matches"`
	Completed  int       `json:"completed_matches"`
	CurrentStage string  `json:"current_stage"`
}

// SearchResult is a unified search hit.
type SearchResult struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Subtitle string `json:"subtitle"`
	ID      int    `json:"id"`
}

// HeadToHead compares two teams.
type HeadToHead struct {
	TeamA   Team    `json:"team_a"`
	TeamB   Team    `json:"team_b"`
	Matches []Match `json:"matches"`
	WinsA   int     `json:"wins_a"`
	WinsB   int     `json:"wins_b"`
	Draws   int     `json:"draws"`
}

// TeamForm holds recent match results.
type TeamForm struct {
	Team    Team     `json:"team"`
	Results []string `json:"results"` // W, D, L
	Scores  []string `json:"scores"`
}
