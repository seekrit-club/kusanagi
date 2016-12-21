package main

/*
WARNING: HACK APPROACHING!

The suggested way to store squares is with bit fiddling. So a square is a
byte, with the following structure:

000XCPPP

Where 0 is junk, X is validity, C is colour, and PPP is piece data.
*/

const (
	EMPTY byte = iota
	PAWN
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

const (
	WHITE byte = 0x00
	BLACK byte = 0x08
)

const (
	ONBOARD  byte = 0x00
	OFFBOARD byte = 0x10
)

type Board struct {
	/* Mailbox style, 10x12 board. */
	Data [120]byte
}

const (
	A1 int = 21
	H8 int = 98
)

func ClearBoard(b *Board) {
	for i := 0; i < 120; i++ {
		if i < A1 || i > H8 {
			b.Data[i] = OFFBOARD
		} else if i%10 == 0 || i%10 == 9 {
			b.Data[i] = OFFBOARD
		} else {
			b.Data[i] = ONBOARD
		}
	}
}

func InitBoard(b *Board) {
	ClearBoard(b)
}

func PrintBoard(b *Board) string {
	retval := ""
	for i := 0; i < 120; i++ {
		retval += ByteToString(b.Data[i])
		if i%10 == 0 && i > A1 && i < H8 {
			retval += "\n"
		}
	}
	return retval
}

func ByteToString(b byte) string {
	if b&OFFBOARD != 0 {
		return ""
	}
	if b == EMPTY {
		return "."
	}
	return "?"
}
