package main

import (
	"strings"
)

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
	b.Data[A1] = KING | WHITE
	b.Data[A1+10] = PAWN | WHITE
	b.Data[H8] = KING | BLACK
	b.Data[H8-1] = ROOK | BLACK
}

func PrintBoard(b *Board) string {
	retval := ""
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			retval +=
				ByteToString(b.Data[CartesianToIndex(file, rank)])
		}
		retval += "\n"
	}
	return retval
}

func CartesianToIndex(file, rank int) int {
	return 21 + (10 * rank) + file
}

func ByteToString(b byte) string {
	if ByteIsOffboard(b) {
		return ""
	}
	retval := ""
	switch GetPiece(b) {
	case EMPTY:
		return "."
	case PAWN:
		retval = "P"
	case KNIGHT:
		retval = "N"
	case BISHOP:
		retval = "B"
	case ROOK:
		retval = "R"
	case QUEEN:
		retval = "Q"
	case KING:
		retval = "K"
	default:
		return "?"
	}
	if IsBlack(b) {
		return strings.ToLower(retval)
	}
	return retval
}

func ByteIsOffboard(b byte) bool {
	return b&OFFBOARD == OFFBOARD
}

func GetPiece(b byte) byte {
	return b & 0x07
}

func IsBlack(b byte) bool {
	return b&BLACK == BLACK
}
