package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log"
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


	// Get the server information and start the server
	server := getServer(router)
	fmt.Println("Starting server...")
	log.Fatal(server.ListenAndServe())
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
	searchText := mux.Vars(r)["searchText"]
	players := playerRepository.GetPlayersBySearchText(searchText)
	respond.With(w, r, http.StatusOK, players)
}


// Get the server information
func getServer(router http.Handler) http.Server {
	srv := &http.Server {
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return *srv
}