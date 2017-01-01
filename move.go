package main

import (
	"errors"
	"fmt"
)

var Vector [8][8]int = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},               // empty
	{0, 0, 0, 0, 0, 0, 0, 0},               // pawn - handled specially.
	{+21, +12, -8, -19, -21, -12, +8, +19}, // N
	{+11, -9, -11, +9, 0, 0, 0, 0},         // B
	{+10, +1, -10, -1, 0, 0, 0, 0},         // R
	{+10, +11, +1, -9, -10, -11, -1, +9},   // Q
	{+10, +11, +1, -9, -10, -11, -1, +9},   // K
	{0, 0, 0, 0, 0, 0, 0, 0},               // ?
}

var Slide [8]bool = [8]bool{false, false, false, true, true, true, false, false}

const (
	MoveQuiet byte = iota
	MoveDoublePush
	MoveCapture
	MoveEnPassant
)

type Move struct {
	From    byte
	To      byte
	Kind    byte
	Promote byte
	Score   int16
}

type Undo struct {
	ToData    byte
	EnPassant byte
}

func pawnmove(b *Board, i byte, retval []Move) []Move {
	var PawnPush, DoublePush byte
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
	if GetPiece(b.Data[PawnPush]) == EMPTY {
		retval = append(retval, Move{i,
			PawnPush, MoveQuiet, EMPTY, 0})
		if CanDouble && GetPiece(b.Data[DoublePush]) ==
			EMPTY {
			retval = append(retval, Move{i,
				DoublePush, MoveDoublePush, EMPTY, 0})
		}
	}
	retval = pawncap(b, i, retval, PawnPush-1)
	retval = pawncap(b, i, retval, PawnPush+1)
	return retval
}

func pawncap(b *Board, i byte, retval []Move, place byte) []Move {
	if OnBoard(place) && GetPiece(b.Data[place]) != EMPTY &&
		GetSide(b.Data[place]) != b.ToMove {
		retval = append(retval, Move{i,
			place, MoveCapture, EMPTY, 0})
	} else if OnBoard(place) && GetPiece(b.Data[place]) == EMPTY && b.EnPassant == place {
		retval = append(retval, Move{i, place, MoveEnPassant, EMPTY, 0})
	}
	return retval
}

func squareattacked(b *Board, i byte, attacking byte) bool {
	var PawnPush byte
	if attacking == BLACK {
		PawnPush = i + 10
	} else {
		PawnPush = i - 10
	}
	if (GetSide(b.Data[PawnPush-1]) == attacking && GetPiece(b.Data[PawnPush-1]) == PAWN) || (GetSide(b.Data[PawnPush+1]) == attacking && GetPiece(b.Data[PawnPush+1]) == PAWN) {
		return true
	}
	for dir := 0; dir < 8; dir++ {
		from := i
		for {
			to := byte(int(from) + Vector[QUEEN][dir])
			piece := GetPiece(b.Data[to])
			if b.Data[to] == OFFBOARD || (piece != EMPTY && GetSide(b.Data[to]) != attacking) {
				break
			} else if piece == QUEEN && GetSide(b.Data[to]) == attacking {
				return true
			} else if piece == ROOK && GetSide(b.Data[to]) == attacking && (Vector[QUEEN][dir] == 10 || Vector[QUEEN][dir] == -10 || Vector[QUEEN][dir] == 1 || Vector[QUEEN][dir] == -1) {
				return true
			} else if piece == BISHOP && GetSide(b.Data[to]) == attacking && (Vector[QUEEN][dir] == 11 || Vector[QUEEN][dir] == -11 || Vector[QUEEN][dir] == 9 || Vector[QUEEN][dir] == -9) {
				return true
			} else if piece == EMPTY {
				from = to
			} else {
				break
			}
		}
		to := byte(int(i) + Vector[KNIGHT][dir])
		if b.Data[to] != OFFBOARD && GetPiece(b.Data[to]) == KNIGHT && GetSide(b.Data[to]) == attacking {
			return true
		}
	}
	return false
}

func quietmove(b *Board, i byte, retval []Move) []Move {
	piece := GetPiece(b.Data[i])
	for dir := 0; dir < 8; dir++ {
		if Vector[piece][dir] == 0 {
			break
		}
		from := i
		for {
			to := byte(int(from) + Vector[piece][dir])
			if b.Data[to] != OFFBOARD {
				if GetPiece(b.Data[to]) == EMPTY {
					retval = append(retval, Move{i,
						to, MoveQuiet, EMPTY, 0})
					if Slide[piece] {
						from = to
					} else {
						break
					}
				} else if GetSide(b.Data[to]) != b.ToMove {
					retval = append(retval, Move{i,
						to, MoveCapture, EMPTY, 0})
					break
				} else {
					break
				}
			} else {
				break
			}
		}
	}
	return retval
}

func MoveGen(b *Board) []Move {
	retval := make([]Move, 0, 32)
	for i := A1; i <= H8; i++ {
		if !OnBoard(i) || GetPiece(b.Data[i]) == EMPTY || GetSide(b.Data[i]) != b.ToMove {
			continue
		}
		if GetPiece(b.Data[i]) == PAWN {
			retval = pawnmove(b, i, retval)
		} else {
			retval = quietmove(b, i, retval)
		}
	}
	return retval
}

func MakeMove(b *Board, m *Move) *Undo {
	retval := &Undo{b.Data[m.To], b.EnPassant}
	if GetPiece(b.Data[m.From]) == KING {
		if b.ToMove == BLACK {
			b.BlackKing = m.To
		} else {
			b.WhiteKing = m.To
		}
	}
	b.EnPassant = INVALID
	b.Data[m.To] = b.Data[m.From]
	b.Data[m.From] = EMPTY
	switch m.Kind {
	case MoveQuiet:
		/* Do nothing */
	case MoveDoublePush:
		if b.ToMove == BLACK {
			b.EnPassant = m.From - 10
		} else {
			b.EnPassant = m.From + 10
		}
	case MoveEnPassant:
		if b.ToMove == BLACK {
			b.Data[m.To+10] = EMPTY
		} else {
			b.Data[m.To-10] = EMPTY
		}
	}
	b.ToMove ^= BLACK
	return retval
}

func UnmakeMove(b *Board, m *Move, u *Undo) {
	b.Data[m.From] = b.Data[m.To]
	b.Data[m.To] = u.ToData
	b.EnPassant = u.EnPassant
	b.ToMove ^= BLACK
	switch m.Kind {
	case MoveEnPassant:
		if b.ToMove == BLACK {
			b.Data[m.To+10] = (b.ToMove ^ BLACK) | PAWN
		} else {
			b.Data[m.To-10] = (b.ToMove ^ BLACK) | PAWN
		}
	}
	if GetPiece(b.Data[m.From]) == KING {
		if b.ToMove == BLACK {
			b.BlackKing = m.From
		} else {
			b.WhiteKing = m.From
		}
	}
}

func (m Move) String() string {
	return fmt.Sprint("{From: ", IndexToAlgebraic(m.From), " to: ",
		IndexToAlgebraic(m.To), " type: ", m.Kind, "}")
}

func DoDividePerft(depth int) uint64 {
	board, _ := Parse(START)
	return Perft(depth, board, true)
}

func Perft(depth int, board *Board, divide bool) uint64 {
	if depth == 0 {
		return 1
	}
	var nodes uint64 = 0
	moves := MoveGen(board)
	for _, move := range moves {
		undo := MakeMove(board, &move)
		if Illegal(board) {
			UnmakeMove(board, &move, undo)
			continue
		}
		if divide {
			fmt.Printf("%s%s ", IndexToAlgebraic(move.From), IndexToAlgebraic(move.To))
		}
		tmp := Perft(depth-1, board, false)
		nodes += tmp
		if divide {
			fmt.Printf("%d\n", tmp)
		}
		UnmakeMove(board, &move, undo)

	}
	return nodes
}

func ParseMove(m string) (*Move, error) {
	if len(m) != 4 {
		return nil, errors.New("Move is wrong length")
	}
	from := m[:2]
	to := m[2:]
	fromi, err1 := AlgebraicToIndex(from)
	if err1 != nil {
		return nil, err1
	}
	toi, err2 := AlgebraicToIndex(to)
	if err2 != nil {
		return nil, err2
	}
	return &Move{fromi, toi, MoveQuiet, EMPTY, 0}, nil
}
