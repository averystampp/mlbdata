package mlb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func retrievePitcher(ids []int) []PlayerOverview {
	var list string
	for i, id := range ids {
		if i+1 == len(ids) {
			list += strconv.Itoa(id)
		} else {
			list += strconv.Itoa(id) + ","
		}
	}
	resp, err := http.Get(BASE + fmt.Sprintf("people?personIds=%s&hydrate=stats(group=[pitching],type=[season])", list))
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}
	var pitcher PlayerHolder
	err = json.NewDecoder(resp.Body).Decode(&pitcher)
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}

	if len(pitcher.People) == 0 {
		return []PlayerOverview{}
	}

	return pitcher.People
}

func singlePitcher(id int) PlayerOverview {
	resp, err := http.Get(BASE + fmt.Sprintf("people?personIds=%d&hydrate=stats(group=[pitching],type=[season])", id))
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}
	var pitcher PlayerHolder
	err = json.NewDecoder(resp.Body).Decode(&pitcher)
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}

	if len(pitcher.People) == 0 {
		return PlayerOverview{}
	}

	return pitcher.People[0]
}

func retrieveBatters(ids []int) []PlayerOverview {
	var list string
	for i, id := range ids {
		if i+1 == len(ids) {
			list += strconv.Itoa(id)
		} else {
			list += strconv.Itoa(id) + ","
		}
	}
	resp, err := http.Get(BASE + fmt.Sprintf("people?personIds=%s&hydrate=stats(group=[batting],type=[season])", list))
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}
	var players PlayerHolder
	err = json.NewDecoder(resp.Body).Decode(&players)
	if err != nil {
		log.Printf("error retrieving pitching %s", err)
	}

	if len(players.People) == 0 {
		return []PlayerOverview{}
	}

	return players.People
}

type PastPitcher struct {
	InGameSplit
	PlayerOverview
}

// playerType must be "pitching" || "hitting" || "fielding"
func playerGameStats(pitchIds []int, gameId int, status string) []PastPitcher {
	var (
		splHolder []PastPitcher
		p         = struct {
			Copyright string           `json:"copyright"`
			Info      []PlayerOverview `json:"people"`
		}{}
		t = struct {
			Copyright string       `json:"copyright"`
			P         []StatHolder `json:"stats"`
		}{}
	)
	if status == "I" {
		pitchIds = pitchIds[:len(pitchIds)-1]
	}

	for _, id := range pitchIds {
		resp, err := http.Get(BASE + fmt.Sprintf("people/%d/stats/game/%d", id, gameId))
		if err != nil {
			return nil
		}

		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return nil
		}
		respTwo, err := http.Get(BASE + fmt.Sprintf("people/%d", id))
		if err != nil {
			return nil
		}
		defer respTwo.Body.Close()
		err = json.NewDecoder(respTwo.Body).Decode(&p)
		if err != nil {
			return nil
		}
		for _, spl := range t.P {
			for _, stat := range spl.Splits {
				if stat.Group == "pitching" {
					splHolder = append(splHolder, PastPitcher{
						InGameSplit:    stat,
						PlayerOverview: p.Info[0],
					})
				}
			}
		}
	}
	return splHolder
}

func addNums(i, l int) int {
	return i + l
}

func subNums(i, l int) int {
	return i - l
}

func calculateERA(er int, ip string) string {
	ipasFloat, err := strconv.ParseFloat(ip, 64)
	if err != nil {
		return ""
	}

	a := float64((9 * er)) / ipasFloat
	return fmt.Sprintf("%.2f", a)
}

func localTime(t string) string {
	local, err := time.ParseInLocation(time.RFC3339, t, time.Local)
	if err != nil {
		return err.Error()
	}
	local = local.In(time.Local)
	y, m, d := local.Date()

	if local.Day() == time.Now().Day() && local.Year() == time.Now().Year() {
		return "@ " + local.Format(time.Kitchen)
	}

	return fmt.Sprintf("%s %d, %d", m.String(), d, y) + " " + local.Format(time.Kitchen)
}
