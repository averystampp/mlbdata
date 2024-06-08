package mlb

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/averystampp/sesame"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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

	rtr.Get("/static/{name}", func(ctx sesame.Context) error {
		fileName := ctx.Request().PathValue("name")
		if fileName == "" {
			return fmt.Errorf("must have a filename to serve")
		}
		http.ServeFile(ctx.Response(), ctx.Request(), fmt.Sprintf("../static/%s", fileName))
		return nil
	})

	// Experimental
	rtr.Get("/batter/{id}/games/logs", BatterGameLogRoute)

	// Custom Teams
	rtr.Get("/custom/home", CustomTeamHomeRoute)
	rtr.Get("/custom/team/{teamId}", GetCustomTeamRoute)
	rtr.Post("/custom/team/create", CreateCustomTeamRoute)

	// Data Export Routes
	// Pitcher Data
	rtr.Post("/export/pitcher/seasonal/{id}", ExportSeasonalPitcherData)
	rtr.Post("/export/pitcher/career/{id}", ExportCareerPitcherData)
	// Batter Data
	rtr.Post("/export/batter/career/{id}", ExportCareerBatterDataRoute)
	rtr.Post("/export/batter/seasonal/{id}", ExportSeasonalBatterDataRoute)
}

func Index(ctx sesame.Context) error {
	teams, err := AllMLBTeams()
	if err != nil {
		return ErrorRoute(ctx, http.StatusInternalServerError, err)
	}

	tmpl, err := template.New("").Parse(`
		{{range .}}
			<p>
				<a href="/team/{{.ID}}">{{.Name}}</a>
			</p>
		{{end}}
	`)
	if err != nil {
		return ErrorRoute(ctx, http.StatusInternalServerError, err)
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
	p, err := GetManyPitcherData(joinIds(pitchers), "season", "2024")
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

func BatterGameLogRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("id")
	if id == "" {
		return ErrorRoute(ctx, http.StatusBadRequest, fmt.Errorf("must have player id in request"))
	}
	season := ctx.Request().URL.Query().Get("season")
	if season == "" {
		season = strconv.Itoa(time.Now().Year())
	}
	batter, err := GetBatterGameLog(id, season)
	if err != nil {
		return ErrorRoute(ctx, http.StatusBadRequest, err)
	}
	f := excelize.NewFile()
	defer f.Close()
	return batter.WriteBatterLogExcel(f, ctx.Response())
}

func ErrorRoute(ctx sesame.Context, status int, topLevelErr error) error {
	s := `
	<div>
		<h1>Server responded with an error</h1>
		<p>{{.}}</p>
	</div>`

	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return err
	}
	ctx.Response().WriteHeader(status)
	ctx.Response().Header().Set("Content-Type", "text/html")
	return tmpl.Execute(ctx.Response(), topLevelErr)
}

func CustomTeamHomeRoute(ctx sesame.Context) error {
	t, err := GetCustomTeams()
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(ctx.Response(), "custom_teams_home.html", t)
}

func CreateCustomTeamRoute(ctx sesame.Context) error {
	teamName := ctx.Request().PostFormValue("teamName")
	if teamName == "" {
		return ErrorRoute(ctx, http.StatusBadRequest, fmt.Errorf("must have a team name to create a team"))
	}
	_, err := CreateCustomTeam(teamName, "someone")
	if err != nil {
		return ErrorRoute(ctx, http.StatusInternalServerError, err)
	}

	http.Redirect(ctx.Response(), ctx.Request(), "/custom/home", http.StatusSeeOther)
	return nil
}

func GetCustomTeamRoute(ctx sesame.Context) error {
	id := ctx.Request().PathValue("teamId")
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrorRoute(ctx, http.StatusBadRequest, err)
	}

	team, err := GetOneCustomTeam(uid)
	if err != nil {
		return ErrorRoute(ctx, http.StatusInternalServerError, err)
	}
	fmt.Println(*team)
	return nil
}
