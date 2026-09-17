package gordle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

const MaxAttempts = 6

type Game struct {
	reader      *bufio.Reader
	solution    []rune
	maxAttempts int
}

var errInvalidWordLength = fmt.Errorf("invalid guess, word doesn't have the same number of characters as the solution")

func New(playerInput io.Reader, solution string, maxAttempts int) *Game {
	g := &Game{
		reader:      bufio.NewReader(playerInput),
		solution:    splitToUppercaseCharacters(solution),
		maxAttempts: maxAttempts,
	}
	return g
}

func (g *Game) Play() {
	fmt.Println("Welcome to Gordle!")
	for currentAttempt := 1; currentAttempt <= g.maxAttempts; currentAttempt++ {
		guess := g.ask()
		if slices.Equal(guess, g.solution) {
			fmt.Printf("You won! You found it in %d guess(es)! The word was: %s.\n", currentAttempt, string(g.solution))
			return
		}
	}
	fmt.Printf("You've lost! The solution was: %s. \n", string(g.solution))
}

func (g *Game) validateGuess(guess []rune) error {
	if len(guess) != len(g.solution) {
		return fmt.Errorf("expected %d, got %d, %w",
			len(g.solution), len(guess), errInvalidWordLength)
	}
	return nil
}

func (g *Game) ask() []rune {
	fmt.Printf("Enter a %d-character guess:\n", len(g.solution))
	for {
		playerInput, err := g.reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Gordle failed to read your guess: %s\n", err.Error())
			continue
		}
		guess := []rune(string(playerInput))
		err = g.validateGuess(guess)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Your attempt is invalid with Gordle's solution: %s.\n", err.Error())
		} else {
			return guess
		}
	}

}

func splitToUppercaseCharacters(input string) []rune {
	return []rune(strings.ToUpper(input))
}
