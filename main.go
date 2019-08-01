package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"./repository"
	"./update"
	"gopkg.in/matryer/respond.v1"
	"fmt"
)

var (
	playerRepository = repository.NewPlayerSqlRepository()
)

func main() {
	// Initialize ticker for update data process
	go startUpdateDataProcess()

	router := mux.NewRouter()

	/** API Routes **/
	router.HandleFunc("/api/search/player/{searchText}", getPlayersBySearchText)
	router.HandleFunc("/api/player/{playerId}", getPlayerStatsByPlayerId)



	// Get the server information and start the server
	fmt.Println("Starting server...")
	http.ListenAndServe(":8080", router)
}

// Init the ticker so data is updated in intervals
func startUpdateDataProcess() {
	// Run every 12 Hours
	ticker := time.NewTicker(12 * time.Hour)
	for tick := range ticker.C {
		fmt.Println("Update data process started at: ", tick)
		update.StartUpdateDataProcess()
	}
}

// get all players whose first or last names start with the search text
func getPlayersBySearchText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	searchText := mux.Vars(r)["searchText"]
	players := playerRepository.GetPlayersBySearchText(searchText)
	respond.With(w, r, http.StatusOK, players)
}

// get player data for a particular player id
func getPlayerStatsByPlayerId(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	playerId := mux.Vars(r)["playerId"]
	playerData := playerRepository.GetPlayerStatsByPlayerId(playerId)
	respond.With(w, r, http.StatusOK, playerData)

}


// Get the server information
func getServer(router http.Handler) http.Server {
	srv := &http.Server {
		Handler:      router,
		Addr:         ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return *srv
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}