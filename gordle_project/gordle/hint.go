package gordle

import "strings"

type hint byte
type feedback []hint

const (
	absentCharacter hint = iota
	wrongPosition
	correctPosition
)

func (h hint) String() string {
	switch h {
	case absentCharacter:
		return "🩶"
	case wrongPosition:
		return "💛"
	case correctPosition:
		return "💚"
	default:
		return "❤️"
	}
}

func (fb feedback) String() string {
	sb := strings.Builder{}
	for _, h := range fb {
		sb.WriteString(h.String())
	}
	return sb.String()
}

func computeFeedback(guess, solution []rune) feedback {
	result := make(feedback, len(guess))
	used := make([]bool, len(solution))

	for i := range guess {
		result[i] = absentCharacter
	}

	for i := range guess {
		if guess[i] == solution[i] {
			result[i] = correctPosition
			used[i] = true
		}
	}

	for i := range guess {
		if result[i] == correctPosition {
			continue
		}
		for j := range solution {
			if !used[j] && guess[i] == solution[j] {
				result[i] = wrongPosition
				used[j] = true
				break
			}
		}
	}
	return result
}

func (fb feedback) Equal(other feedback) bool {
	if len(fb) != len(other) {
		return false
	}
	for i, value := range fb {
		if value != other[i] {
			return false
		}
	}
	return true
}
