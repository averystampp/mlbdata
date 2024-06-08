package mlb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type CustomTeam struct {
	Name     string
	TeamID   uuid.UUID
	Manager  string
	Batters  []int
	Pitchers []int
}

func CreateCustomTeam(teamName, manager string) (string, error) {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()
	team := CustomTeam{
		Name:    teamName,
		Manager: manager,
		TeamID:  uuid.New(),
	}
	body, err := json.Marshal(&team)
	if err != nil {
		return "", err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("teams"))
		return b.Put(team.TeamID[:], body)
	})
	if err != nil {
		return "", err
	}
	return team.TeamID.String(), nil
}

func GetCustomTeams() ([]CustomTeam, error) {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var teams []CustomTeam
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("teams"))
		err := b.ForEach(func(k, v []byte) error {
			var t CustomTeam
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			teams = append(teams, t)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func GetOneCustomTeam(teamId uuid.UUID) (*CustomTeam, error) {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var t CustomTeam
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("teams"))
		body := b.Get(teamId[:])
		if len(body) == 0 {
			return fmt.Errorf("no team found for %s", teamId.String())
		}
		return json.Unmarshal(body, &t)
	})
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Control flow of drafting a custom team
// 1. Create a team, with a name.
// 2. Pick players
// 3. Export team, with data on the team averages and expected perf
// 4. Does it make sense to limit the amount of players on a team?
