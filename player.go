package main

import (
	"math/rand" // Package for generating random numbers
	"time"      // Package for time-related operations
)

// Player represents an NBA player with all their relevant attributes for the guessing game
type Player struct {
	Name        string   // Full name of the player (e.g., "LeBron James")
	Team        string   // Current team or "Retired" for former players
	Position    string   // Playing position (PG, SG, SF, PF, C)
	Height      string   // Player height in feet and inches (e.g., "6'9\"")
	College     string   // College attended or "None" for international/high school players
	DraftYear   int      // Year the player was drafted into the NBA
	PPG         float64  // Career points per game average
	RPG         float64  // Career rebounds per game average
	APG         float64  // Career assists per game average
	Accolades   []string // List of major achievements (championships, awards, etc.)
	TeamHistory []string // List of all teams the player has played for
}

// Global variable to store all loaded players
var players []Player

// initializePlayers loads player data from the NBA API or falls back to hardcoded data
func initializePlayers() error {
	// First attempt to fetch comprehensive player data from NBA API
	apiPlayers, err := fetchAllPlayers()
	if err != nil {
		// If API fails, use the fallback dataset of notable players
		players = getFallbackPlayers()
		return nil // Return nil since fallback is successful
	}

	// If API succeeds, use the fetched data
	players = apiPlayers
	return nil
}

// getRandomPlayer selects and returns a random player from the loaded dataset
func getRandomPlayer() Player {
	// Ensure players are initialized before selecting random player
	if len(players) == 0 {
		initializePlayers() // Initialize if not already done
	}

	// Seed the random number generator with current time for true randomness
	rand.Seed(time.Now().UnixNano())

	// Return a random player from the slice
	return players[rand.Intn(len(players))]
}

// findPlayerByName searches for a player by exact name match
// Returns pointer to player and boolean indicating if found
func findPlayerByName(name string) (*Player, bool) {
	// Iterate through all players in the database
	for _, player := range players {
		// Check for exact name match (case-sensitive)
		if player.Name == name {
			return &player, true // Return pointer to player and true if found
		}
	}
	// Return nil pointer and false if player not found
	return nil, false
}

// getAllPlayerNames returns a slice containing all player names in the database
func getAllPlayerNames() []string {
	// Create slice with capacity equal to number of players
	names := make([]string, len(players))

	// Extract name from each player and add to names slice
	for i, player := range players {
		names[i] = player.Name
	}

	// Return the complete list of player names
	return names
}
