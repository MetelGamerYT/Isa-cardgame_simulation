package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"net/http"

	_ "net/http/pprof"
)

type GameCard struct {
	Color string
	Value string
}

type RoundAction struct {
	PlayerIndex int
	Card        GameCard
}

type RoundData struct {
	RoundNumber        int
	StartingPlayer     int
	StartingCard       GameCard
	PlayerHands        [][]GameCard
	Actions            []RoundAction
	CurrentPlayerIndex int
	CurrentCard        GameCard
	GameStack          []GameCard
	WinProbabilities   []string
}

func createCardDeck() []GameCard {
	deck := make([]GameCard, 0)

	colors := []string{"Blue", "Green", "Red", "Yellow"}
	values := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "Draw Two"}

	for _, color := range colors {
		for _, value := range values {
			deck = append(deck, GameCard{color, value})
			deck = append(deck, GameCard{color, value})
		}
	}

	for i := 0; i < 4; i++ {
		deck = append(deck, GameCard{"", "Color Choice"})
		deck = append(deck, GameCard{"", "Draw Four"})
	}

	return deck
}

func shuffleDeck(deck []GameCard) []GameCard {
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func drawCards(deck []GameCard, numCards int) ([]GameCard, []GameCard) {
	hand := make([]GameCard, 0)
	remainingDeck := make([]GameCard, len(deck))
	copy(remainingDeck, deck)

	for i := 0; i < numCards; i++ {
		if len(remainingDeck) == 0 {
			// Deck is empty, generate a new deck
			remainingDeck = createCardDeck()
			remainingDeck = shuffleDeck(remainingDeck)
		}

		cardIndex := rand.Intn(len(remainingDeck))
		hand = append(hand, remainingDeck[cardIndex])
		remainingDeck = append(remainingDeck[:cardIndex], remainingDeck[cardIndex+1:]...)
	}
	return hand, remainingDeck
}

func drawInitialCard(deck []GameCard) (GameCard, []GameCard) {
	for {
		// Draw a random card from the remaining deck
		index := rand.Intn(len(deck))
		card := deck[index]

		// Check that the card is not one of the forbidden cards
		if card.Value != "Draw Four" && card.Value != "Color Choice" &&
			card.Value != "Reverse" && card.Value != "Skip" && card.Value != "Draw Two" {
			// Remove card from the deck
			deck = append(deck[:index], deck[index+1:]...)
			return card, deck
		}
	}
}

func canPlayCard(currentCard, playedCard GameCard) bool {
	return currentCard.Color == playedCard.Color || currentCard.Value == playedCard.Value || playedCard.Value == "Color Choice" || playedCard.Value == "Draw Four"
}

func nextPlayerIndex(currentPlayerIndex, numPlayers int) int {
	return (currentPlayerIndex + 1) % numPlayers
}

func calculateWinProbability(players [][]GameCard, _ int) []string {
	numPlayers := len(players)
	winProbabilities := make([]string, numPlayers)

	totalCards := 0
	for _, hand := range players {
		totalCards += len(hand)
	}

	for i, hand := range players {
		probab := float64(len(hand)) / float64(totalCards) * 100.0
		winProbabilities[i] = fmt.Sprintf("Player %d: %.2f%%", i+1, probab)
	}

	return winProbabilities
}

func main() {
	go func() {
		fmt.Println("pprof server started at :6060")
		http.ListenAndServe(":6060", nil)
	}()

	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile:", err)
		return
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println("Error starting CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().UnixNano())
	gameDeck := createCardDeck()
	gameDeck = shuffleDeck(gameDeck)

	startCard, gameDeck := drawInitialCard(gameDeck)
	fmt.Println("The starting card is:", startCard.Color, "-", startCard.Value)
	gameStapel := []GameCard{startCard}

	numPlayers := 4
	players := make([][]GameCard, numPlayers)

	for i := 0; i < numPlayers; i++ {
		players[i], gameDeck = drawCards(gameDeck, 7)
		fmt.Printf("Player %d Cards: %v\n", i+1, players[i])
	}

	startPlayerIndex := rand.Intn(numPlayers)
	currentPlayerIndex := startPlayerIndex

	fmt.Println("The game begins with Player", currentPlayerIndex+1)

	gameOver := false
	roundNumber := 1
	var rounds []RoundData

	for !gameOver {
		currentPlayer := currentPlayerIndex
		fmt.Println("\nCurrent Player:", currentPlayer+1)

		currentCard := gameStapel[len(gameStapel)-1]
		fmt.Println("Current card on the stack:", currentCard.Color, "-", currentCard.Value)

		playerHand := make([]GameCard, len(players[currentPlayer]))
		copy(playerHand, players[currentPlayer])
		fmt.Println("Cards in hand:", playerHand)

		playedCardIndex := -1
		for i, card := range playerHand {
			if canPlayCard(currentCard, card) {
				playedCardIndex = i
				break
			}
		}

		if playedCardIndex != -1 {
			// Play card
			playedCard := playerHand[playedCardIndex]
			fmt.Println("Player", currentPlayer+1, "plays a card:", playedCard.Color, "-", playedCard.Value)
			gameStapel = append(gameStapel, playedCard)

			playerHand = append(playerHand[:playedCardIndex], playerHand[playedCardIndex+1:]...)
			players[currentPlayer] = playerHand

			nextPlayer := nextPlayerIndex(currentPlayer, numPlayers)
			if playedCard.Value == "Draw Four" {
				fmt.Println("Player", nextPlayer+1, "has to draw 4 cards.")
				drawnCards, remainingDeck := drawCards(gameDeck, 4)
				players[nextPlayer] = append(players[nextPlayer], drawnCards...)
				gameDeck = remainingDeck

				if gameStapel[len(gameStapel)-1].Value != "Draw Four" && gameStapel[len(gameStapel)-1].Value != "Color Choice" {
					chosenColor := playerHand[0].Color
					fmt.Println("Player", currentPlayer+1, "chooses the color:", chosenColor)
					gameStapel[len(gameStapel)-1].Color = chosenColor
				} else {
					colors := []string{"Blue", "Green", "Red", "Yellow"}
					chosenColor := colors[rand.Intn(len(colors))]
					fmt.Println("Player", currentPlayer+1, "chooses the color:", chosenColor)
					gameStapel[len(gameStapel)-1].Color = chosenColor
				}
			} else if playedCard.Value == "Color Choice" {
				if gameStapel[len(gameStapel)-1].Value != "Draw Four" && gameStapel[len(gameStapel)-1].Value != "Color Choice" {
					chosenColor := playerHand[0].Color
					fmt.Println("Player", currentPlayer+1, "chooses the color:", chosenColor)
					gameStapel[len(gameStapel)-1].Color = chosenColor
				} else {
					colors := []string{"Blue", "Green", "Red", "Yellow"}
					chosenColor := colors[rand.Intn(len(colors))]
					fmt.Println("Player", currentPlayer+1, "chooses the color:", chosenColor)
					gameStapel[len(gameStapel)-1].Color = chosenColor
				}
			} else if playedCard.Value == "Draw Two" {
				nextPlayer := nextPlayerIndex(currentPlayer, numPlayers)
				fmt.Println("Player", nextPlayer+1, "has to draw 2 cards.")
				drawnCards, remainingDeck := drawCards(gameDeck, 2)
				players[nextPlayer] = append(players[nextPlayer], drawnCards...)
				gameDeck = remainingDeck
			}
		} else {
			drawnCard, remainingDeck := drawCards(gameDeck, 1)
			fmt.Println("Player", currentPlayer+1, "must draw:", drawnCard[0].Color, "-", drawnCard[0].Value)
			players[currentPlayer] = append(playerHand, drawnCard[0])
			gameDeck = remainingDeck
		}

		if len(players[currentPlayer]) == 0 {
			fmt.Println("Player", currentPlayer+1, "has no more cards and has won! The game is over")
			gameOver = true
		} else {
			currentPlayerIndex = nextPlayerIndex(currentPlayerIndex, numPlayers)
		}
		winProbabilities := calculateWinProbability(players, currentPlayerIndex)

		roundActions := []RoundAction{}
		for i := 0; i < numPlayers; i++ {
			var lastCard GameCard
			if len(players[i]) > 0 {
				lastCard = players[i][len(players[i])-1]
			} else {
				lastCard = GameCard{}
			}
			roundActions = append(roundActions, RoundAction{PlayerIndex: i, Card: lastCard})
		}

		playerHands := make([][]GameCard, len(players))
		for i, hand := range players {
			playerHands[i] = make([]GameCard, len(hand))
			copy(playerHands[i], hand)
		}

		// Make a copy of the current game stack so that previous rounds
		// remain unchanged even when the stack grows in later rounds.
		gameStackCopy := make([]GameCard, len(gameStapel))
		copy(gameStackCopy, gameStapel)

		roundData := RoundData{
			RoundNumber:        roundNumber,
			StartingPlayer:     startPlayerIndex,
			StartingCard:       startCard,
			PlayerHands:        playerHands,
			Actions:            roundActions,
			CurrentPlayerIndex: currentPlayerIndex,
			CurrentCard:        gameStapel[len(gameStapel)-1],
			GameStack:          gameStackCopy,
			WinProbabilities:   winProbabilities,
		}

		rounds = append(rounds, roundData)

		roundNumber++
	}
	filename := fmt.Sprintf("Round-%s.json", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(rounds)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
