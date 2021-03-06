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
	MovePromote
	MoveCapPromote
	MoveCastle
)

type Move struct {
	From    byte
	To      byte
	Kind    byte
	Promote byte
	Score   int
}

type Undo struct {
	ToData    byte
	EnPassant byte
	Castle    byte
	Index     byte
}

func pawnmove(b *Board, i byte, retval []Move) []Move {
	var PawnPush, DoublePush byte
	CanDouble := false
	CanPromote := false
	if b.ToMove == BLACK {
		PawnPush = i - 10
		DoublePush = i - 20
		CanDouble = i/10 == 8
		CanPromote = i/10 == 3
	} else {
		PawnPush = i + 10
		DoublePush = i + 20
		CanDouble = i/10 == 3
		CanPromote = i/10 == 8
	}
	if GetPiece(b.Data[PawnPush]) == EMPTY {
		if CanPromote {
			retval = append(retval, Move{i,
				PawnPush, MovePromote, QUEEN, 0})
			retval = append(retval, Move{i,
				PawnPush, MovePromote, ROOK, 0})
			retval = append(retval, Move{i,
				PawnPush, MovePromote, BISHOP, 0})
			retval = append(retval, Move{i,
				PawnPush, MovePromote, KNIGHT, 0})
		} else {
			retval = append(retval, Move{i,
				PawnPush, MoveQuiet, EMPTY, 0})
		}
		if CanDouble && GetPiece(b.Data[DoublePush]) ==
			EMPTY {
			retval = append(retval, Move{i,
				DoublePush, MoveDoublePush, EMPTY, 0})
		}
	}
	retval = pawncap(b, i, retval, PawnPush-1, CanPromote)
	retval = pawncap(b, i, retval, PawnPush+1, CanPromote)
	return retval
}

func pawncap(b *Board, i byte, retval []Move, place byte, CanPromote bool) []Move {
	if OnBoard(place) && GetPiece(b.Data[place]) != EMPTY &&
		GetSide(b.Data[place]) != b.ToMove {
		if CanPromote {
			retval = append(retval, Move{i,
				place, MoveCapPromote, QUEEN, 0})
			retval = append(retval, Move{i,
				place, MoveCapPromote, ROOK, 0})
			retval = append(retval, Move{i,
				place, MoveCapPromote, BISHOP, 0})
			retval = append(retval, Move{i,
				place, MoveCapPromote, KNIGHT, 0})
		} else {
			retval = append(retval, Move{i,
				place, MoveCapture, EMPTY, 0})
		}
	} else if OnBoard(place) && GetPiece(b.Data[place]) == EMPTY && b.EnPassant == place {
		retval = append(retval, Move{i, place, MoveEnPassant, EMPTY, 0})
	}
	return retval
}

func InCheck(b *Board) bool {
	if b.ToMove == BLACK && squareattacked(b, b.BlackKing, WHITE) {
		return true
	} else if b.ToMove == WHITE && squareattacked(b, b.WhiteKing, BLACK) {
		return true
	}
	return false
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
		to = byte(int(i) + Vector[KING][dir])
		if b.Data[to] != OFFBOARD && GetPiece(b.Data[to]) == KING && GetSide(b.Data[to]) == attacking {
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

func castle(b *Board, retval []Move, file1, file2, file3, file4 byte) []Move {
	var rank byte
	if b.ToMove == BLACK {
		rank = 7
	} else {
		rank = 0
	}
	sq1 := CartesianToIndex(file1, rank)
	sq2 := CartesianToIndex(file2, rank)
	sq3 := CartesianToIndex(file3, rank)
	sq4 := CartesianToIndex(file4, rank)
	if GetPiece(b.Data[sq2]) != EMPTY || GetPiece(b.Data[sq3]) != EMPTY {
		return retval
	}
	if file4 != 0 && GetPiece(b.Data[sq4]) != EMPTY {
		return retval
	}
	enemy := b.ToMove ^ BLACK
	if !squareattacked(b, sq1, enemy) && !squareattacked(b, sq2, enemy) && !squareattacked(b, sq3, enemy) {
		retval = append(retval, Move{sq1, sq3, MoveCastle, EMPTY, 0})
	}
	return retval
}

func qscastle(b *Board, retval []Move) []Move {
	return castle(b, retval, 4, 3, 2, 1)
}

func kscastle(b *Board, retval []Move) []Move {
	return castle(b, retval, 4, 5, 6, 0)
}

func MoveGen(b *Board) []Move {
	retval := make([]Move, 0, 32)
	if CanCastle(b, b.ToMove, QUEEN) {
		retval = qscastle(b, retval)
	}
	if CanCastle(b, b.ToMove, KING) {
		retval = kscastle(b, retval)
	}
	for _, i := range b.PieceList {
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
	retval := &Undo{b.Data[m.To], b.EnPassant, b.Castle, OFFBOARD}
	if GetPiece(b.Data[m.From]) == KING {
		if b.ToMove == BLACK {
			b.BlackKing = m.To
		} else {
			b.WhiteKing = m.To
		}
	}
	if m.Kind == MoveCapture || m.Kind == MoveCapPromote {
		retval.Index, _ = FindPiece(b, m.To)
		b.PieceList[retval.Index] = OFFBOARD
	}
	b.EnPassant = INVALID
	b.Data[m.To] = b.Data[m.From]
	b.Data[m.From] = EMPTY
	idx, _ := FindPiece(b, m.From)
	b.PieceList[idx] = m.To
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
			retval.Index, _ = FindPiece(b, m.To+10)
			b.Data[m.To+10] = EMPTY
		} else {
			retval.Index, _ = FindPiece(b, m.To-10)
			b.Data[m.To-10] = EMPTY
		}
		b.PieceList[retval.Index] = OFFBOARD
	case MoveCapPromote:
		fallthrough
	case MovePromote:
		b.Data[m.To] = b.ToMove | m.Promote
	case MoveCastle:
		if m.To < m.From {
			/* Queenside */
			retval.Index, _ = FindPiece(b, m.To-2)
			b.PieceList[retval.Index] = m.To + 1
			b.Data[m.To+1] = b.Data[m.To-2]
			b.Data[m.To-2] = EMPTY
		} else {
			/* Kingside */
			retval.Index, _ = FindPiece(b, m.To+1)
			b.PieceList[retval.Index] = m.To - 1
			b.Data[m.To-1] = b.Data[m.To+1]
			b.Data[m.To+1] = EMPTY
		}
	}
	b.Castle &= CASTLEMASK[m.From] & CASTLEMASK[m.To]
	b.ToMove ^= BLACK
	return retval
}

func UnmakeMove(b *Board, m *Move, u *Undo) {
	b.Data[m.From] = b.Data[m.To]
	b.Data[m.To] = u.ToData
	idx, _ := FindPiece(b, m.To)
	b.PieceList[idx] = m.From
	b.EnPassant = u.EnPassant
	b.Castle = u.Castle
	b.ToMove ^= BLACK
	switch m.Kind {
	case MoveCapture:
		b.PieceList[u.Index] = m.To
	case MoveEnPassant:
		if b.ToMove == BLACK {
			b.PieceList[u.Index] = m.To + 10
			b.Data[m.To+10] = (b.ToMove ^ BLACK) | PAWN
		} else {
			b.PieceList[u.Index] = m.To - 10
			b.Data[m.To-10] = (b.ToMove ^ BLACK) | PAWN
		}
	case MoveCapPromote:
		b.PieceList[u.Index] = m.To
		fallthrough
	case MovePromote:
		b.Data[m.From] = b.ToMove | PAWN
	case MoveCastle:
		if m.To < m.From {
			/* Queenside */
			b.PieceList[u.Index] = m.To - 2
			b.Data[m.To-2] = b.Data[m.To+1]
			b.Data[m.To+1] = EMPTY
		} else {
			/* Kingside */
			b.PieceList[u.Index] = m.To + 1
			b.Data[m.To+1] = b.Data[m.To-1]
			b.Data[m.To-1] = EMPTY
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
			fmt.Print(MoveToLongAlgebraic(&move))
		}
		tmp := Perft(depth-1, board, false)
		nodes += tmp
		if divide {
			fmt.Printf(" %d\n", tmp)
		}
		UnmakeMove(board, &move, undo)

	}
	return nodes
}

func MoveToLongAlgebraic(move *Move) string {
	promote := ""
	if move.Kind == MovePromote || move.Kind == MoveCapPromote {
		switch move.Promote {
		case QUEEN:
			promote = "q"
		case ROOK:
			promote = "r"
		case KNIGHT:
			promote = "n"
		case BISHOP:
			promote = "b"
		}
	}
	return fmt.Sprintf("%s%s%s", IndexToAlgebraic(move.From), IndexToAlgebraic(move.To), promote)
}

func ParseMove(b *Board, m string) (*Move, error) {
	moves := MoveGen(b)
	for _, move := range moves {
		if MoveToLongAlgebraic(&move) == m {
			return &move, nil
		}
	}
	return nil, errors.New("Move not legal")
}

func FilterCaptures(movelist []Move) []Move {
	captures := make([]Move, 0, 32)
	for _, move := range movelist {
		if move.Kind == MoveCapture || move.Kind == MoveCapPromote {
			captures = append(captures, move)
		}
	}
	return captures
}
