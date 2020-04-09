package script

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ArgParseFunc func(string) (reflect.Value, error)

type ArgParser struct {
	Fun  ArgParseFunc
	Skip bool
}

func ParseInt(s string) (reflect.Value, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(i), nil
}

func ParseString(s string) (reflect.Value, error) {
	return reflect.ValueOf(s), nil
}

func ParseBool(s string) (reflect.Value, error) {
	var b bool
	switch strings.ToLower(s) {
	case "false", "no", "0":
		b = false
	case "true", "yes", "1":
		b = true
	default:
		return reflect.Value{}, errors.New("not bool:" + s)
	}
	return reflect.ValueOf(b), nil
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

func ParseCompare(s string) (reflect.Value, error) {
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
