---
Title: "!Filter Tag"
Slug: tag-filter
Short: |
  ```
  !Filter
    test: <Predicate>
    over: <List or Dict>
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
# `!Filter` Tag

The `!Filter` tag in Emrichen is used to filter elements of a list or dictionary based on a given predicate. It takes
two arguments: `test`, which defines the filtering condition, and `over`, which is the list or dictionary to be
filtered. The iterated value is exposed as the `item` variable.

```yaml
!Filter
  test: <Predicate>
  over: <List or Dict>
```

## Examples

### Basic Number Filtering

Filter a list of numbers to include only those greater than 3:

```yaml
!Defaults
numbers: [1, 2, 3, 4, 5, 6]
---
filteredNumbers: !Filter
  test: !Op
    a: !Var item
    op: gt
    b: 3
  over: !Var numbers
```

**Output:**

```yaml
filteredNumbers: [4, 5, 6]
```

### String Filtering with Custom Predicate

Filter a list of strings to include only those starting with 'A':

```yaml
!Defaults
strings: ["Apple", "Banana", "Avocado", "Grape"]
---
filteredStrings: !Filter
  test: !Op
    a: !Var item
    op: startswith
    b: "A"
  over: !Var strings
```

**Output:**

```yaml
filteredStrings: ["Apple", "Avocado"]
```

### Filtering Nested Lists in Dictionaries

Filter a list of dictionaries, then filter a nested list within each dictionary:

```yaml
!Defaults
people:
  - name: "Alice"
    hobbies: ["Swimming", "Cycling", "Reading"]
  - name: "Bob"
    hobbies: ["Hiking", "Reading", "Gaming"]
  - name: "Charlie"
    hobbies: ["Gaming", "Cycling", "Cooking"]
---
filteredPeople: !Filter
  test: !Any,Loop
    over: !Lookup item.hobbies
    template: !Op
      a: !Var item
      op: startswith
      b: "C"
  over: !Var people
```

**Output:**

```yaml
filteredPeople:
  - name: "Alice"
    hobbies: ["Swimming", "Cycling", "Reading"]
  - name: "Charlie"
    hobbies: ["Gaming", "Cycling", "Cooking"]
```

## Notes

- The `test` argument in `!Filter` must be a predicate that returns a boolean value.
- `!Filter` can operate over both lists and dictionaries.
- When filtering dictionaries, the `test` is applied to the dictionary's values.