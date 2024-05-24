package mlb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func addNums(i, l int) int {
	return i + l
}
