package main

import (
	"bufio"         // Package for reading files line by line
	"encoding/json" // Package for JSON encoding and decoding
	"fmt"           // Package for formatted I/O operations
	"io"            // Package for I/O primitives
	"net/http"      // Package for HTTP client and server implementations
	"os"            // Package for file operations
	"strings"       // Package for string manipulation functions
	"time"          // Package for time-related operations
)

// Constants for NBA API configuration
const (
	NBA_API_BASE = "https://api.balldontlie.io/v1"                                // Base URL for Ball Don't Lie API
	USER_AGENT   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" // User agent string to mimic browser requests
)

// APIPlayer represents the raw player data structure returned by the NBA API
type APIPlayer struct {
	ID           int    `json:"id"`            // Unique player identifier
	FirstName    string `json:"first_name"`    // Player's first name
	LastName     string `json:"last_name"`     // Player's last name
	Position     string `json:"position"`      // Playing position
	Height       string `json:"height"`        // Player height (e.g., "6-2")
	Weight       string `json:"weight"`        // Player weight in pounds
	JerseyNumber string `json:"jersey_number"` // Jersey number
	College      string `json:"college"`       // College attended
	Country      string `json:"country"`       // Country of origin
	DraftYear    *int   `json:"draft_year"`    // Draft year (pointer to handle null)
	DraftRound   *int   `json:"draft_round"`   // Draft round (pointer to handle null)
	DraftNumber  *int   `json:"draft_number"`  // Draft number (pointer to handle null)
	Team         struct {
		ID           int    `json:"id"`           // Team identifier
		Conference   string `json:"conference"`   // Team conference
		Division     string `json:"division"`     // Team division
		City         string `json:"city"`         // Team city
		Name         string `json:"name"`         // Team name
		FullName     string `json:"full_name"`    // Full team name
		Abbreviation string `json:"abbreviation"` // Team abbreviation
	} `json:"team"` // Current team information
}

// APIResponse represents the structure of API responses
type APIResponse struct {
	Data []APIPlayer `json:"data"` // Array of player data
	Meta struct {
		NextCursor *int `json:"next_cursor"` // Next cursor for pagination (pointer to handle null)
		PerPage    int  `json:"per_page"`    // Items per page
	} `json:"meta"` // Pagination metadata
}

// Global cache variables to store API responses and reduce repeated requests
var playerCache = make(map[int]*Player) // Cache individual player data by ID
var allPlayersCache []Player            // Cache all players list
var cacheExpiry time.Time               // Timestamp when cache expires

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	file, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist, which is okay
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Set environment variable
		os.Setenv(key, value)
	}

	return scanner.Err()
}

// getAPIKey retrieves the API key from .env file or environment variable
func getAPIKey() string {
	// First, try to load from .env file
	loadEnvFile()

	// Get API key from environment variable
	apiKey := os.Getenv("BALLDONTLIE_API_KEY")

	// Check if it's the placeholder value
	if apiKey == "your_api_key_here" || apiKey == "" {
		return ""
	}

	return apiKey
}

// makeAPIRequest performs HTTP GET request to NBA API with proper headers and authentication
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
	req.Header.Set("Connection", "keep-alive")

	// Check for API key and set Authorization header
	apiKey := getAPIKey()
	if apiKey != "" {
		// Use the correct Authorization header format for Ball Don't Lie API
		req.Header.Set("Authorization", apiKey)
		fmt.Printf("DEBUG: Using API key for authentication (key: %s...)\n", apiKey[:8])
	} else {
		fmt.Printf("DEBUG: No API key found - API may return limited data or require authentication\n")
		fmt.Printf("DEBUG: To get an API key, visit: https://app.balldontlie.io\n")
		fmt.Printf("DEBUG: Then add it to your .env file: BALLDONTLIE_API_KEY=your_actual_key\n")
	}

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // Return error if request fails
	}
	defer resp.Body.Close() // Ensure response body is closed when function exits

	// Read response body first to check content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if response status is successful (200 OK)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Check if response is HTML (indicating authentication required)
	if strings.Contains(string(body), "<!DOCTYPE html>") || strings.Contains(string(body), "<html>") {
		return nil, fmt.Errorf("API returned HTML documentation page - authentication required or invalid API key. Please check your API key in the .env file")
	}

	return body, nil
}

// fetchAllPlayers retrieves comprehensive player data from NBA API
func fetchAllPlayers() ([]Player, error) {
	// Check if cached data is still valid (within 1 hour)
	if time.Now().Before(cacheExpiry) && len(allPlayersCache) > 0 {
		return allPlayersCache, nil // Return cached data if still valid
	}

	// Inform user that API fetch is starting
	fmt.Println("Fetching NBA players from Ball Don't Lie API...")

	// Check if API key is available
	apiKey := getAPIKey()
	if apiKey == "" {
		fmt.Println("Note: No API key found in .env file. You can:")
		fmt.Println("1. Get a free API key from https://app.balldontlie.io")
		fmt.Println("2. Add it to your .env file: BALLDONTLIE_API_KEY=your_actual_key")
		fmt.Println("3. The game will use a fallback list of legendary players")
		return nil, fmt.Errorf("no API key provided")
	} else {
		fmt.Println("Note: Using API key from .env file for full player database access.")
	}

	// Initialize slice to store all players
	var allPlayers []Player

	// Start with cursor 0 and continue until we reach the end or hit our limit
	cursor := 0
	maxPages := 10 // Fetch more pages since we have authentication
	pageCount := 0

	for pageCount < maxPages {
		// Construct API URL for current cursor with 100 players per page (max allowed)
		var url string
		if cursor == 0 {
			url = fmt.Sprintf("%s/players?per_page=100", NBA_API_BASE)
		} else {
			url = fmt.Sprintf("%s/players?cursor=%d&per_page=100", NBA_API_BASE, cursor)
		}

		// Make API request for current page
		data, err := makeAPIRequest(url)
		if err != nil {
			// If we have some players already, return them instead of failing completely
			if len(allPlayers) > 0 {
				break
			}
			return nil, fmt.Errorf("failed to fetch players: %v", err)
		}

		// Parse JSON response
		var response APIResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse API response: %v", err)
		}

		// Process each player from current page
		playersProcessed := 0
		for _, apiPlayer := range response.Data {
			// Skip players with missing essential data
			if apiPlayer.FirstName == "" || apiPlayer.LastName == "" {
				continue
			}

			// Create full player name
			playerName := fmt.Sprintf("%s %s", apiPlayer.FirstName, apiPlayer.LastName)
			currentTeam := getTeamName(apiPlayer)

			// Create Player struct with available information from API
			player := Player{
				Name:         playerName,                              // Combine first and last name
				Team:         currentTeam,                             // Extract team name
				Position:     getPosition(apiPlayer.Position),         // Extract and validate position
				Height:       formatHeightFromAPI(apiPlayer.Height),   // Format height from API
				College:      getCollege(apiPlayer.College),           // Get college info
				DraftYear:    getDraftYear(apiPlayer.DraftYear),       // Get draft year
				DraftRound:   getDraftRound(apiPlayer.DraftRound),     // Get draft round
				DraftNumber:  getDraftNumber(apiPlayer.DraftNumber),   // Get draft number
				JerseyNumber: getJerseyNumber(apiPlayer.JerseyNumber), // Get jersey number
				Country:      getCountry(apiPlayer.Country),           // Get country
			}

			// Add completed player to results
			allPlayers = append(allPlayers, player)
			playersProcessed++
		}

		// Show progress to user
		fmt.Printf("Loaded %d players so far... (processed %d from cursor %d)\n",
			len(allPlayers), playersProcessed, cursor)

		// Check if we've reached the last page
		if response.Meta.NextCursor == nil || len(response.Data) < 100 {
			fmt.Printf("Reached end of data at cursor %d (NextCursor: %v, DataCount: %d)\n",
				cursor, response.Meta.NextCursor, len(response.Data))
			break
		}

		// Add delay between requests to be respectful to the API
		time.Sleep(1 * time.Second) // Increased delay for rate limiting

		// Move to next cursor
		cursor = *response.Meta.NextCursor
		pageCount++
	}

	// If we didn't get any players from API, return error
	if len(allPlayers) == 0 {
		return nil, fmt.Errorf("no players retrieved from API - authentication may be required")
	}

	// Cache the results for 1 hour to improve performance
	allPlayersCache = allPlayers
	cacheExpiry = time.Now().Add(1 * time.Hour)

	// Inform user of successful completion
	fmt.Printf("Successfully loaded %d NBA players from API!\n", len(allPlayers))
	return allPlayers, nil
}

// getTeamName safely extracts team name from API player data
func getTeamName(apiPlayer APIPlayer) string {
	// Check if team information is available
	if apiPlayer.Team.FullName != "" {
		return apiPlayer.Team.FullName // Return full team name if available
	} else if apiPlayer.Team.Name != "" {
		return apiPlayer.Team.Name // Return team name if available
	}
	return "Free Agent" // Default value if no team information
}

// getPosition validates and returns a clean position string
func getPosition(position string) string {
	// Clean up position string and return standard abbreviations
	pos := strings.TrimSpace(strings.ToUpper(position))

	// Map common variations to standard positions
	switch pos {
	case "POINT GUARD", "PG":
		return "PG"
	case "SHOOTING GUARD", "SG":
		return "SG"
	case "SMALL FORWARD", "SF":
		return "SF"
	case "POWER FORWARD", "PF":
		return "PF"
	case "CENTER", "C":
		return "C"
	case "GUARD", "G":
		return "G" // Generic guard
	case "FORWARD", "F":
		return "F" // Generic forward
	default:
		if pos != "" {
			return pos // Return as-is if not empty
		}
		return "Unknown" // Default for empty positions
	}
}

// formatHeightFromAPI converts height from API format (e.g., "6-2") to readable format
func formatHeightFromAPI(height string) string {
	if height == "" {
		return "Unknown"
	}

	// API returns height in format like "6-2" (feet-inches)
	parts := strings.Split(height, "-")
	if len(parts) == 2 {
		return fmt.Sprintf("%s'%s\"", parts[0], parts[1])
	}

	return height // Return as-is if format is unexpected
}

// getCollege returns college information or default
func getCollege(college string) string {
	if college == "" {
		return "Unknown"
	}
	return college
}

// getDraftYear returns draft year or estimated year
func getDraftYear(draftYear *int) int {
	if draftYear != nil {
		return *draftYear
	}

	// Return a reasonable default for unknown draft years
	return 2020 // Default to recent year
}

// getDraftRound returns draft round or 0 for undrafted
func getDraftRound(draftRound *int) int {
	if draftRound != nil {
		return *draftRound
	}
	return 0 // 0 indicates undrafted
}

// getDraftNumber returns draft number or 0 for undrafted
func getDraftNumber(draftNumber *int) int {
	if draftNumber != nil {
		return *draftNumber
	}
	return 0 // 0 indicates undrafted
}

// getJerseyNumber returns jersey number or "Unknown" if not available
func getJerseyNumber(jerseyNumber string) string {
	if jerseyNumber == "" {
		return "Unknown"
	}
	return jerseyNumber
}

// getCountry returns country or "Unknown" if not available
func getCountry(country string) string {
	if country == "" {
		return "Unknown"
	}
	return country
}

// getFallbackPlayers provides a curated list of notable players when API is unavailable
func getFallbackPlayers() []Player {
	// Return expanded list of famous NBA players with complete data
	return []Player{
		{
			Name:         "LeBron James",       // Current Lakers superstar
			Team:         "Los Angeles Lakers", // Current team
			Position:     "SF",                 // Small Forward
			Height:       "6'9\"",              // Height in feet/inches
			College:      "None",               // Straight from high school
			DraftYear:    2003,                 // Draft year
			DraftRound:   1,                    // First round
			DraftNumber:  1,                    // First overall pick
			JerseyNumber: "6",                  // Current jersey number
			Country:      "USA",                // Country of origin
		},
		{
			Name:         "Michael Jordan", // Basketball legend
			Team:         "Retired",        // No longer active
			Position:     "SG",             // Shooting Guard
			Height:       "6'6\"",          // Height in feet/inches
			College:      "North Carolina", // College attended
			DraftYear:    1984,             // Draft year
			DraftRound:   1,                // First round
			DraftNumber:  3,                // Third overall pick
			JerseyNumber: "23",             // Famous jersey number
			Country:      "USA",            // Country of origin
		},
		{
			Name:         "Kobe Bryant", // Lakers legend
			Team:         "Retired",     // No longer active
			Position:     "SG",          // Shooting Guard
			Height:       "6'6\"",       // Height in feet/inches
			College:      "None",        // Straight from high school
			DraftYear:    1996,          // Draft year
			DraftRound:   1,             // First round
			DraftNumber:  13,            // 13th overall pick
			JerseyNumber: "24",          // Later jersey number
			Country:      "USA",         // Country of origin
		},
		{
			Name:         "Stephen Curry",         // Warriors superstar
			Team:         "Golden State Warriors", // Current team
			Position:     "PG",                    // Point Guard
			Height:       "6'2\"",                 // Height in feet/inches
			College:      "Davidson",              // College attended
			DraftYear:    2009,                    // Draft year
			DraftRound:   1,                       // First round
			DraftNumber:  7,                       // 7th overall pick
			JerseyNumber: "30",                    // Jersey number
			Country:      "USA",                   // Country of origin
		},
		{
			Name:         "Kevin Durant", // Current Suns star
			Team:         "Phoenix Suns", // Current team
			Position:     "SF",           // Small Forward
			Height:       "6'10\"",       // Height in feet/inches
			College:      "Texas",        // College attended
			DraftYear:    2007,           // Draft year
			DraftRound:   1,              // First round
			DraftNumber:  2,              // Second overall pick
			JerseyNumber: "35",           // Jersey number
			Country:      "USA",          // Country of origin
		},
		{
			Name:         "Giannis Antetokounmpo", // Bucks superstar
			Team:         "Milwaukee Bucks",       // Current team
			Position:     "PF",                    // Power Forward
			Height:       "6'11\"",                // Height in feet/inches
			College:      "None",                  // International player
			DraftYear:    2013,                    // Draft year
			DraftRound:   1,                       // First round
			DraftNumber:  15,                      // 15th overall pick
			JerseyNumber: "34",                    // Jersey number
			Country:      "Greece",                // Country of origin
		},
		{
			Name:         "Luka Doncic",      // Mavericks star
			Team:         "Dallas Mavericks", // Current team
			Position:     "PG",               // Point Guard
			Height:       "6'7\"",            // Height in feet/inches
			College:      "None",             // International player
			DraftYear:    2018,               // Draft year
			DraftRound:   1,                  // First round
			DraftNumber:  3,                  // Third overall pick
			JerseyNumber: "77",               // Jersey number
			Country:      "Slovenia",         // Country of origin
		},
		{
			Name:         "Joel Embiid",        // 76ers center
			Team:         "Philadelphia 76ers", // Current team
			Position:     "C",                  // Center
			Height:       "7'0\"",              // Height in feet/inches
			College:      "Kansas",             // College attended
			DraftYear:    2014,                 // Draft year
			DraftRound:   1,                    // First round
			DraftNumber:  3,                    // Third overall pick
			JerseyNumber: "21",                 // Jersey number
			Country:      "Cameroon",           // Country of origin
		},
		{
			Name:         "Nikola Jokic",   // Nuggets center
			Team:         "Denver Nuggets", // Current team
			Position:     "C",              // Center
			Height:       "6'11\"",         // Height in feet/inches
			College:      "None",           // International player
			DraftYear:    2014,             // Draft year
			DraftRound:   2,                // Second round
			DraftNumber:  41,               // 41st overall pick
			JerseyNumber: "15",             // Jersey number
			Country:      "Serbia",         // Country of origin
		},
		{
			Name:         "Jayson Tatum",   // Celtics forward
			Team:         "Boston Celtics", // Current team
			Position:     "SF",             // Small Forward
			Height:       "6'8\"",          // Height in feet/inches
			College:      "Duke",           // College attended
			DraftYear:    2017,             // Draft year
			DraftRound:   1,                // First round
			DraftNumber:  3,                // Third overall pick
			JerseyNumber: "0",              // Jersey number
			Country:      "USA",            // Country of origin
		},
	}
}
