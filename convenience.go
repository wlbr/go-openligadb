package openligadb

import (
	"context"
	"fmt"
	"strings"
)

// GetLeagueByShortcut returns all league entries matching the given shortcut (e.g. "bl1").
func (c *Client) GetLeagueByShortcut(ctx context.Context, shortcut string) ([]League, error) {
	leagues, err := c.GetAvailableLeagues(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching leagues: %w", err)
	}
	foundLeagues := []League{}
	for _, league := range leagues {
		if league.LeagueShortcut == shortcut {
			foundLeagues = append(foundLeagues, league)
		}
	}
	if len(foundLeagues) == 0 {
		return nil, fmt.Errorf("league with shortcut %s not found", shortcut)
	}
	return foundLeagues, nil
}

// GetLeagueByShortcutInSeason returns leagues matching the given shortcut and season year.
func (c *Client) GetLeagueByShortcutInSeason(ctx context.Context, shortcut string, season int) ([]League, error) {
	leagues, err := c.GetLeagueByShortcut(ctx, shortcut)
	if err != nil {
		return nil, fmt.Errorf("fetching league by shortcut: %w", err)
	}
	foundLeagues := []League{}
	for _, league := range leagues {
		if league.LeagueSeason == fmt.Sprintf("%d", season) {
			foundLeagues = append(foundLeagues, league)
		}
	}
	if len(foundLeagues) == 0 {
		return nil, fmt.Errorf("league with shortcut %s and season %d not found", shortcut, season)
	}
	if len(foundLeagues) > 1 {
		return foundLeagues, fmt.Errorf("ambiguous result for league shortcut %s season %d", shortcut, season)
	}
	return foundLeagues, nil
}

// GetTeamByName returns the team matching the exact name within a league and season.
func (c *Client) GetTeamByName(ctx context.Context, leagueShortcut string, leagueSeason int, teamname string) ([]Team, error) {
	teams, err := c.GetAvailableTeams(ctx, leagueShortcut, leagueSeason)
	if err != nil {
		return nil, fmt.Errorf("fetching teams: %w", err)
	}
	foundTeams := []Team{}
	teamnamelower := strings.ToLower(teamname)

	for _, t := range teams {
		tteamnamelower := strings.ToLower(t.TeamName)
		if strings.Contains(tteamnamelower, teamnamelower) {
			foundTeams = append(foundTeams, t)
		}
	}
	if len(foundTeams) == 0 {
		return nil, fmt.Errorf("team with name %s not found", teamname)
	}
	if len(foundTeams) > 1 {
		return foundTeams, fmt.Errorf("ambiguous result for team name %s", teamname)
	}
	return foundTeams, nil
}

// GetTeamByID returns the team matching the given ID within a league and season.
func (c *Client) GetTeamByID(ctx context.Context, leagueShortcut string, leagueSeason int, teamID int) (Team, error) {
	teams, err := c.GetAvailableTeams(ctx, leagueShortcut, leagueSeason)
	if err != nil {
		return Team{}, fmt.Errorf("fetching teams: %w", err)
	}
	for _, t := range teams {
		if t.TeamID == teamID {
			return t, nil
		}
	}

	return Team{}, fmt.Errorf("team with ID %d not found", teamID)
}
