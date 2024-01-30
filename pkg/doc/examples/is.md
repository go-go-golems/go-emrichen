---
Title: "Is* Tags for Type Checking"
Slug: is-tags
Short: |
  ```
  !IsBoolean, !IsDict, !IsInteger, !IsList, !IsNone, !IsNumber, !IsString
  Data to typecheck.
  ```
Command:
  - emrichen
Topics:
  - tags
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Example
---
# `Is*` Tags for Type Checking

The `Is*` tags in Emrichen are a set of utilities designed for type checking within templates. These tags allow you to
verify the type of a given value, returning `True` if the value matches the expected type, and `False` otherwise. Each
tag accepts a single parameter:

- `Data to typecheck`: The value you wish to check the type of.

Supported tags include `!IsBoolean`, `!IsDict`, `!IsInteger`, `!IsList`, `!IsNone`, `!IsNumber`, and `!IsString`. These
tags are useful for conditional rendering based on the type of data or for validating data structures.

## Examples

### Checking if a Value is a Boolean

```yaml
isTrue: !IsBoolean true
isString: !IsBoolean "true"
```

**Output:**

```yaml
isTrue: true
isString: false
```

### Checking if a Value is a List

```yaml
isList: !IsList [1, 2, 3]
isDict: !IsList {key: "value"}
```

**Output:**

```yaml
isList: true
isDict: false
```

### Composite Type Check

Checking if a value is not `None` and then checking if it's a string.

```yaml
value: "Hello, World!"
isNotNullAndIsString: !Not,IsNone,Var value
isString: !IsString,Var value
```

**Output:**

```yaml
value: "Hello, World!"
isNotNullAndIsString: true
isString: true
```
