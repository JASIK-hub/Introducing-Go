package main

import "testing"

var (
	handmansTale = Book{Author: "Margaret Atwood", Title: "The Handmaid's Tale"}
	oryxAndCrake = Book{
		Author: "Margaret Atwood", Title: "Oryx and Crake",
	}
	theBellJar = Book{Author: "Sylvia Plath", Title: "The Bell Jar"}
	janeEyre   = Book{Author: "Charlotte Brontë", Title: "Jane Eyre"}
)

func TestLoadBookworms(t *testing.T) {
	type testCase struct {
		filePath string
		want     []Bookworm
		wantErr  bool
	}

	tests := map[string]testCase{
		"file exists and contains bookworms": {
			filePath: "data/bookworms.json",
			want: []Bookworm{
				{
					Name:  "Fadi",
					Books: []Book{handmansTale, theBellJar},
				},
				{
					Name:  "Peggy",
					Books: []Book{oryxAndCrake, handmansTale, janeEyre},
				},
			},
			wantErr: false,
		},
		"file does not exist": {
			filePath: "testdata/no_file.json",
			want:     nil,
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := loadBookworms(tc.filePath)
			if (err != nil) != tc.wantErr {
				t.Fatalf("несовпадение по ошибке: got error = %v, wantErr %v", err, tc.wantErr)
			}
			if !equalBookworms(t, got, tc.want) {
				t.Fatalf("different result: got %v, expected %v",
					got, tc.want)
			}
		})
	}
}

func equalBookworms(t *testing.T, bookworms, target []Bookworm) bool {
	t.Helper()

	if len(bookworms) != len(target) {
		return false
	}

	for i := range bookworms {
		if bookworms[i].Name != target[i].Name {
			return false
		}
		if !equalBooks(t, bookworms[i].Books, target[i].Books) {
			return false
		}
	}
	return true
}

func equalBooks(t *testing.T, books, target []Book) bool {
	t.Helper()

	if len(books) != len(target) {
		return false
	}

	count := make(map[Book]int)
	for _, book := range books {
		count[book]++
	}

	for _, book := range target {
		if count[book] == 0 {
			return false
		}
		count[book]--
	}

	return true
}

func equalBooksCount(t *testing.T, got, want map[Book]uint) bool {
	t.Helper()

	if len(got) != len(want) {
		return false
	}

	for book, bookCount := range want {
		count, ok := got[book]
		if !ok || bookCount != count {
			return false
		}
	}
	return true
}

func TestBooksCount(t *testing.T) {
	tcs := map[string]struct {
		input []Bookworm
		want  map[Book]uint
	}{
		"standart scenario": {
			input: []Bookworm{{Name: "Ali", Books: []Book{handmansTale, oryxAndCrake}}},
			want: map[Book]uint{
				handmansTale: 1,
				oryxAndCrake: 1,
			},
		},
		"no bookWorms": {
			input: []Bookworm{},
			want:  map[Book]uint{},
		},
		"bookWorm without books": {
			input: []Bookworm{},
			want:  map[Book]uint{},
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := booksCount(tc.input)
			if !equalBooksCount(t, tc.want, got) {
				t.Fatalf("got a different list of books: %v,expected : %v", got, tc.want)
			}
		})
	}
}

func TestCommonBooks(t *testing.T) {
	tcs := map[string]struct {
		input []Bookworm
		want  []Book
	}{
		"Everyone has read the same books": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmansTale, theBellJar}},
				{Name: "Peggy", Books: []Book{handmansTale, theBellJar}},
			},
			want: []Book{handmansTale, theBellJar},
		},
		"People have completely different lists": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmansTale, theBellJar}},
				{Name: "Peggy", Books: []Book{janeEyre, oryxAndCrake}}},
			want: []Book{},
		},
		"One bookworm has no books": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{}},
				{Name: "Peggy", Books: []Book{handmansTale, theBellJar}},
			},
			want: []Book{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := findCommonBooks(tc.input)
			if !equalBooks(t, got, tc.want) {
				t.Fatalf("got a different list of books: %v,expected : %v", got, tc.want)
			}
		})
	}
}
