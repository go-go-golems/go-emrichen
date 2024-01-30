package emrichen

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"reflect"
	"regexp"
	"strings"
)

func (ei *Interpreter) handleOp(node *yaml.Node) (*yaml.Node, error) {
	args, err := ei.parseArgs(node, []parsedVariable{
		{Name: "op", Required: true, Expand: true},
		{Name: "a", Required: true, Expand: true},
		{Name: "b", Required: true, Expand: true},
	})
	if err != nil {
		return nil, err
	}

	opNode := args["op"]
	if opNode.Kind != yaml.ScalarNode {
		return nil, errors.New("!Op 'op' argument must be a scalar")
	}

	aProcessed, bProcessed := args["a"], args["b"]

	isNumberOperation := false
	switch opNode.Value {
	case "+", "plus", "add", "-", "minus", "sub", "subtract", "*", "×", "mul", "times", "/", "÷", "div", "divide", "truediv", "//", "floordiv",
		"<", "lt", ">", "gt", "<=", "le", "lte", ">=", "ge", "gte", "%", "mod", "modulo":
		isNumberOperation = true
	default:
	}

	isStringOperation := false
	switch opNode.Value {
	case "contains", "startswith", "endswith", "matches":
		isStringOperation = true
	default:
	}

	a, ok := NodeToFloat(aProcessed)
	if isNumberOperation && !ok {
		return nil, errors.New("could not convert first argument to float")
	}
	b, ok := NodeToFloat(bProcessed)
	if isNumberOperation && !ok {
		return nil, errors.New("could not convert second argument to float")
	}

	bothInts := aProcessed.Tag == "!!int" && bProcessed.Tag == "!!int"

	// Handle different operators
	switch opNode.Value {
	case "=", "==", "===":
		if isNumberOperation {
			return makeBool(a == b), nil
		}
		aVal, ok := NodeToInterface(aProcessed)
		if !ok {
			return nil, errors.Errorf("could not convert first argument to interface: %v", aProcessed)
		}
		bVal, ok := NodeToInterface(bProcessed)
		if !ok {
			return nil, errors.Errorf("could not convert second argument to interface: %v", bProcessed)
		}
		return makeBool(reflect.DeepEqual(aVal, bVal)), nil
	case "≠", "!=", "!==", "ne":
		if isNumberOperation {
			return makeBool(a != b), nil
		}
		aVal, ok := NodeToInterface(aProcessed)
		if !ok {
			return nil, errors.Errorf("could not convert first argument to interface: %v", aProcessed)
		}
		bVal, ok := NodeToInterface(bProcessed)
		if !ok {
			return nil, errors.Errorf("could not convert second argument to interface: %v", bProcessed)
		}
		return makeBool(!reflect.DeepEqual(aVal, bVal)), nil

	// Less than, Greater than, Less than or equal to, Greater than or equal to
	case "<", "lt":
		return makeBool(a < b), nil
	case ">", "gt":
		return makeBool(a > b), nil
	case "<=", "le", "lte":
		return makeBool(a <= b), nil
	case ">=", "ge", "gte":
		return makeBool(a >= b), nil

	// Arithmetic operations
	case "+", "plus", "add":
		if bothInts {
			return makeInt(int(a) + int(b)), nil
		}
		return makeFloat(a + b), nil

	case "-", "minus", "sub", "subtract":
		if bothInts {
			return makeInt(int(a) - int(b)), nil
		}
		return makeFloat(a - b), nil

	case "*", "×", "mul", "times":
		if bothInts {
			return makeInt(int(a) * int(b)), nil
		}
		return makeFloat(a * b), nil
	case "/", "÷", "div", "divide", "truediv":
		result := a / b
		if bothInts && result == float64(int(result)) {
			return makeInt(int(result)), nil
		}
		return makeFloat(result), nil
	case "//", "floordiv":
		return makeInt(int(a) / int(b)), nil

	case "%", "mod", "modulo":
		return makeInt(int(a) % int(b)), nil

	case "contains":
		if !isStringOperation {
			return nil, errors.New("contains operator requires string arguments")
		}
		return makeBool(strings.Contains(aProcessed.Value, bProcessed.Value)), nil

	case "startswith":
		if !isStringOperation {
			return nil, errors.New("startswith operator requires string arguments")
		}
		return makeBool(strings.HasPrefix(aProcessed.Value, bProcessed.Value)), nil

	case "endswith":
		if !isStringOperation {
			return nil, errors.New("endswith operator requires string arguments")
		}

		return makeBool(strings.HasSuffix(aProcessed.Value, bProcessed.Value)), nil

	case "matches":
		if !isStringOperation {
			return nil, errors.New("matches operator requires string arguments")
		}
		// do regexp match check
		r, err := regexp.Compile(bProcessed.Value)
		if err != nil {
			return nil, errors.New("invalid regexp")
		}
		return makeBool(r.MatchString(aProcessed.Value)), nil

	// Membership tests
	case "in", "∈":
		r, err := ei.opIn(bProcessed, aProcessed)
		if err != nil {
			return nil, err
		}
		return makeBool(r), nil

	case "not in", "∉":
		r, err := ei.opIn(bProcessed, aProcessed)
		if err != nil {
			return nil, err
		}
		return makeBool(!r), nil

	default:
		return nil, fmt.Errorf("unsupported operator: %s", opNode.Value)
	}
}

func (ei *Interpreter) opIn(bProcessed *yaml.Node, aProcessed *yaml.Node) (bool, error) {
	// check that b is a sequence
	if bProcessed.Kind != yaml.SequenceNode {
		return false, errors.New("in operator requires a sequence as second argument")
	}
	a, ok := NodeToInterface(aProcessed)
	if !ok {
		return false, errors.New("could not convert first argument to interface")
	}
	// check that a is in b
	for _, item := range bProcessed.Content {
		b, ok := NodeToInterface(item)
		if !ok {
			return false, errors.New("could not convert second argument to interface")
		}
		if reflect.DeepEqual(a, b) {
			return true, nil
		}
	}
	return false, nil
}
