---
Title: "!Exists Tag"
Slug: tag-exists
Short: |
  ```
  !Exists JSONPath
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
# `!Exists` Tag

The `!Exists` tag in Emrichen is used to check the existence of a variable or path in the given context. It takes a
JSONPath expression as its argument and returns `true` if the JSONPath expression returns one or more matches, `false`
otherwise.

```yaml
!Exists JSONPath
```

## Examples

### Basic Existence Check

Check if a top-level variable exists:

```yaml
!Defaults
foo: bar
---
existsFoo: !Exists foo
```

**Output:**

```yaml
existsFoo: true
```

### Nested Structure Existence

Check if a nested key exists within a structure:

```yaml
!Defaults
nested:
  key: value
---
existsNestedKey: !Exists nested.key
```

**Output:**

```yaml
existsNestedKey: true
```

### Existence in Array

Check for the existence of an element at a specific index in an array:

```yaml
!Defaults
list:
  - item1
  - item2
---
existsListItem: !Exists list[0]
```

**Output:**

```yaml
existsListItem: true
```