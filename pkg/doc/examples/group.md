---
Title: "!Group Tag"
Slug: tag-group
Short: |
  ```
  !Group
  over: [List or Dict]
  as: item
  by: [Grouping Expression]
  template: [Template (Optional)]
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
# `!Group` Tag

The `!Group` tag in Emrichen is used to group items in a list or dict into a dictionary based on a specified key. It
requires the `over` argument to specify the list or dict to group, the `by` argument to define how to group items, and
optionally a `template` to format the grouped items.

Optionally, the `as` argument can be used to define the name of the variable exposed to the `by` and `template`.

## Examples

### Basic Grouping by Name

Group a list of people by their name:

```yaml
!Group
  over:
    - name: Alice
      age: 30
    - name: Bob
      age: 25
    - name: Alice
      age: 32
  as: person
  by: !Lookup person.name
  template: !Lookup person.age
```

**Output:**

```yaml
Alice: [30, 32]
Bob: [25]
```

### Grouping with Custom Key

Group items and generate a custom key:

```yaml
!Group
  over: [1, 2, 2, 3, 3, 3]
  as: number
  by: !Format "Group_{number}"
```

**Output:**

```yaml
Group_1: [1]
Group_2: [2, 2]
Group_3: [3, 3, 3]
```

### Grouping Without a Template

Group objects by a nested field without specifying a template:

```yaml
!Group
  over:
    - category: Hardware
      item: Keyboard
    - category: Hardware
      item: Mouse
    - category: Software
      item: Operating System
  as: product
  by: !Lookup product.category
```

**Output:**

```yaml
Hardware: 
  - category: Hardware
    item: Keyboard
  - category: Hardware
    item: Mouse
Software: 
  - category: Software
    item: Operating System
```

## Notes

- The `!Group` tag's `over` argument must be a list or dict.
- The `by` argument is required and determines how items are grouped.
- The `template` argument is optional and used to format each item in the group.
- Grouped items are always returned as a dictionary.