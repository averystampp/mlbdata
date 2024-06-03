package mlb

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Batter struct {
	ID            int    `json:"id"`
	FullName      string `json:"fullName"`
	Link          string `json:"link"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	PrimaryNumber string `json:"primaryNumber"`
	BirthDate     string `json:"birthDate"`
	CurrentAge    int    `json:"currentAge"`
	BirthCity     string `json:"birthCity"`
	BirthCountry  string `json:"birthCountry"`
	Height        string `json:"height"`
	Weight        int    `json:"weight"`
	Active        bool   `json:"active"`
	CurrentTeam   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"currentTeam"`
	PrimaryPosition struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"primaryPosition"`
	UseName         string `json:"useName"`
	UseLastName     string `json:"useLastName"`
	BoxscoreName    string `json:"boxscoreName"`
	Gender          string `json:"gender"`
	NameMatrilineal string `json:"nameMatrilineal"`
	IsPlayer        bool   `json:"isPlayer"`
	IsVerified      bool   `json:"isVerified"`
	Pronunciation   string `json:"pronunciation"`
	Stats           []struct {
		Type struct {
			DisplayName string `json:"displayName"`
		} `json:"type"`
		Group struct {
			DisplayName string `json:"displayName"`
		} `json:"group"`
		Exemptions []interface{} `json:"exemptions"`
		Splits     []struct {
			Season       string `json:"season"`
			BattingStats `json:"stat"`
			Team         struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"team"`
			Player struct {
				ID       int    `json:"id"`
				FullName string `json:"fullName"`
				Link     string `json:"link"`
			} `json:"player"`
			League struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"league"`
			Sport struct {
				ID           int    `json:"id"`
				Link         string `json:"link"`
				Abbreviation string `json:"abbreviation"`
			} `json:"sport"`
			GameType string `json:"gameType"`
		} `json:"splits"`
	} `json:"stats"`
	MlbDebutDate string `json:"mlbDebutDate"`
	BatSide      struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"batSide"`
	PitchHand struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"pitchHand"`
	NameFirstLast    string  `json:"nameFirstLast"`
	NameSlug         string  `json:"nameSlug"`
	FirstLastName    string  `json:"firstLastName"`
	LastFirstName    string  `json:"lastFirstName"`
	LastInitName     string  `json:"lastInitName"`
	InitLastName     string  `json:"initLastName"`
	FullFMLName      string  `json:"fullFMLName"`
	FullLFMName      string  `json:"fullLFMName"`
	StrikeZoneTop    float64 `json:"strikeZoneTop"`
	StrikeZoneBottom float64 `json:"strikeZoneBottom"`
}

type BattingStats struct {
	GamesPlayed          int    `json:"gamesPlayed"`
	GroundOuts           int    `json:"groundOuts"`
	AirOuts              int    `json:"airOuts"`
	Runs                 int    `json:"runs"`
	Doubles              int    `json:"doubles"`
	Triples              int    `json:"triples"`
	HomeRuns             int    `json:"homeRuns"`
	StrikeOuts           int    `json:"strikeOuts"`
	BaseOnBalls          int    `json:"baseOnBalls"`
	IntentionalWalks     int    `json:"intentionalWalks"`
	Hits                 int    `json:"hits"`
	HitByPitch           int    `json:"hitByPitch"`
	Avg                  string `json:"avg"`
	AtBats               int    `json:"atBats"`
	Obp                  string `json:"obp"`
	Slg                  string `json:"slg"`
	Ops                  string `json:"ops"`
	CaughtStealing       int    `json:"caughtStealing"`
	StolenBases          int    `json:"stolenBases"`
	StolenBasePercentage string `json:"stolenBasePercentage"`
	GroundIntoDoublePlay int    `json:"groundIntoDoublePlay"`
	NumberOfPitches      int    `json:"numberOfPitches"`
	PlateAppearances     int    `json:"plateAppearances"`
	TotalBases           int    `json:"totalBases"`
	Rbi                  int    `json:"rbi"`
	LeftOnBase           int    `json:"leftOnBase"`
	SacBunts             int    `json:"sacBunts"`
	SacFlies             int    `json:"sacFlies"`
	Babip                string `json:"babip"`
	GroundOutsToAirouts  string `json:"groundOutsToAirouts"`
	CatchersInterference int    `json:"catchersInterference"`
	AtBatsPerHomeRun     string `json:"atBatsPerHomeRun"`
}

// yearByYear,career
func GetOneBatterData(id string, span string, season string) (Batter, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people/%s/?hydrate=stats(group=[hitting],type=[%s],seasons=[%s]),currentTeam", id, span, season))
	if err != nil {
		return Batter{}, err
	}
	d := struct {
		Copyright string   `json:"copyright"`
		Pitcher   []Batter `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return Batter{}, err
	}
	if len(d.Pitcher) == 0 {
		return Batter{}, fmt.Errorf("no data marshaled")
	}
	return d.Pitcher[0], nil
}

func GetManyBatterData(ids string) ([]Batter, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people?personIds=%s&hydrate=stats(group=[hitting],type=[season]),currentTeam", ids))
	if err != nil {
		return nil, err
	}
	d := struct {
		Copyright string   `json:"copyright"`
		Batters   []Batter `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d.Batters, nil
}

func (b *Batter) WriteToExcel(resp http.ResponseWriter) error {
	f := excelize.NewFile()
	_, err := f.NewSheet(b.FullName)
	if err != nil {
		return err
	}
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return err
	}
	err = f.SetSheetRow(b.FullName, "A1", &[]interface{}{
		"Team",
		"Season",
		"G",
		"AB",
		"PA",
		"AVG",
		"BABIP",
		"OBP",
		"SLG",
		"OPS",
		"SO",
		"BB",
		"H",
		"2B",
		"3B",
		"HR",
		"RBI",
		"SB",
		"CS",
		"SB%",
	})
	if err != nil {
		return err
	}
	for _, splits := range b.Stats {
		for i, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			if stat.Season != "" {
				err = f.SetSheetRow(b.FullName, fmt.Sprintf("A%d", i+2), &[]interface{}{
					stat.Team.Name,
					stat.Season,
					stat.GamesPlayed,
					stat.AtBats,
					stat.PlateAppearances,
					stat.Avg,
					stat.Babip,
					stat.Obp,
					stat.Slg,
					stat.Ops,
					stat.StrikeOuts,
					stat.BaseOnBalls,
					stat.Hits,
					stat.Doubles,
					stat.Triples,
					stat.HomeRuns,
					stat.Rbi,
					stat.StolenBases,
					stat.CaughtStealing,
					stat.StolenBasePercentage,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	a := fmt.Sprintf("A%d", len(b.Stats[0].Splits)+2)
	err = f.SetSheetRow(b.FullName, a, &[]interface{}{
		"Total",
		b.Stats[1].Splits[0].Season,
		b.Stats[1].Splits[0].GamesPlayed,
		b.Stats[1].Splits[0].AtBats,
		b.Stats[1].Splits[0].PlateAppearances,
		b.Stats[1].Splits[0].Avg,
		b.Stats[1].Splits[0].Babip,
		b.Stats[1].Splits[0].Obp,
		b.Stats[1].Splits[0].Slg,
		b.Stats[1].Splits[0].Ops,
		b.Stats[1].Splits[0].StrikeOuts,
		b.Stats[1].Splits[0].BaseOnBalls,
		b.Stats[1].Splits[0].Hits,
		b.Stats[1].Splits[0].Doubles,
		b.Stats[1].Splits[0].Triples,
		b.Stats[1].Splits[0].HomeRuns,
		b.Stats[1].Splits[0].Rbi,
		b.Stats[1].Splits[0].StolenBases,
		b.Stats[1].Splits[0].CaughtStealing,
		b.Stats[1].Splits[0].StolenBasePercentage,
	})
	if err != nil {
		return err
	}
	resp.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", b.FullName))
	f.Write(resp)
	return nil
}

func (b *Batter) WriteToExcelSeasonal(resp http.ResponseWriter) error {
	f := excelize.NewFile()
	_, err := f.NewSheet(b.FullName)
	if err != nil {
		return err
	}
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return err
	}
	err = f.SetSheetRow(b.FullName, "A1", &[]interface{}{
		"Team",
		"Season",
		"G",
		"AB",
		"PA",
		"AVG",
		"BABIP",
		"OBP",
		"SLG",
		"OPS",
		"SO",
		"BB",
		"H",
		"2B",
		"3B",
		"HR",
		"RBI",
		"SB",
		"CS",
		"SB%",
	})
	if err != nil {
		return err
	}
	for _, splits := range b.Stats {
		for i, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			if stat.Season != "" {
				err = f.SetSheetRow(b.FullName, fmt.Sprintf("A%d", i+2), &[]interface{}{
					stat.Team.Name,
					stat.Season,
					stat.GamesPlayed,
					stat.AtBats,
					stat.PlateAppearances,
					stat.Avg,
					stat.Babip,
					stat.Ops,
					stat.Slg,
					stat.Ops,
					stat.StrikeOuts,
					stat.BaseOnBalls,
					stat.Hits,
					stat.Doubles,
					stat.Triples,
					stat.HomeRuns,
					stat.Rbi,
					stat.StolenBases,
					stat.CaughtStealing,
					stat.StolenBasePercentage,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	if err := f.SetColWidth(b.FullName, "A", "AZ", 5); err != nil {
		return err
	}

	resp.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", b.FullName))
	f.Write(resp)
	return nil
}

func (b *Batter) WriteToJSON(resp http.ResponseWriter) error {
	var p1 Batter = *b
	body, err := json.Marshal(&p1)
	if err != nil {
		return err
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", b.FullName))
	resp.Write(body)
	return nil
}

func (b *Batter) WriteToCSV(resp http.ResponseWriter) error {
	var fileContent [][]string
	head := []string{
		"Team",
		"Season",
		"G",
		"AB",
		"PA",
		"AVG",
		"BABIP",
		"OBP",
		"SLG",
		"OPS",
		"SO",
		"BB",
		"H",
		"2B",
		"3B",
		"HR",
		"RBI",
		"SB",
		"CS",
		"SB%",
	}
	fileContent = append(fileContent, head)
	for _, splits := range b.Stats {
		for _, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			body := []string{
				stat.Team.Name,
				stat.Season,
				strconv.Itoa(stat.GamesPlayed),
				strconv.Itoa(stat.AtBats),
				strconv.Itoa(stat.PlateAppearances),
				stat.Avg,
				stat.Babip,
				stat.Obp,
				stat.Slg,
				stat.Ops,
				strconv.Itoa(stat.StrikeOuts),
				strconv.Itoa(stat.BaseOnBalls),
				strconv.Itoa(stat.Hits),
				strconv.Itoa(stat.Doubles),
				strconv.Itoa(stat.Triples),
				strconv.Itoa(stat.HomeRuns),
				strconv.Itoa(stat.Rbi),
				strconv.Itoa(stat.StolenBases),
				strconv.Itoa(stat.CaughtStealing),
				stat.StolenBasePercentage,
			}

			fileContent = append(fileContent, body)
		}
	}
	resp.Header().Set("Content-Type", "text/csv")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", b.FullName))
	return csv.NewWriter(resp).WriteAll(fileContent)
}
