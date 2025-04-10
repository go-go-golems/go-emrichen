# Emrichen: YAML Template Engine in Go

Emrichen is a YAML template engine that provides structured data manipulation and configuration generation capabilities. This is a Go implementation of the [Python Emrichen](https://github.com/con2/emrichen) project, with additional Go-specific features.

## Core Features

- **Structured YAML Processing**: Maintains YAML structure during templating
- **Rich Tag System**: Over 30 built-in tags for data manipulation
- **Go Template Integration**: Adds Go templating functionality to the Format tag
- **Extensible Design**: Support for programmatically adding custom tags and operators
- **Strong Type Safety**: Type-aware operations and validations
- **Comprehensive Testing**: Extensive test suite and examples

## Installation

```bash
go get github.com/go-go-golems/go-emrichen
```

## Quick Examples

### Basic Variable Substitution

```yaml
!Defaults
name: John Doe
---
greeting: !Format "Hello, {name}!"
```

Output:

```yaml
greeting: "Hello, John Doe!"
```

### Conditional Logic

```yaml
!Defaults
isAdmin: true
---
accessLevel: !If
  test: !Var isAdmin
  then: "Full Access"
  else: "Restricted Access"
```

Output:

```yaml
accessLevel: "Full Access"
```

### Data Transformation

```yaml
!Defaults
ports: [80, 443]
---
containerPorts: !Loop
  over: !Var ports
  template:
    port: !Lookup item
    protocol: TCP
```

Output:

```yaml
containerPorts:
  - port: 80
    protocol: TCP
  - port: 443
    protocol: TCP
```

## Core Tags Overview

### Data Access and Manipulation

- `!Var`: Access variables defined in !Defaults or environment
- `!Lookup` / `!LookupAll`: JSONPath-style data access
- `!Format`: String formatting with Go template support
- `!Loop`: Iterate over collections with rich control options
- `!With`: Create local variable scopes

### Logic and Control Flow

- `!If`: Conditional branching
- `!All` / `!Any`: Logical operations on collections
- `!Not`: Boolean negation
- `!Op`: Comparison and arithmetic operations
- `!Filter`: Filter collections based on predicates

### Data Structure Operations

- `!Concat`: Combine multiple lists
- `!Merge`: Deep merge dictionaries
- `!Group`: Group items by key
- `!Index`: Create lookup tables
- `!Join`: Join list elements with separator

### Type Operations

- `!IsBoolean`, `!IsDict`, `!IsInteger`, `!IsList`, `!IsNone`, `!IsNumber`, `!IsString`: Type checking
- `!Base64`: Base64 encoding
- `!URLEncode`: URL encoding
- `!MD5`, `!SHA1`, `!SHA256`: Cryptographic hashing

### File Operations

- `!Include`: Include YAML files
- `!IncludeText`: Include text files
- `!IncludeBinary`: Include binary files
- `!IncludeBase64`: Include files as Base64
- `!IncludeGlob`: Include multiple files using glob patterns

### Debugging and Error Handling

- `!Debug`: Output debug information
- `!Error`: Raise errors with custom messages
- `!Exists`: Check for existence of values
- `!Void`: Explicitly remove values from output

## Advanced Example: Kubernetes Configuration

```yaml
!Defaults
app:
  name: myapp
  image: myapp:latest
  replicas: 3
  ports:
    - port: 80
      protocol: TCP
    - port: 443
      protocol: TCP
  env:
    - name: ENVIRONMENT
      value: production
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: !Format "{app.name}-deployment"
spec:
  replicas: !Var app.replicas
  selector:
    matchLabels:
      app: !Var app.name
  template:
    metadata:
      labels:
        app: !Var app.name
    spec:
      containers:
        - name: !Var app.name
          image: !Var app.image
          ports: !Loop
            over: !Var app.ports
            template:
              containerPort: !Lookup item.port
              protocol: !Lookup item.protocol
          env: !Loop
            over: !Var app.env
            template:
              name: !Lookup item.name
              value: !Lookup item.value
```

## VSCode Integration

Add this to your VSCode settings.json for proper syntax highlighting:

```json
{
  "yaml.customTags": [
    "!All sequence",
    "!Any sequence",
    "!Base64 scalar",
    "!Concat sequence",
    "!Debug scalar",
    "!Debug mapping",
    "!Defaults mapping",
    "!Error scalar",
    "!Exists scalar",
    "!Filter mapping",
    "!Format scalar",
    "!Group mapping",
    "!If mapping",
    "!Include scalar",
    "!IncludeBase64 scalar",
    "!IncludeBinary scalar",
    "!IncludeGlob scalar",
    "!IncludeGlob sequence",
    "!IncludeText scalar",
    "!Index mapping",
    "!IsBoolean scalar",
    "!IsDict scalar",
    "!IsInteger scalar",
    "!IsList scalar",
    "!IsNone scalar",
    "!IsNumber scalar",
    "!IsString scalar",
    "!Join scalar",
    "!Join mapping",
    "!Lookup scalar",
    "!LookupAll scalar",
    "!Loop mapping",
    "!MD5 scalar",
    "!Merge sequence",
    "!Not scalar",
    "!Op mapping",
    "!SHA1 scalar",
    "!SHA256 scalar",
    "!URLEncode scalar",
    "!URLEncode mapping",
    "!Var scalar",
    "!Void",
    "!With mapping"
  ]
}
```

## Documentation

Detailed documentation for each tag and feature is available in:

- [Tag Documentation](pkg/doc/examples/)
- [Emrichen Specification](emrichen-spec.md)
- [Practical Guide](emrichen-in-practice.md)

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests to our GitHub repository.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
