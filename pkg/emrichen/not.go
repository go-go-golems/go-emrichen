package emrichen

import "gopkg.in/yaml.v3"

func (ei *Interpreter) handleNot(node *yaml.Node) (*yaml.Node, error) {
	return makeBool(!isTruthy(node)), nil
}
