package openligadb

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetLeagueByShortcut(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]League{
			{LeagueID: 1, LeagueShortcut: "bl1", LeagueName: "1. Bundesliga"},
			{LeagueID: 2, LeagueShortcut: "bl2", LeagueName: "2. Bundesliga"},
			{LeagueID: 3, LeagueShortcut: "bl1", LeagueName: "1. Bundesliga 2024"},
		})
	})

	leagues, err := c.GetLeagueByShortcut(context.Background(), "bl1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(leagues) != 2 {
		t.Errorf("expected 2 leagues, got %d", len(leagues))
	}

	_, err = c.GetLeagueByShortcut(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent shortcut, got nil")
	}
}

func TestGetLeagueByShortcutInSeason(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]League{
			{LeagueID: 1, LeagueShortcut: "bl1", LeagueSeason: "2025", LeagueName: "1. Bundesliga 2025"},
			{LeagueID: 2, LeagueShortcut: "bl1", LeagueSeason: "2024", LeagueName: "1. Bundesliga 2024"},
		})
	})

	leagues, err := c.GetLeagueByShortcutInSeason(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(leagues) != 1 || leagues[0].LeagueSeason != "2025" {
		t.Errorf("unexpected leagues: %+v", leagues)
	}

	_, err = c.GetLeagueByShortcutInSeason(context.Background(), "bl1", 2023)
	if err == nil {
		t.Error("expected error for nonexistent season, got nil")
	}
}

func TestGetTeamByName(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Team{
			{TeamID: 1, TeamName: "FC Bayern"},
			{TeamID: 2, TeamName: "BVB"},
		})
	})

	teams, err := c.GetTeamByName(context.Background(), "bl1", 2025, "FC Bayern")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 1 || teams[0].TeamID != 1 {
		t.Errorf("unexpected teams: %+v", teams)
	}

	_, err = c.GetTeamByName(context.Background(), "bl1", 2025, "Schalke")
	if err == nil {
		t.Error("expected error for nonexistent team, got nil")
	}
}

func TestGetTeamByID(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Team{
			{TeamID: 1, TeamName: "FC Bayern"},
			{TeamID: 2, TeamName: "BVB"},
		})
	})

	team, err := c.GetTeamByID(context.Background(), "bl1", 2025, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.TeamName != "BVB" {
		t.Errorf("expected BVB, got %s", team.TeamName)
	}

	_, err = c.GetTeamByID(context.Background(), "bl1", 2025, 99)
	if err == nil {
		t.Error("expected error for nonexistent team ID, got nil")
	}
}
