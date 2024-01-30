---
Title: "!With Tag"
Slug: tag-with
Short: |
  ```
  !With
  vars: {Variable definitions}
  template: {Template to process}
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
# `!With` Tag

The `!With` tag in Emrichen is designed for defining a scope within which local variables are available. It is
particularly useful for reducing repetition in templates and making complex templates more readable.

The tag accepts two parameters:

- `vars`: A mapping node that defines the local variables. Each key-value pair in this mapping defines a variable name
  and its value.
- `template`: The template to process with the variables defined in `vars`. This template can use the defined variables
  as if they were globally defined.

## Examples

### Basic Usage

Define a local variable and use it within a template.

```yaml
!With
  vars:
    localVar: "Hello, World!"
  template: !Var localVar
```

**Output:**

```yaml
"Hello, World!"
```

### Nested `!With` Blocks

Demonstrate the use of nested `!With` blocks, each with its own local variables.

```yaml
!With
  vars:
    outerVar: "Outer"
  template: !With
    vars:
      innerVar: "Inner"
    template: !Join { items: [!Var outerVar, !Var innerVar], separator: ", " }
```

**Output:**

```yaml
"Outer, Inner"
```

### Interaction with Other Tags

Show how `!With` can be used in conjunction with `!Loop` to iterate over a list of items defined locally.

```yaml
!With
  vars:
    items: [1, 2, 3]
  template: !Loop
    over: !Var items
    as: item
    template: !Format "Item: {item}"
```

**Output:**

```yaml
- "Item: 1"
- "Item: 2"
- "Item: 3"
```
