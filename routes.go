package mlb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/averystampp/sesame"
)

func StartMLBService(rtr *sesame.Router) {
	rtr.Get("/", Index)
	rtr.Get("/team/overview/{id}", TeamOverviewRoute)
	rtr.Get("/player/stats/{id}", PlayerRoute)
	rtr.Get("/players", AllPlayers)
	rtr.Get("/player/search", PlayerSearchRoute)
	rtr.Get("/player/search/name", PlayerResponseRoute)
	rtr.Get("/game/{game}", GameRoute)
	rtr.Get("/team/schedule/{id}", TeamGames)
	rtr.Get("/help", helper)

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
				<a href="/team/overview/{{.ID}}">{{.Name}}</a>
			</p>
		{{end}}
	`)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "text/html")
	return tmpl.Execute(ctx.Response(), teams)
}

func TeamOverviewRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	roster, err := Roster(id)
	if err != nil {
		return err
	}
	team, err := TeambyID(id)
	if err != nil {
		return err
	}
	today := time.Now().Format(time.DateOnly)
	sevenDays := time.Now().Add((time.Hour * 24) * 7).Format(time.DateOnly)
	gd, err := GetGames(id, today, sevenDays)
	if err != nil {
		return err
	}

	rosterTemp := struct {
		Players []Player
		Team    Team
		Games   GameDay
	}{
		Players: roster,
		Team:    team,
		Games:   gd,
	}

	b, err := os.ReadFile("../pages/roster.html")
	if err != nil {
		return err
	}
	fmap := template.FuncMap{
		"time": func(t time.Time) string {
			t = t.In(time.Local)
			y, m, d := t.Date()

			return fmt.Sprintf("%s %d, %d", m.String(), d, y) + " " + t.Format(time.Kitchen)
		},
	}
	tmpl, err := template.New("").Funcs(fmap).Parse(string(b))
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("Content-Type", "text/html")

	return tmpl.Execute(ctx.Response(), rosterTemp)
}

func PlayerRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	var isPitcher bool
	var hydrate string
	player, err := getOnePlayer(id)
	if err != nil {
		return err
	}
	if player.PrimaryPosition.Abbreviation == "P" {
		isPitcher = true
		hydrate = fmt.Sprintf("people?personIds=%s&hydrate=stats(group=[pitching],type=[yearByYear,career])", id)
	} else {
		hydrate = fmt.Sprintf("people?personIds=%s&hydrate=stats(group=[hitting],type=[yearByYear,career])", id)
	}
	resp, err := http.Get(BASE + hydrate)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var pData PlayerHolder
	err = json.Unmarshal(b, &pData)
	if err != nil {
		return err
	}

	data := struct {
		Player      PlayerOverview
		TeamID      string
		Id          string
		Callback    string
		DataRequest string
	}{
		Player:      pData.People[0],
		TeamID:      strconv.Itoa(player.CurrentTeam.ID),
		Id:          id,
		Callback:    ctx.Request().URL.Path,
		DataRequest: hydrate,
	}
	tmpl := template.New("player")
	if isPitcher {
		b, err := os.ReadFile("../pages/pitcher.html")
		if err != nil {
			return err
		}
		tmpl, err := tmpl.Parse(string(b))
		if err != nil {
			return err
		}
		ctx.Response().Header().Set("Content-Type", "text/html")
		return tmpl.ExecuteTemplate(ctx.Response(), "player", data)
	}
	parse, err := os.ReadFile("../pages/player.html")
	if err != nil {
		return err
	}
	tmpl, err = tmpl.Parse(string(parse))
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(ctx.Response(), "player", data)
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

func GameRoute(ctx sesame.Context) error {
	gameId := ctx.Request().PathValue("game")
	tmpl, err := template.New("game.html").Funcs(template.FuncMap{
		"pitchers":     retrievePitcher,
		"batters":      retrieveBatters,
		"add":          addNums,
		"sub":          subNums,
		"onePitcher":   singlePitcher,
		"pastPitchers": playerGameStats,
		"era":          calculateERA,
		"localTime":    localTime,
	}).ParseFiles("../pages/game.html", "../components/game_batter.html", "../components/game_pitcher.html")
	if err != nil {
		return err
	}
	resp, err := http.Get(GAMEURL + fmt.Sprintf("game/%s/feed/live", gameId))
	if err != nil {
		return err
	}
	var game Game
	err = json.NewDecoder(resp.Body).Decode(&game)
	if err != nil {
		return err
	}
	return tmpl.Execute(ctx.Response(), game)
}

func TeamGames(ctx sesame.Context) error {
	start := ctx.Request().FormValue("start")
	end := ctx.Request().FormValue("end")
	id := ctx.Request().PathValue("id")
	if start == "" {
		start = time.Now().Format(time.DateOnly)
	}
	if end == "" {
		end = time.Now().Add(time.Hour * 24).Format(time.DateOnly)
	}
	games, err := GetGames(id, start, end)
	if err != nil {
		return err
	}

	tmpl, err := template.New("teamgames.html").Funcs(template.FuncMap{
		"time": func(t time.Time) string {
			t = t.In(time.Local)
			y, m, d := t.Date()

			return fmt.Sprintf("%s %d, %d", m.String(), d, y) + " " + t.Format(time.Kitchen)
		},
	}).ParseFiles("../pages/teamgames.html")
	if err != nil {
		return err
	}
	var data = struct {
		Data   GameDay
		TeamID string
	}{
		Data:   games,
		TeamID: id,
	}

	return tmpl.Execute(ctx.Response(), data)
}

func helper(ctx sesame.Context) error {
	ls := []int{434378, 650556, 519151, 623352}

	pp := playerGameStats(ls, 745663, "")
	for _, p := range pp {
		fmt.Println(p)
	}
	return nil
}
