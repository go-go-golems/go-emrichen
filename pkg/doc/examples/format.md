---
Title: "!Format Tag"
Slug: tag-format
Short: |
  ```
  !Format "format string"
  ```
Command:
  - emrichen
Topics:
  - tags
IsTemplate: false
IsTopLevel: true
ShowPerDefault: false
SectionType: Example
---
# `!Format` Tag

The `!Format` tag in Emrichen is used for string formatting. It interpolates variables into a format string, supporting
both basic string formatting and Go template syntax, including JSONPath lookup.

```yaml
!Format "format string with {jsonPath.expression} and {{.GoSyntaxVariable}}"
```

## Examples

### Basic String Formatting

Interpolate variables into a basic format string:

```yaml
!Defaults
name: "John"
greeting: !Format "Hello, {name}!"
```

**Output:**

```yaml
greeting: "Hello, John!"
```

### Go Template Syntax

Using Go template syntax for more complex interpolations (using the sprig functions):

```yaml
!Defaults
user: 
  name: "Jane"
  age: 25
message: !Format "Hello, {{.user.name}}! You will be {{add (lookup user.age) 1}} next year."
```

**Output:**

```yaml
message: "Hello, Jane! You will be 26 next year."
```

### JSONPath Lookup in Format String

Embedding JSONPath lookup in the format string:

```yaml
!Defaults
users:
  - name: "Alice"
    age: 30
  - name: "Bob"
    age: 35
userMessage: !Format "User: {users[0].name} is {users[0].age} years old."
```

**Output:**

```yaml
userMessage: "User: Alice is 30 years old."
```

## Notes

- Basic string formatting uses `{}` to enclose variables.
- Go template syntax is enclosed in `{{}}` and supports functions like `add`.
- JSONPath lookup can be used for accessing nested data structures in the basic string formatting. 
- Use the `lookup` function when using the go template syntax.
- The go templates provide the sprig template functions.
