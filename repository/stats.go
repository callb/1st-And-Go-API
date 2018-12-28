package repository

import (
	"../domain"

	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"../utils"
	"database/sql"
	"context"
)

type StatsSqlRepository struct {
	config Configuration
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
