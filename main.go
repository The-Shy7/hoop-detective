package main

import (
	"bufio"     // Package for buffered I/O operations, used for reading user input
	"fmt"       // Package for formatted I/O operations like printing to console
	"math/rand" // Package for generating random numbers
	"os"        // Package for operating system interface, used for standard input
	"strings"   // Package for string manipulation functions
	"time"      // Package for time-related operations
)

// main is the entry point of the program
func main() {
	// Initialize players from API
	fmt.Println("🏀 HOOP DETECTIVE 🏀")
	fmt.Println("Loading NBA player database...")

	// Attempt to load player data from NBA API or fallback to hardcoded data
	err := initializePlayers()
	if err != nil {
		// Display warning if API loading failed but continue with fallback data
		fmt.Printf("Warning: Could not load full player database: %v\n", err)
		fmt.Println("Using fallback player data...")
	}

	// Initialize game variables
	target := getRandomPlayer()                 // Select a random player as the mystery player to guess
	attempts := 0                               // Counter for number of guesses made
	maxAttempts := 8                            // Maximum number of guesses allowed
	hintsUsed := 0                              // Counter for number of hints used
	maxHints := 3                               // Maximum number of hints allowed
	usedHintAttributes := make(map[string]bool) // Track which attributes have been revealed
	scanner := bufio.NewScanner(os.Stdin)       // Create scanner to read user input from terminal

	// Timer setup
	gameStartTime := time.Now()
	gameDuration := 6 * time.Minute // 6 minutes total
	gameEndTime := gameStartTime.Add(gameDuration)

	// Print game instructions and setup information
	printInstructions()
	fmt.Printf("\nYou have %d attempts and 6 minutes to guess the mystery NBA player!\n", maxAttempts)
	fmt.Printf("You can use up to %d hints by typing 'hint'.\n", maxHints)
	fmt.Println("Type 'quit' to exit the game.")
	fmt.Printf("⏰ Game started at: %s\n", gameStartTime.Format("15:04:05"))
	fmt.Printf("⏰ Time limit: %s\n", gameEndTime.Format("15:04:05"))

	// Print header row for the comparison results table
	printHeader()

	// Main game loop - continues until max attempts reached, time runs out, or player guesses correctly
	for attempts < maxAttempts {
		// Check if time has run out
		currentTime := time.Now()
		if currentTime.After(gameEndTime) {
			fmt.Printf("\n⏰ TIME'S UP! You ran out of time after %s.\n", formatDuration(currentTime.Sub(gameStartTime)))
			fmt.Printf("The mystery player was: %s\n", target.Name)
			printPlayerDetails(target)
			return
		}

		// Calculate and display remaining time
		timeRemaining := gameEndTime.Sub(currentTime)

		// Display current attempt number, time remaining, and prompt for user input
		fmt.Printf("\nAttempt %d/%d - Time remaining: %s - Enter your guess: ",
			attempts+1, maxAttempts, formatTimeRemaining(timeRemaining))

		// Read user input from terminal with timeout handling
		inputChan := make(chan string, 1)
		go func() {
			if scanner.Scan() {
				inputChan <- scanner.Text()
			} else {
				inputChan <- ""
			}
		}()

		// Wait for input or timeout
		select {
		case guess := <-inputChan:
			// Process the user's input
			guess = strings.TrimSpace(guess)

			// Check if user wants to quit the game
			if strings.ToLower(guess) == "quit" {
				fmt.Println("\nThanks for playing! The mystery player was:", target.Name)
				return // Exit the program
			}

			// Check if user wants to use a hint
			if strings.ToLower(guess) == "hint" {
				if hintsUsed >= maxHints {
					fmt.Printf("❌ You've already used all %d hints!\n", maxHints)
					continue // Don't count this as an attempt, go to next iteration
				}

				// Show a unique random attribute hint
				hintGiven := showUniqueRandomAttributeHint(target, hintsUsed+1, usedHintAttributes)
				if hintGiven {
					hintsUsed++
					fmt.Printf("💡 Hints remaining: %d\n", maxHints-hintsUsed)
				} else {
					fmt.Printf("❌ All available attributes have already been revealed!\n")
				}
				continue // Don't count this as an attempt, go to next iteration
			}

			// Search for the guessed player in the database
			guessedPlayer, found := findPlayerByName(guess)
			if !found {
				// Player not found in database - show error and continue without counting attempt
				fmt.Printf("❌ Player '%s' not found. Please check the spelling.\n", guess)
				fmt.Printf("💡 Tip: Type 'hint' to get a clue about the mystery player (%d hints remaining)\n", maxHints-hintsUsed)
				continue // Don't increment attempts counter
			}

			// Increment attempts counter since we have a valid guess
			attempts++

			// Compare the guessed player with the target player and display results
			result := compareWithTarget(*guessedPlayer, target)
			fmt.Println(result) // Print the color-coded comparison results

			// Check if the guess is correct (exact name match)
			if guessedPlayer.Name == target.Name {
				// Player guessed correctly - show victory message and exit
				elapsedTime := time.Now().Sub(gameStartTime)
				fmt.Printf("\n🎉 CONGRATULATIONS! 🎉\n")
				fmt.Printf("You guessed correctly in %d attempts and %s!\n", attempts, formatDuration(elapsedTime))
				if hintsUsed > 0 {
					fmt.Printf("You used %d hint(s) to help you.\n", hintsUsed)
				}
				fmt.Printf("The mystery player was: %s\n", target.Name)
				printPlayerDetails(target) // Show detailed information about the target player
				return                     // Exit the program
			}

			// Check if player has used all attempts
			if attempts == maxAttempts {
				// Game over - show failure message and reveal answer
				elapsedTime := time.Now().Sub(gameStartTime)
				fmt.Printf("\n💔 Game Over! You've used all %d attempts in %s.\n", maxAttempts, formatDuration(elapsedTime))
				fmt.Printf("The mystery player was: %s\n", target.Name)
				printPlayerDetails(target) // Show detailed information about the target player
				return                     // Exit the program
			}

			// Provide name hints at specific attempts
			if attempts == 4 {
				// After 4 attempts, reveal first letter of first name and last name
				nameHint := getNameHint(target.Name, 1) // Get first letter hint
				fmt.Printf("💡 Hint: The player's name starts with: %s\n", nameHint)
			} else if attempts == 6 {
				// After 6 attempts, reveal more of the player's name
				nameHint := getNameHint(target.Name, 2) // Get more detailed name hint
				fmt.Printf("💡 Hint: The player's name pattern: %s\n", nameHint)
			}

		case <-time.After(time.Until(gameEndTime)):
			// Time ran out while waiting for input
			fmt.Printf("\n⏰ TIME'S UP! You ran out of time.\n")
			fmt.Printf("The mystery player was: %s\n", target.Name)
			printPlayerDetails(target)
			return
		}
	}
}

// formatTimeRemaining formats the remaining time in a user-friendly way
func formatTimeRemaining(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// formatDuration formats elapsed time in a user-friendly way
func formatDuration(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	if minutes > 0 {
		return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
	}
	return fmt.Sprintf("%d seconds", seconds)
}

// showUniqueRandomAttributeHint displays a unique random attribute of the target player
// Returns true if a hint was given, false if all attributes have been used
func showUniqueRandomAttributeHint(target Player, hintNumber int, usedAttributes map[string]bool) bool {
	// List of all available attributes to hint about
	allAttributes := []string{"team", "position", "height", "college", "draftyear", "draftround", "draftnumber", "jerseynumber", "country"}

	// Filter out already used attributes
	var availableAttributes []string
	for _, attr := range allAttributes {
		if !usedAttributes[attr] {
			availableAttributes = append(availableAttributes, attr)
		}
	}

	// Check if any attributes are still available
	if len(availableAttributes) == 0 {
		return false // No more unique attributes available
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Select a random attribute from available ones
	selectedAttribute := availableAttributes[rand.Intn(len(availableAttributes))]

	// Mark this attribute as used
	usedAttributes[selectedAttribute] = true

	fmt.Printf("💡 Hint #%d: ", hintNumber)

	switch selectedAttribute {
	case "team":
		fmt.Printf("The player's current team is: %s\n", target.Team)
	case "position":
		fmt.Printf("The player's position is: %s\n", target.Position)
	case "height":
		fmt.Printf("The player's height is: %s\n", target.Height)
	case "college":
		if target.College == "None" || target.College == "Unknown" {
			fmt.Printf("The player did not attend college (international or straight from high school)\n")
		} else {
			fmt.Printf("The player attended: %s\n", target.College)
		}
	case "draftyear":
		fmt.Printf("The player was drafted in: %d\n", target.DraftYear)
	case "draftround":
		if target.DraftRound == 0 {
			fmt.Printf("The player was undrafted\n")
		} else {
			fmt.Printf("The player was drafted in round: %d\n", target.DraftRound)
		}
	case "draftnumber":
		if target.DraftNumber == 0 {
			fmt.Printf("The player was undrafted (no draft pick number)\n")
		} else {
			fmt.Printf("The player was the #%d overall pick\n", target.DraftNumber)
		}
	case "jerseynumber":
		if target.JerseyNumber == "Unknown" {
			fmt.Printf("The player's jersey number is not available\n")
		} else {
			fmt.Printf("The player's jersey number is: #%s\n", target.JerseyNumber)
		}
	case "country":
		fmt.Printf("The player is from: %s\n", target.Country)
	}

	return true // Hint was successfully given
}

// getNameHint returns a partial hint of the player's name based on the hint level
func getNameHint(fullName string, hintLevel int) string {
	// Split the name into parts (first name, last name, etc.)
	nameParts := strings.Fields(fullName)

	if len(nameParts) == 0 {
		return "Unknown"
	}

	switch hintLevel {
	case 1:
		// Level 1: Show first letter of each name part
		var hints []string
		for _, part := range nameParts {
			if len(part) > 0 {
				hints = append(hints, string(part[0])+"_")
			}
		}
		return strings.Join(hints, " ")

	case 2:
		// Level 2: Show first 2-3 letters of each name part and length
		var hints []string
		for _, part := range nameParts {
			if len(part) == 0 {
				continue
			}

			var hint string
			if len(part) <= 3 {
				// For very short names, show first letter only
				hint = string(part[0]) + strings.Repeat("_", len(part)-1)
			} else if len(part) <= 5 {
				// For short names, show first 2 letters
				hint = part[:2] + strings.Repeat("_", len(part)-2)
			} else {
				// For longer names, show first 3 letters
				hint = part[:3] + strings.Repeat("_", len(part)-3)
			}

			// Add length indicator
			hint += fmt.Sprintf(" (%d letters)", len(part))
			hints = append(hints, hint)
		}
		return strings.Join(hints, " ")

	default:
		// Default case: just show first letters
		var hints []string
		for _, part := range nameParts {
			if len(part) > 0 {
				hints = append(hints, string(part[0])+"_")
			}
		}
		return strings.Join(hints, " ")
	}
}

// printPlayerDetails displays comprehensive information about a player
func printPlayerDetails(player Player) {
	// Print decorative separator line
	fmt.Println("\n" + strings.Repeat("-", 50))

	// Display basic player information
	fmt.Printf("Name: %s\n", player.Name)
	fmt.Printf("Team: %s\n", player.Team)
	fmt.Printf("Position: %s\n", player.Position)
	fmt.Printf("Height: %s\n", player.Height)
	fmt.Printf("College: %s\n", player.College)
	fmt.Printf("Draft Year: %d\n", player.DraftYear)

	// Display draft information
	if player.DraftRound == 0 {
		fmt.Printf("Draft Status: Undrafted\n")
	} else {
		fmt.Printf("Draft Round: %d\n", player.DraftRound)
		fmt.Printf("Draft Pick: %d\n", player.DraftNumber)
	}

	// Display jersey number and country
	fmt.Printf("Jersey Number: %s\n", player.JerseyNumber)
	fmt.Printf("Country: %s\n", player.Country)

	// Print closing decorative separator line
	fmt.Println(strings.Repeat("-", 50))
}
