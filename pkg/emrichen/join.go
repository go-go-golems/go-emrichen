package emrichen

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (ei *Interpreter) handleJoin(node *yaml.Node) (*yaml.Node, error) {
	separator := " " // Default separator
	itemsNode := node

	switch node.Kind {
	case yaml.MappingNode:
		args, err := ei.ParseArgs(node, []ParsedVariable{
			// don't expand items here yet
			{Name: "items", Required: true, Expand: true},
			{Name: "separator", Expand: true},
		})
		if err != nil {
			return nil, err
		}

		var ok bool
		itemsNode, ok = args["items"]
		if !ok || itemsNode.Kind != yaml.SequenceNode {
			return nil, errors.New("!Join requires a sequence node for 'items'")
		}

		if sepNode, ok := args["separator"]; ok && sepNode.Kind == yaml.ScalarNode {
			separator = sepNode.Value
		}
	case yaml.SequenceNode:
		itemsNode = &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: node.Content,
			Tag:     "!!seq",
		}
		var err error
		itemsNode, err = ei.Process(itemsNode)
		if err != nil {
			return nil, err
		}
	case yaml.DocumentNode, yaml.ScalarNode, yaml.AliasNode:
		// Or handle other kinds if necessary, otherwise this case might not be needed
		// if only MappingNode and SequenceNode are valid inputs for !Join
		return nil, errors.Errorf("!Join expects a Mapping or Sequence node, got %v", node.Kind)
	}

	var items []string
	for _, itemNode := range itemsNode.Content {
		if itemNode.Kind != yaml.ScalarNode {
			return nil, errors.New("!Join items must be scalar values")
		}
		if itemNode.Tag == "!!null" {
			continue
		}
		items = append(items, itemNode.Value)
	}

	joinedStr := strings.Join(items, separator)

	return ValueToNode(joinedStr)
}
