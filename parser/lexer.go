package parser

import (
	"strconv"
	"strings"
)

/*

Grammatic kinda:

ORDER = { SIDE ~ TICKER ~ AMOUNT ~  PRICE}  !
PRICE ::= {LADERED | FLOAT | PERCENT | DIFF | CLADERED} !
SIDE ::= {BUY | SELL} !
VARIABLE ::= {TEMP | FUNCTION | STRING | AMOUNT}
STRING ::= { CONSTANT | CONSTANT ~ STRING | LADERED | HIGH | LOW}
CONSTANT ::= {FLOAT | DIFF | PERCENT | MARKET}
AMOUNT ::= {FLOAT | UFLOAT | PERCENT}
FUNCTION ::={"func" ~ "(" ~ TEMPVAL ~ ")" ~ VARIABLE }
ASSIGN ::= { STRING ~ "="  ~ VARIABLE}
LADERED ::=  { "-l" ~ INTEGER ~ LADEREDVAR }
LADEREDVAR ::=  {FLOAT ~ FLOAT | UFLOAT ~ UFLOAT | PERCENT ~ PERCENT}
CLADERED ::= { "-l" ~ (high,low) ~ DURATION ~ LADEREDVAR}

SIDE   VARIABLE   PERCENT FLOAT UFLOAT VARIABLE
[buy,sell] [ticker] [10%, 0.1, u200] [(market, 38100, d100, 2%), (-(high,low) 3h [d100, 2%]), (-l [38100 38300, d100 d300, 1% 3%])
stop [ticker]


also

cancle -all
cancle btc-perp
cancle -buy

funding -position //funding rate der aktuellen positionen
funding -highest 20 //funding der highest 20 coins


stop buy btc-perp [position, u100,0.1] -low 5h
*/

type TokenType int

const (
	VARIABLE TokenType = iota
	SIDE
	STOP
	FLOAT
	UFLOAT
	DFLOAT
	PERCENT
	ASSIGN
	FLAG
	FUNC
	DURATION
	LBRACKET
	RBRACKET
	SOURCE
	LOAD
	CANCEL
	FUNDING
	POSITION
	FUNDINGRATES
)

type Token struct {
	Type TokenType
	Text string
}

//Lexer reads the input checks it for some errors and outputs a Token list
func Lexer(inputS string) (t []Token, err error) {

	input := strings.Split(inputS, " ")
	for _, s := range input {

		if len(s) == 0 {
			continue
		}
		last := len(s) - 1
		switch s {
		case "buy", "Buy":
			t = append(t, Token{SIDE, "buy"})
		case "sell", "Sell":
			t = append(t, Token{SIDE, "sell"})
		case "Stop", "stop":
			t = append(t, Token{STOP, "stop"})
		case "=":
			t = append(t, Token{ASSIGN, "="})
		case "func":
			t = append(t, Token{FUNC, "func"})
		case "(":
			t = append(t, Token{LBRACKET, ""})
		case ")":
			t = append(t, Token{RBRACKET, ""})
		case "cancel":
			t = append(t, Token{CANCEL, "cancel"})
		case "fundingpay":
			t = append(t, Token{FUNDING, "fundus"})
		case "position":
			t = append(t, Token{POSITION, ""})
		case "fundingrates":
			t = append(t, Token{FUNDINGRATES, ""})
		case "load":
			t = append(t, Token{LOAD, ""})
		default:
			if (s[last] == 'h' || s[last] == 'm' || s[last] == 'd') && len(s) > 1 {
				_, err := strconv.Atoi(s[:last])
				if err == nil {
					t = append(t, Token{DURATION, s})
					continue
				}
			}

			if len(s) > 6 {
				if s[:4] == "func" {
					t = append(t, Token{FUNC, "func"})
					t = append(t, lexFunc([]byte(s[4:]))...)
					continue
				}
			}

			if s[0] == '-' {
				_, err := strconv.ParseFloat(s[1:], 64)

				if err == nil {
					t = append(t, Token{DFLOAT, s[1:]})
				} else {
					ss := s[1:]
					if ss == "low" || ss == "high" {
						t = append(t, Token{SOURCE, ss})
					} else {
						t = append(t, Token{FLAG, ss})
					}
				}
				continue
			}

			if s[0] == 'u' && len(s) > 1 {
				_, err := strconv.Atoi(s[1:])
				if err == nil {
					t = append(t, Token{UFLOAT, s[1:]})
				} else {
					t = append(t, Token{VARIABLE, s})
				}
				continue
			}

			if s[last] == '%' {
				_, err := strconv.ParseFloat(s[:last], 64)
				if err != nil {
					return t, err
				}
				t = append(t, Token{PERCENT, s[:len(s)-1]})
				continue
			}

			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				t = append(t, Token{FLOAT, s})
				continue
			}

			t = append(t, lexVariable([]byte(s))...)
		}
	}
	return t, nil
}

//lexVariable parses Functions
func lexVariable(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{LBRACKET, ""})
			temp = []byte("")
		case ')':
			l, _ := Lexer(string(temp))
			tk = append(tk, l...)
			tk = append(tk, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			temp = append(temp, ' ')
		default:
			temp = append(temp, v)
		}
	}

	if len(temp) != 0 {
		tk = append(tk, Token{VARIABLE, string(temp)})
	}
	return tk
}

func lexFunc(s []byte) []Token {
	var temp []byte
	var tk []Token

	for _, v := range s {
		switch v {
		case '(':
			tk = append(tk, Token{LBRACKET, ""})
		case ')':
			tk = append(tk, Token{VARIABLE, string(temp)}, Token{RBRACKET, ""})
			temp = []byte("")
		case ',':
			tk = append(tk, Token{VARIABLE, string(temp)})
			temp = []byte("")
		default:
			temp = append(temp, v)
		}
	}

	return tk
}

func (t TokenType) String() string {
	var s string
	switch t {
	case SIDE:
		s = "side"
	case VARIABLE:
		s = "var"
	case STOP:
		s = "stop"
	case FLOAT:
		s = "float"
	case UFLOAT:
		s = "ufloat"
	case DFLOAT:
		s = "dfloat"
	case PERCENT:
		s = "percent"
	case ASSIGN:
		s = "assign"
	case FLAG:
		s = "flag"
	case FUNC:
		s = "func"
	case DURATION:
		s = "duration"
	case LBRACKET:
		s = "("
	case RBRACKET:
		s = ")"
	case SOURCE:
		s = "source"
	case CANCEL:
		s = "cancel"
	case FUNDING:
		s = "funding"
	case POSITION:
		s = "position"
	case LOAD:
		s = "LOAD"
	}

	return s
}
