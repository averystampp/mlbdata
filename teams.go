package mlb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
