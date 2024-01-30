---
Title: "!Concat Tag"
Slug: tag-concat
Short: |
  ```
  !Concat [List1, List2, ..., ListN]
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
# `!Concat` Tag 

The `!Concat` tag in Emrichen is used to concatenate lists.
It takes a list of lists as its argument and returns a single list containing all the elements of the input lists,
in the order they were provided.

```yaml
!Concat [List1, List2, ..., ListN]
```

## Examples

### Basic Concatenation

Concatenate two simple lists:

```yaml
- !Concat [[1, 2, 3], [4, 5]]
```

**Output:**

```yaml
[1, 2, 3, 4, 5]
```

### Concatenation with Variables

Concatenate lists stored in variables:

```yaml
!Defaults
list1: [1, 2, 3]
list2: [4, 5]
---
concatenated: !Concat [!Var list1, !Var list2]
```

**Output:**

```yaml
concatenated: [1, 2, 3, 4, 5]
```

### Concatenation with Different Types

Concatenate lists containing different types of elements:

```yaml
- !Concat [[1, 'hello'], [true, 3.14]]
```

**Output:**

```yaml
[1, 'hello', true, 3.14]
```

## Notes

- The `!Concat` tag requires that its arguments are lists. Concatenating non-list items will result in an error.
- It's possible to concatenate lists containing different types of elements, including strings, numbers, booleans, etc.