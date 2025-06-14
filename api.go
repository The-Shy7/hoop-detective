package main

import (
	// Package for JSON encoding and decoding
	"fmt"      // Package for formatted I/O operations
	"io"       // Package for I/O primitives
	"net/http" // Package for HTTP client and server implementations
	// Package for string conversions
	"strings" // Package for string manipulation functions
	"time"    // Package for time-related operations

	"github.com/tidwall/gjson" // Third-party package for JSON parsing
)

// Constants for NBA API configuration
const (
	NBA_API_BASE = "https://stats.nba.com/stats"                                  // Base URL for NBA Stats API
	USER_AGENT   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" // User agent string to mimic browser requests
)

// APIPlayer represents the raw player data structure returned by the NBA API
type APIPlayer struct {
	ID          int    `json:"PERSON_ID"`          // Unique player identifier
	FirstName   string `json:"FIRST_NAME"`         // Player's first name
	LastName    string `json:"LAST_NAME"`          // Player's last name
	DisplayName string `json:"DISPLAY_FIRST_LAST"` // Full display name
	TeamID      int    `json:"TEAM_ID"`            // Current team identifier
	TeamName    string `json:"TEAM_NAME"`          // Current team name
	TeamCity    string `json:"TEAM_CITY"`          // Current team city
	Position    string `json:"POSITION"`           // Playing position
	Height      string `json:"HEIGHT"`             // Player height
	Weight      string `json:"WEIGHT"`             // Player weight
	College     string `json:"SCHOOL"`             // College attended
	Country     string `json:"COUNTRY"`            // Country of origin
	DraftYear   int    `json:"DRAFT_YEAR"`         // Year drafted
	DraftRound  int    `json:"DRAFT_ROUND"`        // Draft round
	DraftNumber int    `json:"DRAFT_NUMBER"`       // Draft pick number
	Experience  int    `json:"EXPERIENCE"`         // Years of NBA experience
	Jersey      string `json:"JERSEY"`             // Jersey number
	IsActive    bool   `json:"ROSTERSTATUS"`       // Whether player is currently active
}

// PlayerStats holds career statistical averages for a player
type PlayerStats struct {
	PPG float64 // Points per game average
	RPG float64 // Rebounds per game average
	APG float64 // Assists per game average
}

// PlayerCareerInfo holds additional career information not available from basic API
type PlayerCareerInfo struct {
	Accolades     []string // List of major achievements and awards
	TeamHistory   []string // List of all teams played for
	AllStarYears  []int    // Years selected as All-Star
	Championships []int    // Years won championships
}

// Global cache variables to store API responses and reduce repeated requests
var playerCache = make(map[int]*Player) // Cache individual player data by ID
var allPlayersCache []Player            // Cache all players list
var cacheExpiry time.Time               // Timestamp when cache expires

// makeAPIRequest performs HTTP GET request to NBA API with proper headers
func makeAPIRequest(url string) ([]byte, error) {
	// Create HTTP client with 30-second timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err // Return error if request creation fails
	}

	// Set headers to mimic browser request and avoid blocking
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // Return error if request fails
	}
	defer resp.Body.Close() // Ensure response body is closed when function exits

	// Check if response status is successful (200 OK)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Read and return response body
	return io.ReadAll(resp.Body)
}

// fetchAllPlayers retrieves comprehensive player data from NBA API
func fetchAllPlayers() ([]Player, error) {
	// Check if cached data is still valid (within 1 hour)
	if time.Now().Before(cacheExpiry) && len(allPlayersCache) > 0 {
		return allPlayersCache, nil // Return cached data if still valid
	}

	// Inform user that API fetch is starting
	fmt.Println("Fetching NBA players from API...")

	// Construct API URL for all players (current and historical)
	url := fmt.Sprintf("%s/commonallplayers?LeagueID=00&Season=2023-24&IsOnlyCurrentSeason=0", NBA_API_BASE)

	// Make API request
	data, err := makeAPIRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch players: %v", err)
	}

	// Parse JSON response using gjson for efficient extraction
	result := gjson.GetBytes(data, "resultSets.0") // Get first result set
	headers := result.Get("headers").Array()       // Extract column headers
	rows := result.Get("rowSet").Array()           // Extract data rows

	// Initialize slice to store processed players
	var players []Player

	// Create map for quick header index lookup
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header.String()] = i // Map header name to column index
	}

	// Process each player row from API response
	for _, row := range rows {
		rowData := row.Array() // Convert row to array of values
		if len(rowData) == 0 {
			continue // Skip empty rows
		}

		// Extract essential player information
		playerID := int(rowData[headerMap["PERSON_ID"]].Int()) // Get player ID
		firstName := rowData[headerMap["FIRST_NAME"]].String() // Get first name
		lastName := rowData[headerMap["LAST_NAME"]].String()   // Get last name

		// Skip players with missing essential data
		if firstName == "" || lastName == "" {
			continue
		}

		// Create Player struct with basic information
		player := Player{
			Name:      fmt.Sprintf("%s %s", firstName, lastName),      // Combine first and last name
			Team:      getTeamName(rowData, headerMap),                // Extract team name
			Position:  getStringValue(rowData, headerMap, "POSITION"), // Extract position
			Height:    getStringValue(rowData, headerMap, "HEIGHT"),   // Extract height
			College:   getStringValue(rowData, headerMap, "SCHOOL"),   // Extract college
			DraftYear: int(rowData[headerMap["FROM_YEAR"]].Int()),     // Extract draft year
		}

		// Fetch detailed career statistics for this player
		stats, err := fetchPlayerStats(playerID)
		if err == nil {
			// If stats fetch successful, add to player data
			player.PPG = stats.PPG
			player.RPG = stats.RPG
			player.APG = stats.APG
		}

		// Get additional career information (accolades, team history)
		careerInfo := getPlayerCareerInfo(player.Name)
		player.Accolades = careerInfo.Accolades
		player.TeamHistory = careerInfo.TeamHistory

		// Add completed player to results
		players = append(players, player)

		// Add small delay to avoid overwhelming the API with requests
		time.Sleep(50 * time.Millisecond)

		// Limit total players to prevent excessive API calls during development
		if len(players) >= 500 {
			break // Stop after 500 players
		}
	}

	// Cache the results for 1 hour to improve performance
	allPlayersCache = players
	cacheExpiry = time.Now().Add(1 * time.Hour)

	// Inform user of successful completion
	fmt.Printf("Successfully loaded %d NBA players!\n", len(players))
	return players, nil
}

// fetchPlayerStats retrieves career statistics for a specific player
func fetchPlayerStats(playerID int) (*PlayerStats, error) {
	// Construct API URL for player career stats
	url := fmt.Sprintf("%s/playercareerstats?PlayerID=%d", NBA_API_BASE, playerID)

	// Make API request for player stats
	data, err := makeAPIRequest(url)
	if err != nil {
		return nil, err // Return error if request fails
	}

	// Parse career totals from response (resultSets.1 contains career summary)
	result := gjson.GetBytes(data, "resultSets.1")
	if !result.Exists() {
		return nil, fmt.Errorf("no career stats found") // Error if no stats data
	}

	// Extract rows from career totals
	rows := result.Get("rowSet").Array()
	if len(rows) == 0 {
		return nil, fmt.Errorf("no career stats data") // Error if no data rows
	}

	// Get the career totals row (typically the last row in the dataset)
	careerRow := rows[len(rows)-1].Array()

	// Create PlayerStats struct with extracted values
	stats := &PlayerStats{
		PPG: careerRow[26].Float(), // Points per game (column 26)
		RPG: careerRow[23].Float(), // Rebounds per game (column 23)
		APG: careerRow[24].Float(), // Assists per game (column 24)
	}

	return stats, nil
}

// getTeamName safely extracts team name from player data row
func getTeamName(rowData []gjson.Result, headerMap map[string]int) string {
	// Check if TEAM_NAME column exists and is within bounds
	if idx, exists := headerMap["TEAM_NAME"]; exists && idx < len(rowData) {
		teamName := rowData[idx].String() // Extract team name
		if teamName != "" {
			return teamName // Return team name if not empty
		}
	}
	return "Free Agent" // Default value if no team or empty
}

// getStringValue safely extracts string value from player data row
func getStringValue(rowData []gjson.Result, headerMap map[string]int, key string) string {
	// Check if requested column exists and is within bounds
	if idx, exists := headerMap[key]; exists && idx < len(rowData) {
		value := rowData[idx].String() // Extract string value
		if value == "" {
			return "Unknown" // Return "Unknown" if empty
		}
		return value // Return actual value if not empty
	}
	return "Unknown" // Default value if column doesn't exist
}

// getPlayerCareerInfo provides additional career information for well-known players
// This is a simplified implementation - a full version would query additional APIs
func getPlayerCareerInfo(playerName string) PlayerCareerInfo {
	// Initialize empty career info structure
	info := PlayerCareerInfo{
		Accolades:   []string{}, // Empty accolades list
		TeamHistory: []string{}, // Empty team history list
	}

	// Add hardcoded accolades and team history for famous players
	// In a production system, this would query additional API endpoints
	switch playerName {
	case "LeBron James":
		info.Accolades = []string{"4x NBA Champion", "4x Finals MVP", "4x MVP", "19x All-Star"}
		info.TeamHistory = []string{"Cleveland Cavaliers", "Miami Heat", "Los Angeles Lakers"}
	case "Michael Jordan":
		info.Accolades = []string{"6x NBA Champion", "6x Finals MVP", "5x MVP", "14x All-Star"}
		info.TeamHistory = []string{"Chicago Bulls", "Washington Wizards"}
	case "Kobe Bryant":
		info.Accolades = []string{"5x NBA Champion", "2x Finals MVP", "1x MVP", "18x All-Star"}
		info.TeamHistory = []string{"Los Angeles Lakers"}
	case "Stephen Curry":
		info.Accolades = []string{"4x NBA Champion", "1x Finals MVP", "2x MVP", "9x All-Star"}
		info.TeamHistory = []string{"Golden State Warriors"}
	case "Kevin Durant":
		info.Accolades = []string{"2x NBA Champion", "2x Finals MVP", "1x MVP", "14x All-Star"}
		info.TeamHistory = []string{"Seattle SuperSonics", "Oklahoma City Thunder", "Golden State Warriors", "Brooklyn Nets", "Phoenix Suns"}
	default:
		// For other players, add generic accolades based on name patterns
		if strings.Contains(playerName, "All-Star") {
			info.Accolades = append(info.Accolades, "All-Star")
		}
	}

	return info // Return career information
}

// getFallbackPlayers provides a curated list of notable players when API is unavailable
func getFallbackPlayers() []Player {
	// Return hardcoded list of famous NBA players with complete data
	return []Player{
		{
			Name:        "LeBron James",                                                         // Current Lakers superstar
			Team:        "Los Angeles Lakers",                                                   // Current team
			Position:    "SF",                                                                   // Small Forward
			Height:      "6'9\"",                                                                // Height in feet/inches
			College:     "None",                                                                 // Straight from high school
			DraftYear:   2003,                                                                   // Draft year
			PPG:         27.2,                                                                   // Career points per game
			RPG:         7.5,                                                                    // Career rebounds per game
			APG:         7.3,                                                                    // Career assists per game
			Accolades:   []string{"4x NBA Champion", "4x Finals MVP", "4x MVP", "19x All-Star"}, // Major achievements
			TeamHistory: []string{"Cleveland Cavaliers", "Miami Heat", "Los Angeles Lakers"},    // Teams played for
		},
		{
			Name:        "Michael Jordan",                                                       // Basketball legend
			Team:        "Retired",                                                              // No longer active
			Position:    "SG",                                                                   // Shooting Guard
			Height:      "6'6\"",                                                                // Height in feet/inches
			College:     "North Carolina",                                                       // College attended
			DraftYear:   1984,                                                                   // Draft year
			PPG:         30.1,                                                                   // Career points per game
			RPG:         6.2,                                                                    // Career rebounds per game
			APG:         5.3,                                                                    // Career assists per game
			Accolades:   []string{"6x NBA Champion", "6x Finals MVP", "5x MVP", "14x All-Star"}, // Major achievements
			TeamHistory: []string{"Chicago Bulls", "Washington Wizards"},                        // Teams played for
		},
		{
			Name:        "Kobe Bryant",                                                          // Lakers legend
			Team:        "Retired",                                                              // No longer active
			Position:    "SG",                                                                   // Shooting Guard
			Height:      "6'6\"",                                                                // Height in feet/inches
			College:     "None",                                                                 // Straight from high school
			DraftYear:   1996,                                                                   // Draft year
			PPG:         25.0,                                                                   // Career points per game
			RPG:         5.2,                                                                    // Career rebounds per game
			APG:         4.7,                                                                    // Career assists per game
			Accolades:   []string{"5x NBA Champion", "2x Finals MVP", "1x MVP", "18x All-Star"}, // Major achievements
			TeamHistory: []string{"Los Angeles Lakers"},                                         // Teams played for
		},
		{
			Name:        "Stephen Curry",                                                       // Warriors superstar
			Team:        "Golden State Warriors",                                               // Current team
			Position:    "PG",                                                                  // Point Guard
			Height:      "6'2\"",                                                               // Height in feet/inches
			College:     "Davidson",                                                            // College attended
			DraftYear:   2009,                                                                  // Draft year
			PPG:         24.3,                                                                  // Career points per game
			RPG:         4.6,                                                                   // Career rebounds per game
			APG:         6.5,                                                                   // Career assists per game
			Accolades:   []string{"4x NBA Champion", "1x Finals MVP", "2x MVP", "9x All-Star"}, // Major achievements
			TeamHistory: []string{"Golden State Warriors"},                                     // Teams played for
		},
		{
			Name:        "Kevin Durant",                                                                                                     // Current Suns star
			Team:        "Phoenix Suns",                                                                                                     // Current team
			Position:    "SF",                                                                                                               // Small Forward
			Height:      "6'10\"",                                                                                                           // Height in feet/inches
			College:     "Texas",                                                                                                            // College attended
			DraftYear:   2007,                                                                                                               // Draft year
			PPG:         27.3,                                                                                                               // Career points per game
			RPG:         7.1,                                                                                                                // Career rebounds per game
			APG:         4.3,                                                                                                                // Career assists per game
			Accolades:   []string{"2x NBA Champion", "2x Finals MVP", "1x MVP", "14x All-Star"},                                             // Major achievements
			TeamHistory: []string{"Seattle SuperSonics", "Oklahoma City Thunder", "Golden State Warriors", "Brooklyn Nets", "Phoenix Suns"}, // Teams played for
		},
	}
}
