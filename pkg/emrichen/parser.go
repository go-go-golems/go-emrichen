package emrichen

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ParsedVariable struct {
	Name     string
	Expand   bool
	Required bool
}

// ParseArgs processes a YAML mapping node according to a list of variable specifications.
// It's a core utility function used by various Emrichen tags to parse their arguments
// in a consistent way.
//
// Parameters:
// - node: A pointer to a yaml.Node that must be a mapping node containing key-value pairs
// - variables: A slice of parsedVariable structs that specify:
//   - Name: The expected argument name
//   - Required: Whether the argument must be present
//   - Expand: Whether to process the value through the Emrichen interpreter
//
// Returns:
// - map[string]*yaml.Node: A map of processed arguments where:
//   - Keys are the argument names
//   - Values are the processed YAML nodes (expanded if specified)
//
// - error: Returns an error if:
//   - The input node is not a mapping node
//   - An unknown argument key is encountered
//   - A required argument is missing
//   - A key is not a scalar value
//   - Value expansion fails
//
// Example usage:
//
//	args, err := ei.ParseArgs(node, []parsedVariable{
//	  {Name: "test", Required: true, Expand: true},
//	  {Name: "then", Required: true, Expand: false},
//	  {Name: "else", Required: false, Expand: false},
//	})
func (ei *Interpreter) ParseArgs(
	node *yaml.Node,
	variables []ParsedVariable,
) (map[string]*yaml.Node, error) {
	argsMap := make(map[string]*yaml.Node)
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("expected a mapping node")
	}

	varMap := make(map[string]ParsedVariable)
	for _, v := range variables {
		varMap[v.Name] = v
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		parsedVar, ok := varMap[keyNode.Value]
		if !ok {
			return nil, errors.Errorf("unknown key '%s'", keyNode.Value)
		}
		key, ok := NodeToString(keyNode)
		if !ok {
			return nil, errors.Errorf("expected scalar key '%s'", keyNode.Value)
		}

		if parsedVar.Expand {
			value, err := ei.Process(valueNode)
			if err != nil {
				return nil, err
			}
			valueNode = value
		}
		argsMap[key] = valueNode
	}

	for _, v := range variables {
		if v.Required {
			if _, ok := argsMap[v.Name]; !ok {
				return nil, errors.Errorf("required key '%s' not found", v.Name)
			}
		}
	}

	return argsMap, nil
}

// parseURLEncodeArgs extracts 'url' and 'query' parameters from a YAML node and organizes them suitably for URL encoding.
// This function is specifically tailored for extracting URL and query parameters for URL encoding purposes.
//
// Parameters:
// - node: A pointer to a yaml.Node that should contain the 'url' and 'query' parameters in a mapping structure.
//
// Returns:
// - A string representing the URL extracted from the node.
// - A map[string]string where keys and values represent query parameters.
// - An error if the node doesn't contain the necessary structure or required keys ('url', and optionally 'query').
//
// Note: The 'query' parameter is optional and can be a mapping node containing key-value pairs of query parameters.
func (ei *Interpreter) parseURLEncodeArgs(node *yaml.Node) (string, map[string]interface{}, error) {
	args, err := ei.ParseArgs(node, []ParsedVariable{
		{Name: "url", Required: true, Expand: true},
		{Name: "query", Expand: true},
	})
	if err != nil {
		return "", nil, err
	}

	urlStr, ok := NodeToString(args["url"])
	if !ok {
		return "", nil, errors.New("url must be a string")
	}

	// TODO need to process node
	queryParams := make(map[string]interface{})
	if queryNode, ok := args["query"]; ok && queryNode.Kind == yaml.MappingNode {
		for i := 0; i < len(queryNode.Content); i += 2 {
			paramKey := queryNode.Content[i].Value
			param, err := ei.Process(queryNode.Content[i+1])
			if err != nil {
				return "", nil, err
			}
			paramValue, ok := NodeToScalarInterface(param)
			if !ok {
				return "", nil, errors.New("query parameter value must be a scalar")
			}
			queryParams[paramKey] = paramValue
		}
	}

	return urlStr, queryParams, nil
}
