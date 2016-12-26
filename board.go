package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
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

const (
	CASTLEWK byte = 0x01
	CASTLEWQ byte = 0x02
	CASTLEBK byte = 0x04
	CASTLEBQ byte = 0x08
)

type Board struct {
	/* Mailbox style, 10x12 board. */
	Data      [120]byte
	ToMove    byte
	Castle    byte
	EnPassant byte
}

const (
	A1 byte = 21
	H8 byte = 98
)

const START string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const INVALID byte = 0

func ClearBoard(b *Board) {
	b.EnPassant = 1
	var i byte
	for i = 0; i < 120; i++ {
		if OnBoard(i) {
			b.Data[i] = ONBOARD
		} else {
			b.Data[i] = OFFBOARD
		}
	}
}

func Parse(fen string) (*Board, error) {
	b := new(Board)
	ClearBoard(b)
	var rank byte = 7
	var file byte = 0
	var eprank int = 0 /* Deliberate; this stores Atoi's result later on. */
	var epfile byte = 0
	stage := 0
	for _, runeValue := range fen {
		switch stage {
		case 0:
			/* Fill the data */
			if runeValue >= '1' && runeValue <= '8' {
				inc, _ := strconv.Atoi(string(runeValue))
				file += byte(inc)
			} else if runeValue == '/' {
				rank -= 1
				file = 0
			} else if runeValue == ' ' {
				stage++
			} else {
				switch unicode.ToUpper(runeValue) {
				case 'P':
					b.Data[CartesianToIndex(file, rank)] |= PAWN
				case 'N':
					b.Data[CartesianToIndex(file, rank)] |= KNIGHT
				case 'B':
					b.Data[CartesianToIndex(file, rank)] |= BISHOP
				case 'R':
					b.Data[CartesianToIndex(file, rank)] |= ROOK
				case 'Q':
					b.Data[CartesianToIndex(file, rank)] |= QUEEN
				case 'K':
					b.Data[CartesianToIndex(file, rank)] |= KING
				default:
					return nil, errors.New("Unexpected character in board data")
				}
				if unicode.IsLower(runeValue) {
					b.Data[CartesianToIndex(file, rank)] |= BLACK
				}
				file += 1
			}
		case 1:
			/* Get who's to play next */
			switch runeValue {
			case 'w':
				b.ToMove = WHITE
			case 'b':
				b.ToMove = BLACK
			case ' ':
				stage++
			default:
				return nil, errors.New("Unexpected character for active colour")
			}
		case 2:
			/* Castling */
			switch runeValue {
			case '-':
				/* Do nothing */
			case 'K':
				b.Castle |= CASTLEWK
			case 'Q':
				b.Castle |= CASTLEWQ
			case 'k':
				b.Castle |= CASTLEBK
			case 'q':
				b.Castle |= CASTLEBQ
			case ' ':
				stage++
			default:
				return nil, errors.New("Unexpected character for castling")
			}
		case 3:
			/* En-passant */
			if runeValue >= '1' && runeValue <= '8' {
				eprank, _ = strconv.Atoi(string(runeValue))
				eprank--
			} else if runeValue >= 'a' && runeValue <= 'h' {
				epfile = byte(runeValue - 'a')
			} else if runeValue == ' ' {
				if b.EnPassant != INVALID {
					b.EnPassant = CartesianToIndex(epfile, byte(eprank))
				}
				stage++
			} else if runeValue == '-' {
				b.EnPassant = INVALID
			} else {
				return nil, errors.New("Unexpected character for en passant")
			}
		}
	}
	return b, nil
}

func PrintBoard(b *Board) string {
	retval := ""
	var rank, file byte
	for rank = 7; rank != 255; rank-- {
		for file = 0; file < 8; file++ {
			retval +=
				ByteToString(b.Data[CartesianToIndex(file, rank)])
		}
		retval += "\n"
	}
	return retval
}

func CartesianToIndex(file, rank byte) byte {
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

func GetSide(b byte) byte {
	return b & BLACK
}

func IsBlack(b byte) bool {
	return GetSide(b) == BLACK
}

func OnBoard(i byte) bool {
	return i >= A1 && i <= H8 && !(i%10 == 0 || i%10 == 9)
}

func IndexToCartesian(index byte) (byte, byte) {
	var file, rank byte
	file = (index % 10) - 1
	rank = (index / 10) - 2
	return file, rank
}

func CartesianToAlgebraic(file, rank byte) string {
	afile := rune(byte('a') + file)
	arank := rank + 1
	return fmt.Sprintf("%c%d", afile, arank)
}

func IndexToAlgebraic(i byte) string {
	return CartesianToAlgebraic(IndexToCartesian(i))
}

func FindKing(b *Board, colour byte) (byte, error) {
	for king := A1; king <= H8; king++ {
		if b.Data[king] != OFFBOARD && GetPiece(b.Data[king]) == KING &&
			GetSide(b.Data[king]) == colour {
			return king, nil
		}
	}
	return INVALID, errors.New("Couldn't find the king")
}
