package emrichen

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (ei *Interpreter) handleVar(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind == yaml.ScalarNode {
		varName := node.Value
		varValue, ok := ei.env.GetVar(varName)
		if !ok {
			return nil, errors.Errorf("variable %s not found", varName)
		}
		v, err := ValueToNode(varValue)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, errors.New("variable definition must be !Var variable name")
}
