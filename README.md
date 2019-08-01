# 1st-And-Go-API
1st-And-Go

Public API to get NFL player data
- Uses data service under /update folder to update data in the db
- Endpoints defined in main.go file

API Documentation

GET /api/search/player/{searchText}
- Gets all players that have the search text contained in their name
- Response includes the internal playerId needed for /api/player/{playerId}

GET /api/player/{playerId}
- Gets all game stats for a player given their playerId
