package gordle

type hint byte

const (
	absentCharacter hint = iota
	wrongPosiotion
	correctPosition
)
