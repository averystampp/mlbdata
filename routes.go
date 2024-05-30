package mlb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/averystampp/sesame"
)

func StartMLBService(rtr *sesame.Router) {
	rtr.Get("/", Index)
	// remove when possible
	rtr.Get("/players", AllPlayers)
	rtr.Get("/player/search", PlayerSearchRoute)
	rtr.Get("/player/search/name", PlayerResponseRoute)

	// new
	rtr.Get("/team/{teamId}", TeamRosterRoute)
	rtr.Get("/pitcher/{id}", PitcherRoute)
	rtr.Post("/export/pitcher/{id}", ExportPitcherDataRoute)
	rtr.Post("/test", t)

	// dont touch yet
	rtr.Post("/player/export", ExportPlayerData)

}

func Index(ctx sesame.Context) error {
	teams, err := AllMLBTeams()
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(`
		{{range .}}
			<p>
				<a href="/team/{{.ID}}">{{.Name}}</a>
			</p>
		{{end}}
	`)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "text/html")
	return tmpl.Execute(ctx.Response(), teams)
}

func AllPlayers(ctx sesame.Context) error {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var players []PlayerSearch
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		if err != nil {
			return err
		}
		err = b.ForEach(func(k, v []byte) error {
			id, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				return err
			}
			p := PlayerSearch{
				Name: string(k),
				ID:   int(id),
			}
			fmt.Println(p)
			players = append(players, p)
			return nil

		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(&players)
	ctx.Response().Write(body)
	return nil
}

func PlayerSearchRoute(ctx sesame.Context) error {
	tmpl, err := template.ParseFiles("../pages/builder.html")
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "text/html")
	return tmpl.Execute(ctx.Response(), nil)
}

func PlayerResponseRoute(ctx sesame.Context) error {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		ctx.Response().Write([]byte("{}"))
		return nil
	}
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var players []PlayerSearch
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		if err != nil {
			return err
		}
		err = b.ForEach(func(k, v []byte) error {
			name = strings.ToTitle(string(name[0])) + name[1:]
			if strings.HasPrefix(string(k), name) {
				id, err := strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					return err
				}

				p := PlayerSearch{
					Name: string(k),
					ID:   int(id),
				}
				players = append(players, p)
				return nil
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "application/json")
	if len(players) > 5 {
		players = players[:5]
	}
	body, err := json.Marshal(&players)
	ctx.Response().Write(body)
	return nil
}

func PitcherRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to continue")
	}

	p, err := GetOnePitcherData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
	if err != nil {
		return err
	}
	tmpl, err := template.New("player.html").Funcs(template.FuncMap{
		"removeDups": func(p Pitcher) []string {
			var seasons []string
			for _, stat := range p.Stats {
				for _, spl := range stat.Splits {
					if len(seasons) == 0 {
						seasons = append(seasons, spl.Season)
					}
					if spl.Season == seasons[len(seasons)-1] || spl.Season == "" {
						continue
					}
					seasons = append(seasons, spl.Season)
				}
			}
			return seasons
		},
	}).ParseFiles("../pages/player.html")
	if err != nil {
		return err
	}

	d := struct {
		PlayerType string
		Data       Pitcher
	}{
		PlayerType: "pitcher",
		Data:       p,
	}

	return tmpl.Execute(ctx.Response(), d)
}

func ExportPitcherDataRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to process request")
	}
	span := ctx.Request().PostFormValue("span")
	if span == "career" {
		p, err := GetOnePitcherData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
		if err != nil {
			return err
		}
		b, err := json.Marshal(&p)
		if err != nil {
			return err
		}
		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.Response().Write(b)
		return nil
	}
	if span == "season" {
		err := ctx.Request().ParseForm()
		if err != nil {
			return err
		}
		var s string
		for k, v := range ctx.Request().PostForm {
			if k == "season" {
				s = joinSeasons(v)
			}
		}
		p, err := GetOnePitcherData(id, "season", s)
		if err != nil {
			return err
		}
		b, err := json.Marshal(&p)
		if err != nil {
			return err
		}
		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.Response().Write(b)
		return nil
	}

	return fmt.Errorf("must have a span of \"career\" or \"season\"")
}

func TeamRosterRoute(ctx sesame.Context) error {
	now := time.Now()
	defer func() {
		log.Printf("path=%s took=%s\n", ctx.Request().URL.Path, time.Since(now))
	}()
	teamId := ctx.Request().PathValue("teamId")
	if teamId == "" {
		return fmt.Errorf("need a team id to process request")
	}

	team, err := TeambyID(teamId)
	if err != nil {
		return err
	}
	pitchers, batters, err := Roster(teamId)
	if err != nil {
		return err
	}

	p, err := GetManyPitcherData(joinIds(pitchers))
	if err != nil {
		return err
	}
	b, err := GetManyBatterData(joinIds(batters))
	if err != nil {
		return err
	}
	d := struct {
		Team     Team
		Pitchers []Pitcher
		Batters  []Batter
	}{
		Team:     team,
		Pitchers: p,
		Batters:  b,
	}

	tmpl, err := template.ParseFiles("../pages/teamRoster.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(ctx.Response(), d)
}

func joinIds(players []FortyManSearch) string {
	var ids string
	for i, p := range players {
		if i != len(players)-1 {
			ids += strconv.Itoa(p.Person.ID) + ","
		} else {
			ids += strconv.Itoa(p.Person.ID)
		}

	}
	return ids
}

func t(ctx sesame.Context) error {
	err := ctx.Request().ParseForm()
	if err != nil {
		return err
	}
	for k, v := range ctx.Request().PostForm {
		if k == "season" {
			s := joinSeasons(v)
			fmt.Println(s)
		}
	}
	return nil
}

func joinSeasons(s []string) string {
	var seasons string
	for i, p := range s {
		if i != len(s)-1 {
			seasons += p + ","
		} else {
			seasons += p
		}

	}
	return seasons
}
