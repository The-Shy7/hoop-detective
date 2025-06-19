package main

import (
	"bufio" // Package for buffered I/O operations, used for reading user input
	"fmt"   // Package for formatted I/O operations like printing to console
	"math/rand"
	"os"      // Package for operating system interface, used for standard input
	"strings" // Package for string manipulation functions
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
	target := getRandomPlayer()           // Select a random player as the mystery player to guess
	attempts := 0                         // Counter for number of guesses made
	maxAttempts := 8                      // Maximum number of guesses allowed
	scanner := bufio.NewScanner(os.Stdin) // Create scanner to read user input from terminal

	// Print game instructions and setup information
	printInstructions()
	fmt.Printf("\nYou have %d attempts to guess the mystery NBA player!\n", maxAttempts)
	fmt.Println("Type 'quit' to exit the game.")
	fmt.Println("Type 'hint' to get a hint about available players.")

	// Print header row for the comparison results table
	printHeader()

	// Main game loop - continues until max attempts reached or player guesses correctly
	for attempts < maxAttempts {
		// Display current attempt number and prompt for user input
		fmt.Printf("\nAttempt %d/%d - Enter your guess: ", attempts+1, maxAttempts)

		// Read user input from terminal
		if !scanner.Scan() {
			break // Exit if there's an error reading input
		}

		// Get the user's guess and remove leading/trailing whitespace
		guess := strings.TrimSpace(scanner.Text())

		// Check if user wants to quit the game
		if strings.ToLower(guess) == "quit" {
			fmt.Println("\nThanks for playing! The mystery player was:", target.Name)
			return // Exit the program
		}

		// Check if user wants to see a hint of available players
		if strings.ToLower(guess) == "hint" {
			showPlayerHint() // Display sample of available players
			continue         // Don't count this as an attempt, go to next iteration
		}

		// Search for the guessed player in the database
		guessedPlayer, found := findPlayerByName(guess)
		if !found {
			// Player not found in database - show error and continue without counting attempt
			fmt.Printf("❌ Player '%s' not found. Please check the spelling.\n", guess)
			fmt.Println("💡 Tip: Type 'hint' to see some available players")
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
			fmt.Printf("\n🎉 CONGRATULATIONS! 🎉\n")
			fmt.Printf("You guessed correctly in %d attempts!\n", attempts)
			fmt.Printf("The mystery player was: %s\n", target.Name)
			printPlayerDetails(target) // Show detailed information about the target player
			return                     // Exit the program
		}

		// Check if player has used all attempts
		if attempts == maxAttempts {
			// Game over - show failure message and reveal answer
			fmt.Printf("\n💔 Game Over! You've used all %d attempts.\n", maxAttempts)
			fmt.Printf("The mystery player was: %s\n", target.Name)
			printPlayerDetails(target) // Show detailed information about the target player
			return                     // Exit the program
		}

		// Provide progressive hints based on number of attempts made
		if attempts == 3 {
			// After 3 attempts, reveal the player's position
			fmt.Printf("\n💡 Hint: The player's position is %s\n", target.Position)
		} else if attempts == 5 {
			// After 5 attempts, reveal the player's draft year
			fmt.Printf("💡 Hint: The player was drafted in %d\n", target.DraftYear)
		} else if attempts == 7 {
			// After 7 attempts, reveal the player's current team
			fmt.Printf("💡 Hint: The player's current team is %s\n", target.Team)
		}
	}
}

// showPlayerHint displays a random sample of available players to help the user
func showPlayerHint() {
	fmt.Println("\n🔍 Here are some available players to help you get started:")

	// Determine how many players to show (15 or total if less than 15)
	sampleSize := 15
	if len(players) < sampleSize {
		sampleSize = len(players)
	}

	// Create slice of indices for all players
	indices := make([]int, len(players))
	for i := range indices {
		indices[i] = i // Fill with sequential numbers 0, 1, 2, ...
	}

	// Shuffle the indices using Fisher-Yates algorithm to get random sample
	for i := len(indices) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)                           // Generate random index from 0 to i
		indices[i], indices[j] = indices[j], indices[i] // Swap elements
	}

	// Display the random sample of players
	for i := 0; i < sampleSize; i++ {
		fmt.Printf("- %s\n", players[indices[i]].Name)
	}

	// Show how many more players are available
	fmt.Printf("\n... and %d more players in the database!\n", len(players)-sampleSize)
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
