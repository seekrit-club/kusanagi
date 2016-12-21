package main

const (
	MoveQuiet byte = iota
)

type Move struct {
	From    byte
	To      byte
	Kind    byte
	Promote byte
	Score   int16
}

func MoveGen(b *Board) []Move {
	retval := make([]Move, 0, 32)
	for i := A1; i <= H8; i++ {
		if !OnBoard(i) || GetPiece(b.Data[i]) == EMPTY {
			continue
		}
		if GetPiece(b.Data[i]) == PAWN {
			PawnPush := i + 10
			if b.ToMove == BLACK {
				PawnPush = i - 10
			}
			if GetPiece(b.Data[PawnPush]) == EMPTY &&
				GetSide(b.Data[i]) == b.ToMove {
				retval = append(retval, Move{byte(i),
					byte(PawnPush), MoveQuiet, EMPTY, 0})
			}
		}
	}
	return retval
}

func MakeMove(b *Board, m *Move) {
	b.Data[m.To] = b.Data[m.From]
	b.Data[m.From] = EMPTY
	b.ToMove ^= BLACK
}
