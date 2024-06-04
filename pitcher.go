package mlb

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Team_ struct {
	ID int `json:"id"`
}

type FortyManSearch struct {
	Person struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
		Link     string `json:"link"`
	} `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
	Status struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	ParentTeamID int `json:"parentTeamId"`
}

type Pitcher struct {
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
			Season        string `json:"season"`
			PitchingStats `json:"stat"`
			Team          struct {
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

type PitchingStats struct {
	GamesPlayed            int    `json:"gamesPlayed"`
	GamesStarted           int    `json:"gamesStarted"`
	GroundOuts             int    `json:"groundOuts"`
	AirOuts                int    `json:"airOuts"`
	Runs                   int    `json:"runs"`
	Doubles                int    `json:"doubles"`
	Triples                int    `json:"triples"`
	HomeRuns               int    `json:"homeRuns"`
	StrikeOuts             int    `json:"strikeOuts"`
	BaseOnBalls            int    `json:"baseOnBalls"`
	IntentionalWalks       int    `json:"intentionalWalks"`
	Hits                   int    `json:"hits"`
	HitByPitch             int    `json:"hitByPitch"`
	Avg                    string `json:"avg"`
	AtBats                 int    `json:"atBats"`
	Obp                    string `json:"obp"`
	Slg                    string `json:"slg"`
	Ops                    string `json:"ops"`
	CaughtStealing         int    `json:"caughtStealing"`
	StolenBases            int    `json:"stolenBases"`
	StolenBasePercentage   string `json:"stolenBasePercentage"`
	GroundIntoDoublePlay   int    `json:"groundIntoDoublePlay"`
	NumberOfPitches        int    `json:"numberOfPitches"`
	Era                    string `json:"era"`
	InningsPitched         string `json:"inningsPitched"`
	Wins                   int    `json:"wins"`
	Losses                 int    `json:"losses"`
	Saves                  int    `json:"saves"`
	SaveOpportunities      int    `json:"saveOpportunities"`
	Holds                  int    `json:"holds"`
	BlownSaves             int    `json:"blownSaves"`
	EarnedRuns             int    `json:"earnedRuns"`
	Whip                   string `json:"whip"`
	BattersFaced           int    `json:"battersFaced"`
	Outs                   int    `json:"outs"`
	GamesPitched           int    `json:"gamesPitched"`
	CompleteGames          int    `json:"completeGames"`
	Shutouts               int    `json:"shutouts"`
	Strikes                int    `json:"strikes"`
	StrikePercentage       string `json:"strikePercentage"`
	HitBatsmen             int    `json:"hitBatsmen"`
	Balks                  int    `json:"balks"`
	WildPitches            int    `json:"wildPitches"`
	Pickoffs               int    `json:"pickoffs"`
	TotalBases             int    `json:"totalBases"`
	GroundOutsToAirouts    string `json:"groundOutsToAirouts"`
	WinPercentage          string `json:"winPercentage"`
	PitchesPerInning       string `json:"pitchesPerInning"`
	GamesFinished          int    `json:"gamesFinished"`
	StrikeoutWalkRatio     string `json:"strikeoutWalkRatio"`
	StrikeoutsPer9Inn      string `json:"strikeoutsPer9Inn"`
	WalksPer9Inn           string `json:"walksPer9Inn"`
	HitsPer9Inn            string `json:"hitsPer9Inn"`
	RunsScoredPer9         string `json:"runsScoredPer9"`
	HomeRunsPer9           string `json:"homeRunsPer9"`
	InheritedRunners       int    `json:"inheritedRunners"`
	InheritedRunnersScored int    `json:"inheritedRunnersScored"`
	CatchersInterference   int    `json:"catchersInterference"`
	SacBunts               int    `json:"sacBunts"`
	SacFlies               int    `json:"sacFlies"`
}

// .people.stats.splits.stat.metric.averageValue
type PitcherMetrics struct {
	ID                 int    `json:"id"`
	FullName           string `json:"fullName"`
	Link               string `json:"link"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	PrimaryNumber      string `json:"primaryNumber"`
	BirthDate          string `json:"birthDate"`
	CurrentAge         int    `json:"currentAge"`
	BirthCity          string `json:"birthCity"`
	BirthStateProvince string `json:"birthStateProvince"`
	BirthCountry       string `json:"birthCountry"`
	Height             string `json:"height"`
	Weight             int    `json:"weight"`
	Active             bool   `json:"active"`
	PrimaryPosition    struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"primaryPosition"`
	UseName      string `json:"useName"`
	UseLastName  string `json:"useLastName"`
	MiddleName   string `json:"middleName"`
	BoxscoreName string `json:"boxscoreName"`
	NickName     string `json:"nickName"`
	Gender       string `json:"gender"`
	IsPlayer     bool   `json:"isPlayer"`
	IsVerified   bool   `json:"isVerified"`
	DraftYear    int    `json:"draftYear"`
	Stats        []struct {
		Type struct {
			DisplayName string `json:"displayName"`
		} `json:"type"`
		Group struct {
			DisplayName string `json:"displayName"`
		} `json:"group"`
		Exemptions []interface{} `json:"exemptions"`
		Splits     []struct {
			Season string `json:"season"`
			Stat   struct {
				Metric struct {
					Group        string  `json:"group"`
					Name         string  `json:"name"`
					AverageValue float64 `json:"averageValue"`
					MinValue     float64 `json:"minValue"`
					MaxValue     float64 `json:"maxValue"`
					Unit         string  `json:"unit"`
					MetricID     int     `json:"metricId"`
				} `json:"metric"`
				Event struct {
					Details struct {
						Type struct {
							Code        string `json:"code"`
							Description string `json:"description"`
						} `json:"type"`
					} `json:"details"`
					Count struct {
					} `json:"count"`
				} `json:"event"`
			} `json:"stat"`
			Player struct {
				ID              int    `json:"id"`
				FullName        string `json:"fullName"`
				Link            string `json:"link"`
				FirstName       string `json:"firstName"`
				LastName        string `json:"lastName"`
				PrimaryPosition struct {
					Code         string `json:"code"`
					Name         string `json:"name"`
					Type         string `json:"type"`
					Abbreviation string `json:"abbreviation"`
				} `json:"primaryPosition"`
			} `json:"player"`
			GameType       string `json:"gameType"`
			NumOccurrences int    `json:"numOccurrences"`
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

// yearByYear,career
func GetOnePitcherData(id string, span string, season string) (Pitcher, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people/%s/?hydrate=stats(group=[pitching],type=[%s],seasons=[%s]),currentTeam", id, span, season))
	if err != nil {
		return Pitcher{}, err
	}
	d := struct {
		Copyright string    `json:"copyright"`
		Pitcher   []Pitcher `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return Pitcher{}, err
	}
	if len(d.Pitcher) == 0 {
		return Pitcher{}, fmt.Errorf("no data marshaled")
	}
	return d.Pitcher[0], nil
}

func GetManyPitcherData(ids, span, season string) ([]Pitcher, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people?personIds=%s&hydrate=stats(group=[pitching],type=[%s],seasons=[%s]),currentTeam", ids, span, season))
	if err != nil {
		return nil, err
	}
	d := struct {
		Copyright string    `json:"copyright"`
		Pitcher   []Pitcher `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return nil, err
	}

	return d.Pitcher, nil
}

func (p *Pitcher) WriteToJSON(resp http.ResponseWriter) error {
	var p1 Pitcher = *p
	b, err := json.Marshal(&p1)
	if err != nil {
		return err
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", p.FullName))
	resp.Write(b)
	return nil
}

func (p *Pitcher) WriteToCSV(resp http.ResponseWriter) error {
	var fileContent [][]string
	head := []string{
		"Team",
		"Season",
		"G",
		"IP",
		"NOP",
		"BF",
		"ERA",
		"HR",
		"BB",
		"K",
		"2B",
		"3B",
		"R",
		"ER",
		"HBP",
		"AO",
		"GO",
		"O",
		"BK",
		"WP",
		"PO",
		"WHIP",
		"AVG",
		"OBP",
		"SLG",
		"OPS",
		"SO9",
		"WP9",
		"HP9",
		"RP9",
		"HRP9",
	}
	fileContent = append(fileContent, head)
	for _, splits := range p.Stats {
		for _, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			body := []string{
				stat.Team.Name,
				stat.Season,
				strconv.Itoa(stat.GamesPlayed),
				stat.InningsPitched,
				strconv.Itoa(stat.NumberOfPitches),
				strconv.Itoa(stat.BattersFaced),
				stat.Era,
				strconv.Itoa(stat.HomeRuns),
				strconv.Itoa(stat.BaseOnBalls),
				strconv.Itoa(stat.StrikeOuts),
				strconv.Itoa(stat.Doubles),
				strconv.Itoa(stat.Triples),
				strconv.Itoa(stat.Runs),
				strconv.Itoa(stat.EarnedRuns),
				strconv.Itoa(stat.HitByPitch),
				strconv.Itoa(stat.AirOuts),
				strconv.Itoa(stat.GroundOuts),
				strconv.Itoa(stat.Outs),
				strconv.Itoa(stat.Balks),
				strconv.Itoa(stat.WildPitches),
				strconv.Itoa(stat.Pickoffs),
				stat.Whip,
				stat.Avg,
				stat.Obp,
				stat.Slg,
				stat.Ops,
				stat.StrikeoutsPer9Inn,
				stat.WalksPer9Inn,
				stat.HitsPer9Inn,
				stat.RunsScoredPer9,
				stat.HomeRunsPer9,
			}

			fileContent = append(fileContent, body)
		}
	}
	resp.Header().Set("Content-Type", "text/csv")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", p.FullName))
	return csv.NewWriter(resp).WriteAll(fileContent)
}

func (p *Pitcher) WriteToExcel(resp http.ResponseWriter) error {
	f := excelize.NewFile()
	_, err := f.NewSheet(p.FullName)
	if err != nil {
		return err
	}
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return err
	}
	err = f.SetSheetRow(p.FullName, "A1", &[]interface{}{
		"Team",
		"Season",
		"G",
		"IP",
		"NOP",
		"BF",
		"ERA",
		"HR",
		"BB",
		"K",
		"2B",
		"3B",
		"R",
		"ER",
		"HBP",
		"AO",
		"GO",
		"O",
		"BK",
		"WP",
		"PO",
		"WHIP",
		"AVG",
		"OBP",
		"SLG",
		"OPS",
		"SO9",
		"WP9",
		"HP9",
		"RP9",
		"HRP9",
	})
	if err != nil {
		return err
	}
	for _, splits := range p.Stats {
		for i, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			if stat.Season != "" {

				err = f.SetSheetRow(p.FullName, fmt.Sprintf("A%d", i+2), &[]interface{}{
					stat.Team.Name,
					stat.Season,
					stat.GamesPlayed,
					stat.InningsPitched,
					stat.NumberOfPitches,
					stat.BattersFaced,
					stat.Era,
					stat.HomeRuns,
					stat.BaseOnBalls,
					stat.StrikeOuts,
					stat.Doubles,
					stat.Triples,
					stat.Runs,
					stat.EarnedRuns,
					stat.HitByPitch,
					stat.AirOuts,
					stat.GroundOuts,
					stat.Outs,
					stat.Balks,
					stat.WildPitches,
					stat.Pickoffs,
					stat.Whip,
					stat.Avg,
					stat.Obp,
					stat.Slg,
					stat.Ops,
					stat.StrikeoutsPer9Inn,
					stat.WalksPer9Inn,
					stat.HitsPer9Inn,
					stat.RunsScoredPer9,
					stat.HomeRunsPer9,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	a := fmt.Sprintf("A%d", len(p.Stats[0].Splits)+2)
	err = f.SetSheetRow(p.FullName, a, &[]interface{}{
		p.Stats[1].Splits[0].Team.Name,
		p.Stats[1].Splits[0].Season,
		p.Stats[1].Splits[0].GamesPlayed,
		p.Stats[1].Splits[0].InningsPitched,
		p.Stats[1].Splits[0].NumberOfPitches,
		p.Stats[1].Splits[0].BattersFaced,
		p.Stats[1].Splits[0].Era,
		p.Stats[1].Splits[0].HomeRuns,
		p.Stats[1].Splits[0].BaseOnBalls,
		p.Stats[1].Splits[0].StrikeOuts,
		p.Stats[1].Splits[0].Doubles,
		p.Stats[1].Splits[0].Triples,
		p.Stats[1].Splits[0].Runs,
		p.Stats[1].Splits[0].EarnedRuns,
		p.Stats[1].Splits[0].HitByPitch,
		p.Stats[1].Splits[0].AirOuts,
		p.Stats[1].Splits[0].GroundOuts,
		p.Stats[1].Splits[0].Outs,
		p.Stats[1].Splits[0].Balks,
		p.Stats[1].Splits[0].WildPitches,
		p.Stats[1].Splits[0].Pickoffs,
		p.Stats[1].Splits[0].Whip,
		p.Stats[1].Splits[0].Avg,
		p.Stats[1].Splits[0].Obp,
		p.Stats[1].Splits[0].Slg,
		p.Stats[1].Splits[0].Ops,
		p.Stats[1].Splits[0].StrikeoutsPer9Inn,
		p.Stats[1].Splits[0].WalksPer9Inn,
		p.Stats[1].Splits[0].HitsPer9Inn,
		p.Stats[1].Splits[0].RunsScoredPer9,
		p.Stats[1].Splits[0].HomeRunsPer9,
	})
	if err != nil {
		return err
	}
	resp.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", p.FullName))
	f.Write(resp)
	return nil
}

func (p *Pitcher) WriteToExcelSeasonal(resp http.ResponseWriter) error {
	f := excelize.NewFile()
	_, err := f.NewSheet(p.FullName)
	if err != nil {
		return err
	}
	if err = f.DeleteSheet("Sheet1"); err != nil {
		return err
	}
	err = f.SetSheetRow(p.FullName, "A1", &[]interface{}{
		"Team",
		"Season",
		"G",
		"IP",
		"NOP",
		"BF",
		"ERA",
		"HR",
		"BB",
		"K",
		"2B",
		"3B",
		"R",
		"ER",
		"HBP",
		"AO",
		"GO",
		"O",
		"BK",
		"WP",
		"PO",
		"WHIP",
		"AVG",
		"OBP",
		"SLG",
		"OPS",
		"SO9",
		"WP9",
		"HP9",
		"RP9",
		"HRP9",
	})
	if err != nil {
		return err
	}
	for _, splits := range p.Stats {
		for i, stat := range splits.Splits {
			if stat.Team.Name == "" {
				stat.Team.Name = "Total " + stat.Season
			}
			if stat.Season != "" {
				err = f.SetSheetRow(p.FullName, fmt.Sprintf("A%d", i+2), &[]interface{}{
					stat.Team.Name,
					stat.Season,
					stat.GamesPlayed,
					stat.InningsPitched,
					stat.NumberOfPitches,
					stat.BattersFaced,
					stat.Era,
					stat.HomeRuns,
					stat.BaseOnBalls,
					stat.StrikeOuts,
					stat.Doubles,
					stat.Triples,
					stat.Runs,
					stat.EarnedRuns,
					stat.HitByPitch,
					stat.AirOuts,
					stat.GroundOuts,
					stat.Outs,
					stat.Balks,
					stat.WildPitches,
					stat.Pickoffs,
					stat.Whip,
					stat.Avg,
					stat.Obp,
					stat.Slg,
					stat.Ops,
					stat.StrikeoutsPer9Inn,
					stat.WalksPer9Inn,
					stat.HitsPer9Inn,
					stat.RunsScoredPer9,
					stat.HomeRunsPer9,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	if err := f.SetColWidth(p.FullName, "A", "AZ", 5); err != nil {
		return err
	}

	resp.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", p.FullName))
	f.Write(resp)
	return nil
}

func GetPitcherMetrics(id int) (PitcherMetrics, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people/%d?hydrate=stats(type=[metricAverages],metrics=[releaseSpinRate,releaseSpeed])", id))
	if err != nil {
		return PitcherMetrics{}, err
	}
	d := struct {
		Copyright string           `json:"copyright"`
		Metrics   []PitcherMetrics `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return PitcherMetrics{}, err
	}
	if len(d.Metrics) == 0 {
		return PitcherMetrics{}, fmt.Errorf("no data marshaled")
	}
	return d.Metrics[0], nil
}
