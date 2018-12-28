package domain

import (
	"strings"
	"time"
)

// All information for a particular player
type PlayerStats struct {
	Name string
	TeamAbbr string
	GameDate time.Time
	PassingStats PassingStats
	RushingStats RushingStats
	ReceivingStats ReceivingStats
}


// Passing Stats for a player
type PassingStats struct {
	Attempts int `json:"att"`
	Completions int `json:"cmp"`
	Yards int `json:"yds"`
	Touchdowns int `json:"tds"`
	Interceptions int `json:"ints"`
	TwoPointAttempts int `json:"twopta"`
	TwoPointSuccesses int `json:"twoptm"`
}

// Rushing stats for a player
type RushingStats struct {
	Attempts int `json:"att"`
	Yards int `json:"yds"`
	Touchdowns int `json:"tds"`
	Longest int `json:"lng"`
	LongestTouchdown int `json:"lngtd"`
	TwoPointAttempts int `json:"twopta"`
	TwoPointSuccesses int `json:"twoptm"`
}

// Receiving stats for a player
type ReceivingStats struct {
	Receptions int `json:"rec"`
	Yards int `json:"yds"`
	Touchdowns int `json:"tds"`
	Longest int `json:"lng"`
	LongestTouchdown int `json:"lngtd"`
	TwoPointAttempts int `json:"twopta"`
	TwoPointSuccesses int `json:"twoptm"`
}

// Return the stats label for a given stats type
func GetStatsLabels(statsType string) []string {
	switch strings.ToLower(statsType) {
	case "passing":
		return []string {
			"name",
			"att",
			"cmp",
			"yds",
			"tds",
			"ints",
			"twopta",
			"twoptm",
		}
	case "rushing":
		return []string {
			"name",
			"att",
			"yds",
			"tds",
			"lng",
			"lngtd",
			"twopta",
			"twoptm",
		}
	case "receiving":
		return []string {
			"name",
			"rec",
			"yds",
			"tds",
			"lng",
			"lngtd",
			"twopta",
			"twoptm",
		}

	}

	return nil
}

func NewPassingStats(stats map[string]interface{}) PassingStats {
	att, okAtt := stats["att"].(float64)
	cmp, okCmp := stats["cmp"].(float64)
	yds, okYds := stats["yds"].(float64)
	tds, okTds := stats["tds"].(float64)
	ints, okInts := stats["ints"].(float64)
	twopta, okTwoPta := stats["twopta"].(float64)
	twoptm, okTwoPtm := stats["twoptm"].(float64)


	// if can't type assert a field, return an empty Passing stats object
	if !okAtt || !okCmp || !okYds || !okTds || !okInts || !okTwoPta || !okTwoPtm {
		return PassingStats{}
	}

	return PassingStats {
		int(att),
		int(cmp),
		int(yds),
		int(tds),
		int(ints),
		int(twopta),
		int(twoptm),
	}
}

func NewRushingStats(stats map[string]interface{}) RushingStats {
	att, okAtt := stats["att"].(float64)
	yds, okYds := stats["yds"].(float64)
	tds, okTds := stats["tds"].(float64)
	lng, okLng := stats["lng"].(float64)
	lngtd, okLngtd := stats["lngtd"].(float64)
	twopta, okTwoPta := stats["twopta"].(float64)
	twoptm, okTwoPtm := stats["twoptm"].(float64)


	// if can't type assert a field, return an empty Rushing stats object
	if !okAtt || !okYds || !okTds || !okLng || !okLngtd || !okTwoPta || !okTwoPtm {
		return RushingStats{}
	}

	return RushingStats {
		int(att),
		int(yds),
		int(tds),
		int(lng),
		int(lngtd),
		int(twopta),
		int(twoptm),
	}
}

func NewReceivingStats(stats map[string]interface{}) ReceivingStats {
	rec, okRec := stats["rec"].(float64)
	yds, okYds := stats["yds"].(float64)
	tds, okTds := stats["tds"].(float64)
	lng, okLng := stats["lng"].(float64)
	lngtd, okLngtd := stats["lngtd"].(float64)
	twopta, okTwoPta := stats["twopta"].(float64)
	twoptm, okTwoPtm := stats["twoptm"].(float64)


	// if can't type assert a field, return an empty Rushing stats object
	if !okRec || !okYds || !okTds || !okLng || !okLngtd || !okTwoPta || !okTwoPtm {
		return ReceivingStats{}
	}

	return ReceivingStats {
		int(rec),
		int(yds),
		int(tds),
		int(lng),
		int(lngtd),
		int(twopta),
		int(twoptm),
	}
}
