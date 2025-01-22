# Emrichen: A Go Implementation

Emrichen is a powerful templating engine designed for generating YAML configurations with ease and precision.
This is a heavily LLM-supported implementation of the super-cool [Python implementation](https://github.com/con2/emrichen) (why don't more people use this?).
This Go version brings the same flexibility and robustness to Go developers,
allowing them to dynamically generate configuration files for a wide range of applications, including Kubernetes
deployments, configuration management, and more.

Emrichen stands out by understanding the structure of YAML, enabling users to avoid common pitfalls associated
with text-based templating systems, such as indentation errors and type mismatches. With its rich set of tags, Emrichen
offers a pragmatic and powerful way to template complex configurations.

This go implementation does not support JSON, but it adds go templating functionality to Format,
as well as additional operators. 

Additional tags and template operators can be added programmatically.

You can find detailed documentation for each tag in the [doc section](pkg/doc/examples)
as well as an exhaustive list of examples in [the examples yamls](test-data)
and [in the go unit tests](pkg/emrichen/).

## Overview of Emrichen

### Introduction

Emrichen is a powerful templating engine designed for generating YAML configurations with ease and precision. Built with Go, Emrichen offers flexibility and robustness to developers, enabling the dynamic creation of configuration files for a wide range of applications, including Kubernetes deployments, configuration management, and more.

### Key Features

- **Structured Templating**: Emrichen understands YAML's structure, reducing common pitfalls associated with text-based templating systems.
- **Rich Tag Set**: A comprehensive set of tags allows for complex data manipulation, conditional logic, loops, and more.
- **Variable Management**: Define default values and reuse variables throughout your templates effortlessly.
- **Error Handling and Debugging**: Utilize tags like `!Error` and `!Debug` to manage errors and debug templates effectively.
- **Extensibility**: Additional tags and template operators can be added programmatically, allowing for customization to fit specific needs.

### Concise Examples

#### Example 1: Basic Variable Substitution

```yaml
!Defaults
name: John Doe
---
greeting: !Format "Hello, {name}!"
```

**Output:**

```yaml
greeting: "Hello, John Doe!"
```

#### Example 2: Conditional Logic with `!If`

```yaml
!Defaults
isAdmin: true
---
accessLevel: !If
  test: !Var isAdmin
  then: "Full Access"
  else: "Restricted Access"
```

**Output:**

```yaml
accessLevel: "Full Access"
```

#### Example 3: Looping with `!Loop`

```yaml
!Defaults
ports: [80, 443]
---
containerPorts: !Loop
  over: !Var ports
  template: !Lookup "item"
```

**Output:**

```yaml
containerPorts:
  - 80
  - 443
```

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

This reference aims to provide a quick overview of the capabilities and parameters of each tag supported by Emrichen.
For detailed examples and advanced usage, refer to the specific documentation for each tag.