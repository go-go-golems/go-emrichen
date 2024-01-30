---
Title: "!If Tag"
Slug: tag-if
Short: |
  ```
  !If {test: condition, then: value_if_true, else: value_if_false}
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
# `!If` Tag 

The `!If` tag in Emrichen is used for conditional logic. It evaluates a condition and returns one of two values based on that condition.

```yaml
!If {test: condition, then: value_if_true, else: value_if_false}
```

## Examples

### Basic Usage

A simple condition to check a boolean value:

```yaml
!Defaults
isAdult: true
---
ageCheck: !If
  test: !Var isAdult
  then: 'Adult'
  else: 'Child'
```

**Output:**

```yaml
ageCheck: 'Adult'
```

### Nested Conditions

Using nested `!If` statements for complex logic:

```yaml
!Defaults
age: 20
---
ageCategory: !If
  test: !Op {a: !Var age, op: '>=', b: 18}
  then: !If {test: !Op {a: !Var age, op: '<', b: 65}, then: 'Adult', else: 'Senior'}
  else: 'Child'
```

**Output:**

```yaml
ageCategory: 'Adult'
```

### Using with Other Tags

Combining `!If` with other tags like `!Var` and `!Op`:

```yaml
!Defaults
temperature: 30
---
feeling: !If
  test: !Op {a: !Var temperature, op: '>', b: 25}
  then: 'Hot'
  else: 'Comfortable'
```

**Output:**

```yaml
feeling: 'Hot'
```

## Notes

- The `test` argument must be a condition that evaluates to a boolean value.
- The `then` and `else` arguments define the values to return based on the `test` result. They are optional.