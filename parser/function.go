package parser

import (
	"errors"
	"strconv"
)

type Func struct {
	Name               string  //Name of the Variable the function is assigned too
	NumberOfParameters int     //Number of Parameters the functions has
	ParameterPosition  []int   //position where the Parameter exist
	FunctionTokens     []Token //Unparsed TokenList, which the parameter get inserted and afterwards parsed
}

func ParseFunc(v interface{}, tl []Token) ([]Token, error) {
	fun, ok := v.(Func)

	if !ok {
		return tl, errors.New("This variable is not a function")
	}
	e, err := fun.Parse(tl)
	if err != nil {
		return tl, err
	}

	return e, nil
}

func (f *Func) Parse(tl []Token) ([]Token, error) {
	if tl[0].Type != LBRACKET {
		return tl, errors.New("Syntax error, no bracket" + tl[0].Text)
	}
	var s []Token
	for _, v := range tl[1:] {
		switch v.Type {
		case RBRACKET:
			break
		default:
			s = append(s, v)
		}
	}

	if len(s) != f.NumberOfParameters {
		return tl, errors.New("To much Parameters, amount has to be" + strconv.Itoa(f.NumberOfParameters))
	}

	for i, t := range s {
		n := f.ParameterPosition[i]
		f.FunctionTokens[n] = t
	}

	return f.FunctionTokens, nil
}
