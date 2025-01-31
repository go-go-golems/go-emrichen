---
Title: Adding Custom Tags
Slug: adding-custom-tags
Topics:
  - development
  - customization
Commands:
  - help
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: DevelopmentTopic
---

# Adding Custom Tags to Go-Emrichen

Go-Emrichen's functionality can be extended through custom tags, allowing you to add domain-specific templating capabilities. This guide explains how to implement custom tags and integrate them with the Go-Emrichen interpreter.

## Understanding Tag Processing

At its core, Go-Emrichen treats YAML as a tree of nodes that can be transformed through tag handlers. Each tag represents a specific transformation operation on these nodes. When the interpreter encounters a tag, it delegates the processing to the corresponding handler function, which can:

- Transform the node's value
- Create new nodes
- Access and modify variables in the environment
- Process child nodes recursively

This node-based approach provides several advantages:
- Type safety through YAML's typing system
- Preservation of YAML structure during transformations
- Clear separation between parsing and processing
- Ability to handle complex nested structures

## Basic Tag Implementation

Custom tags in Go-Emrichen are implemented as functions that process YAML nodes. The basic signature for a tag handler is:

```go
import (
    "github.com/go-go-golems/go-emrichen/pkg/emrichen"
    "gopkg.in/yaml.v3"
)

type TagFunc func(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error)
```

This signature provides access to both:
- The interpreter instance (`ei`) for environment access and node processing
- The YAML node to be processed (`node`)

Here's a simple example of a custom tag that converts text to uppercase:

```go
import (
    "strings"
    "github.com/go-go-golems/go-emrichen/pkg/emrichen"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v3"
)

func handleUppercase(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    if node.Kind != yaml.ScalarNode {
        return nil, errors.New("!Uppercase requires a scalar node")
    }
    return &yaml.Node{
        Kind:  yaml.ScalarNode,
        Value: strings.ToUpper(node.Value),
    }, nil
}
```

### Registering Custom Tags

Custom tags are registered when creating a new interpreter. The registration process maps tag names to their handler functions, allowing the interpreter to look up the appropriate handler when processing YAML:

```go
customTags := map[string]emrichen.TagFunc{
    "!Uppercase": handleUppercase,
}

interpreter, err := emrichen.NewInterpreter(
    emrichen.WithAdditionalTags(customTags),
)
```

Note that tag handlers are pure functions - they don't modify the interpreter state directly. Instead, they receive the interpreter instance as an argument, which provides access to the environment and utility functions.

## Working with Arguments

Many tags need to handle complex arguments. Go-Emrichen provides a structured approach to argument handling that ensures:
- Clear validation of required arguments
- Type safety for argument values
- Support for optional arguments with defaults
- Variable expansion when needed
- Consistent error handling
- Recursive processing of nested structures

This structured approach helps prevent runtime errors and provides clear error messages when arguments are missing or invalid.

### Core Principles

When implementing tag handlers, follow these core principles:
1. **Always Use ParseArgs**: Use the `ParseArgs` method to handle arguments instead of processing them manually. This ensures consistent behavior and proper error handling.
2. **Process Recursively**: All sub-fields should be processed recursively using the interpreter's `Process` method to ensure proper handling of nested tags and variables.
3. **Handle Expansion**: Use the `Expand` flag in `ParsedVariable` to control when variable substitution occurs.

### Argument Definition

Arguments are defined using the `ParsedVariable` struct, which encapsulates the validation rules for each argument:

```go
type ParsedVariable struct {
    Name     string  // The argument name in the YAML
    Required bool    // Whether the argument must be provided
    Expand   bool    // Whether to process variables in the argument value
}
```

The `Expand` field is particularly important as it determines whether the argument value should be processed for variable substitution before being used. This allows for dynamic argument values while maintaining control over when expansion occurs.

### Using ParseArgs

The `ParseArgs` method is a core utility function that processes YAML mapping nodes according to a list of variable specifications. It provides:

1. **Consistent Argument Handling**:
   - Validates the input node is a mapping node
   - Checks for unknown argument keys
   - Ensures required arguments are present
   - Verifies keys are scalar values

2. **Variable Expansion Control**:
   - Processes values through the interpreter when `Expand` is true
   - Preserves raw values when `Expand` is false
   - Handles nested structures properly

3. **Type Safety**:
   - Maintains YAML node types
   - Provides clear error messages for type mismatches
   - Preserves node metadata

Here's a complete example of proper argument handling:

```go
func handleCustomTag(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    args, err := ei.ParseArgs(node, []emrichen.ParsedVariable{
        {Name: "input", Required: true, Expand: true},    // Will be processed for variables
        {Name: "format", Required: false},                // Optional, no expansion
        {Name: "options", Expand: true},                  // Optional with expansion
    })
    if err != nil {
        return nil, err
    }
    
    // Access and validate arguments
    inputNode := args["input"]
    if inputNode.Kind != yaml.ScalarNode {
        return nil, errors.New("input must be a scalar value")
    }

    // Optional argument handling
    var format string
    if formatNode, ok := args["format"]; ok {
        formatStr, ok := emrichen.NodeToString(formatNode)
        if !ok {
            return nil, errors.New("format must be a string")
        }
        format = formatStr
    }

    // Process nested structures
    if optionsNode, ok := args["options"]; ok {
        // Options were already expanded by ParseArgs if found
        // Process them according to your tag's logic
        result, err := processOptions(ei, optionsNode)
        if err != nil {
            return nil, errors.Wrap(err, "failed to process options")
        }
        // Use the processed result
    }

    // Return processed result
    return processedNode, nil
}
```

### Recursive Processing

When working with nested structures, always ensure proper recursive processing:

```go
func handleNestedStructure(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    args, err := ei.ParseArgs(node, []emrichen.ParsedVariable{
        {Name: "data", Required: true, Expand: true},
        {Name: "template", Required: true},  // Not expanded initially
    })
    if err != nil {
        return nil, err
    }

    // Process the template with the expanded data
    dataNode := args["data"]
    templateNode := args["template"]

    // Create a new scope with the processed data
    result := &yaml.Node{}
    err = ei.env.With(map[string]interface{}{
        "data": emrichen.NodeToInterface(dataNode),
    }, func() error {
        // Now process the template in the new scope
        processed, err := ei.Process(templateNode)
        if err != nil {
            return err
        }
        *result = *processed
        return nil
    })
    
    return result, err
}
```

### Error Handling in Argument Processing

When using ParseArgs, handle errors appropriately:

```go
func handleWithErrorChecking(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    args, err := ei.ParseArgs(node, []emrichen.ParsedVariable{
        {Name: "source", Required: true, Expand: true},
    })
    if err != nil {
        // ParseArgs already provides detailed error messages
        return nil, errors.Wrap(err, "invalid arguments")
    }

    // Additional type-specific validation
    sourceNode := args["source"]
    switch sourceNode.Kind {
    case yaml.ScalarNode:
        // Handle scalar case
    case yaml.SequenceNode:
        // Handle sequence case
    case yaml.MappingNode:
        // Handle mapping case
    default:
        return nil, errors.Errorf("unexpected node kind for source: %v", sourceNode.Kind)
    }

    // Process the node
    return processedNode, nil
}
```

## Environment Interaction

The environment in Go-Emrichen serves as a shared context for variable storage and retrieval. It's designed to support:
- Hierarchical scoping of variables
- Safe concurrent access
- Temporary variable overrides
- Clean variable cleanup

Understanding how to interact with the environment is crucial for creating tags that work well with variables and maintain proper scoping.

### Accessing Variables

Variable access is mediated through the environment to ensure consistency and proper scoping:

```go
func handleVarAccess(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    varName := node.Value
    varValue, exists := ei.env.GetVar(varName)
    if !exists {
        return nil, errors.Errorf("variable %s not found", varName)
    }
    return emrichen.ValueToNode(varValue)
}
```

### Scoped Variables

The `With` method provides a powerful way to create temporary variable scopes. This is essential for:
- Isolating variable changes
- Preventing variable leakage
- Managing cleanup automatically
- Supporting nested scopes

```go
func handleScopedOperation(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    tempVars := map[string]interface{}{
        "tempVar": "value",
    }
    
    var result *yaml.Node
    err := ei.env.With(tempVars, func() error {
        var err error
        result, err = ei.Process(node)
        return err
    })
    return result, err
}
```

## Node Processing

YAML nodes are the fundamental data structure in Go-Emrichen. Understanding how to work with them effectively is crucial for creating robust tags. The node processing utilities provide:
- Type-safe conversions
- Structured node creation
- Error handling
- Performance optimization

### Converting Values

Value conversion is a common operation that needs to be handled carefully to maintain type safety and provide clear error messages:

```go
// Go value to YAML node
value := "example"
node, err := emrichen.ValueToNode(value)

// YAML node to interface{}
if val, ok := emrichen.NodeToInterface(node); ok {
    // Use val
}

// Type-specific conversions
if intVal, ok := emrichen.NodeToInt(node); ok {
    // Use intVal
}
```

### Creating Nodes

Node creation follows specific patterns to ensure consistency and proper YAML structure:

```go
// Scalar node
scalarNode := &yaml.Node{
    Kind:  yaml.ScalarNode,
    Value: "value",
    Tag:   "!!str",
}

// Sequence node
seqNode := &yaml.Node{
    Kind:    yaml.SequenceNode,
    Tag:     "!!seq",
    Content: []*yaml.Node{...},
}

// Mapping node
mapNode := &yaml.Node{
    Kind:    yaml.MappingNode,
    Tag:     "!!map",
    Content: []*yaml.Node{...},
}
```

## Advanced Features

Advanced features in Go-Emrichen build upon the basic concepts to provide powerful functionality. Understanding these features helps create more sophisticated tags.

### Template Processing

Template processing is a recursive operation that allows tags to work with complex nested structures:

```go
func handleTemplate(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    // Process template with current environment
    result, err := ei.Process(templateNode)
    if err != nil {
        return nil, err
    }
    return result, nil
}
```

### Conditional Evaluation

Conditional logic in tags should be implemented carefully to ensure proper evaluation of conditions and handling of both branches:

```go
func handleCondition(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    if isTruthy(conditionNode) {
        return ei.Process(thenNode)
    }
    return ei.Process(elseNode)
}
```

### Error Handling

Proper error handling is crucial for debugging and maintaining templates. Errors should be:
- Descriptive and actionable
- Properly wrapped for context
- Handled at appropriate levels
- Clear about what went wrong

```go
func handleComplexOperation(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    // Validate input
    if node.Kind != yaml.MappingNode {
        return nil, errors.New("expected mapping node")
    }
    
    // Handle nested errors
    result, err := someOperation()
    if err != nil {
        return nil, errors.Wrap(err, "operation failed")
    }
    
    // Return processed result
    return emrichen.ValueToNode(result)
}
```

## Complete Example

The following example demonstrates how these concepts come together in a practical tag implementation. It shows:
- Proper argument handling
- Clear error messages
- Type safety
- Clean code structure

```go
func handlePrefix(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    args, err := ei.ParseArgs(node, []emrichen.ParsedVariable{
        {Name: "text", Required: true},
        {Name: "prefix", Required: true},
    })
    if err != nil {
        return nil, err
    }
    
    // Get text value
    text, ok := emrichen.NodeToString(args["text"])
    if !ok {
        return nil, errors.New("text must be a string")
    }
    
    // Get prefix value
    prefix, ok := emrichen.NodeToString(args["prefix"])
    if !ok {
        return nil, errors.New("prefix must be a string")
    }
    
    // Create result node
    return &yaml.Node{
        Kind:  yaml.ScalarNode,
        Tag:   "!!str",
        Value: prefix + text,
    }, nil
}
```

Usage in YAML:

```yaml
formatted: !Prefix
  text: Hello World
  prefix: ">> "
```

## Best Practices

These best practices are derived from real-world experience with template processing and aim to create robust, maintainable tags:

1. **Input Validation**: Always validate input node types and arguments before processing.
2. **Error Messages**: Provide clear error messages that help users understand what went wrong.
3. **Documentation**: Document your custom tags with examples and parameter descriptions.
4. **Type Safety**: Use type conversion utilities to safely handle different value types.
5. **Scoped Changes**: Use `With` for temporary environment modifications.
6. **Testing**: Write comprehensive tests for your custom tags.

## Testing Custom Tags

Testing is crucial for ensuring tag reliability. A good test suite should cover:
- Basic functionality
- Edge cases
- Error conditions
- Complex inputs
- Performance characteristics

```go
func TestCustomTag(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "basic usage",
            input:    "!CustomTag value",
            expected: "processed value",
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "!CustomTag {invalid: input}",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Parse input YAML
            var node yaml.Node
            err := yaml.Unmarshal([]byte(tt.input), &node)
            require.NoError(t, err)
            
            // Create interpreter with custom tag
            interpreter, err := emrichen.NewInterpreter(
                emrichen.WithAdditionalTags(map[string]emrichen.TagFunc{
                    "!CustomTag": handleCustomTag,
                }),
            )
            require.NoError(t, err)
            
            // Process node
            result, err := interpreter.Process(&node)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result.Value)
        })
    }
}
```

## Common Patterns

These patterns emerge from common use cases and provide tested solutions to recurring problems.

### Handling Collections

Collection handling requires careful attention to:
- Maintaining order
- Proper error propagation
- Memory efficiency
- Type consistency

```go
func handleCollection(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    switch node.Kind {
    case yaml.SequenceNode:
        // Process sequence items
        var results []*yaml.Node
        for _, item := range node.Content {
            result, err := processItem(ei, item)
            if err != nil {
                return nil, err
            }
            results = append(results, result)
        }
        return &yaml.Node{
            Kind:    yaml.SequenceNode,
            Tag:     "!!seq",
            Content: results,
        }, nil
        
    case yaml.MappingNode:
        // Process key-value pairs
        var results []*yaml.Node
        for i := 0; i < len(node.Content); i += 2 {
            key := node.Content[i]
            value := node.Content[i+1]
            // Process key-value pair
        }
        return &yaml.Node{
            Kind:    yaml.MappingNode,
            Tag:     "!!map",
            Content: results,
        }, nil
        
    default:
        return nil, errors.New("expected sequence or mapping node")
    }
}
```

### Variable Expansion

Variable expansion is a common requirement that needs careful handling to:
- Maintain proper scoping
- Handle circular references
- Provide clear error messages
- Support nested expansion

```go
func handleWithExpansion(ei *emrichen.Interpreter, node *yaml.Node) (*yaml.Node, error) {
    // First process the node to expand any variables
    expanded, err := ei.Process(node)
    if err != nil {
        return nil, err
    }
    
    // Then process the expanded result
    return processExpanded(ei, expanded)
}
```

## Conclusion

Custom tags provide a powerful way to extend Go-Emrichen's functionality. By following these patterns and best practices, you can create robust and reusable tags that integrate seamlessly with the templating system.

The key to successful tag implementation lies in:
- Understanding the YAML node structure
- Proper error handling and validation
- Clean and maintainable code
- Comprehensive testing
- Clear documentation

Remember to:
- Validate inputs thoroughly
- Handle errors gracefully
- Document your tags clearly
- Test edge cases
- Consider performance implications

For more examples, look at the built-in tag implementations in the Go-Emrichen source code. 
