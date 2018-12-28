package repository

import (
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"database/sql"
	"../utils"
	"context"
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


func (repo PlayerSqlRepository) GetPlayersBySearchText(searchText string) []domain.Player {
	db := repo.getDbConn()
	ctx := context.Background()
	query := fmt.Sprintf("select name, teamAbbr from nfldata.dbo.PlayerStats where name like '%%%v%%'", searchText)
	rows, err := db.QueryContext(ctx, query,
		sql.Named("searchText", searchText))
	utils.CheckForError(err)

	var players []domain.Player
	for rows.Next() {
		var currPlayer domain.Player
		rows.Scan(
			&currPlayer.Name,
			&currPlayer.Team,
		)
		// append the player from the row to the list of players returned
		players = append(players, currPlayer)
	}
	return players
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
