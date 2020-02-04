package script

import (
	"errors"
	"reflect"
	"strconv"
)

func parseInt(s string) (reflect.Value, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(i), nil
}

func parseString(s string) (reflect.Value, error) {
	return reflect.ValueOf(s), nil
}

type CompareOp int

const (
	GT  CompareOp = iota // >
	GTE                  // >=
	LT                   // <
	LTE                  // <=
	EQ                   // ==
	NEQ                  // !=
)

func parseCompare(s string) (reflect.Value, error) {
	var op CompareOp
	switch s {
	case "<":
		op = LT
	case ">":
		op = GT
	case ">=":
		op = GTE
	case "<=":
		op = LTE
	case "!=":
		op = NEQ
	case "==":
		op = EQ
	default:
		return reflect.Value{}, errors.New("not support:" + s)
	}

	return reflect.ValueOf(op), nil
}

func CompareInt(op CompareOp, a, b int) bool {
	switch op {
	case GT:
		return a > b
	case GTE:
		return a >= b
	case LT:
		return a < b
	case LTE:
		return a <= b
	case EQ:
		return a == b
	case NEQ:
		return a != b
	}
	return false
}

type ArgParseFunc func(string) (reflect.Value, error)
