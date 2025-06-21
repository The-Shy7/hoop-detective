package main

import (
	"fmt"     // Package for formatted I/O operations
	"strings" // Package for string manipulation functions
)

// ComparisonResult holds the formatted comparison results between guess and target player
type ComparisonResult struct {
	Name         string // Formatted name with color indicator
	Team         string // Formatted team with color indicator
	Position     string // Formatted position with color indicator
	Height       string // Formatted height with color indicator
	College      string // Formatted college with color indicator
	DraftYear    string // Formatted draft year with color indicator
	DraftRound   string // Formatted draft round with color indicator
	DraftNumber  string // Formatted draft number with color indicator
	JerseyNumber string // Formatted jersey number with color indicator
	Country      string // Formatted country with color indicator
}

// String method formats ComparisonResult for display in tabular format
func (cr ComparisonResult) String() string {
	// Return formatted string with fixed-width columns for aligned display
	return fmt.Sprintf("%-20s | %-20s | %-8s | %-6s | %-15s | %-9s | %-5s | %-6s | %-6s | %-12s",
		cr.Name, cr.Team, cr.Position, cr.Height, cr.College, cr.DraftYear, cr.DraftRound, cr.DraftNumber, cr.JerseyNumber, cr.Country)
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

	// Compare Draft Round - exact match required for green
	if guess.DraftRound == target.DraftRound {
		if guess.DraftRound == 0 {
			result.DraftRound = "🟢 Undrafted" // Special case for undrafted players
		} else {
			result.DraftRound = "🟢 " + fmt.Sprintf("%d", guess.DraftRound)
		}
	} else {
		if guess.DraftRound == 0 {
			result.DraftRound = "🔴 Undrafted" // Special case for undrafted players
		} else {
			result.DraftRound = "🔴 " + fmt.Sprintf("%d", guess.DraftRound)
		}
	}

	// Compare Draft Number with tolerance for close matches
	if guess.DraftNumber == target.DraftNumber {
		if guess.DraftNumber == 0 {
			result.DraftNumber = "🟢 N/A" // Special case for undrafted players
		} else {
			result.DraftNumber = "🟢 " + fmt.Sprintf("%d", guess.DraftNumber)
		}
	} else if guess.DraftNumber != 0 && target.DraftNumber != 0 && abs(guess.DraftNumber-target.DraftNumber) <= 5 {
		// Within 5 picks gets yellow (close match) - only for drafted players
		result.DraftNumber = "🟡 " + fmt.Sprintf("%d", guess.DraftNumber)
	} else {
		if guess.DraftNumber == 0 {
			result.DraftNumber = "🔴 N/A" // Special case for undrafted players
		} else {
			result.DraftNumber = "🔴 " + fmt.Sprintf("%d", guess.DraftNumber)
		}
	}

	// Compare Jersey Number - exact match required for green
	if guess.JerseyNumber == target.JerseyNumber {
		result.JerseyNumber = "🟢 " + guess.JerseyNumber // Green circle for exact match
	} else {
		result.JerseyNumber = "🔴 " + guess.JerseyNumber // Red circle for no match
	}

	// Compare Country - exact match required for green
	if guess.Country == target.Country {
		result.Country = "🟢 " + guess.Country // Green circle for exact match
	} else {
		result.Country = "🔴 " + guess.Country // Red circle for no match
	}

	// Return the complete comparison result
	return result
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x // Return positive version of negative number
	}
	return x // Return unchanged if already positive
}

// printHeader displays the column headers for the comparison results table
func printHeader() {
	// Print separator line of equal signs
	fmt.Println(strings.Repeat("=", 150))

	// Print column headers with fixed widths for alignment
	fmt.Printf("%-20s | %-20s | %-8s | %-6s | %-15s | %-9s | %-5s | %-6s | %-6s | %-12s\n",
		"NAME", "TEAM", "POSITION", "HEIGHT", "COLLEGE", "DRAFT YR", "ROUND", "PICK", "JERSEY", "COUNTRY")

	// Print another separator line
	fmt.Println(strings.Repeat("=", 150))
}

// printInstructions displays the game rules and setup information
func printInstructions() {
	// Print game rules and instructions
	fmt.Println("\nHow to play:")
	fmt.Println("- Guess NBA players by typing their full name")
	fmt.Println("- You have 8 attempts and 6 minutes to guess correctly")
	fmt.Println("- 🟢 Green = Exact match")
	fmt.Println("- 🟡 Yellow = Close match (within range for numbers)")
	fmt.Println("- 🔴 Red = No match")

	// Display information about the player database size
	fmt.Printf("\nDatabase contains %d NBA players from throughout history!\n", len(players))
	fmt.Println("Type 'hint' during the game to get clues about the mystery player (limited to 3 hints).")
	fmt.Println("⏰ Race against time - you only have 6 minutes!")

	// Print decorative separator line
	fmt.Println(strings.Repeat("=", 80))
}
