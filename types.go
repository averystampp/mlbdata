package mlb

const (
	BASE string = "https://statsapi.mlb.com/api/v1/"
)

// Team types
// _____________________________________________________________________________________________
type Team struct {
	SpringLeague    `json:"springLeague"`
	AllStarStatus   string `json:"allStarStatus"`
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Link            string `json:"link"`
	Season          int    `json:"season"`
	Venue           `json:"venue"`
	SpringVenue     `json:"springVenue"`
	TeamCode        string `json:"teamCode"`
	FileCode        string `json:"fileCode"`
	Abbreviation    string `json:"abbreviation"`
	TeamName        string `json:"teamName"`
	LocationName    string `json:"locationName"`
	FirstYearOfPlay string `json:"firstYearOfPlay"`
	League          `json:"league"`
	Division        `json:"division"`
	Sport           `json:"sport"`
	ShortName       string `json:"shortName"`
	Record          `json:"record"`
	FranchiseName   string `json:"franchiseName"`
	ClubName        string `json:"clubName"`
	Active          bool   `json:"active"`
}

type Record struct {
	GamesPlayed           int    `json:"gamesPlayed"`
	WildCardGamesBack     string `json:"wildCardGamesBack"`
	LeagueGamesBack       string `json:"leagueGamesBack"`
	SpringLeagueGamesBack string `json:"springLeagueGamesBack"`
	SportGamesBack        string `json:"sportGamesBack"`
	DivisionGamesBack     string `json:"divisionGamesBack"`
	ConferenceGamesBack   string `json:"conferenceGamesBack"`
	LeagueRecord          `json:"leagueRecord"`
}

type SpringLeague struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
}

type Venue struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type SpringVenue struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Division struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Sport struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
	Name string `json:"name"`
}

// Player types
// _____________________________________________________________________________________________
type Player_ struct {
	Person       `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     `json:"position"`
	Status       `json:"status"`
	ParentTeamID int `json:"parentTeamId"`
}

type Person struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type Position struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type Status struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Game types
// _____________________________________________________________________________________________

type GameTeam struct {
	LeagueRecord   `json:"leagueRecord"`
	Score          int `json:"score"`
	GameTeamDetail `json:"team"`
	IsWinner       bool `json:"isWinner"`
	SplitSquad     bool `json:"splitSquad"`
	SeriesNumber   int  `json:"seriesNumber"`
}

type LeagueRecord struct {
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Ties   int    `json:"ties"`
	Pct    string `json:"pct"`
}

type GameTeamDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type GameTeamHolder struct {
	Away GameTeam `json:"away"`
	Home GameTeam `json:"home"`
}

type Content struct {
	Link string `json:"link"`
}
