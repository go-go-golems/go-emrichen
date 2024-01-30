---
Title: "!All and !Any Tags"
Slug: tag-all-any
Short: |
  ```
  !All [Condition1, Condition2, ..., ConditionN]
  !Any [Condition1, Condition2, ..., ConditionN]
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
# `!All` and `!Any` Tags

The `!All` and `!Any` tags in Emrichen are used to evaluate a list of conditions, returning `true` or `false` based on
their collective truthiness. `!All` returns `true` if all conditions are truthy, while `!Any` returns `true` if at least
one condition is truthy.

## Examples

### Basic Usage

Check if all conditions are true:

```yaml
- all_true: !All [true, true, true]
```

**Output:**

```yaml
all_true: true
```

Check if any condition is true:

```yaml
- any_true: !Any [false, false, true]
```

**Output:**

```yaml
any_true: true
```

### Using Variables

Given variables `is_admin: true` and `feature_enabled: false`:

```yaml
- all_conditions: !All [!Var is_admin, !Var feature_enabled]
- any_condition: !Any [!Var is_admin, !Var feature_enabled]
```

**Output:**

```yaml
all_conditions: false
any_condition: true
```

### Nested Lists and Mappings

Check conditions within nested structures:

```yaml
- all_nested: !All [[true, true], {key: true}]
- any_nested: !Any [[false, false], {key: false}]
```

**Output:**

```yaml
all_nested: true
any_nested: false
```
