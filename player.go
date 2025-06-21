package main

import (
	"math/rand" // Package for generating random numbers
	"strings"   // Package for string manipulation functions
	"time"      // Package for time-related operations
)

// Player represents an NBA player with all their relevant attributes for the guessing game
type Player struct {
	Name         string // Full name of the player (e.g., "LeBron James")
	Team         string // Current team or "Retired" for former players
	Position     string // Playing position (PG, SG, SF, PF, C)
	Height       string // Player height in feet and inches (e.g., "6'9\"")
	College      string // College attended or "None" for international/high school players
	DraftYear    int    // Year the player was drafted into the NBA
	DraftRound   int    // Round the player was drafted in (1-2, or 0 for undrafted)
	DraftNumber  int    // Overall pick number in the draft (1-60, or 0 for undrafted)
	JerseyNumber string // Current jersey number (or "Unknown" if not available)
	Country      string // Country of origin
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

// findPlayerByName searches for a player by case-insensitive name match
// Returns pointer to player and boolean indicating if found
func findPlayerByName(name string) (*Player, bool) {
	// Convert input name to lowercase for case-insensitive comparison
	lowerName := strings.ToLower(strings.TrimSpace(name))

	// Iterate through all players in the database
	for _, player := range players {
		// Check for case-insensitive name match
		if strings.ToLower(player.Name) == lowerName {
			return &player, true // Return pointer to player and true if found
		}
	}

	// If exact match not found, try partial matching for common variations
	for _, player := range players {
		playerLower := strings.ToLower(player.Name)

		// Check if the input matches any part of the player's name (for nicknames or partial names)
		if strings.Contains(playerLower, lowerName) && len(lowerName) >= 3 {
			// Only match if the input is at least 3 characters to avoid too many false positives
			return &player, true
		}

		// Check if player name contains the input (reverse check for partial matches)
		if strings.Contains(lowerName, playerLower) && len(playerLower) >= 3 {
			return &player, true
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
