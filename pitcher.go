package mlb

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// yearByYear,career
func GetOnePitcherData(id string, span string, season string) (Pitcher, error) {
	fmt.Println(season)
	fmt.Println(fmt.Sprintf(BASE+"people/%s/?hydrate=stats(group=[pitching],type=[%s],seasons=%s),currentTeam", id, span, season))
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

func GetManyPitcherData(ids string) ([]Pitcher, error) {
	resp, err := http.Get(fmt.Sprintf(BASE+"people?personIds=%s&hydrate=stats(group=[pitching],type=[season]),currentTeam", ids))
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
