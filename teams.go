package mlb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func Roster(id string) ([]Player, error) {
	resp, err := http.Get(BASE + fmt.Sprintf("teams/%s/roster?rosterType=40Man", id))
	if err != nil {
		return nil, err
	}
	var roster = struct {
		Copyright string   `json:"copyright"`
		Roster    []Player `json:"roster"`
	}{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &roster)
	if err != nil {
		return nil, err
	}
	return roster.Roster, nil
}

func GetGames(id, start, end string) (GameDay, error) {
	if start == "" || end == "" {
		start = time.Now().Format(time.DateOnly)
		end = time.Now().Format(time.DateOnly)
	}
	full := fmt.Sprintf(BASE+"schedule?sportId=1&startDate=%s&endDate=%s&teamId=%s", start, end, id)
	resp, err := http.Get(full)
	if err != nil {
		return GameDay{}, err
	}

	var gd GameDay
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&gd)
	if err != nil {
		return GameDay{}, err
	}

	return gd, nil
}
