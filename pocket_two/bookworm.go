package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Bookworm struct {
	Name  string `json:"name"`
	Books []Book `json:"books"`
}
type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}
type byAuthor []Book
type bookRecommendations map[Book]map[Book]uint

func loadBookworms(filePath string) ([]Bookworm, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var bookworms []Bookworm

	err = json.NewDecoder(f).Decode(&bookworms)
	if err != nil {
		return nil, err
	}
	return bookworms, nil
}

func booksCount(bookWorms []Bookworm) map[Book]uint {
	count := make(map[Book]uint)

	for _, bookWorm := range bookWorms {
		for _, book := range bookWorm.Books {
			count[book]++
		}
	}
	return count
}

func findCommonBooks(bookWorms []Bookworm) []Book {
	booksOnShelves := booksCount(bookWorms)
	var commonBooks []Book
	for book, count := range booksOnShelves {
		if count > 1 {
			commonBooks = append(commonBooks, book)
		}
	}
	return sortBooks(commonBooks)
}

func sortBooks(books []Book) []Book {
	sort.Sort(byAuthor(books))
	return books
}

func displayBooks(books []Book) {
	for _, book := range books {
		fmt.Println("-", book.Title, "by", book.Author)
	}
}

func (b byAuthor) Len() int { return len(b) }

func (b byAuthor) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byAuthor) Less(i, j int) bool {
	if b[i].Author != b[j].Author {
		return b[i].Author < b[j].Author
	}
	return b[i].Title < b[j].Title
}

func recommendOtherBooks(bookWorms []Bookworm) []Bookworm {

	sb := make(bookRecommendations)

	for _, bookWorm := range bookWorms {
		for i, book := range bookWorm.Books {
			otherBooksOnShelves := listOtherBooksOnShelves(i, bookWorm.Books)
			registerBookRecommendations(sb, book, otherBooksOnShelves)
		}
	}
	result := make([]Bookworm, 0, len(bookWorms))
	for i, bookWorm := range bookWorms {
		result[i] = Bookworm{Name: bookWorm.Name, Books: recommendBooks(sb, bookWorm.Books)}
	}
	return result
}

func listOtherBooksOnShelves(i int, books []Book) []Book {
	otherBooks := make([]Book, 0, len(books)-1)
	for idx, book := range books {
		if idx == i {
			continue
		}
		otherBooks = append(otherBooks, book)
	}
	return otherBooks
}

func registerBookRecommendations(sb bookRecommendations, target Book, otherBooksOnShelves []Book) {
	if sb[target] == nil {
		sb[target] = make(map[Book]uint)
	}
	for _, otherBook := range otherBooksOnShelves {
		sb[target][otherBook]++
	}
}

func recommendBooks(sb bookRecommendations, userBooks []Book) []Book {
	counts := make(map[Book]uint)
	for _, userBook := range userBooks {
		for relatedBook, score := range sb[userBook] {
			counts[relatedBook] += score
		}
	}

	for _, userBook := range userBooks {
		delete(counts, userBook)
	}

	type recommendation struct {
		book  Book
		count uint
	}

	list := make([]recommendation, 0, len(counts))
	for book, count := range counts {
		list = append(list, recommendation{book: book, count: count})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].count > list[j].count
	})

	result := make([]Book, 0, len(list))
	for _, rec := range list {
		result = append(result, rec.book)
	}
	return result

}
