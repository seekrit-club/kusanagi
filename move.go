package main

const (
	MoveQuiet byte = iota
	MoveDoublePush
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
			var PawnPush, DoublePush int
			CanDouble := false
			if b.ToMove == BLACK {
				PawnPush = i - 10
				DoublePush = i - 20
				CanDouble = i/10 == 8
			} else {
				PawnPush = i + 10
				DoublePush = i + 20
				CanDouble = i/10 == 3
			}
			if GetPiece(b.Data[PawnPush]) == EMPTY &&
				GetSide(b.Data[i]) == b.ToMove {
				retval = append(retval, Move{byte(i),
					byte(PawnPush), MoveQuiet, EMPTY, 0})
				if CanDouble && GetPiece(b.Data[DoublePush]) ==
					EMPTY {
					retval = append(retval, Move{byte(i),
						byte(DoublePush), MoveDoublePush, EMPTY, 0})
				}
			}
		}
	}
	return retval
}

func MakeMove(b *Board, m *Move) {
	b.EnPassant = INVALID
	b.Data[m.To] = b.Data[m.From]
	b.Data[m.From] = EMPTY
	switch m.Kind {
	case MoveQuiet:
		/* Do nothing */
	case MoveDoublePush:
		if b.ToMove == BLACK {
			b.EnPassant = int(m.From - 10)
		} else {
			b.EnPassant = int(m.From + 10)
		}
	}
	b.ToMove ^= BLACK
}
