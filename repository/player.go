package repository

import (
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"database/sql"
	"../utils"
	"../domain"
	"fmt"
)

type PlayerSqlRepository struct {
	config Configuration
}

func NewPlayerSqlRepository() PlayerSqlRepository {
	repo := PlayerSqlRepository{}
	repo.config = Configuration {
		"den1.mssql7.gear.host",
		0,
		"nfldata",
		"database!",
	}
	return repo
}

// Get players that have the search text in their name
func (repo PlayerSqlRepository) GetPlayersBySearchText(searchText string) []domain.Player {
	db := repo.getDbConn()
	query := fmt.Sprintf("select id, name, teamAbbr from nfldata.dbo.Player where name like '%%%v%%'", searchText)
	rows, err := db.Query(query)
	utils.CheckForError(err)

	var players []domain.Player
	for rows.Next() {
		var currPlayer domain.Player
		rows.Scan(
			&currPlayer.Id,
			&currPlayer.Name,
			&currPlayer.Team,
		)
		// append the player from the row to the list of players returned
		players = append(players, currPlayer)
	}
	return players
}

// Get stats for a particular player
func (repo PlayerSqlRepository) GetPlayerStatsByPlayerId(playerId string) []domain.PlayerStats {
	db := repo.getDbConn()
	query := fmt.Sprintf(
		"select p.name, p.teamAbbr, ps.gamedate, " +
	"ps.att PassAtt, ps.cmp, ps.yds PassYards, ps.tds PassTds, ps.ints," +
	"ps.twopta PassTwoPta, ps.twoptm PassTwoPta, " +
	"rus.att RushAtt, rus.yds RushYds, rus.tds RushTds," +
	"rus.lng RushLng, rus.lngtd RushLngTd," +
	"rus.twopta RushTwoPta, rus.twoptm RushTwoPtm," +
	"rs.rec, rs.yds RecYds, rs.tds RecTds," +
	"rs.lng RecLng, rs.lngtd RecLngTd," +
	"rs.twopta RecTwoPta, rs.twoptm RecTwoPtm " +
	"from Player p " +
	"join ReceivingStats rs " +
	"on p.nflid = rs.playerid " +
	"join RushingStats rus " +
	"on p.nflid = rus.playerid " +
	"and rs.gamedate = rus.gamedate " +
	"join PassingStats ps " +
	"on p.nflid = ps.playerid " +
	"and rs.gamedate = ps.gamedate " +
	"where p.id = %v " +
			"order by rs.gamedate", playerId)
	fmt.Println(query)

	rows, err := db.Query(query)
	utils.CheckForError(err)

	var playerStats []domain.PlayerStats
	for rows.Next() {
		var currPlayerStats domain.PlayerStats
		rows.Scan(
			&currPlayerStats.Name,
			&currPlayerStats.TeamAbbr,
			&currPlayerStats.GameDate,
			&currPlayerStats.PassingStats.Attempts,
			&currPlayerStats.PassingStats.Completions,
			&currPlayerStats.PassingStats.Yards,
			&currPlayerStats.PassingStats.Touchdowns,
			&currPlayerStats.PassingStats.Interceptions,
			&currPlayerStats.PassingStats.TwoPointAttempts,
			&currPlayerStats.PassingStats.TwoPointSuccesses,
			&currPlayerStats.RushingStats.Attempts,
			&currPlayerStats.RushingStats.Yards,
			&currPlayerStats.RushingStats.Touchdowns,
			&currPlayerStats.RushingStats.Longest,
			&currPlayerStats.RushingStats.LongestTouchdown,
			&currPlayerStats.RushingStats.TwoPointAttempts,
			&currPlayerStats.RushingStats.TwoPointSuccesses,
			&currPlayerStats.ReceivingStats.Receptions,
			&currPlayerStats.ReceivingStats.Yards,
			&currPlayerStats.ReceivingStats.Touchdowns,
			&currPlayerStats.ReceivingStats.Longest,
			&currPlayerStats.ReceivingStats.LongestTouchdown,
			&currPlayerStats.ReceivingStats.TwoPointAttempts,
			&currPlayerStats.ReceivingStats.TwoPointSuccesses,
		)
		playerStats = append(playerStats, currPlayerStats)
	}

	return playerStats
}

// Get the database connection
func (repo PlayerSqlRepository) getDbConn() *sql.DB {
	config := repo.config
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(config.Username, config.Password),
		Host:	  config.Host,
		//Host:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		// Path:  instance, // if connecting to an instance instead of a port
	}
	conn, err := sql.Open("sqlserver", u.String())
	utils.CheckForError(err)

	return conn
}
