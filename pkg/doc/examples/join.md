---
Title: "!Join Tag"
Slug: tag-join
Short: |
  ```
  !Join [Item1, Item2, ..., ItemN]
  ```
  or
  ```
  !Join { items: [Item1, Item2, ..., ItemN], separator: "Separator" }
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
# `!Join` Tag

The `!Join` tag in Emrichen is used to concatenate a list of items into a single string, optionally using a specified
separator.

## Examples

### Basic Join with Default Separator

Join a list of words using the default space separator:

```yaml
joinedWords: !Join [hello, world]
```

**Output:**

```yaml
joinedWords: "hello world"
```

### Join with Custom Separator

Join a list of words using a custom separator:

```yaml
joinedWords: !Join { items: [hello, world], separator: ", " }
```

**Output:**

```yaml
joinedWords: "hello, world"
```

### Join with Var

Join a list of words using a custom separator and demonstrate use with `!Var`:

```yaml
!Defaults
words: [hello, world]
separator: ", "
---
joinedWords: !Join { items: !Var words, separator: !Var separator }
```

**Output:**

```yaml
joinedWords: "hello, world"
```

## Notes

- The `!Join` tag can concatenate any list of scalar values (strings, numbers, booleans).
- If `separator` is not specified, a single space (`" "`) is used as the default separator.
- To use multiple tags in sequence, such as `!Not !Var`, compose them with a comma, e.g., `!Not,Var`, to maintain valid YAML syntax.