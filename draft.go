package mlb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/averystampp/sesame"
	bolt "go.etcd.io/bbolt"
)

type DraftTeam struct {
	Name        string
	DateCreated string
	CreatedBy   string
	Players     []int
}

func DraftService(rtr *sesame.Router) {
	rtr.Get("/draft/teams", DraftTeamViews)
	rtr.Post("/draft/team/create", CreateTeam)
}

func DraftTeamViews(ctx sesame.Context) error {
	tmpl, err := template.ParseFiles("../pages/draftteams.html")
	if err != nil {
		return err
	}

	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	var teamHolder []DraftTeam
	db.View(func(tx *bolt.Tx) error {
		teams := tx.Bucket([]byte("teams"))
		if err != nil {
			return err
		}
		teams.ForEach(func(k, v []byte) error {
			var t DraftTeam
			err := json.Unmarshal(v, &t)
			if err != nil {
				return err
			}
			teamHolder = append(teamHolder, t)
			return nil
		})
		return nil
	})

	return tmpl.Execute(ctx.Response(), teamHolder)
}

func CreateTeam(ctx sesame.Context) error {
	team := ctx.Request().FormValue("team-name")
	if team == "" {
		return fmt.Errorf("must have a team name")
	}

	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("teams"))
		if err != nil {
			return err
		}
		b, err := json.Marshal(DraftTeam{
			Name:        team,
			DateCreated: time.Now().String(),
			CreatedBy:   "Avery",
		})
		if err != nil {
			return err
		}
		return bucket.Put([]byte(team), b)
	})
	if err != nil {
		return err
	}
	db.Close()
	http.Redirect(ctx.Response(), ctx.Request(), "/draft/teams", http.StatusSeeOther)
	return nil
}

func CreateTeamView(ctx sesame.Context) error {
	return nil
}

// db, err := bolt.Open("../players.db", 0600, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// Bolt schema for draft: Player Name -> []DraftTeams
