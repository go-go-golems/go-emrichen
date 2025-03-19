package emrichen

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-go-golems/go-emrichen/pkg/env"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Interpreter struct {
	env            *env.Env
	additionalTags map[string]TagFunc
	funcmaps       []template.FuncMap
}

type InterpreterOption func(*Interpreter) error

func WithVars(vars map[string]interface{}) InterpreterOption {
	return func(ei *Interpreter) error {
		ei.env.Push(vars)
		return nil
	}
}

func WithFuncMap(funcmap ...template.FuncMap) InterpreterOption {
	return func(ei *Interpreter) error {
		ei.funcmaps = append(ei.funcmaps, funcmap...)
		return nil
	}
}

type TagFunc func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error)
type TagFuncMap map[string]TagFunc

var defaultHandlers = TagFuncMap{
	"!Defaults": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind == yaml.MappingNode {
			err := ei.updateVars(node.Content)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	},
	"!All": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleAll(node)
	},
	"!And": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleAll(node)
	},
	"!Any": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleAny(node)
	},
	"!Or": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleAny(node)
	},
	"!Base64": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind != yaml.ScalarNode {
			return nil, errors.New("!Base64 requires a scalar value")
		}
		return makeString(base64.StdEncoding.EncodeToString([]byte(node.Value))), nil
	},
	"!Concat": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleConcat(node)
	},
	"!Debug": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		// need to remove debug tag
		switch node.Kind {
		case yaml.SequenceNode:
			node.Tag = "!!seq"
		case yaml.MappingNode:
			node.Tag = "!!map"
		case yaml.ScalarNode:
			node.Tag = "!!str"
		case yaml.DocumentNode:
			node.Tag = "!!doc"
		case yaml.AliasNode:
			node.Tag = "!!alias"
		}
		v, err := ei.Process(node)
		if err != nil {
			return nil, err
		}
		toInterface, _ := NodeToInterface(v)
		fmt.Printf("DEBUG: %s\n", toInterface)
		return v, nil
	},
	"!Error": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind != yaml.ScalarNode {
			return nil, errors.New("!Error tag requires a scalar value for the error message")
		}
		errorString, err := ei.renderFormatString(node.Value)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorString)
	},
	"!Exists": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleExists(node)
	},
	"!Format": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleFormat(node)
	},
	"!Filter": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleFilter(node)
	},
	"!Group": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleGroup(node)
	},
	"!If": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIf(node)
	},
	"!Include": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleInclude(node)
	},
	"!IncludeBase64": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIncludeBase64(node)
	},
	"!IncludeBinary": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIncludeBinary(node)
	},
	"!IncludeGlob": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIncludeGlob(node)
	},
	"!IncludeText": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIncludeText(node)
	},
	"!Index": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleIndex(node)
	},
	"!IsBoolean": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return makeBool(node.Kind == yaml.ScalarNode && (node.Value == "true" || node.Value == "false")), nil
	},
	"!IsDict": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return makeBool(node.Kind == yaml.MappingNode), nil
	},
	"!IsInteger": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		_, err := strconv.Atoi(node.Value)
		return makeBool(err == nil && node.Kind == yaml.ScalarNode), nil
	},
	"!IsList": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return makeBool(node.Kind == yaml.SequenceNode), nil
	},
	"!IsNone": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return makeBool(node.Tag == "!!null" || node.Value == "null"), nil
	},
	"!IsNumber": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		_, err := strconv.ParseFloat(node.Value, 64)
		return makeBool(err == nil && node.Kind == yaml.ScalarNode), nil
	},
	"!IsString": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return makeBool(node.Kind == yaml.ScalarNode), nil
	},
	"!Join": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleJoin(node)
	},
	"!Loop": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleLoop(node)
	},
	"!Lookup": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleLookup(node)
	},
	"!LookupAll": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleLookupAll(node)
	},
	"!MD5": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind != yaml.ScalarNode {
			return nil, errors.New("!MD5 requires a scalar value")
		}
		hash := md5.Sum([]byte(node.Value))
		return makeString(hex.EncodeToString(hash[:])), nil
	},
	"!Merge": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleMerge(node)
	},
	"!Not": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleNot(node)
	},
	"!Op": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleOp(node)
	},
	"!SHA1": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind != yaml.ScalarNode {
			return nil, errors.New("!SHA1 requires a scalar value")
		}
		hash := sha1.Sum([]byte(node.Value))
		return makeString(hex.EncodeToString(hash[:])), nil
	},
	"!SHA256": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		if node.Kind != yaml.ScalarNode {
			return nil, errors.New("!SHA256 requires a scalar value")
		}
		hash := sha256.Sum256([]byte(node.Value))
		return makeString(hex.EncodeToString(hash[:])), nil
	},
	"!Var": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleVar(node)
	},
	"!URLEncode": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleURLEncode(node)
	},
	"!Void": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return nil, nil
	},
	"!With": func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return ei.handleWith(node)
	},
}

func WithAdditionalTags(tags TagFuncMap) InterpreterOption {
	return func(ei *Interpreter) error {
		for k, v := range tags {
			if _, ok := ei.additionalTags[k]; ok {
				return errors.Errorf("tag %s already exists", k)
			}
			ei.additionalTags[k] = v
		}
		return nil
	}
}

func NewInterpreter(options ...InterpreterOption) (*Interpreter, error) {
	ret := &Interpreter{
		env:            env.NewEnv(),
		additionalTags: map[string]TagFunc{},
	}

	// Copy default handlers
	for k, v := range defaultHandlers {
		ret.additionalTags[k] = v
	}

	for _, option := range options {
		err := option(ret)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

type interpretHelper struct {
	target      interface{}
	interpreter *Interpreter
}

type rawInterpretHelper struct {
	target      *yaml.Node
	interpreter *Interpreter
}

func (ei *interpretHelper) UnmarshalYAML(value *yaml.Node) error {
	resolved, err := ei.interpreter.Process(value)
	if err != nil {
		return err
	}
	if resolved == nil {
		return nil
	}
	return resolved.Decode(ei.target)
}

func (ei *rawInterpretHelper) UnmarshalYAML(value *yaml.Node) error {
	resolved, err := ei.interpreter.Process(value)
	if err != nil {
		return err
	}
	if resolved != nil {
		*ei.target = *resolved
	}
	return nil
}

func (ei *Interpreter) CreateDecoder(target interface{}) *interpretHelper {
	return &interpretHelper{
		target:      target,
		interpreter: ei,
	}
}

func (ei *Interpreter) CreateRawDecoder(target *yaml.Node) *rawInterpretHelper {
	return &rawInterpretHelper{
		target:      target,
		interpreter: ei,
	}
}

func (ei *Interpreter) RegisterTag(tag string, f func(node *yaml.Node) (*yaml.Node, error)) error {
	if _, ok := ei.additionalTags[tag]; ok {
		return errors.Errorf("tag %s already exists", tag)
	}
	// Wrap the old-style function to match the new signature
	ei.additionalTags[tag] = func(ei_ *Interpreter, node *yaml.Node) (*yaml.Node, error) {
		return f(node)
	}
	return nil
}

func (ei *Interpreter) LookupFirst(jsonPath string) (*yaml.Node, error) {
	v, err := ei.env.LookupFirst("$." + jsonPath)
	if err != nil {
		return nil, err
	}
	node, err := ValueToNode(v)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (ei *Interpreter) LookupAll(jsonPath string) (*yaml.Node, error) {
	v, err := ei.env.LookupAll("$."+jsonPath, true)
	if err != nil {
		return nil, err
	}
	node, err := ValueToNode(v)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (ei *Interpreter) Process(node *yaml.Node) (*yaml.Node, error) {
	tag := node.Tag
	ss := strings.Split(tag, ",")
	if len(ss) == 0 {
		return nil, errors.New("custom tag is empty")
	}

	for i, s := range ss[1:] {
		if !strings.HasPrefix(s, "!") {
			ss[i+1] = "!" + s
		}
	}

	// reverse ss
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}

	for _, verb := range ss {
		ret, err := func() (*yaml.Node, error) {
			// we allow overriding our own tags
			if f, ok := ei.additionalTags[verb]; ok {
				return f(ei, node)
			}

			// If no handler is found, process the node based on its kind
			switch node.Kind {
			case yaml.SequenceNode:
				retContent := make([]*yaml.Node, 0)
				for i := range node.Content {
					v, err := ei.Process(node.Content[i])
					if err != nil {
						return nil, err
					}
					if v == nil {
						continue
					}
					retContent = append(retContent, v)
				}
				return &yaml.Node{
					Kind:    yaml.SequenceNode,
					Content: retContent,
					Tag:     "!!seq",
				}, nil
			case yaml.MappingNode:
				retContent := make([]*yaml.Node, 0)
				for i := 0; i < len(node.Content); i += 2 {
					key := node.Content[i]
					value := node.Content[i+1]

					v, err := ei.Process(value)
					if err != nil {
						return nil, err
					}
					if v == nil {
						continue
					}
					retContent = append(retContent, key, v)
				}
				return &yaml.Node{
					Kind:    yaml.MappingNode,
					Content: retContent,
					Tag:     "!!map",
				}, nil
			case yaml.ScalarNode:
				return node, nil
			case yaml.AliasNode:
				return nil, errors.New("alias nodes are not supported")
			case yaml.DocumentNode:
				if len(node.Content) == 1 {
					return ei.Process(node.Content[0])
				}
				return ei.Process(node.Content[0])
			default:
				return nil, errors.Errorf("unknown node kind: %v", node.Kind)
			}
		}()

		if err != nil {
			return nil, err
		}

		node = ret
	}

	return node, nil
}

func (ei *Interpreter) updateVars(content []*yaml.Node) error {
	name := ""
	vars := map[string]interface{}{}
	for i := range content {
		if i%2 == 0 {
			name = content[i].Value
			continue
		}
		node, err := ei.Process(content[i])
		if err != nil {
			return err
		}
		v, ok := NodeToInterface(node)
		if !ok {
			return errors.New("could not get node value")
		}
		vars[name] = v
	}

	ei.env.Push(vars)

	return nil
}
