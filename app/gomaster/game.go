package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

type Colour string

var allColours = [...]Colour{"red", "pink", "orange", "yellow", "green", "blue", "purple", "black", "gray"}

type GameConfig struct {
	NumberOfColours   int
	NumberOfPositions int
}

type Game struct {
	NumberOfPositions int
	Colours           []Colour
	Complete          bool
	Answer            []Colour
	Attempts          []Attempt
}

type Attempt struct {
	Positions             []Colour
	Index                 int
	Complete              bool
	RightColourRightPlace int
	RightColourWrongPlace int
	WrongColour           int
}

func StartGame(config GameConfig) (*Game, error) {

	numberOfColours := Min(config.NumberOfColours, len(allColours)-1)

	if numberOfColours <= 0 || numberOfColours >= len(allColours) {
		return nil, errors.New(fmt.Sprintf("Invalid number of colours: %v (max is %v)", numberOfColours, len(allColours)))
	}

	numberOfPositions := Min(config.NumberOfPositions, numberOfColours)

	if numberOfPositions <= 0 || numberOfPositions >= len(allColours) {
		return nil, errors.New(fmt.Sprintf("Invalid number of positions: %v (max is %v)", numberOfPositions, len(allColours)))
	}

	log.Printf("colours %v positions %v", numberOfColours, numberOfPositions)

	colours := func() []Colour {
		colours := make([]Colour, len(allColours))
		copy(colours, allColours[:])
		rand.Shuffle(len(colours), func(i int, j int) {
			colours[i], colours[j] = colours[j], colours[i]
		})
		return colours[:numberOfColours]
	}()

	answer := func() []Colour {
		answer := make([]Colour, numberOfPositions)
		copy(answer, colours[:numberOfPositions])
		rand.Shuffle(len(answer), func(i int, j int) {
			answer[i], answer[j] = answer[j], answer[i]
		})
		return answer[:]
	}()

	return &Game{NumberOfPositions: numberOfPositions, Colours: colours, Answer: answer}, nil
}

func (game *Game) SubmitGuess(positions []Colour) (*Attempt, error) {

	if game.Complete {
		return nil, errors.New("game already complete")
	}

	if len(positions) != game.NumberOfPositions {
		return nil, errors.New(fmt.Sprintf("invalid guess size: %v (expected: %v)", len(positions), game.NumberOfPositions))
	}

	for _, p := range positions {
		if !Contains(game.Colours, p) {
			return nil, errors.New(fmt.Sprintf("Colour not supported: %v", p))
		}
	}

	attempt := Attempt{Positions: positions, Index: len(game.Attempts)}

	for k, v := range positions {
		if game.Answer[k] == v {
			attempt.RightColourRightPlace++
		} else if Contains(game.Answer, v) {
			attempt.RightColourWrongPlace++
		} else {
			attempt.WrongColour++
		}
	}

	if attempt.RightColourRightPlace == len(positions) {
		attempt.Complete = true
		game.Complete = true
	}

	game.Attempts = append(game.Attempts, attempt)

	return &attempt, nil
}
