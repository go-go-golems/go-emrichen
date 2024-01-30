---
Title: "!Loop Tag"
Slug: tag-loop
Short: |
  ```
  !Loop
  over: [List or Dict]
  as: item (optional, default `item`)
  index_as: index (optional)
  index_start: start_index (optional, default `0`)
  previous_as: previous_item (optional)
  template: !Var item
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
# `!Loop` Tag

The `!Loop` tag in Emrichen is a versatile construct designed for iterating over collections—either lists or
dictionaries—and applying a specified template to each element of the collection.

The tag supports several keys to control the iteration process and how each element is processed:

- `over`: (Required) Specifies the collection to iterate over. This can be a direct list or dictionary, or a
  variable (`!Var`) pointing to a list or dictionary.
- `as`: (Optional, default `item`) Defines the variable name to represent the current element in the collection during
  each iteration of the loop.
- `index_as`: (Optional) When provided, assigns the current index (for lists) or key (for dictionaries) to a variable
  with the given name, making it accessible within the template.
- `index_start`: (Optional, default `0`) Sets the starting index for the loop. This is useful for skipping a certain
  number of elements at the beginning of the collection. Note that this applies to sequences only and is ignored when
  iterating over dictionaries.
- `previous_as`: (Optional) Allows access to the previous element in the collection by assigning it to a variable with
  the specified name. On the first iteration, the previous element is considered `null`.
- `template`: (Required) The template to be applied to each element of the collection. This template can utilize the
  variables defined by `as`, `index_as`, and `previous_as`.
- `as_documents`: (Not supported yet) If set to `true`, each iteration's output is treated as a separate YAML document. This is
  primarily useful at the top level of a template to generate multiple documents from a single loop. 

## Examples

### Basic Loop Over List

Loop over a list of numbers, printing each number.

```yaml
- !Loop
    over: [1, 2, 3]
    template: !Format "Number: {item}"
```

**Output:**

```yaml
- "Number: 1"
- "Number: 2"
- "Number: 3"
```

### Loop Over Dictionary

Loop over a dictionary, printing each key-value pair.

```yaml
- !Loop
    over: {a: 1, b: 2}
    index_as: key
    template: !Format "{key}: {item}"
```

**Output:**

```yaml
- "a: 1"
- "b: 2"
```

### Nested Loops

Demonstrate nested loops by iterating over a list of dictionaries, each containing a list.

```yaml
- !Loop
    over:
      - name: "list1"
        items: [1, 2]
      - name: "list2"
        items: [3, 4]
    as: list
    template:
      name: !Lookup list.name
      items: !Loop
        over: !Lookup list.items
        template: !Format "Item: {item}"
```

**Output:**

```yaml
- name: "list1"
  items:
    - "Item: 1"
    - "Item: 2"
- name: "list2"
  items:
    - "Item: 3"
    - "Item: 4"
```
