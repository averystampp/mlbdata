package mlb

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/averystampp/sesame"
)

var tmpl *template.Template = template.Must(template.New("player.html").Funcs(template.FuncMap{
	"metrics":          pitchingMetrics,
	"removeDups":       RemoveDups,
	"removeDupsBatter": RemoveDupsBatter,
}).ParseGlob("../static/*.html"))

func StartMLBService(rtr *sesame.Router) {
	rtr.Get("/", Index)

	// new
	rtr.Get("/team/{teamId}", TeamRosterRoute)
	rtr.Get("/pitcher/{id}", PitcherRoute)
	rtr.Get("/batter/{id}", BatterRoute)

	rtr.Get("/test", func(ctx sesame.Context) error {

		metrics, err := GetPitcherMetrics(573186)
		if err != nil {
			return err
		}
		fmt.Println(metrics)
		return nil
	})

	// exports pitcher
	rtr.Post("/export/pitcher/seasonal/{id}", ExportSeasonalPitcherData)
	rtr.Post("/export/pitcher/career/{id}", ExportCareerPitcherData)

	rtr.Post("/export/batter/career/{id}", ExportCareerBatterDataRoute)
	rtr.Post("/export/batter/seasonal/{id}", ExportSeasonalBatterDataRoute)

	// test
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

func PitcherRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to continue")
	}

	p, err := GetOnePitcherData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
	if err != nil {
		return err
	}

	metrics, err := GetPitcherMetrics(p.ID)
	if err != nil {
		return err
	}

	d := struct {
		PlayerType string
		Data       Pitcher
		Metrics    PitcherMetrics
	}{
		PlayerType: "pitcher",
		Data:       p,
		Metrics:    metrics,
	}

	return tmpl.ExecuteTemplate(ctx.Response(), "player.html", d)
}

func BatterRoute(ctx sesame.Context) error {

	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to continue")
	}

	batter, err := GetOneBatterData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
	if err != nil {
		return err
	}

	d := struct {
		PlayerType string
		Data       Batter
	}{
		PlayerType: "batter",
		Data:       batter,
	}

	return tmpl.ExecuteTemplate(ctx.Response(), "player.html", d)
}

func ExportSeasonalPitcherData(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
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

	exportType := ctx.Request().PostFormValue("exportType")
	switch exportType {
	case "xlsx":
		return p.WriteToExcelSeasonal(ctx.Response())
	case "json":
		return p.WriteToJSON(ctx.Response())
	case "csv":
		return p.WriteToCSV(ctx.Response())
	default:
		return fmt.Errorf("must have a supported export type (json, xlsx, txt)")
	}
}

func ExportCareerPitcherData(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to process request")
	}

	p, err := GetOnePitcherData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
	if err != nil {
		return err
	}

	exportType := ctx.Request().PostFormValue("exportType")
	switch exportType {
	case "xlsx":
		return p.WriteToExcel(ctx.Response())
	case "json":
		return p.WriteToJSON(ctx.Response())
	case "csv":
		return p.WriteToCSV(ctx.Response())
	default:
		return fmt.Errorf("must have a supported export type (json, xlsx, txt)")
	}
}

func ExportCareerBatterDataRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return fmt.Errorf("must have id to process request")
	}

	b, err := GetOneBatterData(id, "yearByYear,career", strconv.Itoa(time.Now().Year()))
	if err != nil {
		return err
	}

	exportType := ctx.Request().PostFormValue("exportType")
	switch exportType {
	case "xlsx":
		return b.WriteToExcel(ctx.Response())
	case "json":
		return b.WriteToJSON(ctx.Response())
	case "csv":
		return b.WriteToCSV(ctx.Response())
	default:
		return fmt.Errorf("must have a supported export type (json, xlsx, txt)")
	}
}

func ExportSeasonalBatterDataRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
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
	b, err := GetOneBatterData(id, "season", s)
	if err != nil {
		return err
	}

	exportType := ctx.Request().PostFormValue("exportType")
	switch exportType {
	case "xlsx":
		return b.WriteToExcelSeasonal(ctx.Response())
	case "json":
		return b.WriteToJSON(ctx.Response())
	case "csv":
		return b.WriteToCSV(ctx.Response())
	default:
		return fmt.Errorf("must have a supported export type (json, xlsx, txt)")
	}
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

	record, err := TeamRecord(teamId)
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
		Record   SimpleRecord
	}{
		Team:     team,
		Pitchers: p,
		Batters:  b,
		Record:   record,
	}

	return tmpl.ExecuteTemplate(ctx.Response(), "teamRoster.html", d)
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
