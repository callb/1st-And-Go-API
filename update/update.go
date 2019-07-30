package update

import (
	"encoding/json"
	"../domain"
	"../repository"
	"../utils"
	"strings"
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
)

var (
	// ** Modify year, month, day to run through the specified date up through today **
	year = 2018
	month = 8
	day = 1
	gameNum = 0
	gameDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
)


func StartUpdateDataProcess() {
	// ** Comment this out and edit the global variables to run for more days than today **
	//setDateAsToday()

	// Run for each day up through the current date
	for !isDateTomorrow(gameDate) {
		// Update data for the particular day and game num
		isAbleToGetData := updateDataForGameKey()

		// If no more game data for that day, go to the next day
		if !isAbleToGetData {
			gameDate = gameDate.AddDate(0, 0, 1)
			gameNum = 0
			continue
		}

		// If game data was found, go to the next game on the same date
		gameNum++
	}
}

// Set the date as today
func setDateAsToday() {
	todayYear, todayMonth, todayDay := time.Now().Date()
	year = todayYear
	month = int(todayMonth)
	day = todayDay
	gameNum = 0
	gameDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Is the given date the day after today?
func isDateTomorrow(date time.Time) bool {
	year, month, day := date.Date()
	tomYear, tomMonth, tomDay := time.Now().AddDate(0, 0, 1).Date()
	return year == tomYear &&
			month == tomMonth &&
			day == tomDay
}

// Update the game data for the given gameKey (Year, month, day, gameNumber)
func updateDataForGameKey() bool {
	gameDateKeyForUrl := getGameDateKey(gameDate, gameNum)
	url := fmt.Sprintf("http://www.nfl.com/liveupdate/game-center/%s/%s_gtd.json",
		gameDateKeyForUrl,
		gameDateKeyForUrl)
	fmt.Println("Updating stats for game key: " + gameDateKeyForUrl)
	response, err := http.Get(url)

	if response.StatusCode == http.StatusNotFound {
		fmt.Println("Data not found for game key: " + gameDateKeyForUrl)
		return false
	}

	utils.CheckForError(err)
	bytes, err := ioutil.ReadAll(response.Body)
	utils.CheckForError(err)
	parseJson(bytes)

	return true
}

func getGameDateKey(dateTime time.Time, gameNum int) string {
	year, month, day := dateTime.Date()
	yearFormatted := addZeroToSingleDigit(year)
	monthFormatted := addZeroToSingleDigit(int(month))
	dayFormatted := addZeroToSingleDigit(int(day))
	gameNumFormatted := addZeroToSingleDigit(gameNum)
	return yearFormatted +
		monthFormatted +
		dayFormatted +
		gameNumFormatted

}

func addZeroToSingleDigit(num int) string {
	numAsString := strconv.Itoa(num)
	if len(numAsString) == 1 {
		numAsString = "0" + numAsString
	}
	return numAsString


}

func parseJson(bytes []byte) {
	var unmarshalled interface{}
	err := json.Unmarshal(bytes, &unmarshalled)
	utils.CheckForError(err)
	data := assertToMap(unmarshalled)

	gameDateKey := getGameDateKey(gameDate, gameNum)

	if len(data) == 0 || !containsKey(data, gameDateKey) {
		return
	}

	value := assertToMap(data[gameDateKey])

	// Parse the Game data
	saveGameData(value)
}

func saveGameData(data map[string]interface{}) {
	if !containsKey(data, "home") || !containsKey(data, "away") {
		return
	}

	homeMap := assertToMap(data["home"])
	homeGameData := getGameDataForTeam(homeMap)

	saveStatsToDb(homeGameData)

	awayMap := assertToMap(data["away"])
	awayGameData := getGameDataForTeam(awayMap)

	saveStatsToDb(awayGameData)
}

func saveStatsToDb(statsMap map[string]domain.PlayerStats) {
	statsRepository := repository.NewStatsSqlRepository()
	/**
	for playerKey, playerData := range statsMap {
		statsRepository.SavePlayerStats(playerKey, playerData)
	}
	**/
	statsRepository.SavePlayerStatsBatch(statsMap)
}

// Get data for the team; home or away fields
func getGameDataForTeam(teamData map[string]interface{}) map[string]domain.PlayerStats {
	if !containsKey(teamData, "abbr") || !containsKey(teamData, "stats") {
		return make(map[string]domain.PlayerStats)
	}

	teamAbbr, ok := teamData["abbr"].(string)

	// Check if team abbreviation was able to be type asserted
	if !ok {
		return make(map[string]domain.PlayerStats)
	}

	teamStats := assertToMap(teamData["stats"])

	passingStatsKey := "passing"
	rushingStatsKey := "rushing"
	receivingStatsKey := "receiving"

	playerData := make(map[string]domain.PlayerStats)
	hasTeamPassingStats := containsKey(teamStats, passingStatsKey)
	hasTeamRushingStats := containsKey(teamStats, rushingStatsKey)
	hasTeamReceivingStats := containsKey(teamStats, receivingStatsKey)

	if hasTeamPassingStats {
		teamPassingStats := assertToMap(teamStats[passingStatsKey])
		playerData = addTeamStatsToPlayerData(teamPassingStats, playerData, passingStatsKey, teamAbbr)
	}

	if hasTeamRushingStats {
		teamRushingStats := assertToMap(teamStats[rushingStatsKey])
		playerData = addTeamStatsToPlayerData(teamRushingStats, playerData, rushingStatsKey, teamAbbr)
	}

	if hasTeamReceivingStats {
		teamReceivingStats := assertToMap(teamStats[receivingStatsKey])
		playerData = addTeamStatsToPlayerData(teamReceivingStats, playerData, receivingStatsKey, teamAbbr)
	}

	return playerData
}

func addTeamStatsToPlayerData(teamStats map[string]interface{}, playerData map[string]domain.PlayerStats,
	statsType string, teamAbbr string) map[string]domain.PlayerStats {
	for playerKey, value := range teamStats {
		statsMap := assertToMap(value)
		statsLabels := domain.GetStatsLabels(statsType)

		// Check if all the stat values exist in the data
		if !containsKeys(statsMap, statsLabels) {
			continue
		}

		var player domain.PlayerStats
		var name string

		name, ok := statsMap["name"].(string)
		if !ok {
			continue
		}


		switch strings.ToLower(statsType) {
			case "passing":
				passingStats := domain.NewPassingStats(statsMap)
				player = getPlayerWithPassingStats(passingStats, playerData, playerKey)
				break
			case "rushing":
				rushingStats := domain.NewRushingStats(statsMap)
				player = getPlayerWithRushingStats(rushingStats, playerData, playerKey)
				break
			case "receiving":
				receivingStats := domain.NewReceivingStats(statsMap)
				player = getPlayerWithReceivingStats(receivingStats, playerData, playerKey)
				break
			default:
				continue

		}

		player.Name = name
		player.TeamAbbr = teamAbbr
		player.GameDate = time.Date(gameDate.Year(),
			gameDate.Month(),
			gameDate.Day(),
			0,
			0,
			0,
			0,
			time.UTC)
		playerData[playerKey] = player
	}

	return playerData
}





// Return a player with their passing stats added to the object
func getPlayerWithPassingStats(passingStats domain.PassingStats,
	playerData map[string]domain.PlayerStats, playerKey string) domain.PlayerStats {
		player := domain.PlayerStats{}

		if containsPlayerKey(playerData, playerKey) {
			player = playerData[playerKey]
		}

		player.PassingStats = passingStats
		return player
}

// Return a player with their rushing stats added to the object
func getPlayerWithRushingStats(rushingStats domain.RushingStats,
	playerData map[string]domain.PlayerStats, playerKey string) domain.PlayerStats {
	player := domain.PlayerStats{}

	if containsPlayerKey(playerData, playerKey) {
		player = playerData[playerKey]
	}

	player.RushingStats = rushingStats
	return player
}

// Return a player with their receiving stats added to the object
func getPlayerWithReceivingStats(receivingStats domain.ReceivingStats,
	playerData map[string]domain.PlayerStats, playerKey string) domain.PlayerStats {
	player := domain.PlayerStats{}

	if containsPlayerKey(playerData, playerKey) {
		player = playerData[playerKey]
	}

	player.ReceivingStats = receivingStats
	return player
}



// Assert an object to a map if it's possible
func assertToMap(i interface{}) map[string]interface{} {
	m, ok := i.(map[string]interface{})

	if !ok {
		return make(map[string]interface{})
	}

	return m
}

// Check if a dictionary contains the given key
func containsKey(dict map[string]interface{}, keyToCheck string) bool {
	for key := range dict {
		if key == keyToCheck {
			return true
		}
	}
	return false
}

// Check if a player map has the given player key string
func containsPlayerKey(dict map[string]domain.PlayerStats, keyToCheck string) bool {
	for key := range dict {
		if key == keyToCheck {
			return true
		}
	}
	return false
}

// check if a dictionary contains the given keys
func containsKeys(dict map[string]interface{}, keysToCheck []string) bool {
	// Iterate through each key in the given list
	for _, keyToCheck := range keysToCheck {
		isFound := false
		// interate through the map to see if the key to check is contained in it
		for dictKey := range dict {
			// if the key is found, set the bool to true
			if dictKey == keyToCheck {
				isFound = true
			}
		}
		// If the isFound bool was never set to true, the map doesn't contain one of the keys
		if !isFound {
			return false
		}
	}

	// All keys were found in the map
	return true
}
