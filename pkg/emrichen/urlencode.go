package emrichen

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (ei *Interpreter) handleURLEncode(node *yaml.Node) (*yaml.Node, error) {
	switch node.Kind {
	case yaml.ScalarNode:
		// Simple string encoding
		return makeString(url.QueryEscape(node.Value)), nil
	case yaml.MappingNode:
		urlStr, queryParams, err := ei.parseURLEncodeArgs(node)
		if err != nil {
			return nil, err
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing URL in !URLEncode")
		}

		query := parsedURL.Query()
		for k, v := range queryParams {
			query.Set(k, fmt.Sprintf("%s", v))
		}
		parsedURL.RawQuery = query.Encode()

		return makeString(parsedURL.String()), nil

	case yaml.DocumentNode, yaml.SequenceNode, yaml.AliasNode:
		return nil, errors.New("!URLEncode requires a scalar or mapping node")
	}

	return nil, errors.New("!URLEncode requires a scalar or mapping node")
}
