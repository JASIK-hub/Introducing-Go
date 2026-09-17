package main

import (
	"learngo-pockets/gordle/gordle"
	"os"
)

func main() {
	solution := "Hello"
	g := gordle.New(os.Stdin, solution, gordle.MaxAttempts)
	g.Play()
}
