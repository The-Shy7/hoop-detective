# Hoop Detective

A comprehensive CLI-based guessing game where you try to identify NBA players from throughout basketball history based on their attributes!

## Features

- **Complete NBA Database**: Access to hundreds of NBA players from throughout history via Ball Don't Lie API
- **Smart Comparison System**: Intelligent matching with color-coded feedback
- **Progressive Hints**: Get helpful hints as you play
- **Fallback System**: Works even if API is unavailable

## Quick Setup

1. **Get your API key** (optional but recommended):
   - Create a free account at [https://app.balldontlie.io](https://app.balldontlie.io)
   - Copy your API key from the dashboard

2. **Configure the game**:
   - Create `.env` file in the project directory
   - Add your API key to the `.env` file in the following format:
     ```
     BALLDONTLIE_API_KEY=your_actual_api_key_here
     ```

3. **Run the game**:
   ```bash
   go run .
   ```

**Without an API key**, the game will automatically fall back to a curated list of 10 legendary players.

## Project Structure

```
hoop-detective/
├── main.go          # Main entry point and game loop
├── player.go        # Player data structures and management
├── game.go          # Game logic and comparison algorithms
├── api.go           # Ball Don't Lie API integration and data fetching
├── .env             # Environment variables (API key)
├── .gitignore       # Git ignore file
├── go.mod           # Go module dependencies
└── README.md        # Project documentation
```

### File Descriptions

#### `main.go`
**Purpose**: Main entry point and game orchestration
- **Game Loop**: Manages the 8-attempt guessing cycle
- **User Input**: Handles player guesses, commands ('hint', 'quit')
- **Progressive Hints**: Reveals position (attempt 3), draft year (attempt 5), team (attempt 7)
- **Victory/Defeat Logic**: Determines win conditions and displays results
- **Player Hint System**: Shows random sample of available players
- **Game Flow Control**: Coordinates between all other modules

#### `player.go`
**Purpose**: Player data structures and core player management
- **Player Struct**: Defines the complete player data model with 10 attributes
- **Database Management**: Handles player loading from API or fallback data
- **Random Selection**: Implements mystery player selection algorithm
- **Player Search**: Provides exact name matching functionality
- **Data Initialization**: Coordinates between API and fallback systems
- **Global State**: Manages the in-memory player database

#### `game.go`
**Purpose**: Game mechanics and comparison logic
- **Comparison Engine**: Implements sophisticated attribute matching with tolerance ranges
- **Color Coding System**: 
  - 🟢 Green for exact matches
  - 🟡 Yellow for close matches (±2 years for draft year, ±5 picks for draft number)
  - 🔴 Red for no matches
- **Draft Information**: Handles special cases for undrafted players
- **Display Formatting**: Creates aligned tabular output for results
- **Game Instructions**: Provides user guidance and rule explanations

#### `api.go`
**Purpose**: Ball Don't Lie API integration and external data management
- **Ball Don't Lie NBA API**: Fetches comprehensive player data from api.balldontlie.io
- **Environment File Loading**: Reads API key from .env file
- **Authentication**: Handles API key authentication via Authorization header
- **HTTP Client**: Handles API requests with proper headers and timeouts
- **JSON Parsing**: Uses standard library for efficient data extraction
- **Caching System**: 1-hour cache to reduce API calls and improve performance
- **Rate Limiting**: Built-in delays to respect API usage limits
- **Cursor-based Pagination**: Handles the new pagination system
- **Fallback Data**: Provides curated list of legendary players when API fails
- **Error Handling**: Graceful degradation when external services are unavailable or require authentication
- **Data Processing**: Extracts all available player information from API responses

#### `.env`
**Purpose**: Environment configuration
- **API Key Storage**: Securely stores your Ball Don't Lie API key
- **Easy Configuration**: Simple key=value format for easy editing
- **Template Provided**: Includes placeholder for your actual API key

#### `go.mod`
**Purpose**: Dependency management
- **Module Definition**: Defines the project as `hoop-detective` module
- **Go Version**: Requires Go 1.21 or higher
- **No External Dependencies**: Uses only Go standard library for maximum compatibility

## How to Play

1. The game selects a random NBA player as the mystery player
2. You have 8 attempts to guess the correct player
3. After each guess, you'll receive feedback using color-coded indicators:
   - 🟢 **Green**: Exact match
   - 🟡 **Yellow**: Close match (within range for numbers)
   - 🔴 **Red**: No match

## Game Features

- **Comprehensive Player Database**: Fetches data from Ball Don't Lie NBA API for hundreds of players
- **Smart Comparison**: 
  - Draft years within 2 years show as yellow
  - Draft picks within 5 positions show as yellow
  - Special handling for undrafted players
- **Progressive Hints**: Get hints after attempts 3, 5, and 7
- **Interactive Help**: Type 'hint' to see available players
- **Detailed Results**: See complete player information after the game

## Player Attributes Compared

- **Name**: Full player name
- **Current Team**: Current team or "Free Agent" for players without teams
- **Position**: Primary playing position (PG, SG, SF, PF, C)
- **Height**: Player height in feet and inches
- **College**: College attended (from API when available)
- **Draft Year**: Year entered NBA (from API when available)
- **Draft Round**: Round drafted in (1-2, or "Undrafted")
- **Draft Number**: Overall pick number (1-60, or "N/A" for undrafted)
- **Jersey Number**: Current jersey number
- **Country**: Country of origin

## Architecture Overview

### Data Flow
1. **Initialization**: `main.go` calls `initializePlayers()` in `player.go`
2. **Environment Loading**: `api.go` loads `.env` file and reads API key
3. **API Fetching**: `player.go` calls `fetchAllPlayers()` in `api.go`
4. **Authentication**: `api.go` uses API key from .env file for authentication
5. **Data Processing**: `api.go` processes Ball Don't Lie NBA API responses and caches results
6. **Game Setup**: `main.go` selects random target player and starts game loop
7. **User Interaction**: `main.go` handles input and calls comparison functions
8. **Comparison**: `game.go` performs attribute matching with tolerance logic
9. **Display**: Results formatted and displayed with color coding

### Key Design Patterns
- **Separation of Concerns**: Each file has a distinct responsibility
- **Configuration Management**: .env file for easy API key management
- **Caching Strategy**: API responses cached for 1 hour to improve performance
- **Graceful Degradation**: Fallback to curated data if API unavailable or unauthenticated
- **Progressive Enhancement**: Hints become more specific as attempts increase
- **Data Integrity**: All player information comes directly from the API

## Installation & Setup

1. Make sure you have Go installed on your system
2. Clone or download the game files
3. Navigate to the game directory
4. (Optional) Set up your API key:
   - Get a free API key from [https://app.balldontlie.io](https://app.balldontlie.io)
   - Edit the `.env` file and replace `your_api_key_here` with your actual key
5. Run the game:

```bash
go run .
```

Or build and run:

```bash
go build -o hoop-detective
./hoop-detective
```

## API Integration

The game uses the **Ball Don't Lie NBA API** (api.balldontlie.io) to fetch comprehensive player data including:
- Player biographical information (name, position, height, weight)
- Current team affiliations
- College information
- Draft year, round, and pick number
- Jersey numbers
- Country of origin

**API Features:**
- **Authentication Required** - Free API key required from app.balldontlie.io
- **Comprehensive data** - Covers current and historical NBA players
- **Cursor-based pagination** - Efficiently handles large datasets
- **Rate limiting friendly** - Built-in delays respect API limits

**Configuration:**
- **Easy Setup**: API key stored in `.env` file
- **Automatic Loading**: Game automatically reads configuration on startup
- **Secure**: .env file can be excluded from version control

If the API is unavailable or unauthenticated, the game falls back to a curated list of 10 legendary players with complete data.

## Commands During Game

- **Player Name**: Guess a player by typing their full name
- **'hint'**: See a random sample of available players
- **'quit'**: Exit the game

## Example Gameplay

```
🏀 HOOP DETECTIVE 🏀
Loading NBA player database...
Fetching NBA players from Ball Don't Lie API...
Note: Using API key from .env file for full player database access.
Successfully loaded 500 NBA players from API!

Database contains 500 NBA players from throughout history!

Attempt 1/8 - Enter your guess: Michael Jordan
🔴 Michael Jordan    | 🔴 Retired           | 🔴 SG     | 🔴 6'6" | 🔴 North Carolina   | 🟡 1984    | 🔴 1     | 🔴 3      | 🔴 23    | 🟢 USA

💡 Hint: The player's position is PG

Attempt 2/8 - Enter your guess: hint

🔍 Here are some available players to help you get started:
- LeBron James
- Stephen Curry
- Kevin Durant
- Giannis Antetokounmpo
- Luka Doncic
- Jayson Tatum
- Joel Embiid
- Nikola Jokic
- Damian Lillard
- Russell Westbrook
- Chris Paul
- James Harden
- Kawhi Leonard
- Anthony Davis
- Jimmy Butler

... and 485 more players in the database!

Attempt 2/8 - Enter your guess: Magic Johnson
🟢 Magic Johnson     | 🟢 Retired           | 🟢 PG     | 🟢 6'9" | 🟢 Michigan State   | 🟢 1979    | 🟢 1     | 🟢 1      | 🟢 32    | 🟢 USA

🎉 CONGRATULATIONS! 🎉
You guessed correctly in 2 attempts!
```

## Technical Details

- **Language**: Go 1.21+
- **Dependencies**: None (uses only Go standard library)
- **API**: Ball Don't Lie NBA API (api.balldontlie.io)
- **Configuration**: .env file for API key management
- **Authentication**: API key via Authorization header
- **Caching**: 1-hour cache for API responses to improve performance
- **Rate Limiting**: Built-in delays to respect API limits
- **Data Processing**: JSON parsing with Go's encoding/json package

## Code Organization

The codebase follows Go best practices with clear separation of concerns:

- **Modular Design**: Each file handles a specific aspect of the game
- **Configuration Management**: Simple .env file for API key storage
- **Error Handling**: Comprehensive error handling with graceful fallbacks
- **Performance**: Caching and rate limiting for optimal API usage
- **Maintainability**: Extensive commenting and clear function naming
- **Extensibility**: Easy to add new player attributes or comparison logic
- **Zero Dependencies**: Uses only Go standard library for maximum portability

## Troubleshooting

If you encounter issues:

1. **API Authentication Problems**: 
   - Make sure you have a valid API key from [https://app.balldontlie.io](https://app.balldontlie.io)
   - Check that your `.env` file contains the correct API key
   - Ensure the API key is not the placeholder `your_api_key_here`
   - The game will automatically fall back to a curated player list if authentication fails

2. **Player Not Found**: Use the 'hint' command to see available players

3. **Slow Loading**: Initial load may take time as player data is fetched from the API

4. **Network Issues**: The game includes retry logic and graceful degradation

5. **Configuration Issues**: 
   - Make sure the `.env` file is in the same directory as the game
   - Check that the file format is correct: `BALLDONTLIE_API_KEY=your_key`

## Future Enhancements

Potential improvements for the game:
- **Real Statistics**: Integrate with additional APIs for actual career stats (PPG, RPG, APG)
- **Team History**: Fetch complete team history for all players (requires different data source)
- **Accolades**: Add major achievements and awards (requires different data source)
- **Difficulty Levels**: Add easy/medium/hard modes
- **Multiplayer**: Add competitive multiplayer functionality
- **Player Photos**: Add visual elements to the game
- **Configuration UI**: Add interactive setup for API key configuration

Enjoy testing your NBA knowledge across all eras of basketball history!