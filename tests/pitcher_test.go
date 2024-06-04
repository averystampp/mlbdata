package mlbtests

import (
	"strconv"
	"strings"
	"testing"

	"github.com/averystampp/mlb"
)

func TestGetOnePitcher(t *testing.T) {
	pitcher, err := mlb.GetOnePitcherData("669923", "career", "")
	if err != nil {
		t.Fatal(err)
	}

	if pitcher.FullName != "George Kirby" {
		t.Fatal("pitcher did return correct name")
	}
}

func TestGetManyPitchers(t *testing.T) {
	pitchers := []string{"669923", "605135", "669203"}
	joined := strings.Join(pitchers, ",")
	returnedPitchers, err := mlb.GetManyPitcherData(joined, "season", "2024")
	if err != nil {
		t.Fatal(err)
	}
	for _, pitcher := range returnedPitchers {
		for i, p := range pitchers {
			if strconv.Itoa(pitcher.ID) == p {
				pitchers = RemoveIndex(pitchers, i)
			}
		}

	}

	if len(pitchers) != 0 {
		t.Fatal("not all pitchers were found in list \"pitchers\"")
	}
}

func RemoveIndex(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
