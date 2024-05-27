package mlb

const (
	GAMEURL string = "https://statsapi.mlb.com/api/v1.1/"
	// game/745335/feed/live
)

type Game struct {
	Copyright string `json:"copyright"`
	Pk        int    `json:"gamePk"`
	Link      string `json:"link"`
	GameData  `json:"gameData"`
	LiveData  `json:"liveData"`
}

type GameData struct {
	GameInfo        `json:"game"`
	TeamHolder      `json:"teams"`
	Players         map[string]PlayerOverview `json:"players"`
	ProbablePitcher `json:"probablePitchers"`
	GameStatus      `json:"status"`
	GameDate        `json:"datetime"`
}

type GameDate struct {
	DateTime string `json:"datetime"`
}

type GameInfo struct {
	Pk              int    `json:"pk"`
	Type            string `json:"type"`
	DoubleHeader    string `json:"doubleHeader"`
	Id              string
	GameType        string
	TieBreaker      string
	GameNumber      int
	CalendarEventID string
}

type ProbablePitcher struct {
	Away PitcherInfo `json:"away"`
	Home PitcherInfo `json:"home"`
}

type PitcherInfo struct {
	Id   int    `json:"id"`
	Name string `json:"fullName"`
	Link string `json:"link"`
}

type TeamHolder struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type LiveData struct {
	Boxscore  `json:"boxscore"`
	LineScore `json:"linescore"`
}

type Boxscore struct {
	BoxscoreTeamHolder `json:"teams"`
}

type BoxscoreTeam struct {
	BoxscoreTeamStats `json:"teamStats"`
	Pitchers          []int `json:"pitchers"`
	BullPen           []int `json:"bullpen"`
	Batters           []int `json:"batters"`
	Bench             []int `json:"bench"`
	BattingOrder      []int `json:"battingOrder"`
}

type BoxscoreTeamStats struct {
	BoxscoreTeamBatting  `json:"batting"`
	BoxscoreTeamPitching `json:"pitching"`
}

type BoxscoreTeamBatting struct {
	FlyOuts              int    `json:"flyOuts"`
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
	GroundIntoTriplePlay int    `json:"groundIntoTriplePlay"`
	PlateAppearances     int    `json:"plateAppearances"`
	TotalBases           int    `json:"totalBases"`
	Rbi                  int    `json:"rbi"`
	LeftOnBase           int    `json:"leftOnBase"`
	SacBunts             int    `json:"sacBunts"`
	SacFlies             int    `json:"sacFlies"`
	CatchersInterference int    `json:"catchersInterference"`
	Pickoffs             int    `json:"pickoffs"`
	AtBatsPerHomeRun     string `json:"atBatsPerHomeRun"`
	PopOuts              int    `json:"popOuts"`
	LineOuts             int    `json:"lineOuts"`
}

type BoxscoreTeamPitching struct {
}

type BoxscoreTeamHolder struct {
	Away BoxscoreTeam `json:"away"`
	Home BoxscoreTeam `json:"home"`
}

// in game stats

type StatHolder struct {
	Splits []InGameSplit `json:"splits"`
}

type InGameSplit struct {
	InGameStat `json:"stat"`
	Type       string `json:"type"`
	Group      string `json:"group"`
}

type InGameStat struct {
	Summary                string `json:"summary"`
	GamesPlayed            int    `json:"gamesPlayed"`
	FlyOuts                int    `json:"flyOuts"`
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
	AtBats                 int    `json:"atBats"`
	CaughtStealing         int    `json:"caughtStealing"`
	StolenBases            int    `json:"stolenBases"`
	StolenBasePercentage   string `json:"stolenBasePercentage"`
	GroundIntoDoublePlay   int    `json:"groundIntoDoublePlay"`
	GroundIntoTriplePlay   int    `json:"groundIntoTriplePlay"`
	PlateAppearances       int    `json:"plateAppearances"`
	TotalBases             int    `json:"totalBases"`
	Rbi                    int    `json:"rbi"`
	LeftOnBase             int    `json:"leftOnBase"`
	SacBunts               int    `json:"sacBunts"`
	SacFlies               int    `json:"sacFlies"`
	CatchersInterference   int    `json:"catchersInterference"`
	Pickoffs               int    `json:"pickoffs"`
	AtBatsPerHomeRun       string `json:"atBatsPerHomeRun"`
	Note                   string `json:"note"`
	GamesStarted           int    `json:"gamesStarted"`
	NumberOfPitches        int    `json:"numberOfPitches"`
	InningsPitched         string `json:"inningsPitched"`
	Wins                   int    `json:"wins"`
	Losses                 int    `json:"losses"`
	Saves                  int    `json:"saves"`
	SaveOpportunities      int    `json:"saveOpportunities"`
	Holds                  int    `json:"holds"`
	BlownSaves             int    `json:"blownSaves"`
	EarnedRuns             int    `json:"earnedRuns"`
	BattersFaced           int    `json:"battersFaced"`
	Outs                   int    `json:"outs"`
	GamesPitched           int    `json:"gamesPitched"`
	CompleteGames          int    `json:"completeGames"`
	Shutouts               int    `json:"shutouts"`
	PitchesThrown          int    `json:"pitchesThrown"`
	Balls                  int    `json:"balls"`
	Strikes                int    `json:"strikes"`
	StrikePercentage       string `json:"strikePercentage"`
	HitBatsmen             int    `json:"hitBatsmen"`
	Balks                  int    `json:"balks"`
	WildPitches            int    `json:"wildPitches"`
	GamesFinished          int    `json:"gamesFinished"`
	RunsScoredPer9         string `json:"runsScoredPer9"`
	HomeRunsPer9           string `json:"homeRunsPer9"`
	InheritedRunners       int    `json:"inheritedRunners"`
	InheritedRunnersScored int    `json:"inheritedRunnersScored"`
	PassedBall             int    `json:"passedBall"`
}

type LineScore struct {
	CurrentInning        int    `json:"currentInning"`
	CurrentInningOrdinal string `json:"currentInningOrdinal"`
	InningState          string `json:"inningState"`
	InningHalf           string `json:"inningHalf"`
	IsTopInning          bool   `json:"isTopInning"`
	ScheduledInnings     int    `json:"scheduledInnings"`
}
