package main

import (
	"flag"
	"fmt"
)

func main() {
	var lang string
	flag.StringVar(&lang, "lang", "en", "The required language, e.g. en, ur...")
	flag.Parse()
	greeting := greet(language(lang))
	fmt.Println(greeting)
}

type language string

var phraseBook = map[language]string{
	"el": "Χαίρετε Κόσμε",
	"en": "Hello world",
	"fr": "Bonjour le monde",
	"he": "שלום עולם",
	"ur": "یلو دنیا",
	"vi": "Xin chào Thế Giới",
}

func greet(l language) string {
	greeting, ok := phraseBook[l]

	if !ok {
		return fmt.Sprintf("unsupported language: %q", l)
	}
	return greeting
}
