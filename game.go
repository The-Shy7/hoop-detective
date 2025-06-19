package main

import (
	"fmt"     // Package for formatted I/O operations
	"strings" // Package for string manipulation functions
)

// ComparisonResult holds the formatted comparison results between guess and target player
type ComparisonResult struct {
	Name        string // Formatted name with color indicator
	Team        string // Formatted team with color indicator
	Position    string // Formatted position with color indicator
	Height      string // Formatted height with color indicator
	College     string // Formatted college with color indicator
	DraftYear   string // Formatted draft year with color indicator
	PPG         string // Formatted points per game with color indicator
	RPG         string // Formatted rebounds per game with color indicator
	APG         string // Formatted assists per game with color indicator
	Accolades   string // Formatted accolades with color indicator
	TeamHistory string // Formatted team history with color indicator
}

// String method formats ComparisonResult for display in tabular format
func (cr ComparisonResult) String() string {
	// Return formatted string with fixed-width columns for aligned display
	return fmt.Sprintf("%-20s | %-20s | %-8s | %-6s | %-15s | %-9s | %-6s | %-6s | %-6s | %-30s | %-30s",
		cr.Name, cr.Team, cr.Position, cr.Height, cr.College, cr.DraftYear, cr.PPG, cr.RPG, cr.APG, cr.Accolades, cr.TeamHistory)
}

// compareWithTarget compares a guessed player with the target player and returns color-coded results
func compareWithTarget(guess, target Player) ComparisonResult {
	// Initialize empty result structure
	result := ComparisonResult{}

	// Compare Name - exact match required for green
	if guess.Name == target.Name {
		result.Name = "🟢 " + guess.Name // Green circle for exact match
	} else {
		result.Name = "🔴 " + guess.Name // Red circle for no match
	}

	// Compare Team - exact match required for green
	if guess.Team == target.Team {
		result.Team = "🟢 " + guess.Team // Green circle for exact match
	} else {
		result.Team = "🔴 " + guess.Team // Red circle for no match
	}

	// Compare Position - exact match required for green
	if guess.Position == target.Position {
		result.Position = "🟢 " + guess.Position // Green circle for exact match
	} else {
		result.Position = "🔴 " + guess.Position // Red circle for no match
	}

	// Compare Height - exact match required for green
	if guess.Height == target.Height {
		result.Height = "🟢 " + guess.Height // Green circle for exact match
	} else {
		result.Height = "🔴 " + guess.Height // Red circle for no match
	}

	// Compare College - exact match required for green
	if guess.College == target.College {
		result.College = "🟢 " + guess.College // Green circle for exact match
	} else {
		result.College = "🔴 " + guess.College // Red circle for no match
	}

	// Compare Draft Year with tolerance for close matches
	if guess.DraftYear == target.DraftYear {
		// Exact match gets green
		result.DraftYear = "🟢 " + fmt.Sprintf("%d", guess.DraftYear)
	} else if abs(guess.DraftYear-target.DraftYear) <= 2 {
		// Within 2 years gets yellow (close match)
		result.DraftYear = "🟡 " + fmt.Sprintf("%d", guess.DraftYear)
	} else {
		// More than 2 years difference gets red
		result.DraftYear = "🔴 " + fmt.Sprintf("%d", guess.DraftYear)
	}

	// Compare Points Per Game with tolerance ranges
	if abs64(guess.PPG-target.PPG) < 1.0 {
		// Within 1.0 point gets green (very close)
		result.PPG = "🟢 " + fmt.Sprintf("%.1f", guess.PPG)
	} else if abs64(guess.PPG-target.PPG) < 3.0 {
		// Within 3.0 points gets yellow (close)
		result.PPG = "🟡 " + fmt.Sprintf("%.1f", guess.PPG)
	} else {
		// More than 3.0 points difference gets red
		result.PPG = "🔴 " + fmt.Sprintf("%.1f", guess.PPG)
	}

	// Compare Rebounds Per Game with tolerance ranges
	if abs64(guess.RPG-target.RPG) < 1.0 {
		// Within 1.0 rebound gets green (very close)
		result.RPG = "🟢 " + fmt.Sprintf("%.1f", guess.RPG)
	} else if abs64(guess.RPG-target.RPG) < 2.0 {
		// Within 2.0 rebounds gets yellow (close)
		result.RPG = "🟡 " + fmt.Sprintf("%.1f", guess.RPG)
	} else {
		// More than 2.0 rebounds difference gets red
		result.RPG = "🔴 " + fmt.Sprintf("%.1f", guess.RPG)
	}

	// Compare Assists Per Game with tolerance ranges
	if abs64(guess.APG-target.APG) < 1.0 {
		// Within 1.0 assist gets green (very close)
		result.APG = "🟢 " + fmt.Sprintf("%.1f", guess.APG)
	} else if abs64(guess.APG-target.APG) < 2.0 {
		// Within 2.0 assists gets yellow (close)
		result.APG = "🟡 " + fmt.Sprintf("%.1f", guess.APG)
	} else {
		// More than 2.0 assists difference gets red
		result.APG = "🔴 " + fmt.Sprintf("%.1f", guess.APG)
	}

	// Compare Accolades - check for common elements between lists
	commonAccolades := findCommonElements(guess.Accolades, target.Accolades)
	if len(commonAccolades) > 0 {
		// Check if accolades match exactly (same length and all common)
		if len(commonAccolades) == len(target.Accolades) && len(guess.Accolades) == len(target.Accolades) {
			result.Accolades = "🟢 " + strings.Join(guess.Accolades, ", ") // Perfect match
		} else {
			result.Accolades = "🟡 " + strings.Join(guess.Accolades, ", ") // Partial match
		}
	} else {
		result.Accolades = "🔴 " + strings.Join(guess.Accolades, ", ") // No common accolades
	}

	// Compare Team History - check for common teams between lists
	commonTeams := findCommonElements(guess.TeamHistory, target.TeamHistory)
	if len(commonTeams) > 0 {
		// Check if team histories match exactly (same length and all common)
		if len(commonTeams) == len(target.TeamHistory) && len(guess.TeamHistory) == len(target.TeamHistory) {
			result.TeamHistory = "🟢 " + strings.Join(guess.TeamHistory, ", ") // Perfect match
		} else {
			result.TeamHistory = "🟡 " + strings.Join(guess.TeamHistory, ", ") // Partial match
		}
	} else {
		result.TeamHistory = "🔴 " + strings.Join(guess.TeamHistory, ", ") // No common teams
	}

	// Return the complete comparison result
	return result
}

// findCommonElements finds elements that appear in both string slices
func findCommonElements(slice1, slice2 []string) []string {
	// Initialize empty slice to store common elements
	common := []string{}

	// Iterate through first slice
	for _, item1 := range slice1 {
		// For each item in first slice, check if it exists in second slice
		for _, item2 := range slice2 {
			if item1 == item2 {
				// If match found, add to common elements and break inner loop
				common = append(common, item1)
				break // Avoid duplicates by breaking after first match
			}
		}
	}

	// Return slice of common elements
	return common
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x // Return positive version of negative number
	}
	return x // Return unchanged if already positive
}

// abs64 returns the absolute value of a float64
func abs64(x float64) float64 {
	if x < 0 {
		return -x // Return positive version of negative number
	}
	return x // Return unchanged if already positive
}

// printHeader displays the column headers for the comparison results table
func printHeader() {
	// Print separator line of equal signs
	fmt.Println(strings.Repeat("=", 200))

	// Print column headers with fixed widths for alignment
	fmt.Printf("%-20s | %-20s | %-8s | %-6s | %-15s | %-9s | %-6s | %-6s | %-6s | %-30s | %-30s\n",
		"NAME", "TEAM", "POSITION", "HEIGHT", "COLLEGE", "DRAFT YR", "PPG", "RPG", "APG", "ACCOLADES", "TEAM HISTORY")

	// Print another separator line
	fmt.Println(strings.Repeat("=", 200))
}

// printInstructions displays the game rules and setup information
func printInstructions() {
	// Print game title with basketball emoji
	fmt.Println("\n🏀 NBA PLAYER GUESSING GAME 🏀")

	// Print game rules and instructions
	fmt.Println("\nHow to play:")
	fmt.Println("- Guess NBA players by typing their full name")
	fmt.Println("- 🟢 Green = Exact match")
	fmt.Println("- 🟡 Yellow = Close match (within range for numbers, partial match for lists)")
	fmt.Println("- 🔴 Red = No match")

	// Display information about the player database size
	fmt.Printf("\nDatabase contains %d NBA players from throughout history!\n", len(players))
	fmt.Println("Type 'hint' during the game to see some available players.")

	// Print decorative separator line
	fmt.Println(strings.Repeat("=", 80))
}
