---
Title: "!Var Tag"
Slug: tag-var
Short: |
  ```
  !Var variable_name
  ```
Command:
  - emrichen
Topics:
  - tags
IsTemplate: true
IsTopLevel: true
ShowPerDefault: false
SectionType: Example
---

# `!Var` Tag Documentation

The `!Var` tag allows the user to define variables that can be reused throughout the template.
This tag dynamically inserts the value of a variable wherever it is referenced.

## Syntax

```yaml
!Var variable_name
```

## Examples

### Basic Variable Substitution

Substitute a simple string variable:

```yaml
!Defaults
greeting: Hello, World!
---
message: !Var greeting
```

**Output:**

```yaml
message: Hello, World!
```

### Substitution with Nested Structures

Use `!Var` to insert complex data structures like lists or maps:

```yaml
!Defaults
userInfo:
  name: Jane Doe
  age: 28
---
user: !Var userInfo
```

**Output:**

```yaml
user:
  name: Jane Doe
  age: 28
```

### Conditional Substitution

Combine `!Var` with `!If` for conditional variable substitution:

```yaml
!Defaults
environment: production
---
path: !If
  test: !Var environment
  then: /var/www/production
  else: /var/www/staging
```

**Output (if `environment` is `production`):**

```yaml
path: /var/www/production
```

## Notes

- The `!Var` tag is essential for creating dynamic and reusable templates.
- It can be used with any data type supported by YAML, including strings, numbers, lists, and maps.