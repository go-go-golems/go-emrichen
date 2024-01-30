# Emrichen: A Go Implementation

Emrichen is a powerful templating engine designed for generating YAML configurations with ease and precision.
Originating from a Python implementation, this Go version brings the same flexibility and robustness to Go developers,
allowing them to dynamically generate configuration files for a wide range of applications, including Kubernetes
deployments, configuration management, and more.

Emrichen stands out by understanding the structure of YAML, enabling users to avoid common pitfalls associated
with text-based templating systems, such as indentation errors and type mismatches. With its rich set of tags, Emrichen
offers a pragmatic and powerful way to template complex configurations.

This go implementation does not support JSON, but it adds go templating functionality to Format,
as well as additional operators. 

Additional tags and template operators can be added programmatically.

## Kubernetes Deployment Example

Below is an example of a Kubernetes deployment template using Emrichen, showcasing advanced features such as conditional
logic, loops, and variable substitution.

```yaml
!Defaults
app_name: myapp
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
  - name: DEBUG
    value: "false"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: !Format "{app_name}-deployment"
spec:
  replicas: !Var replicas
  selector:
    matchLabels:
      app: !Var app_name
  template:
    metadata:
      labels:
        app: !Var app_name
    spec:
      containers:
        - name: !Var app_name
          image: !Var image
          ports: !Loop
            over: !Var ports
            template:
              containerPort: !Lookup item.port
              protocol: !Lookup item.protocol
          env: !Loop
            over: !Filter
              test: !Op
                a: !Lookup item.name
                op: ne
                b: "DEBUG"  # Assuming we want to filter out the DEBUG environment variable in production
              over: !Var env
            template:
              name: !Lookup item.name
              value: !Lookup item.value
```

This template demonstrates the use of:

- `!Defaults` for setting default values
- `!Var` for variable substitution
- `!Loop` for iterating over lists to dynamically generate port and environment variable configurations
- `!Filter` for filtering out environment variables based on a condition
- `!Format` for formatting the deployment name based on a variable
- `!Lookup` for performing JSONPath lookups

## Emrichen Tags Reference

Emrichen is a powerful template engine designed for generating YAML and JSON configurations. It supports a variety of tags to manipulate and generate structured data dynamically. Below is a reference table of all supported tags along with their parameters and a brief description.

| Tag            | Parameters                                                                 | Description                                                                                   |
|----------------|----------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `!All`         | `iterable`                                                                 | Returns `true` if all items in the iterable are truthy.                                       |
| `!Any`         | `iterable`                                                                 | Returns `true` if at least one item in the iterable is truthy.                                |
| `!Base64`      | `value`                                                                    | Encodes the given value into Base64.                                                          |
| `!Concat`      | `lists`                                                                    | Concatenates lists.                                                                           |
| `!Debug`       | `value`                                                                    | Outputs the value to stderr for debugging purposes.                                           |
| `!Defaults`    | `variables`                                                                | Defines default values for variables.                                                         |
| `!Error`       | `message`                                                                  | Outputs an error message and halts processing if encountered.                                 |
| `!Exists`      | `JSONPath`                                                                 | Checks for the existence of a path or variable.                                               |
| `!Filter`      | `test`, `over`                                                             | Filters elements based on a predicate.                                                        |
| `!Format`      | `format_string`                                                            | Formats a string using variables and expressions.                                             |
| `!Group`       | `over`, `by`, `template` (optional), `result_as` (optional)                 | Groups items based on a key.                                                                  |
| `!If`          | `test`, `then`, `else`                                                     | Conditional logic to return values based on a test.                                          |
| `!Include`     | `path`                                                                     | Includes and processes another template file.                                                 |
| `!IncludeBase64`| `path`                                                                     | Includes a binary file as a Base64-encoded string.                                            |
| `!IncludeBinary`| `path`                                                                     | Includes the contents of a binary file.                                                       |
| `!IncludeGlob` | `patterns`                                                                 | Includes and processes multiple files matching glob patterns.                                 |
| `!IncludeText` | `path`                                                                     | Includes the contents of a text file.                                                         |
| `!Index`       | `over`, `by`, `template` (optional), `duplicates` (optional)               | Creates a dictionary out of a list based on a key.                                            |
| `!IsBoolean`   | `value`                                                                    | Checks if the value is a boolean.                                                             |
| `!IsDict`      | `value`                                                                    | Checks if the value is a dictionary.                                                          |
| `!IsInteger`   | `value`                                                                    | Checks if the value is an integer.                                                            |
| `!IsList`      | `value`                                                                    | Checks if the value is a list.                                                                |
| `!IsNone`      | `value`                                                                    | Checks if the value is `None` (null).                                                         |
| `!IsNumber`    | `value`                                                                    | Checks if the value is a number.                                                              |
| `!IsString`    | `value`                                                                    | Checks if the value is a string.                                                              |
| `!Join`        | `items`, `separator` (optional)                                            | Joins a list of items with a separator.                                                       |
| `!Lookup`      | `JSONPath`                                                                 | Performs a JSONPath lookup.                                                                   |
| `!LookupAll`   | `JSONPath`                                                                 | Performs a JSONPath lookup, returning all matches.                                            |
| `!Loop` | `over`: (Required) Collection to iterate over. <br> `as`: (Optional, default `item`) Variable name for the current element. <br> `index_as`: (Optional) Variable name for the current index or key. <br> `index_start`: (Optional, default `0`) Starting index for the loop. <br> `previous_as`: (Optional) Variable name for the previous element. <br> `template`: (Required) Template applied to each element. <br> `as_documents`: (Optional, default `false`) Treats each iteration's output as a separate YAML document. <br> `filter`: (Optional) Predicate to filter items to loop over. <br> `sort_by`: (Optional) Key or function to sort items before looping. <br> `reverse`: (Optional, default `false`) Reverses the order of items before looping. | Iterates over collections, applying a specified template to each element, with extensive control over the iteration process. |
| `!MD5`         | `data`                                                                     | Hashes the given data using the MD5 algorithm.                                                |
| `!Merge`       | `dicts`                                                                    | Merges dictionaries, with later values overriding earlier ones.                               |
| `!Not`         | `value`                                                                    | Negates a boolean value.                                                                      |
| `!Op`          | `a`, `op`, `b`                                                             | Performs a binary operation between two values.                                               |
| `!SHA1`        | `data`                                                                     | Hashes the given data using the SHA1 algorithm.                                               |
| `!SHA256`      | `data`                                                                     | Hashes the given data using the SHA256 algorithm.                                             |
| `!URLEncode`   | `string` or `url`, `query`                                                 | Encodes a string for URL inclusion or combines a URL with query parameters.                   |
| `!Var`         | `name`                                                                     | Substitutes the value of a variable.                                                          |
| `!Void`        | -                                                                          | Used to remove items from the output.                                                         |
| `!With`        | `vars`, `template`                                                         | Defines a scope with local variables for a template.                                          |

This reference aims to provide a quick overview of the capabilities and parameters of each tag supported by Emrichen. For detailed examples and advanced usage, refer to the specific documentation for each tag.


You can find detailed documentation for each tag in the [doc section](pkg/doc/examples)
as well as an exhaustive list of examples in [the examples yamls](test-data)
and [in the go unit tests](pkg/emrichen/).


Below is a structured approach to documenting each Emrichen tag with a short description and an example. This format is designed to be concise yet informative, providing users with a quick understanding of each tag's functionality and usage.

## `!All`

Evaluates if all items in the iterable are truthy.

**Example:**
```yaml
allTrue: !All [true, true, true]
```

## `!Any`

Evaluates if at least one item in the iterable is truthy.

**Example:**
```yaml
anyTrue: !Any [false, false, true]
```

## `!Base64`

Encodes the given value into Base64.

**Example:**
```yaml
encoded: !Base64 "Hello, World!"
```

## `!Concat`

Concatenates lists.

**Example:**
```yaml
concatenated: !Concat [[1, 2], [3, 4]]
```

## `!Debug`

Outputs the value to stderr for debugging purposes.

**Example:**
```yaml
debugged: !Debug "Debug this value"
```

## `!Defaults`

Defines default values for variables.

**Example:**
```yaml
!Defaults
defaultVar: "Default Value"
---
defaulted: !Var defaultVar
```

## `!Error`

Outputs an error message and halts processing if encountered.

**Example:**
```yaml
!Error "An error occurred"
```

## `!Exists`

Checks for the existence of a path or variable.

**Example:**
```yaml
existsVar: !Exists varName
```

## `!Filter`

Filters elements based on a predicate.

**Example:**
```yaml
filtered: !Filter {test: !Op {a: !Var item, op: ">", b: 5}, over: [4, 5, 6, 7]}
```

## `!Format`

Formats a string using variables and expressions.

**Example:**
```yaml
formatted: !Format "Hello, {name}!"
```

## `!Group`

Groups items based on a key.

**Example:**
```yaml
grouped: !Group {over: [item1, item2], by: !Var key}
```

## `!If`

Conditional logic to return values based on a test.

**Example:**
```yaml
conditional: !If {test: !Var condition, then: "Yes", else: "No"}
```

## `!Include`

Includes and processes another template file.

**Example:**
```yaml
included: !Include "path/to/template.yml"
```

## `!IncludeBase64`

Includes a binary file as a Base64-encoded string.

**Example:**
```yaml
includedBase64: !IncludeBase64 "path/to/file.bin"
```

## `!IncludeBinary`

Includes the contents of a binary file.

**Example:**
```yaml
includedBinary: !IncludeBinary "path/to/file.bin"
```

## `!IncludeGlob`

Includes and processes multiple files matching glob patterns.

**Example:**
```yaml
includedGlob: !IncludeGlob "configs/*.yml"
```

## `!IncludeText`

Includes the contents of a text file.

**Example:**
```yaml
includedText: !IncludeText "path/to/text.txt"
```

## `!Index`

Creates a dictionary out of a list based on a key.

**Example:**
```yaml
indexed: !Index {over: [item1, item2], by: !Var key}
```

## `!IsBoolean`

Checks if the value is a boolean.

**Example:**
```yaml
isBoolean: !IsBoolean true
```

## `!IsDict`

Checks if the value is a dictionary.

**Example:**
```yaml
isDict: !IsDict {key: "value"}
```

## `!IsInteger`

Checks if the value is an integer.

**Example:**
```yaml
isInteger: !IsInteger 42
```

## `!IsList`

Checks if the value is a list.

**Example:**
```yaml
isList: !IsList [1, 2, 3]
```

## `!IsNone`

Checks if the value is `None` (null).

**Example:**
```yaml
isNone: !IsNone null
```

## `!IsNumber`

Checks if the value is a number.

**Example:**
```yaml
isNumber: !IsNumber 3.14
```

## `!IsString`

Checks if the value is a string.

**Example:**
```yaml
isString: !IsString "Hello"
```

## `!Join`

Joins a list of items with a separator.

**Example:**
```yaml
joined: !Join {items: [hello, world], separator: ", "}
```

## `!Lookup`

Performs a JSONPath lookup.

**Example:**
```yaml
lookup: !Lookup "path.to.value"
```

## `!LookupAll`

Performs a JSONPath lookup, returning all matches.

**Example:**
```yaml
lookupAll: !LookupAll "path.to.values[*]"
```

## `!Loop`

Loops over a list or dict, applying a template to each item.

**Example:**
```yaml
looped: !Loop {over: [1, 2, 3], template: !Format "Number: {item}"}
```

## `!MD5`

Hashes the given data using the MD5 algorithm.

**Example:**
```yaml
hashedMD5: !MD5 "data to hash"
```

## `!Merge`

Merges dictionaries, with later values overriding earlier ones.

**Example:**
```yaml
merged: !Merge [{a: 1}, {b: 2}]
```

## `!Not`

Negates a boolean value.

**Example:**
```yaml
negated: !Not true
```

## `!Op`

Performs a binary operation between two values.

**Example:**
```yaml
operation: !Op {a: 5, op: "+", b: 3}
```

## `!SHA1`

Hashes the given data using the SHA1 algorithm.

**Example:**
```yaml
hashedSHA1: !SHA1 "data to hash"
```

## `!SHA256`

Hashes the given data using the SHA256 algorithm.

**Example:**
```yaml
hashedSHA256: !SHA256 "data to hash"
```

## `!URLEncode`

Encodes a string for URL inclusion or combines a URL with query parameters.

**Example:**
```yaml
urlEncoded: !URLEncode "string to encode"
```

## `!Var`

Substitutes the value of a variable.

**Example:**
```yaml
variable: !Var variableName
```

## `!Void`

Used to remove items from the output.

**Example:**
```yaml
voided: !Void
```

## `!With`

Defines a scope with local variables for a template.

**Example:**
```yaml
withScope: !With {vars: {localVar: "Local Value"}, template: !Var localVar}
```
