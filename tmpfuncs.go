package mlb

import (
	"fmt"
	"strconv"
)

type PitchMetric struct {
	MetricName   string
	AverageValue string
	MinValue     string
	MaxValue     string
	PitchCount   string
}

func pitchingMetrics(pm PitcherMetrics) map[string][]PitchMetric {
	stats := make(map[string][]PitchMetric)
	for _, stat := range pm.Stats {
		for _, split := range stat.Splits {
			_, ok := stats[split.Stat.Event.Details.Type.Description]
			if ok {
				stats[split.Stat.Event.Details.Type.Description] = append(stats[split.Stat.Event.Details.Type.Description],
					PitchMetric{
						MetricName:   split.Stat.Metric.Name,
						AverageValue: fmt.Sprintf("%.2f %s", split.Stat.Metric.AverageValue, split.Stat.Metric.Unit),
						MaxValue:     fmt.Sprintf("%.2f %s", split.Stat.Metric.MaxValue, split.Stat.Metric.Unit),
						MinValue:     fmt.Sprintf("%.2f %s", split.Stat.Metric.MinValue, split.Stat.Metric.Unit),
						PitchCount:   strconv.Itoa(split.NumOccurrences),
					})
			} else {
				stats[split.Stat.Event.Details.Type.Description] = []PitchMetric{
					{
						MetricName:   split.Stat.Metric.Name,
						AverageValue: fmt.Sprintf("%.2f %s", split.Stat.Metric.AverageValue, split.Stat.Metric.Unit),
						MaxValue:     fmt.Sprintf("%.2f %s", split.Stat.Metric.MaxValue, split.Stat.Metric.Unit),
						MinValue:     fmt.Sprintf("%.2f %s", split.Stat.Metric.MinValue, split.Stat.Metric.Unit),
						PitchCount:   strconv.Itoa(split.NumOccurrences),
					},
				}
			}
		}
	}
	return stats
}

func RemoveDups(p Pitcher) []string {
	var seasons []string
	for _, stat := range p.Stats {
		for _, spl := range stat.Splits {
			if len(seasons) == 0 {
				seasons = append(seasons, spl.Season)
			}
			if spl.Season == seasons[len(seasons)-1] || spl.Season == "" {
				continue
			}
			seasons = append(seasons, spl.Season)
		}
	}
	return seasons
}

func RemoveDupsBatter(b Batter) []string {
	var seasons []string
	for _, stat := range b.Stats {
		for _, spl := range stat.Splits {
			if len(seasons) == 0 {
				seasons = append(seasons, spl.Season)
			}
			if spl.Season == seasons[len(seasons)-1] || spl.Season == "" {
				continue
			}
			seasons = append(seasons, spl.Season)
		}
	}
	return seasons
}

func joinIds(players []FortyManSearch) string {
	var ids string
	for i, p := range players {
		if i != len(players)-1 {
			ids += strconv.Itoa(p.Person.ID) + ","
		} else {
			ids += strconv.Itoa(p.Person.ID)
		}
	}
	return ids
}

func joinSeasons(s []string) string {
	var seasons string
	for i, p := range s {
		if i != len(s)-1 {
			seasons += p + ","
		} else {
			seasons += p
		}

	}
	return seasons
}
