---
Title: "!Merge Tag"
Slug: tag-merge
Short: |
  ```
  !Merge [Dict1, Dict2, ..., DictN]
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
# `!Merge` Tag

The `!Merge` tag in Emrichen is used to merge multiple dictionaries into a single dictionary. For overlapping keys, the
value from the last dictionary takes precedence.

## Examples

### Basic Merge

Merge two dictionaries with distinct keys:

```yaml
- !Merge
  - a: 1
  - b: 2
```

**Output:**

```yaml
a: 1
b: 2
```

### Overlapping Keys

Merge dictionaries with overlapping keys:

```yaml
- !Merge
  - a: 1
  - a: 2
  - b: 3
```

**Output:**

```yaml
a: 2
b: 3
```

### Merge with Tag Composition

Merge dictionaries using variables and tag composition:

```yaml
!Defaults
dict1: {a: 1}
dict2: {a: 2, b: !Not,Var truthy}
truthy: true
---
- !Merge
  - !Var dict1
  - !Var dict2
```

**Output:**

```yaml
a: 2
b: false
```

## Notes

- The `!Merge` tag requires that its arguments are dictionaries. Merging non-dictionary items will result in an error.