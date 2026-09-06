package main

import (
	"fmt"
	"os"
)

func main() {
	bookworms, err := loadBookworms(`./data/bookworms.json`)
	if err != nil {
		fmt.Printf(`failed to load bookworms:%v\n`, err)
		os.Exit(1)
	}
	commonBooks := findCommonBooks(bookworms)
	displayBooks(commonBooks)
}
