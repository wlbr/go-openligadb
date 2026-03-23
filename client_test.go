package openligadb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClientDefaults(t *testing.T) {
	c := NewClient()
	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
	}
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := NewClient(WithBaseURL("https://example.com"), WithHTTPClient(custom))
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
	}
	if c.httpClient != custom {
		t.Error("httpClient was not set by WithHTTPClient")
	}
}

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient(WithBaseURL(srv.URL))
	return srv, c
}

func TestGetAvailableLeagues(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getavailableleagues" {
			t.Errorf("path = %q, want /getavailableleagues", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept = %q, want application/json", r.Header.Get("Accept"))
		}
		json.NewEncoder(w).Encode([]League{
			{LeagueID: 1, LeagueName: "Bundesliga", LeagueShortcut: "bl1", LeagueSeason: "2025"},
		})
	})

	leagues, err := c.GetAvailableLeagues(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(leagues) != 1 {
		t.Fatalf("len = %d, want 1", len(leagues))
	}
	if leagues[0].LeagueShortcut != "bl1" {
		t.Errorf("LeagueShortcut = %q, want bl1", leagues[0].LeagueShortcut)
	}
}

func TestGetAvailableSports(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getavailablesports" {
			t.Errorf("path = %q, want /getavailablesports", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Sport{{SportID: 1, SportName: "Fussball"}})
	})

	sports, err := c.GetAvailableSports(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(sports) != 1 || sports[0].SportName != "Fussball" {
		t.Errorf("unexpected sports: %+v", sports)
	}
}

func TestGetMatch(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchdata/12345" {
			t.Errorf("path = %q, want /getmatchdata/12345", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Match{
			MatchID:         12345,
			MatchIsFinished: true,
			Team1:           Team{TeamName: "Team A"},
			Team2:           Team{TeamName: "Team B"},
		})
	})

	m, err := c.GetMatch(context.Background(), 12345)
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 12345 {
		t.Errorf("MatchID = %d, want 12345", m.MatchID)
	}
	if !m.MatchIsFinished {
		t.Error("MatchIsFinished = false, want true")
	}
}

func TestGetMatchesByLeagueSeason(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchdata/bl1/2025" {
			t.Errorf("path = %q, want /getmatchdata/bl1/2025", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 1}, {MatchID: 2}})
	})

	matches, err := c.GetMatchesByLeagueSeason(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Fatalf("len = %d, want 2", len(matches))
	}
}

func TestGetMatchesByLeagueSeasonGroup(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchdata/bl1/2025/3" {
			t.Errorf("path = %q, want /getmatchdata/bl1/2025/3", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 10}})
	})

	matches, err := c.GetMatchesByLeagueSeasonGroup(context.Background(), "bl1", 2025, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("len = %d, want 1", len(matches))
	}
}

func TestGetMatchesByLeagueSeasonTeam(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchdata/bl1/2025/Bayern" {
			t.Errorf("path = %q, want /getmatchdata/bl1/2025/Bayern", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 20}})
	})

	matches, err := c.GetMatchesByLeagueSeasonTeam(context.Background(), "bl1", 2025, "Bayern")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("len = %d, want 1", len(matches))
	}
}

func TestGetMatchesByTeamIDs(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchdata/40/7" {
			t.Errorf("path = %q, want /getmatchdata/40/7", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 99}})
	})

	matches, err := c.GetMatchesByTeamIDs(context.Background(), 40, 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("len = %d, want 1", len(matches))
	}
}

func TestGetLastChangeDate(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getlastchangedate/bl1/2025/1" {
			t.Errorf("path = %q, want /getlastchangedate/bl1/2025/1", r.URL.Path)
		}
		json.NewEncoder(w).Encode("2025-09-15T18:30:00Z")
	})

	dt, err := c.GetLastChangeDate(context.Background(), "bl1", 2025, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2025, 9, 15, 18, 30, 0, 0, time.UTC)
	if !dt.Equal(want) {
		t.Errorf("date = %v, want %v", dt, want)
	}
}

func TestGetNextMatchByLeagueTeam(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getnextmatchbyleagueteam/100/40" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Match{MatchID: 50})
	})

	m, err := c.GetNextMatchByLeagueTeam(context.Background(), 100, 40)
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 50 {
		t.Errorf("MatchID = %d, want 50", m.MatchID)
	}
}

func TestGetNextMatchByLeagueShortcut(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getnextmatchbyleagueshortcut/bl1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Match{MatchID: 51})
	})

	m, err := c.GetNextMatchByLeagueShortcut(context.Background(), "bl1")
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 51 {
		t.Errorf("MatchID = %d, want 51", m.MatchID)
	}
}

func TestGetLastMatchByLeagueShortcut(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getlastmatchbyleagueshortcut/bl2" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Match{MatchID: 60})
	})

	m, err := c.GetLastMatchByLeagueShortcut(context.Background(), "bl2")
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 60 {
		t.Errorf("MatchID = %d, want 60", m.MatchID)
	}
}

func TestGetLastMatchByLeagueTeam(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getlastmatchbyleagueteam/100/7" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Match{MatchID: 61})
	})

	m, err := c.GetLastMatchByLeagueTeam(context.Background(), 100, 7)
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 61 {
		t.Errorf("MatchID = %d, want 61", m.MatchID)
	}
}

func TestGetCurrentGroup(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getcurrentgroup/bl1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Group{GroupName: "5. Spieltag", GroupOrderID: 5, GroupID: 42})
	})

	g, err := c.GetCurrentGroup(context.Background(), "bl1")
	if err != nil {
		t.Fatal(err)
	}
	if g.GroupOrderID != 5 {
		t.Errorf("GroupOrderID = %d, want 5", g.GroupOrderID)
	}
}

func TestGetAvailableGroups(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getavailablegroups/bl1/2025" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Group{{GroupOrderID: 1}, {GroupOrderID: 2}})
	})

	groups, err := c.GetAvailableGroups(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("len = %d, want 2", len(groups))
	}
}

func TestGetGoalGetters(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getgoalgetters/bl1/2025" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]GoalGetter{
			{GoalGetterID: 1, GoalGetterName: "Harry Kane", GoalCount: 20},
		})
	})

	gg, err := c.GetGoalGetters(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(gg) != 1 || gg[0].GoalCount != 20 {
		t.Errorf("unexpected: %+v", gg)
	}
}

func TestGetAvailableTeams(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getavailableteams/bl1/2025" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Team{
			{TeamID: 40, TeamName: "FC Bayern"},
			{TeamID: 7, TeamName: "BVB"},
		})
	})

	teams, err := c.GetAvailableTeams(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 2 {
		t.Fatalf("len = %d, want 2", len(teams))
	}
}

func TestGetBlTable(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getbltable/bl1/2025" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]BlTableTeam{
			{TeamName: "Bayern", Points: 55},
		})
	})

	table, err := c.GetBlTable(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 1 || table[0].Points != 55 {
		t.Errorf("unexpected table: %+v", table)
	}
}

func TestGetGroupTable(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getgrouptable/bl1/2025" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]BlTableTeam{{TeamName: "BVB", Points: 40}})
	})

	table, err := c.GetGroupTable(context.Background(), "bl1", 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(table) != 1 {
		t.Fatalf("len = %d, want 1", len(table))
	}
}

func TestGetMatchesByTeam(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchesbyteam/Bayern/4/4" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 70}})
	})

	matches, err := c.GetMatchesByTeam(context.Background(), "Bayern", 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("len = %d, want 1", len(matches))
	}
}

func TestGetMatchesByTeamID(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getmatchesbyteamid/40/2/2" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Match{{MatchID: 80}})
	})

	matches, err := c.GetMatchesByTeamID(context.Background(), 40, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("len = %d, want 1", len(matches))
	}
}

func TestGetResultInfos(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getresultinfos/4500" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ResultInfo{ID: 1, Name: "Endergebnis"})
	})

	ri, err := c.GetResultInfos(context.Background(), 4500)
	if err != nil {
		t.Fatal(err)
	}
	if ri.Name != "Endergebnis" {
		t.Errorf("Name = %q, want Endergebnis", ri.Name)
	}
}

func TestHTTPError(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := c.GetAvailableLeagues(context.Background())
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}

func TestInvalidJSON(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})

	_, err := c.GetAvailableLeagues(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestContextCanceled(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]League{})
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetAvailableLeagues(ctx)
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestMatchJSONDeserialization(t *testing.T) {
	oneInt := 1
	twoInt := 2
	zeroInt := 0
	falseVal := false
	viewers := 75000
	minute := 23
	wantMatch := Match{
		MatchID:         100,
		MatchIsFinished: true,
		Team1:           Team{TeamID: 1, TeamName: "FC Bayern"},
		Team2:           Team{TeamID: 2, TeamName: "BVB"},
		Goals: []Goal{{
			GoalID: 1, ScoreTeam1: &oneInt, ScoreTeam2: &zeroInt,
			MatchMinute: &minute, GoalGetterName: "Kane",
			IsPenalty: &falseVal, IsOwnGoal: &falseVal, IsOvertime: &falseVal,
		}},
		MatchResults: []MatchResult{{
			ResultID: 1, ResultName: "Endergebnis",
			PointsTeam1: &twoInt, PointsTeam2: &oneInt,
			ResultOrderID: 1, ResultTypeID: 2,
		}},
		Location:        &Location{LocationID: 10, LocationCity: "Munich", LocationStadium: "Allianz Arena"},
		NumberOfViewers: &viewers,
	}

	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(wantMatch)
	})

	m, err := c.GetMatch(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchID != 100 {
		t.Errorf("MatchID = %d, want 100", m.MatchID)
	}
	if m.Team1.TeamName != "FC Bayern" {
		t.Errorf("Team1 = %q, want FC Bayern", m.Team1.TeamName)
	}
	if len(m.Goals) != 1 {
		t.Fatalf("goals len = %d, want 1", len(m.Goals))
	}
	if m.Goals[0].GoalGetterName != "Kane" {
		t.Errorf("GoalGetterName = %q, want Kane", m.Goals[0].GoalGetterName)
	}
	if m.Goals[0].ScoreTeam1 == nil || *m.Goals[0].ScoreTeam1 != 1 {
		t.Error("ScoreTeam1 mismatch")
	}
	if m.Location == nil || m.Location.LocationStadium != "Allianz Arena" {
		t.Error("Location not decoded correctly")
	}
	if m.NumberOfViewers == nil || *m.NumberOfViewers != 75000 {
		t.Error("NumberOfViewers mismatch")
	}
	if len(m.MatchResults) != 1 || m.MatchResults[0].PointsTeam1 == nil || *m.MatchResults[0].PointsTeam1 != 2 {
		t.Error("MatchResults not decoded correctly")
	}
}

func TestMatchWithNullableFields(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Match{MatchID: 200})
	})

	m, err := c.GetMatch(context.Background(), 200)
	if err != nil {
		t.Fatal(err)
	}
	if m.MatchDateTime != nil {
		t.Errorf("MatchDateTime = %v, want nil", m.MatchDateTime)
	}
	if m.Location != nil {
		t.Errorf("Location = %v, want nil", m.Location)
	}
	if m.NumberOfViewers != nil {
		t.Errorf("NumberOfViewers = %v, want nil", m.NumberOfViewers)
	}
}
