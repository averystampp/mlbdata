package mlb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/averystampp/sesame"
	"github.com/xuri/excelize/v2"
)

func ExportPlayerData(ctx sesame.Context) error {
	p := ctx.Request().PostFormValue("callback")
	data := ctx.Request().PostFormValue("data")
	resp, err := http.Get(BASE + data)
	if err != nil {
		log.Printf("error exporting player data: %s", err)
		http.Redirect(ctx.Response(), ctx.Request(), p, http.StatusSeeOther)
		return nil
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
	if pData.People[0].PrimaryPosition.Code == "1" {
		f, err := writePitcherData(pData.People[0].FullName, pData.People[0].Stats)
		if err != nil {
			return err
		}
		ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", pData.People[0].FullName))
		return f.Write(ctx.Response())
	}
	bd, err := writeBatterData(pData.People[0].FullName, pData.People[0].Stats)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", pData.People[0].FullName))
	return bd.Write(ctx.Response())
}

func writePitcherData(name string, stats []Stat) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err := f.NewSheet(name)
	if err != nil {
		return nil, err
	}
	f.SetSheetRow(name, "A1", &[]interface{}{
		"Team Name",
		"Season",
		"G",
		"K",
		"H",
		"BB",
		"HR",
		"OBP",
		"OPS",
		"SLG",
		"AVG",
		"ERA",
		"IP",
		"ER",
		"S",
		"S%",
		"HP9",
		"BB9",
		"HR9",
		"SO9",
	})
	var totalrow []Split
	var i int
	for _, s := range stats {
		if s.Type.DisplayName == "career" {
			totalrow = s.Splits
			continue
		}
		for _, spl := range s.Splits {
			if spl.SplitTeam.Name == "" {
				spl.SplitTeam.Name = fmt.Sprintf("%s Total", spl.Season)
			}
			err := f.SetSheetRow(name, fmt.Sprintf("A%d", i+2), &[]interface{}{
				spl.SplitTeam.Name,
				spl.Season,
				spl.GamesPlayed,
				spl.StrikeOuts,
				spl.Hits,
				spl.BaseOnBalls,
				spl.HomeRuns,
				spl.Obp,
				spl.Ops,
				spl.Slg,
				spl.Avg,
				spl.Era,
				spl.InningsPitched,
				spl.EarnedRuns,
				spl.Strikes,
				spl.StrikePercentage,
				spl.SplitStat.HitsPer9Inn,
				spl.SplitStat.WalksPer9Inn,
				spl.SplitStat.HomeRunsPer9,
				spl.SplitStat.StrikeoutsPer9Inn,
			})
			i += 1
			if err != nil {
				return nil, err
			}
		}
	}
	err = f.SetSheetRow(name, fmt.Sprintf("A%d", i+2), &[]interface{}{
		"Career Total",
		"",
		totalrow[0].GamesPlayed,
		totalrow[0].StrikeOuts,
		totalrow[0].Hits,
		totalrow[0].BaseOnBalls,
		totalrow[0].HomeRuns,
		totalrow[0].Obp,
		totalrow[0].Ops,
		totalrow[0].Slg,
		totalrow[0].Avg,
		totalrow[0].Era,
		totalrow[0].InningsPitched,
		totalrow[0].EarnedRuns,
		totalrow[0].Strikes,
		totalrow[0].StrikePercentage,
		totalrow[0].SplitStat.HitsPer9Inn,
		totalrow[0].SplitStat.WalksPer9Inn,
		totalrow[0].SplitStat.HomeRunsPer9,
		totalrow[0].SplitStat.StrikeoutsPer9Inn,
	})
	if err != nil {
		return nil, err
	}
	if err := f.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}
	styleId, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}
	err = f.SetRowStyle(name, i+2, i+2, styleId)
	if err != nil {
		return nil, err
	}
	if err := f.SetColWidth(name, "A", "T", 5); err != nil {
		return nil, err
	}
	return f, nil
}

func writeBatterData(name string, stats []Stat) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err := f.NewSheet(name)
	if err != nil {
		return nil, err
	}
	f.SetSheetRow(name, "A1", &[]interface{}{
		"Team Name",
		"Season",
		"G",
		"AB",
		"H",
		"BB",
		"RBI",
		"SB",
		"CS",
		"R",
		"HBP",
		"K",
		"AO",
		"GO",
		"2B",
		"3B",
		"TB",
		"HR",
		"OBP",
		"OPS",
		"SLG",
		"AVG",
		"BABIP",
	})
	var totalrow []Split
	var i int
	for _, s := range stats {
		if s.Type.DisplayName == "career" {
			totalrow = s.Splits
			continue
		}
		for _, spl := range s.Splits {
			if spl.SplitTeam.Name == "" {
				spl.SplitTeam.Name = fmt.Sprintf("%s Total", spl.Season)
			}
			err := f.SetSheetRow(name, fmt.Sprintf("A%d", i+2), &[]interface{}{
				spl.SplitTeam.Name,
				spl.Season,
				spl.GamesPlayed,
				spl.AtBats,
				spl.Hits,
				spl.BaseOnBalls,
				spl.Rbi,
				spl.StolenBases,
				spl.CaughtStealing,
				spl.Runs,
				spl.HitByPitch,
				spl.StrikeOuts,
				spl.AirOuts,
				spl.GroundOuts,
				spl.Doubles,
				spl.Triples,
				spl.TotalBases,
				spl.HomeRuns,
				spl.Obp,
				spl.Ops,
				spl.Slg,
				spl.Avg,
				spl.SplitStat.Babip,
			})
			i += 1
			if err != nil {
				return nil, err
			}
		}
	}
	lastRow := fmt.Sprintf("A%d", i+2)
	err = f.SetSheetRow(name, lastRow, &[]interface{}{
		"Career Total",
		"",
		totalrow[0].GamesPlayed,
		totalrow[0].AtBats,
		totalrow[0].Hits,
		totalrow[0].BaseOnBalls,
		totalrow[0].Rbi,
		totalrow[0].StolenBases,
		totalrow[0].CaughtStealing,
		totalrow[0].Runs,
		totalrow[0].HitByPitch,
		totalrow[0].StrikeOuts,
		totalrow[0].AirOuts,
		totalrow[0].GroundOuts,
		totalrow[0].Doubles,
		totalrow[0].Triples,
		totalrow[0].TotalBases,
		totalrow[0].HomeRuns,
		totalrow[0].Obp,
		totalrow[0].Ops,
		totalrow[0].Slg,
		totalrow[0].Avg,
		totalrow[0].Babip,
	})
	if err != nil {
		return nil, err
	}
	styleId, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}
	err = f.SetRowStyle(name, i+2, i+2, styleId)
	if err != nil {
		return nil, err
	}
	if err := f.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}
	if err := f.SetColWidth(name, "A", "T", 5); err != nil {
		return nil, err
	}
	return f, nil
}
