package mlb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type PlayerHolder struct {
	Copyright string           `json:"copyright"`
	People    []PlayerOverview `json:"people"`
}

type PlayerOverview struct {
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
	CurrentTeam        `json:"currentTeam"`
	PrimaryPosition    `json:"primaryPosition"`
	UseName            string `json:"useName"`
	UseLastName        string `json:"useLastName"`
	MiddleName         string `json:"middleName"`
	BoxscoreName       string `json:"boxscoreName"`
	NickName           string `json:"nickName"`
	Gender             string `json:"gender"`
	IsPlayer           bool   `json:"isPlayer"`
	IsVerified         bool   `json:"isVerified"`
	DraftYear          int    `json:"draftYear"`
	Stats              []Stat `json:"stats"`
	MLBDebut           string `json:"mlbDebutDate"`
	Batside            `json:"batSide"`
	PitchHand          `json:"pitchHand"`
	NFL                string  `json:"nameFirstLast"`
	NameSlug           string  `json:"nameSlug"`
	FirstLastName      string  `json:"firstLastName"`
	LastFirstName      string  `json:"lastFirstName"`
	LastInitName       string  `json:"lastInitName"`
	InitLastName       string  `json:"initLastName"`
	FullFMLName        string  `json:"fullFMLName"`
	FullLFM            string  `json:"fullLFMName"`
	StrikeZoneTop      float64 `json:"strikeZoneTop"`
	StrikeZoneBottom   float64 `json:"strikeZoneBottom"`
}

type CurrentTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Batside struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PitchHand struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PrimaryPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type Stat struct {
	Type       `json:"type"`
	Group      `json:"group"`
	Splits     []Split       `json:"splits"`
	Exemptions []interface{} `json:"exemptions"`
}

type Type struct {
	DisplayName string `json:"displayName"`
}

type Group struct {
	DisplayName string `json:"displayName"`
}

type Split struct {
	Season      string `json:"season"`
	SplitStat   `json:"stat"`
	SplitTeam   `json:"team"`
	SplitPlayer `json:"player"`
	SplitLeague `json:"league"`
	SplitSport  `json:"sport"`
	GameType    string `json:"gameType"`
}

type SplitTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type SplitSport struct {
	ID           int    `json:"id"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
}

type SplitLeague struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type SplitPlayer struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type SplitStat struct {
	GamesPlayed            int    `json:"gamesPlayed"`
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
	PlateAppearances       int    `json:"plateAppearances"`
	TotalBases             int    `json:"totalBases"`
	Rbi                    int    `json:"rbi"`
	LeftOnBase             int    `json:"leftOnBase"`
	SacBunts               int    `json:"sacBunts"`
	SacFlies               int    `json:"sacFlies"`
	Babip                  string `json:"babip"`
	GroundOutsToAirouts    string `json:"groundOutsToAirouts"`
	CatchersInterference   int    `json:"catchersInterference"`
	AtBatsPerHomeRun       string `json:"atBatsPerHomeRun"`
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
}

type SimplePlayer struct {
	SimplePlayerInfo     `json:"person"`
	JerseyNumber         string `json:"jerseyNumber"`
	SimplePlayerPosition `json:"position"`
	SimplePlayerStatus   `json:"status"`
	ParentTeamID         int `json:"parentTeamId"`
}

type SimplePlayerInfo struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type SimplePlayerPosition struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type SimplePlayerStatus struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PlayerSearch struct {
	Name string
	ID   int
}

func GetAllPlayers() error {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	team, err := AllMLBTeams()
	if err != nil {
		return err
	}

	for _, t := range team {
		resp, err := http.Get(BASE + fmt.Sprintf("teams/%d/roster?rosterType=40Man", t.ID))
		if err != nil {
			return err
		}
		var p struct {
			Copyright string         `json:"copyright"`
			Players   []SimplePlayer `json:"roster"`
		}
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			return err
		}
		for _, player := range p.Players {
			err = db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("players"))
				if err != nil {
					return err
				}
				err = b.Put([]byte(player.SimplePlayerInfo.FullName), []byte(strconv.Itoa(player.ID)))
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	defer db.Close()
	return nil
}
