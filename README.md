# Uno Game Round Visualizer | Isa-CardGame Simulation

## Description:
This repository contains a Go application that simulates rounds of the Uno card game and generates a visual representation of each round in HTML format. The application utilizes JSON files to store data about each round, including player hands, actions taken, current card on the stack, win probabilities, and more.

## Usage:
1. Clone the repository to your local machine.
2. Build and run the Go application using the provided source code.
3. The application will simulate Uno game rounds and generate a JSON file for each game.
4. Open the `rounds.html` file in your web browser to visualize the Uno game rounds.
5. Optionally, you can load a JSON file by clicking the "Load JSON File" button or drop a JSON file onto the drop area to visualize specific game rounds.
6. Use the "Clear JSON Data" button to remove the loaded JSON data from the visualization.

## JSON Files:
Each JSON file represents a single round of the Uno game and contains the following information:
- `RoundNumber`: The number of the round.
- `StartingPlayer`: The index of the starting player.
- `StartingCard`: The starting card of the round.
- `PlayerHands`: An array of arrays containing the cards in each player's hand.
- `Actions`: An array of objects representing the actions taken by each player in the round.
- `CurrentPlayerIndex`: The index of the current player.
- `CurrentCard`: The current card on the stack.
- `GameStack`: An array containing all the cards played in the round.
- `WinProbabilities`: An array of strings representing the win probabilities for each player.

## CPU Profile and HTML Page:
- The application includes functionality to generate a CPU profile (`cpu_profile.prof`) for performance analysis.
- The HTML page (`rounds.html`) provides an interactive visualization of Uno game rounds using the JSON data.
