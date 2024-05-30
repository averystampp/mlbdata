package mlb

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetOneBatterData(link string) (Batter, error) {
	resp, err := http.Get(fmt.Sprintf("https://statsapi.mlb.com%s?hydrate=stats(group=[hitting],type=[season]),currentTeam", link))
	if err != nil {
		return Batter{}, err
	}
	d := struct {
		Copyright string   `json:"copyright"`
		Batter    []Batter `json:"people"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return Batter{}, err
	}
	if len(d.Batter) == 0 {
		return Batter{}, fmt.Errorf("no data marshaled")
	}
	return d.Batter[0], nil
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
