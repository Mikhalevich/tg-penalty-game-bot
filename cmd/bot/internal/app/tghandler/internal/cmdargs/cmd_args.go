package cmdargs

import (
	"strings"
)

type Operation int

const (
	OperationNope Operation = iota + 1
	OperationJoin
)

const (
	joinPrefix = "join_"
)

// Parse parse command arguments and return operation and raw arguments(withoud operation prefix).
func Parse(args string) (Operation, string) {
	if args == "" {
		return OperationNope, args
	}

	if after, ok := strings.CutPrefix(args, joinPrefix); ok {
		return OperationJoin, after
	}

	return OperationNope, args
}
