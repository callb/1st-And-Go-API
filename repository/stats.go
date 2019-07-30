package repository

import (
	"../domain"

	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"../utils"
	"database/sql"
	"context"
	"fmt"
	"strings"
	"time"
)

type StatsSqlRepository struct {
	config Configuration
}

type tvpSaveData struct {
	playerKey string
	playerData domain.PlayerStats
	statsType string
}

func NewStatsSqlRepository() StatsSqlRepository {
	repo := StatsSqlRepository {}
	repo.config = Configuration {
		"den1.mssql7.gear.host",
		0,
		"nfldata",
		"database!",
	}
	return repo
}

func newTvpSaveData(playerKey string, playerData domain.PlayerStats, statsType string) tvpSaveData {
	return tvpSaveData {
		playerKey: playerKey,
		playerData: playerData,
		statsType: statsType,
	}
}

func (repo StatsSqlRepository) SavePlayerStats(key string, player domain.PlayerStats) {
	conn := repo.getDbConn()
	ctx := context.Background()

	// Save PlayerStats Information
	_, err := conn.ExecContext(ctx, "SavePlayer",
		sql.Named("nflid", key),
		sql.Named("name", player.Name),
		sql.Named("teamAbbr", player.TeamAbbr),
	)
	utils.CheckForError(err)

	// Save Passing Stats
	passingStats := player.PassingStats
	_, err = conn.ExecContext(ctx, "SavePassingStats",
		sql.Named("playerid", key),
		sql.Named("gamedate", player.GameDate),
		sql.Named("att", passingStats.Attempts),
		sql.Named("cmp", passingStats.Completions),
		sql.Named("yds", passingStats.Yards),
		sql.Named("tds", passingStats.Touchdowns),
		sql.Named("ints", passingStats.Interceptions),
		sql.Named("twopta", passingStats.TwoPointAttempts),
		sql.Named("twoptm", passingStats.TwoPointSuccesses),
	)
	utils.CheckForError(err)

	// Save Rushing Stats
	rushingStats := player.RushingStats
	_, err = conn.ExecContext(ctx, "SaveRushingStats",
		sql.Named("playerid", key),
		sql.Named("gamedate", player.GameDate),
		sql.Named("att", rushingStats.Attempts),
		sql.Named("yds", rushingStats.Yards),
		sql.Named("tds", rushingStats.Touchdowns),
		sql.Named("lng", rushingStats.Longest),
		sql.Named("lngtd", rushingStats.LongestTouchdown),
		sql.Named("twopta", rushingStats.TwoPointAttempts),
		sql.Named("twoptm", rushingStats.TwoPointSuccesses),
	)
	utils.CheckForError(err)

	// Save receiving stats
	receivingStats := player.ReceivingStats
	_, err = conn.ExecContext(ctx, "SaveReceivingStats",
		sql.Named("playerid", key),
		sql.Named("gamedate", player.GameDate),
		sql.Named("rec", receivingStats.Receptions),
		sql.Named("yds", receivingStats.Yards),
		sql.Named("tds", receivingStats.Touchdowns),
		sql.Named("lng", receivingStats.Longest),
		sql.Named("lngtd", receivingStats.LongestTouchdown),
		sql.Named("twopta", receivingStats.TwoPointAttempts),
		sql.Named("twoptm", receivingStats.TwoPointSuccesses),
	)
	utils.CheckForError(err)
}

// Save the player stats in the given map of player key/id to player data
func (repo StatsSqlRepository) SavePlayerStatsBatch(statsMap map[string]domain.PlayerStats) {
	playerTvpSaveQuery := ""
	passingTvpSaveQuery := ""
	rushingTvpSaveQuery := ""
	receivingTvpSaveQuery := ""
	// Iterate through each player in the player data and add to the save query
	for playerKey, playerData := range statsMap {
		passingTvpSaveData := newTvpSaveData(playerKey, playerData, "passing")
		passingTvpSaveQuery = addToStatsTvp(passingTvpSaveQuery, passingTvpSaveData)

		rushingTvpSaveData := newTvpSaveData(playerKey, playerData, "rushing")
		rushingTvpSaveQuery = addToStatsTvp(rushingTvpSaveQuery, rushingTvpSaveData)

		receivingTvpSaveData := newTvpSaveData(playerKey, playerData, "receiving")
		receivingTvpSaveQuery = addToStatsTvp(receivingTvpSaveQuery, receivingTvpSaveData)

		playerTvpSaveQuery = addToPlayerTvp(playerTvpSaveQuery, playerKey, playerData)
	}

	playerTvpSaveQuery += "\nexec SavePlayer @records = @r"
	passingTvpSaveQuery += "\nexec SavePassingStats @records = @r"
	rushingTvpSaveQuery += "\nexec SaveRushingStats @records = @r"
	receivingTvpSaveQuery += "\nexec SaveReceivingStats @records = @r"

	conn := repo.getDbConn()
	defer conn.Close()
	executeModifyQuery(*conn, playerTvpSaveQuery)
	executeModifyQuery(*conn, passingTvpSaveQuery)
	executeModifyQuery(*conn, rushingTvpSaveQuery)
	executeModifyQuery(*conn, receivingTvpSaveQuery)

}

// Add the next line of data to the player tvp
func addToPlayerTvp(tvpCurrQuery string, playerKey string, playerData domain.PlayerStats) string {
	newQueryLine := fmt.Sprintf("\nSELECT '%v', '%v', '%v'", playerKey, playerData.Name, playerData.TeamAbbr)

	if len(tvpCurrQuery) == 0 {
		tvpCurrQuery += fmt.Sprintf("\nDECLARE @r PlayerTvp\n")
		tvpCurrQuery += fmt.Sprintf("INSERT INTO @r %v", newQueryLine)
		return tvpCurrQuery
	}

	tvpCurrQuery += fmt.Sprintf(" UNION %v", newQueryLine)
	return tvpCurrQuery
}

// Add the next line of data to the given tvp Query for saving stats
func addToStatsTvp(tvpCurrQuery string, data tvpSaveData) string {
	statsType := strings.ToLower(data.statsType)
	newQueryLine := ""
	gameDate := formatDateForQuery(data.playerData.GameDate)

	// format the new query string line for the given data
	switch statsType {
	case "passing":
		newQueryLine = fmt.Sprintf(
			"select '%v', '%v', %v, %v, %v, %v, %v, %v, %v",
			data.playerKey,
			gameDate,
			data.playerData.PassingStats.Attempts,
			data.playerData.PassingStats.Completions,
			data.playerData.PassingStats.Yards,
			data.playerData.PassingStats.Touchdowns,
			data.playerData.PassingStats.Interceptions,
			data.playerData.PassingStats.TwoPointAttempts,
			data.playerData.PassingStats.TwoPointSuccesses)
		break
	case "rushing":
		newQueryLine = fmt.Sprintf(
			"select '%v', '%v', %v, %v, %v, %v, %v, %v, %v",
			data.playerKey,
			gameDate,
			data.playerData.RushingStats.Attempts,
			data.playerData.RushingStats.Yards,
			data.playerData.RushingStats.Touchdowns,
			data.playerData.RushingStats.Longest,
			data.playerData.RushingStats.LongestTouchdown,
			data.playerData.RushingStats.TwoPointAttempts,
			data.playerData.RushingStats.TwoPointSuccesses)
		break
	case "receiving":
		newQueryLine = fmt.Sprintf(
			"select '%v', '%v', %v, %v, %v, %v, %v, %v, %v",
			data.playerKey,
			gameDate,
			data.playerData.ReceivingStats.Receptions,
			data.playerData.ReceivingStats.Yards,
			data.playerData.ReceivingStats.Touchdowns,
			data.playerData.ReceivingStats.Longest,
			data.playerData.ReceivingStats.LongestTouchdown,
			data.playerData.ReceivingStats.TwoPointAttempts,
			data.playerData.ReceivingStats.TwoPointSuccesses)
	}

	// If the current save query has no data, add initial tvp declaration
	if len(tvpCurrQuery) == 0 {
		tvpCurrQuery += fmt.Sprintf("DECLARE @r %vStatsTvp\n", statsType)
		tvpCurrQuery += fmt.Sprintf("INSERT INTO @r")
		tvpCurrQuery += fmt.Sprintf("\n%v", newQueryLine)
		return tvpCurrQuery
	}

	// If more than 1 insert row, add in UNION then the new data
	tvpCurrQuery += fmt.Sprintf(" UNION\n%v", newQueryLine)
	return tvpCurrQuery

}

// Format the datetime for the save querystring
func formatDateForQuery(dateTime time.Time) string {
	year, month, day := dateTime.Date()
	return fmt.Sprintf("%v %v, %v", month.String(), day, year)
}


func (repo StatsSqlRepository) getDbConn() *sql.DB {
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

// execute a sql query that inserts or updates data and doesn't return any rows
func executeModifyQuery(conn sql.DB, query string) {
	ctx := context.Background()
	if len(query) == 0 {
		return
	}
	conn.QueryContext(ctx, query)
}
