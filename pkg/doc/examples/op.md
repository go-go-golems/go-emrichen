---
Title: "!Op Tag"
Slug: tag-op
Short: |
  ```
  !Op
  a: Value1
  op: Operator
  b: Value2
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
# `!Op` Tag

The `!Op` tag in Emrichen performs binary operations on two values. It supports arithmetic, comparison, and logical
operations. The tag requires three arguments:

- `a`: The first operand.
- `op`: The operation to perform. Supported operations include arithmetic (`+`, `-`, `*`, `/`, `%`),
  comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), and logical (`&&`, `||`).
- `b`: The second operand.

The `!Op` tag in Emrichen supports a wide range of operators for performing arithmetic,
comparison, logical, and string operations.

### Arithmetic Operators
- `+`, `plus`, `add`: Addition of two numbers.
- `-`, `minus`, `sub`, `subtract`: Subtraction of two numbers.
- `*`, `×`, `mul`, `times`: Multiplication of two numbers.
- `/`, `÷`, `div`, `divide`, `truediv`: Division of two numbers, resulting in a float.
- `//`, `floordiv`: Integer division of two numbers, discarding any remainder.
- `%`, `mod`, `modulo`: Modulus operation, returning the remainder of division.

### Comparison Operators (all types)
- `=`, `==`, `===`: Equality check between two values.
- `≠`, `!=`, `!==`, `ne`: Inequality check between two values.
- `<`, `lt`: Less than comparison between two numbers.
- `>`, `gt`: Greater than comparison between two numbers.
- `<=`, `le`, `lte`: Less than or equal to comparison between two numbers.
- `>=`, `ge`, `gte`: Greater than or equal to comparison between two numbers.

### Logical Operators
- `&&`, `and`: Logical AND operation between two boolean values.
- `||`, `or`: Logical OR operation between two boolean values.
- `!`, `not`: Logical NOT operation, negating a boolean value.

### String Operators
- `contains`: Checks if the first string contains the second string.
- `startswith`: Checks if the first string starts with the second string.
- `endswith`: Checks if the first string ends with the second string.
- `matches`: Checks if the first string matches the regular expression provided in the second string.

### Membership Tests (all types)
- `in`, `∈`: Checks if the first value is present in the second value (which must be a sequence).
- `not in`, `∉`: Checks if the first value is not present in the second value (which must be a sequence).

### Special Notes
- The `!Op` tag dynamically determines the operation to perform based on the `op` argument provided.
- For arithmetic operations, if both operands are integers, the result will also be an integer. If at least one operand is a float, the result will be a float.
- Comparison and logical operations return boolean values (`true` or `false`).
- String operations require both operands to be strings, except for the `matches` operation, where the second operand is a regular expression pattern.

## Examples

### Arithmetic Operation: Addition

```yaml
- expression: !Op
    a: 5
    op: "+"
    b: 3
  result: 8
```

### Comparison Operation: Greater Than

```yaml
- expression: !Op
    a: 10
    op: ">"
    b: 5
  result: true
```

### Logical Operation: AND

```yaml
- conditions:
    - conditionA: true
    - conditionB: false
  check: !Op
    a: !Var conditionA
    op: "&&"
    b: !Var conditionB
  result: false
```
