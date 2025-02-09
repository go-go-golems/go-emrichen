package emrichen

import "gopkg.in/yaml.v3"

func (ei *Interpreter) handleIf(node *yaml.Node) (*yaml.Node, error) {
	args, err := ei.ParseArgs(node, []ParsedVariable{
		{Name: "test", Required: true},
		{Name: "then"},
		{Name: "else"},
	})
	if err != nil {
		return nil, err
	}

	testResult, err := ei.Process(args["test"])
	if err != nil {
		return nil, err
	}

	if isTruthy(testResult) {
		if args["then"] == nil {
			return ValueToNode(nil)
		}
		return ei.Process(args["then"])
	} else {
		if args["else"] == nil {
			return ValueToNode(nil)
		}
		return ei.Process(args["else"])
	}
}
