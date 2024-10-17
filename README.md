# Fruit or Doom

Fruit or Doom is a real-time multiplayer web game where two players compete in an intense battle of strategy and reflexes.
Player one collects apples while avoiding bombs, while player two plants bombs to thwart their opponent , apples to gain more score.
The game leverages WebSockets for seamless real-time communication between players, providing a fast-paced and interactive experience.

## Rules of the Game

### Game Modes

Players can choose between two modes of play:

    Mode 1: Player One (the collector) attempts to guess the contents of the cells on an empty grid.
    Mode 2: Player Two (the planter) strategically places apples and bombs on the grid.

### Scoring System

    Both players start with 100 points.

### Mode 1 (Collector)

    Player One interacts with an empty grid and must guess whether each cell contains an apple or a bomb.
    Scoring:
        +5 points for clicking an apple.
        -10 points for clicking a bomb.
    Player One can spend points for additional information:
        Pay 5 points to reveal which cells are currently active.
        Pay 5 points to see the number of apples and bombs on the grid.

### Mode 2 (Planter)

    Player Two can plant apples or bombs on the grid at the cost of 5 points per action.
    Scoring:
        +10 points if an apple remains on the grid for 15 seconds without being clicked by Player One.
        -5 points if a bomb remains on the grid for 15 seconds without being clicked by Player One.
    Every 10 seconds, either an apple or a bomb is automatically planted in a random cell on the grid.

### The game ends when one of the following conditions is met:

    1- One player's score drops to 0.
    2- One player leaves the room.
    3- The game timer reaches 3 minutes.

## Demo

[![Demo](https://img.youtube.com/vi/fyLtcAapIkc/0.jpg)](https://www.youtube.com/watch?v=fyLtcAapIkc)
https://www.youtube.com/watch?v=fyLtcAapIkc

## Run Locally

Clone the project

```bash
git https://github.com/abdulrhmanm03/real_time_board_game.git
```

Go to the project directory

```bash
cd real_time_board_game
```

Go to the backend directory

```bash
cd backend
```

Install dependencies

```bash
go mod tidy
```

Start the server

```bash
go run main.go
```

Go to the client directory

```bash
cd ../client
```

Install dependencies

```bash
npm install
```

Start the server

```bash
npm run dev
```
