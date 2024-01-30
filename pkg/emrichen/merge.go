package emrichen

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (ei *Interpreter) handleMerge(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, errors.New("!Merge requires a sequence of mapping nodes")
	}

	mergedMap := make(map[string]*yaml.Node)
	for _, item := range node.Content {
		var err error
		item, err = ei.Process(item)
		if err != nil {
			return nil, err
		}
		// skip !Void nodes
		if item == nil {
			continue
		}
		if item.Kind != yaml.MappingNode {
			return nil, errors.New("!Merge items must be mapping nodes")
		}

		tempMap := make(map[string]*yaml.Node)
		for i := 0; i < len(item.Content); i += 2 {
			keyNode := item.Content[i]
			valueNode := item.Content[i+1]

			tempMap[keyNode.Value] = valueNode
		}

		for k, v := range tempMap {
			mergedMap[k] = v
		}
	}

	mergedContent := []*yaml.Node{}
	for k, v := range mergedMap {
		mergedContent = append(mergedContent, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: k,
		}, v)
	}

	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: mergedContent,
	}, nil

}
