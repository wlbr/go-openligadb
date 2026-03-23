package openligadb

import (
	"fmt"
	"strings"
	"time"
)

// OpenLigaTime handles the specific "2025-08-22T20:30:00" format
type OpenLigaTime struct {
	time.Time
}

//const layout = "2006-01-02T15:04:05"

// UnmarshalJSON parses the OpenLigaDB date/time format into an OpenLigaTime.
func (ot *OpenLigaTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	layout := "2006-01-02T15:04:05"
	if s[len(s)-1] == 'Z' {
		layout += "Z"
	}
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("failed to parse OpenLigaTime: %w", err)
	}
	ot.Time = t
	return nil
}

// Sport represents a sport discipline (e.g. football).
type Sport struct {
	SportID   int    `json:"sportId"`
	SportName string `json:"sportName"`
}

// League represents a league in a specific season.
type League struct {
	LeagueID       int    `json:"leagueId"`
	LeagueName     string `json:"leagueName"`
	LeagueShortcut string `json:"leagueShortcut"`
	LeagueSeason   string `json:"leagueSeason"`
	Sport          Sport  `json:"sport"`
}

// Team represents a team participating in a league.
type Team struct {
	TeamID        int    `json:"teamId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconUrl"`
	TeamGroupName string `json:"teamGroupName"`
}

// Group represents a matchday or round within a league season.
type Group struct {
	GroupName    string `json:"groupName"`
	GroupOrderID int    `json:"groupOrderID"`
	GroupID      int    `json:"groupID"`
}

// Location represents the venue where a match is played.
type Location struct {
	LocationID      int    `json:"locationID"`
	LocationCity    string `json:"locationCity"`
	LocationStadium string `json:"locationStadium"`
}

// Goal represents a single goal scored during a match.
type Goal struct {
	GoalID         int    `json:"goalID"`
	ScoreTeam1     *int   `json:"scoreTeam1"`
	ScoreTeam2     *int   `json:"scoreTeam2"`
	MatchMinute    *int   `json:"matchMinute"`
	GoalGetterID   int    `json:"goalGetterID"`
	GoalGetterName string `json:"goalGetterName"`
	IsPenalty      *bool  `json:"isPenalty"`
	IsOwnGoal      *bool  `json:"isOwnGoal"`
	IsOvertime     *bool  `json:"isOvertime"`
	Comment        string `json:"comment"`
}

// MatchResult represents a result entry for a match (e.g. halftime, final).
type MatchResult struct {
	ResultID          int    `json:"resultID"`
	ResultName        string `json:"resultName"`
	PointsTeam1       *int   `json:"pointsTeam1"`
	PointsTeam2       *int   `json:"pointsTeam2"`
	ResultOrderID     int    `json:"resultOrderID"`
	ResultTypeID      int    `json:"resultTypeID"`
	ResultDescription string `json:"resultDescription"`
}

// Match represents a single match including teams, goals and results.
type Match struct {
	MatchID            int           `json:"matchID"`
	MatchDateTime      *OpenLigaTime `json:"matchDateTime"`
	TimeZoneID         string        `json:"timeZoneID"`
	LeagueID           int           `json:"leagueId"`
	LeagueName         string        `json:"leagueName"`
	LeagueSeason       int           `json:"leagueSeason"`
	LeagueShortcut     string        `json:"leagueShortcut"`
	MatchDateTimeUTC   *OpenLigaTime `json:"matchDateTimeUTC"`
	Group              Group         `json:"group"`
	Team1              Team          `json:"team1"`
	Team2              Team          `json:"team2"`
	LastUpdateDateTime *OpenLigaTime `json:"lastUpdateDateTime"`
	MatchIsFinished    bool          `json:"matchIsFinished"`
	MatchResults       []MatchResult `json:"matchResults"`
	Goals              []Goal        `json:"goals"`
	Location           *Location     `json:"location"`
	NumberOfViewers    *int          `json:"numberOfViewers"`
}

// GoalGetter represents a top scorer entry for a league season.
type GoalGetter struct {
	GoalGetterID   int    `json:"goalGetterId"`
	GoalGetterName string `json:"goalGetterName"`
	GoalCount      int    `json:"goalCount"`
}

// BlTableTeam represents a team's standing in the league table.
type BlTableTeam struct {
	TeamInfoID    int    `json:"teamInfoId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconUrl"`
	Points        int    `json:"points"`
	OpponentGoals int    `json:"opponentGoals"`
	Goals         int    `json:"goals"`
	Matches       int    `json:"matches"`
	Won           int    `json:"won"`
	Lost          int    `json:"lost"`
	Draw          int    `json:"draw"`
	GoalDiff      int    `json:"goalDiff"`
}

// GlobalResultInfo describes a global result type definition.
type GlobalResultInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ResultInfo describes the configured result types for a league.
type ResultInfo struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	OrderID          *int             `json:"orderId"`
	GlobalResultInfo GlobalResultInfo `json:"globalResultInfo"`
}
