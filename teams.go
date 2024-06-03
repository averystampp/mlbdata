package mlb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TeamRecordType struct {
	Copyright string `json:"copyright"`
	Teams     []struct {
		SpringLeague struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Link         string `json:"link"`
			Abbreviation string `json:"abbreviation"`
		} `json:"springLeague"`
		AllStarStatus string `json:"allStarStatus"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Link          string `json:"link"`
		Season        int    `json:"season"`
		Venue         struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"venue"`
		SpringVenue struct {
			ID   int    `json:"id"`
			Link string `json:"link"`
		} `json:"springVenue"`
		TeamCode        string `json:"teamCode"`
		FileCode        string `json:"fileCode"`
		Abbreviation    string `json:"abbreviation"`
		TeamName        string `json:"teamName"`
		LocationName    string `json:"locationName"`
		FirstYearOfPlay string `json:"firstYearOfPlay"`
		League          struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"league"`
		Division struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Link string `json:"link"`
		} `json:"division"`
		Sport struct {
			ID   int    `json:"id"`
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"sport"`
		ShortName string `json:"shortName"`
		Record    struct {
			Season string `json:"season"`
			Streak struct {
				StreakType   string `json:"streakType"`
				StreakNumber int    `json:"streakNumber"`
				StreakCode   string `json:"streakCode"`
			} `json:"streak"`
			DivisionRank          string `json:"divisionRank"`
			LeagueRank            string `json:"leagueRank"`
			WildCardRank          string `json:"wildCardRank"`
			SportRank             string `json:"sportRank"`
			GamesPlayed           int    `json:"gamesPlayed"`
			GamesBack             string `json:"gamesBack"`
			WildCardGamesBack     string `json:"wildCardGamesBack"`
			LeagueGamesBack       string `json:"leagueGamesBack"`
			SpringLeagueGamesBack string `json:"springLeagueGamesBack"`
			SportGamesBack        string `json:"sportGamesBack"`
			DivisionGamesBack     string `json:"divisionGamesBack"`
			ConferenceGamesBack   string `json:"conferenceGamesBack"`
			LeagueRecord          struct {
				Wins   int    `json:"wins"`
				Losses int    `json:"losses"`
				Ties   int    `json:"ties"`
				Pct    string `json:"pct"`
			} `json:"leagueRecord"`
			LastUpdated time.Time `json:"lastUpdated"`
			Records     struct {
				SplitRecords []struct {
					Wins   int    `json:"wins"`
					Losses int    `json:"losses"`
					Type   string `json:"type"`
					Pct    string `json:"pct"`
				} `json:"splitRecords"`
				DivisionRecords []struct {
					Wins     int    `json:"wins"`
					Losses   int    `json:"losses"`
					Pct      string `json:"pct"`
					Division struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
						Link string `json:"link"`
					} `json:"division"`
				} `json:"divisionRecords"`
				OverallRecords []struct {
					Wins   int    `json:"wins"`
					Losses int    `json:"losses"`
					Type   string `json:"type"`
					Pct    string `json:"pct"`
				} `json:"overallRecords"`
				LeagueRecords []struct {
					Wins   int    `json:"wins"`
					Losses int    `json:"losses"`
					Pct    string `json:"pct"`
					League struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
						Link string `json:"link"`
					} `json:"league"`
				} `json:"leagueRecords"`
				ExpectedRecords []struct {
					Wins   int    `json:"wins"`
					Losses int    `json:"losses"`
					Type   string `json:"type"`
					Pct    string `json:"pct"`
				} `json:"expectedRecords"`
			} `json:"records"`
			RunsAllowed                 int    `json:"runsAllowed"`
			RunsScored                  int    `json:"runsScored"`
			DivisionChamp               bool   `json:"divisionChamp"`
			DivisionLeader              bool   `json:"divisionLeader"`
			HasWildcard                 bool   `json:"hasWildcard"`
			Clinched                    bool   `json:"clinched"`
			EliminationNumber           string `json:"eliminationNumber"`
			EliminationNumberSport      string `json:"eliminationNumberSport"`
			EliminationNumberLeague     string `json:"eliminationNumberLeague"`
			EliminationNumberDivision   string `json:"eliminationNumberDivision"`
			EliminationNumberConference string `json:"eliminationNumberConference"`
			WildCardEliminationNumber   string `json:"wildCardEliminationNumber"`
			Wins                        int    `json:"wins"`
			Losses                      int    `json:"losses"`
			RunDifferential             int    `json:"runDifferential"`
			WinningPercentage           string `json:"winningPercentage"`
		} `json:"record"`
		FranchiseName string `json:"franchiseName"`
		ClubName      string `json:"clubName"`
		Active        bool   `json:"active"`
	} `json:"teams"`
}

func AllMLBTeams() ([]Team, error) {
	resp, err := http.Get(BASE + "teams?sportId=1")
	if err != nil {
		return nil, err
	}
	var teams = struct {
		Copyright string `json:"copyright"`
		Teams     []Team `json:"teams"`
	}{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &teams)
	if err != nil {
		return nil, err
	}

	return teams.Teams, nil
}

func TeambyID(id string) (Team, error) {
	resp, err := http.Get(BASE + fmt.Sprintf("teams/%s", id))
	if err != nil {
		return Team{}, err
	}
	var team = struct {
		Copyright string `json:"copyright"`
		Team      []Team `json:"teams"`
	}{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Team{}, err
	}
	err = json.Unmarshal(b, &team)
	if err != nil {
		return Team{}, err
	}
	return team.Team[0], nil
}

func Roster(id string) ([]FortyManSearch, []FortyManSearch, error) {
	resp, err := http.Get(BASE + fmt.Sprintf("teams/%s/roster?rosterType=40Man", id))
	if err != nil {
		return nil, nil, err
	}
	var roster = struct {
		Copyright string           `json:"copyright"`
		Roster    []FortyManSearch `json:"roster"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&roster)
	if err != nil {
		return nil, nil, err
	}

	var pitchers []FortyManSearch
	var batters []FortyManSearch
	for _, player := range roster.Roster {
		if player.Position.Code == "1" {
			pitchers = append(pitchers, player)
		} else {
			batters = append(batters, player)
		}
	}

	return pitchers, batters, nil
}

type SimpleRecord struct {
	Wins   int
	Losses int
	Pct    string
}

func TeamRecord(id string) (SimpleRecord, error) {
	resp, err := http.Get(fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%s?hydrate=standings", id))
	if err != nil {
		return SimpleRecord{}, err
	}
	var t TeamRecordType

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return SimpleRecord{}, err
	}

	return SimpleRecord{
		Wins:   t.Teams[0].Record.Wins,
		Losses: t.Teams[0].Record.Losses,
		Pct:    t.Teams[0].Record.WinningPercentage,
	}, nil
}
