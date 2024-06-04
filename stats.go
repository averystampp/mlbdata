package mlb

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xuri/excelize/v2"
)

var counter atomic.Int32

// magic number is 737 pitchers
func LeagePitchingStats() error {
	counter.Add(1)
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()
	teams, err := AllMLBTeams()
	if err != nil {
		return err
	}
	pitchers := make(chan []FortyManSearch)
	wg := &sync.WaitGroup{}
	for _, team := range teams {
		wg.Add(1)
		go func(id int) {
			p, _, err := Roster(strconv.Itoa(id))
			if err != nil {
				return
			}
			pitchers <- p
		}(team.ID)
	}
	f := excelize.NewFile()
	defer f.Close()
	go func() {
		for {
			var ids []string
			p := <-pitchers
			for _, person := range p {
				ids = append(ids, strconv.Itoa(person.Person.ID))
			}
			s := strings.Join(ids, ",")
			pitchs, err := GetManyPitcherData(s, "season", "2023")
			if err != nil {
				return
			}
			wg.Add(1)
			go process(pitchs, wg, f)
			wg.Done()
		}
	}()
	wg.Wait()
	counter.Swap(0)

	return f.SaveAs("some.xlsx")
}

func process(pitchers []Pitcher, wg *sync.WaitGroup, f *excelize.File) {
	defer wg.Done()

	for _, pitcher := range pitchers {
		row := pitcher.excelrow()
		if row != nil {
			f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", counter.Add(1)), pitcher.excelrow())
		}
	}

	f.SetSheetRow("Sheet1", "A1", &[]interface{}{
		"Player Name",
		"Season",
		"Team",
		"G",
		"IP",
		"ERA",
	})

}

func (p *Pitcher) excelrow() *[]interface{} {
	row := &[]interface{}{}
	if len(p.Stats) == 0 {
		return nil
	}

	if len(p.Stats[0].Splits) > 0 {
		var t []string
		for _, team := range p.Stats[0].Splits {
			if team.Team.ID != 0 {
				t = append(t, team.Team.Name)
			}
		}
		p.Stats[0].Splits[0].Team.Name = strings.Join(t, ",")

	}

	*row = append(*row,
		p.FullName,
		p.Stats[0].Splits[0].Season,
		p.Stats[0].Splits[0].Team.Name,
		p.Stats[0].Splits[0].GamesPlayed,
		p.Stats[0].Splits[0].InningsPitched,
		p.Stats[0].Splits[0].Era,
		p.Stats[0].Splits[0].StrikeOuts,
	)

	return row

}
