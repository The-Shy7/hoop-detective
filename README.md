# Hoop Detective [IN PROGRESS]

A comprehensive CLI-based guessing game where you try to identify NBA players from throughout basketball history based on their attributes!

## Features

- **Complete NBA Database**: Access to hundreds of NBA players from throughout history via NBA Stats API
- **Smart Comparison System**: Intelligent matching with color-coded feedback
- **Progressive Hints**: Get helpful hints as you play

## Project Structure

```
hoop-detective/
├── main.go          # Main entry point and game loop
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