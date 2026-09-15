package main

import (
	"learngo-pockets/gordle/gordle"
	"os"
)

const maxAttempts = 6

func main() {
	solution := "Hello"
	g := gordle.New(os.Stdin, solution, maxAttempts)
	g.Play()
}
