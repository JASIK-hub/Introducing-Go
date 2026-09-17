package gordle

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestGameAsk(t *testing.T) {
	tt := map[string]struct {
		input string
		want  []rune
	}{
		"5 characters in english": {
			input: "HELLO",
			want:  []rune("HELLO"),
		},
		"5 characters in arabic": {
			input: "مرحبا",
			want:  []rune("مرحبا"),
		},
		"5 characters in japanese": {
			input: "こんにちは",
			want:  []rune("こんにちは"),
		},
		"3 characters in japanese": {
			input: "こんに\nこんにちは",
			want:  []rune("こんにちは"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g := New(strings.NewReader(tc.input), string(tc.want), MaxAttempts)
			got := g.ask()
			if !slices.Equal(got, tc.want) {
				t.Errorf("got = %v, want %v", string(got), string(tc.want))
			}
		})
	}
}

func TestValidateGuess(t *testing.T) {
	tt := map[string]struct {
		input   []rune
		wantErr error
	}{
		"nominal": {
			input:   []rune("GUESS"),
			wantErr: nil,
		},
		"too long guess": {
			input:   []rune("gordle-game"),
			wantErr: errInvalidWordLength,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			g := Game{}
			err := g.validateGuess(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}

}

func TestComputeFeedback(t *testing.T) {
	tt := map[string]struct {
		guess            []rune
		solution         []rune
		expectedFeedback feedback
	}{
		"nominal": {
			guess:    []rune("hello"),
			solution: []rune("hello"),
			expectedFeedback: feedback{
				correctPosition,
				correctPosition,
				correctPosition,
				correctPosition,
				correctPosition,
			},
		},
		"double character": {
			guess:    []rune("hello"),
			solution: []rune("world"),
			expectedFeedback: feedback{
				absentCharacter,
				absentCharacter,
				absentCharacter,
				correctPosition,
				absentCharacter,
			},
		},
		"double character with wrong answer": {
			guess:    []rune("geeks"),
			solution: []rune("hello"),
			expectedFeedback: feedback{
				absentCharacter,
				correctPosition,
				absentCharacter,
				absentCharacter,
				absentCharacter,
			},
		},
		"two identical, but not in the right position": {
			guess:    []rune("hlleo"),
			solution: []rune("hello"),
			expectedFeedback: feedback{
				correctPosition,
				wrongPosition,
				correctPosition,
				wrongPosition,
				correctPosition,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			fb := computeFeedback(tc.guess, tc.solution)
			if !fb.Equal(tc.expectedFeedback) {
				t.Errorf("guess: %q, got the wrong feedback, wanted %v, got %v",
					tc.guess, tc.expectedFeedback, fb)
			}
		})
	}
}
