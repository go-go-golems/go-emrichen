---
Title: "!Lookup Tag"
Slug: tag-lookup
Short: |
  ```
  !Lookup JSONPathExpression
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
# `!Lookup` Tag

The `!Lookup` tag in Emrichen is used to perform a JSONPath lookup, returning the first match for the given expression. If there is no match, an error is raised.

```yaml
!Lookup JSONPathExpression
```

## Examples

### Basic Lookup

Retrieve a single value from a list of dictionaries:

```yaml
!Defaults
people:
  - name: "Alice"
    age: 30
---
personName: !Lookup people[0].name
```

**Output:**

```yaml
personName: "Alice"
```

### Lookup with Error Handling

Attempt to lookup a non-existent index, demonstrating error handling:

```yaml
!Defaults
people:
  - name: "Alice"
    age: 30
---
tryLookupNonexistent: !Lookup people[999].name
```

**Output:**

Error: No match found for JSONPath expression.

### Lookup and Modify

Lookup a value and modify it using tag composition:

```yaml
!Defaults
people:
  - name: "Alice"
    age: 30
---
personNameUppercase: !Format,Lookup people[0].name
```

**Output:**

```yaml
personNameUppercase: "ALICE"
```

## Notes

- The `!Lookup` tag requires a valid JSONPath expression as its argument.