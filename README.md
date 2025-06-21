# Hoop Detective

A comprehensive CLI-based guessing game where you try to identify NBA players from throughout basketball history based on their attributes! Race against time with a 6-minute timer while using strategic hints to narrow down your guesses.

## Features

- **Complete NBA Database**: Access to hundreds of NBA players from throughout history via Ball Don't Lie API
- **Smart Comparison System**: Intelligent matching with color-coded feedback
- **Case-Insensitive Input**: Player names work regardless of capitalization (e.g., "lebron james" = "LeBron James")
- **Dual Challenge System**: 8 attempts AND 6-minute time limit
- **Strategic Hint System**: Get up to 3 unique random attribute hints plus automatic name hints
- **Real-time Timer**: Live countdown showing remaining time during gameplay
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

## Game Rules

### **Dual Challenge System**
- **8 Attempts Maximum**: You have up to 8 guesses to identify the mystery player
- **6-Minute Time Limit**: Complete the challenge before time runs out
- **Win Condition**: Guess correctly within both the attempt limit AND time limit
- **Lose Condition**: Run out of attempts OR time expires

### **User-Friendly Input**
- **Case-Insensitive**: Type player names in any case (e.g., "lebron james", "LEBRON JAMES", "LeBron James")
- **Flexible Matching**: Supports partial name matching for common variations
- **Automatic Trimming**: Extra spaces are automatically removed
- **Smart Recognition**: The system recognizes players even with minor typing variations

### **Timer Features**
- **Real-time Display**: See remaining time with each guess prompt
- **Start/End Times**: Game shows when it started and when it will end
- **Victory Timing**: See exactly how long it took you to win
- **Time Pressure**: Adds excitement and urgency to decision-making

## Project Structure

```
hoop-detective/
├── main.go          # Main entry point, game loop, and timer management
├── player.go        # Player data structures and case-insensitive matching
├── game.go          # Game logic and comparison algorithms
├── api.go           # Ball Don't Lie API integration and data fetching
├── .env             # Environment variables (API key)
├── .gitignore       # Git ignore file
├── go.mod           # Go module dependencies
└── README.md        # Project documentation
```

### File Descriptions

#### `main.go`
**Purpose**: Main entry point, game orchestration, and timer management
- **Game Loop**: Manages the 8-attempt guessing cycle with time constraints
- **Timer System**: Implements 6-minute countdown with real-time display
- **User Input**: Handles player guesses, commands ('hint', 'quit') with timeout detection
- **Case-Insensitive Processing**: Converts input to appropriate format for matching
- **Unique Hint System**: Manages both unique random attribute hints (up to 3) and automatic name hints
- **Victory/Defeat Logic**: Determines win conditions based on attempts AND time
- **Time Formatting**: Provides user-friendly time display (minutes/seconds)
- **Concurrent Input**: Handles user input while monitoring timer expiration
- **Hint Tracking**: Prevents duplicate attribute hints using a tracking map

#### `player.go`
**Purpose**: Player data structures and flexible player matching
- **Player Struct**: Defines the complete player data model with 10 attributes
- **Database Management**: Handles player loading from API or fallback data
- **Random Selection**: Implements mystery player selection algorithm
- **Case-Insensitive Search**: Provides flexible name matching functionality
  - **Exact Matching**: Case-insensitive exact name matching
  - **Partial Matching**: Supports common name variations and nicknames
  - **Smart Filtering**: Prevents false positives with minimum length requirements
- **Data Initialization**: Coordinates between API and fallback systems
- **Global State**: Manages the in-memory player database

#### `game.go`
**Purpose**: Game mechanics, comparison logic, and user interface
- **Comparison Engine**: Implements sophisticated attribute matching with tolerance ranges
- **Color Coding System**: 
  - 🟢 Green for exact matches
  - 🟡 Yellow for close matches (±2 years for draft year, ±5 picks for draft number)
  - 🔴 Red for no matches
- **Draft Information**: Handles special cases for undrafted players
- **Display Formatting**: Creates aligned tabular output for results
- **Game Instructions**: Provides user guidance including timer and input information
- **Timer Integration**: Updated instructions reflect the 6-minute time limit and case-insensitive input

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
2. You have **8 attempts and 6 minutes** to guess the correct player
3. **Type player names in any case** - "lebron james", "LEBRON JAMES", and "LeBron James" all work
4. Each guess shows remaining time and attempt count
5. After each guess, you'll receive feedback using color-coded indicators:
   - 🟢 **Green**: Exact match
   - 🟡 **Yellow**: Close match (within range for numbers)
   - 🔴 **Red**: No match

## Game Features

- **Comprehensive Player Database**: Fetches data from Ball Don't Lie NBA API for hundreds of players
- **Dual Challenge System**: Both attempt-based and time-based constraints
- **Real-time Timer**: Live countdown creates urgency and excitement
- **Case-Insensitive Input**: Type names however you want - the game understands
- **Smart Name Matching**: Handles common variations and partial names intelligently
- **Smart Comparison**: 
  - Draft years within 2 years show as yellow
  - Draft picks within 5 positions show as yellow
  - Special handling for undrafted players
- **Unique Hint System**: 
  - **Random Attribute Hints**: Up to 3 unique hints revealing different player attributes
  - **Automatic Name Hints**: Progressive name hints at attempts 4 and 6
- **Interactive Help**: Strategic hint system to help narrow down possibilities
- **Detailed Results**: See complete player information and timing after the game

## Player Attributes Compared

- **Name**: Full player name (case-insensitive matching)
- **Current Team**: Current team or "Free Agent" for players without teams
- **Position**: Primary playing position (PG, SG, SF, PF, C)
- **Height**: Player height in feet and inches
- **College**: College attended (from API when available)
- **Draft Year**: Year entered NBA (from API when available)
- **Draft Round**: Round drafted in (1-2, or "Undrafted")
- **Draft Number**: Overall pick number (1-60, or "N/A" for undrafted)
- **Jersey Number**: Current jersey number
- **Country**: Country of origin

## Input System

### **Case-Insensitive Matching**
The game accepts player names in any capitalization:
- **"lebron james"** ✅
- **"LEBRON JAMES"** ✅
- **"LeBron James"** ✅
- **"lEbRoN jAmEs"** ✅

### **Smart Name Recognition**
- **Exact Matching**: Perfect case-insensitive name matches
- **Partial Matching**: Recognizes common name variations (minimum 3 characters)
- **Automatic Trimming**: Extra spaces are removed automatically
- **Flexible Input**: Works with various typing styles and preferences

### **Input Tips**
- Don't worry about capitalization - type naturally
- Full names work best for exact matches
- The system is forgiving with common variations
- Use 'hint' command if you're unsure of exact spelling

## Hint System

The game features a comprehensive dual hint system designed to help you strategically narrow down possibilities while racing against time:

### **Unique Random Attribute Hints (Limited to 3)**
Type 'hint' during the game to reveal a **unique** attribute of the mystery player:
- **Team**: Current team or retirement status
- **Position**: Playing position (PG, SG, SF, PF, C)
- **Height**: Player height in feet and inches
- **College**: College attended or international status
- **Draft Year**: Year the player entered the NBA
- **Draft Round**: Round drafted (1-2) or undrafted status
- **Draft Pick**: Overall pick number or undrafted
- **Jersey Number**: Current jersey number
- **Country**: Country of origin

**Key Feature**: Each hint command reveals a **different** attribute - no duplicates! This ensures maximum strategic value from your limited 3 hints.

### **Automatic Name Hints**
Progressive name hints are automatically provided at specific attempts:

#### **Attempt 4**: First Letter Hints
- Shows the first letter of each name part
- Example: "L_ J_" for LeBron James
- Helps eliminate players with different starting letters

#### **Attempt 6**: Detailed Name Pattern
- Shows first 2-3 letters of each name part
- Includes character count for each name
- Example: "Leb____ (6 letters) Jam__ (5 letters)" for LeBron James
- Provides enough information to make educated guesses

### **Strategic Hint Usage with Timer**
- **Early Game**: Use random attribute hints to eliminate large groups of players quickly
- **No Duplicates**: Each hint command guarantees new information
- **Time Management**: Balance hint usage with time remaining
- **Mid Game**: Combine attribute knowledge with first letter hints
- **Late Game**: Use detailed name patterns for final guesses under time pressure
- **Resource Management**: Choose wisely - you only get 3 unique random attribute hints!

### **Hint Tracking System**
- **Unique Guarantee**: The game tracks which attributes have been revealed
- **Smart Selection**: Only unused attributes are available for random selection
- **Maximum Value**: Every hint provides new, strategic information
- **Fallback Protection**: If all attributes are somehow revealed, the system prevents errors

## Timer System

### **6-Minute Challenge**
- **Total Time**: Exactly 6 minutes from game start to completion
- **Real-time Display**: See remaining time with each guess prompt
- **Format**: Displays as "Xm Ys" (e.g., "4m 32s") or just seconds when under 1 minute
- **Start/End Times**: Shows when game started and when it will end
- **Concurrent Monitoring**: Timer runs while waiting for user input

### **Time-based Victory/Defeat**
- **Victory**: Guess correctly within 8 attempts AND 6 minutes
- **Time Victory**: Shows total time taken when you win
- **Time Defeat**: Game ends immediately when timer expires
- **Attempt Defeat**: Game ends when 8 attempts are used (if time remains)

### **Strategic Time Management**
- **Quick Decisions**: Time pressure encourages faster decision-making
- **Hint Timing**: Consider when to use limited hints based on remaining time
- **Name Recognition**: Faster players may rely more on name pattern recognition
- **Pressure Element**: Adds excitement and urgency to the guessing process

## Architecture Overview

### Data Flow
1. **Initialization**: `main.go` calls `initializePlayers()` in `player.go`
2. **Environment Loading**: `api.go` loads `.env` file and reads API key
3. **API Fetching**: `player.go` calls `fetchAllPlayers()` in `api.go`
4. **Authentication**: `api.go` uses API key from .env file for authentication
5. **Data Processing**: `api.go` processes Ball Don't Lie NBA API responses and caches results
6. **Game Setup**: `main.go` selects random target player and starts timer
7. **Timer Management**: `main.go` tracks game start time and 6-minute deadline
8. **User Interaction**: `main.go` handles case-insensitive input with timeout detection
9. **Name Matching**: `player.go` performs flexible case-insensitive player search
10. **Unique Hint Management**: `main.go` tracks used attributes and provides unique reveals
11. **Comparison**: `game.go` performs attribute matching with tolerance logic
12. **Display**: Results formatted and displayed with color coding and time remaining
13. **Victory/Defeat**: Determined by both attempt count AND timer expiration

### Key Design Patterns
- **Separation of Concerns**: Each file has a distinct responsibility
- **Configuration Management**: .env file for easy API key management
- **Caching Strategy**: API responses cached for 1 hour to improve performance
- **Graceful Degradation**: Fallback to curated data if API unavailable or unauthenticated
- **Strategic Hint System**: Unique hint system maximizes strategic value
- **Timer Integration**: Concurrent timer monitoring with user input handling
- **Case-Insensitive Design**: Flexible input handling throughout the system
- **Smart Matching**: Multiple matching strategies for better user experience
- **Data Integrity**: All player information comes directly from the API
- **Resource Management**: Limited unique hints and time encourage strategic thinking
- **State Tracking**: Map-based tracking prevents duplicate hint attributes

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

- **Player Name**: Guess a player by typing their full name (case-insensitive)
- **'hint'**: Get a unique random attribute hint about the mystery player (limited to 3 uses, no duplicates)
- **'quit'**: Exit the game

## Example Gameplay

```
🏀 HOOP DETECTIVE 🏀
Loading NBA player database...
Fetching NBA players from Ball Don't Lie API...
Note: Using API key from .env file for full player database access.
Successfully loaded 500 NBA players from API!

Database contains 500 NBA players from throughout history!
You have 8 attempts and 6 minutes to guess the mystery NBA player!
You can use up to 3 hints by typing 'hint'.
⏰ Game started at: 14:30:15
⏰ Time limit: 14:36:15
💡 Tip: Player names are case-insensitive (e.g., 'lebron james' works)

Attempt 1/8 - Time remaining: 5m 59s - Enter your guess: hint
💡 Hint #1: The player's position is: PG
💡 Hints remaining: 2

Attempt 1/8 - Time remaining: 5m 55s - Enter your guess: hint
💡 Hint #2: The player's current team is: Los Angeles Lakers
💡 Hints remaining: 1

Attempt 1/8 - Time remaining: 5m 50s - Enter your guess: shake milton
🔴 Shake Milton | 🟢 Los Angeles Lakers | 🔴 SG | 🔴 6'5" | 🔴 SMU | 🔴 2018 | 🔴 2 | 🔴 54 | 🔴 20 | 🟢 USA

💡 Hint: The player's name starts with: G_ V_

Attempt 2/8 - Time remaining: 5m 35s - Enter your guess: gabe vincent
🟢 Gabe Vincent | 🟢 Los Angeles Lakers | 🟢 PG | 🟢 6'3" | 🟢 UC Santa Barbara | 🟢 2020 | 🟢 Undrafted | 🟢 N/A | 🟢 7 | 🟢 USA

🎉 CONGRATULATIONS! 🎉
You guessed correctly in 2 attempts and 25 seconds!
You used 2 hint(s) to help you.
The mystery player was: Gabe Vincent
```

## Technical Details

- **Language**: Go 1.21+
- **Dependencies**: None (uses only Go standard library)
- **API**: Ball Don't Lie NBA API (api.balldontlie.io)
- **Configuration**: .env file for API key management
- **Authentication**: API key via Authorization header
- **Caching**: 1-hour cache for API responses to improve performance
- **Rate Limiting**: Built-in delays to respect API limits
- **Timer Implementation**: Concurrent goroutines for input handling and timer monitoring
- **Case-Insensitive Matching**: String manipulation with Go's strings package
- **Hint Tracking**: Map-based system to prevent duplicate attribute hints
- **Data Processing**: JSON parsing with Go's encoding/json package

## Code Organization

The codebase follows Go best practices with clear separation of concerns:

- **Modular Design**: Each file handles a specific aspect of the game
- **Configuration Management**: Simple .env file for API key storage
- **Error Handling**: Comprehensive error handling with graceful fallbacks
- **Performance**: Caching and rate limiting for optimal API usage
- **Timer Integration**: Concurrent programming for real-time timer functionality
- **Case-Insensitive Design**: Flexible input handling throughout the system
- **Unique Hint System**: Map-based tracking ensures no duplicate attribute hints
- **Smart Matching**: Multiple matching strategies for better user experience
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

2. **Player Name Issues**:
   - Don't worry about capitalization - "lebron james" works the same as "LeBron James"
   - Try the full name if partial matching doesn't work
   - Use the hint system if you're unsure of exact spelling
   - The system supports common name variations

3. **Timer Issues**:
   - The timer runs concurrently with user input
   - If you're taking too long to type, the game may end due to time expiration
   - Practice typing player names quickly to maximize your time

4. **Hint System**:
   - Each 'hint' command reveals a unique attribute - no repeats
   - If you somehow exhaust all 9 possible attributes, the system will notify you
   - Use hints strategically early in the game for maximum benefit

5. **Input Recognition**:
   - The game is case-insensitive but still requires reasonably accurate spelling
   - Partial matching works for names with at least 3 characters
   - If a player isn't found, double-check the spelling or use a hint

6. **Slow Loading**: Initial load may take time as player data is fetched from the API

7. **Network Issues**: The game includes retry logic and graceful degradation

8. **Configuration Issues**: 
   - Make sure the `.env` file is in the same directory as the game
   - Check that the file format is correct: `BALLDONTLIE_API_KEY=your_key`

## Future Enhancements

Potential improvements for the game:
- **Fuzzy Matching**: Even more flexible name recognition with typo tolerance
- **Nickname Support**: Recognition of common player nicknames (e.g., "King James" for LeBron)
- **Autocomplete**: Suggest player names as you type
- **Difficulty Levels**: Different time limits (3min/6min/9min) and attempt counts
- **Hint Categories**: Allow players to choose specific categories (physical, career, team)
- **Real Statistics**: Integrate with additional APIs for actual career stats (PPG, RPG, APG)
- **Team History**: Fetch complete team history for all players (requires different data source)
- **Accolades**: Add major achievements and awards (requires different data source)
- **Multiplayer**: Add competitive multiplayer functionality with shared timers
- **Player Photos**: Add visual elements to the game
- **Configuration UI**: Add interactive setup for API key and timer configuration
- **Advanced Hint System**: More sophisticated hint algorithms and weighted categories
- **Smart Hint Ordering**: Prioritize more valuable attributes based on elimination potential
- **Leaderboard**: Track best times and attempt counts
- **Pause Functionality**: Allow players to pause the timer (with limitations)
- **Sound Effects**: Add audio cues for timer warnings and victory

Enjoy testing your NBA knowledge against the clock with this strategic guessing challenge featuring guaranteed unique hints and flexible input handling!