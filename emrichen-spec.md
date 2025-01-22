# Detailed Specification of Emrichen Language

Emrichen is a powerful YAML templating language that introduces a variety of custom tags to facilitate dynamic and flexible configuration generation. This specification provides an in-depth look at each Emrichen tag, detailing the types of its arguments, structure, and usage examples to guide users in effectively leveraging Emrichen for their YAML configurations.

## Table of Contents

1. [`!All` Tag](#all-tag)
2. [`!Any` Tag](#any-tag)
3. [`!Base64` Tag](#base64-tag)
4. [`!Concat` Tag](#concat-tag)
5. [`!Debug` Tag](#debug-tag)
6. [`!Defaults` Tag](#defaults-tag)
7. [`!Error` Tag](#error-tag)
8. [`!Exists` Tag](#exists-tag)
9. [`!Filter` Tag](#filter-tag)
10. [`!Format` Tag](#format-tag)
11. [`!Group` Tag](#group-tag)
12. [`!If` Tag](#if-tag)
13. [`!Include` Tags Family](#include-tags-family)
    - [`!Include` Tag](#include-tag)
    - [`!IncludeBase64` Tag](#includebase64-tag)
    - [`!IncludeBinary` Tag](#includebinary-tag)
    - [`!IncludeGlob` Tag](#includeglob-tag)
    - [`!IncludeText` Tag](#includetext-tag)
14. [`!Index` Tag](#index-tag)
15. [`!IsBoolean` Tag](#isboolean-tag)
16. [`!IsDict` Tag](#isdict-tag)
17. [`!IsInteger` Tag](#isinteger-tag)
18. [`!IsList` Tag](#islist-tag)
19. [`!IsNone` Tag](#isnone-tag)
20. [`!IsNumber` Tag](#isnumber-tag)
21. [`!IsString` Tag](#isstring-tag)
22. [`!Join` Tag](#join-tag)
23. [`!Lookup` Tags](#lookup-tags)
    - [`!Lookup` Tag](#lookup-tag)
    - [`!LookupAll` Tag](#lookupall-tag)
24. [`!Loop` Tag](#loop-tag)
25. [`!MD5` Tag](#md5-tag)
26. [`!Merge` Tag](#merge-tag)
27. [`!Not` Tag](#not-tag)
28. [`!Op` Tag](#op-tag)
29. [`!SHA1` Tag](#sha1-tag)
30. [`!SHA256` Tag](#sha256-tag)
31. [`!URLEncode` Tag](#urlencode-tag)
32. [`!Var` Tag](#var-tag)
33. [`!Void` Tag](#void-tag)
34. [`!With` Tag](#with-tag)

---

## `!All` Tag

**Purpose:**  
Evaluates whether **all** items in a sequence are truthy. Returns `true` if every item evaluates to a truthy value; otherwise, returns `false`.

**Arguments:**  
- **Sequence Node:** A YAML sequence (`[]`) containing expressions or values to evaluate.

**Structure:**
```yaml
allConditionsMet: !All
  - !Var is_admin
  - !Lookup user_roles[0]
  - !Lookup feature_flags.feature_x
```

**Example:**
```yaml
allConditionsMet: !All
  - !Var is_admin
  - !Lookup user_roles[0]  # Should be 'admin'
  - !Lookup feature_flags.feature_x
```

**Behavior:**
- Each item in the sequence is processed.
- If **any** item evaluates to a falsy value, `!All` returns `false`.
- If **all** items are truthy, `!All` returns `true`.

---

## `!Any` Tag

**Purpose:**  
Evaluates whether **at least one** item in a sequence is truthy. Returns `true` if any item evaluates to a truthy value; otherwise, returns `false`.

**Arguments:**  
- **Sequence Node:** A YAML sequence (`[]`) containing expressions or values to evaluate.

**Structure:**
```yaml
anyConditionMet: !Any
  - !Var is_admin
  - !Lookup feature_flags.feature_y
```

**Example:**
```yaml
anyConditionMet: !Any
  - !Var is_admin
  - !Lookup feature_flags.feature_y  # This might be false, but 'is_admin' could be true
```

**Behavior:**
- Each item in the sequence is processed.
- If **any** item evaluates to a truthy value, `!Any` returns `true`.
- If **all** items are falsy, `!Any` returns `false`.

---

## `!Base64` Tag

**Purpose:**  
Encodes a given scalar value into a Base64-encoded string.

**Arguments:**  
- **Scalar Node:** A YAML scalar (`string`, `int`, `bool`, etc.) to be encoded.

**Structure:**
```yaml
encodedString: !Base64 "Hello, World!"
```

**Example:**
```yaml
encodedString: !Base64 "Hello, World!"
```

**Behavior:**
- Takes the scalar value, encodes it using Base64, and returns the encoded string.

**Result:**
```yaml
encodedString: "SGVsbG8sIFdvcmxkIQ=="
```

---

## `!Concat` Tag

**Purpose:**  
Concatenates multiple YAML sequences into a single sequence.

**Arguments:**  
- **Sequence of Sequences:** A YAML sequence where each item is itself a sequence to be concatenated.

**Structure:**
```yaml
concatenatedList: !Concat
  - [1, 2]
  - [3, 4]
  - [5]
```

**Example:**
```yaml
concatenatedList: !Concat
  - ["apple", "banana"]
  - ["orange"]
  - ["grape", "melon"]
```

**Behavior:**
- Each sequence within the `!Concat` tag is processed and concatenated in order.
- The result is a single, flattened sequence containing all elements from the input sequences.

**Result:**
```yaml
concatenatedList:
  - "apple"
  - "banana"
  - "orange"
  - "grape"
  - "melon"
```

---

## `!Debug` Tag

**Purpose:**  
Outputs the value of a node to stderr for debugging purposes without altering the YAML structure.

**Arguments:**  
- **Any Node:** Can be a scalar, sequence, mapping, etc.

**Structure:**
```yaml
debuggedValue: !Debug "This is a debug message"
```

**Example:**
```yaml
debuggedValue: !Debug
  - !Var debugExampleVar
```

**Behavior:**
- Processes the node and prints its value to stderr.
- The original node's value remains unchanged in the final YAML output.

**Note:**  
Useful for developers to inspect intermediate values during the templating process.

---

## `!Defaults` Tag

**Purpose:**  
Defines default values for variables that can be referenced throughout the YAML document using `!Var`.

**Arguments:**  
- **Mapping Node:** A YAML mapping (`{}`) containing key-value pairs that set default variables.

**Structure:**
```yaml
!Defaults
app_name: myapp
image_tag: v1.0.0
replica_count: 2
```

**Example:**
```yaml
!Defaults
name: John Doe
age: 30
hobbies: [Reading, Coding, Hiking]
address:
  street: 123 Main St
  city: Anytown
  zip: 12345
---
person:
  name: !Var name
  age: !Var age
  hobbies: !Var hobbies
  address: !Var address
```

**Behavior:**
- Sets default variables that can be accessed using the `!Var` tag.
- Variables defined under `!Defaults` are available globally within the document unless overridden.

---

## `!Error` Tag

**Purpose:**  
Outputs an error message and halts the processing of the YAML document when encountered.

**Arguments:**  
- **Scalar Node:** A YAML scalar containing the error message to be displayed.

**Structure:**
```yaml
fatalError: !Error "An unrecoverable error occurred."
```

**Example:**
```yaml
conditionalError: !If
  test: !IsNone,Var errorExampleVar
  then: !Error "errorExampleVar must not be null"
  else: !Var errorExampleVar
```

**Behavior:**
- When the `!Error` tag is processed, it immediately stops further processing and returns the specified error message.
- Useful for enforcing mandatory configurations and validating conditions.

---

## `!Exists` Tag

**Purpose:**  
Checks for the existence of a specified path or variable within the current environment.

**Arguments:**  
- **String Node:** A YAML scalar representing the JSONPath or variable name to check for existence.

**Structure:**
```yaml
existsFoo: !Exists foo
existsNestedKey: !Exists nested.key
```

**Example:**
```yaml
!Defaults
foo: bar
nested:
  key: value
list:
  - item1
  - item2
---
existsFoo: !Exists foo
existsBaz: !Exists baz
```

**Behavior:**
- Evaluates to `true` if the specified path or variable exists; otherwise, `false`.
- Handles nested paths and list indices.

**Result:**
```yaml
existsFoo: true
existsBaz: false
```

---

## `!Filter` Tag

**Purpose:**  
Filters elements from a collection based on a provided predicate.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `test` (required): An expression that evaluates to a boolean, determining whether an item should be included.
  - `over` (required): The collection (sequence or mapping) to filter.
  - `as` (optional): The variable name to assign the current item during evaluation.

**Structure:**
```yaml
filteredItems: !Filter
  test: !Op
    a: !Var item
    op: ">"
    b: 5
  over: !Var numbers
  as: num
```

**Example:**
```yaml
!Defaults
numbers: [1, 2, 3, 4, 5, 6]
---
title: "Basic Filter Example"
filteredNumbers: !Filter
  test: !Op
    a: !Var item
    op: ">"
    b: 3
  over: !Var numbers
```

**Behavior:**
- Iterates over each item in the `over` collection.
- Evaluates the `test` expression for each item.
- Includes the item in the result if the `test` evaluates to `true`.
- Supports filtering within mappings and nested structures.

**Result:**
```yaml
filteredNumbers:
  - 4
  - 5
  - 6
```

---

## `!Format` Tag

**Purpose:**  
Formats a string using variables and expressions, supporting both simple and complex interpolations.

**Arguments:**  
- **String Node:** A YAML scalar containing the format string with placeholders.

**Structure:**
```yaml
greeting: !Format "Hello, {name}! You are {age} years old."
```

**Example:**
```yaml
!Defaults
name: John
age: 30
---
person:
  greeting: !Format "Hello, {name}! You are {age} years old."
```

**Behavior:**
- Processes the format string, replacing placeholders with corresponding variable values.
- Supports Go template syntax for more advanced formatting and expressions.

**Advanced Example with Go Template Syntax:**
```yaml
!Defaults
user:
  firstName: Jane
  lastName: Doe
  age: 35
---
userGreeting: !Format "Hello, {{.user.firstName}} {{.user.lastName}}! Next year, you will be {{ add (lookup \"user.age\") 1 }} years old."
```

**Result:**
```yaml
userGreeting: "Hello, Jane Doe! Next year, you will be 36 years old."
```

---

## `!Group` Tag

**Purpose:**  
Groups items from a collection based on a specified key or expression, optionally applying a template to each grouped item.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `over` (required): The collection (sequence or mapping) to group.
  - `by` (required): An expression or key to determine the grouping.
  - `template` (optional): An expression to transform each item before grouping.
  - `as` (optional): The variable name for the current item during evaluation.
  - `result_as` (optional): The variable name for the result of the template.

**Structure:**
```yaml
groupedPeople: !Group
  over: !Var people
  by: !Lookup person.age
  template: !Lookup person.name
  as: person
```

**Example:**
```yaml
!Defaults
people:
  - name: Alice
    age: 30
  - name: Bob
    age: 40
  - name: Carol
    age: 30
---
groupedPeopleByAge: !Group
  over: !Var people
  by: !Lookup person.age
  template: !Lookup person.name
  as: person
```

**Behavior:**
- Iterates over each item in the `over` collection.
- Evaluates the `by` expression to determine the grouping key.
- Applies the `template` expression to transform each item if provided.
- Groups items under their respective keys.
- Handles duplicate keys based on the grouping logic.

**Result:**
```yaml
groupedPeopleByAge:
  "30":
    - "Alice"
    - "Carol"
  "40":
    - "Bob"
```

---

## `!If` Tag

**Purpose:**  
Implements conditional logic to include different values based on the evaluation of a test expression.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `test` (required): An expression that evaluates to a boolean.
  - `then` (optional): The value to return if `test` is `true`.
  - `else` (optional): The value to return if `test` is `false`.

**Structure:**
```yaml
message: !If
  test: !Var is_admin
  then: "You have administrative privileges."
  else: "You are a regular user."
```

**Example:**
```yaml
!Defaults
user:
  isAdmin: true
  name: John
---
message: !If
  test: !Var user.isAdmin
  then: !Format "You are an admin, {{.user.name}}!"
  else: !Format "Hello, {{.user.name}}!"
```

**Behavior:**
- Evaluates the `test` expression.
- If `test` is `true`, processes and returns the `then` value.
- If `test` is `false`, processes and returns the `else` value.
- If `then` or `else` is omitted, the result for that branch is `null`.

**Result:**
```yaml
message: "You are an admin, John!"
```

---

## `!Include` Tags Family

Emrichen provides a family of `!Include` tags to incorporate external files into the YAML document. Each variant handles different types of files and inclusion methods.

### `!Include` Tag

**Purpose:**  
Includes and processes another YAML template file within the current document.

**Arguments:**  
- **Scalar Node:** A YAML scalar containing the file path to include.

**Structure:**
```yaml
includedContent: !Include "path/to/another_template.yml"
```

**Example:**
```yaml
includedContent: !Include,Var singleFile
```

**Behavior:**
- Reads the specified YAML file.
- Processes its content as if it were part of the current document.
- Supports multiple YAML documents within the included file.

---

### `!IncludeBase64` Tag

**Purpose:**  
Includes a binary file and encodes its content into a Base64 string.

**Arguments:**  
- **Scalar Node:** A YAML scalar containing the file path to include.

**Structure:**
```yaml
base64Content: !IncludeBase64 "path/to/file.bin"
```

**Example:**
```yaml
base64Content: !IncludeBase64,Var base64
```

**Behavior:**
- Reads the specified binary file.
- Encodes its content using Base64.
- Inserts the encoded string into the YAML document.

**Result:**
```yaml
base64Content: "SGVsbG8sIFdvcmxkIQ=="
```

---

### `!IncludeBinary` Tag

**Purpose:**  
Includes the raw content of a binary file directly into the YAML document.

**Arguments:**  
- **Scalar Node:** A YAML scalar containing the file path to include.

**Structure:**
```yaml
binaryContent: !IncludeBinary "path/to/file.bin"
```

**Example:**
```yaml
binaryContent: !IncludeBinary,Var binary
```

**Behavior:**
- Reads the specified binary file.
- Inserts its raw content as a string into the YAML document.

**Note:**  
Ensure that the consuming application can handle raw binary data appropriately.

---

### `!IncludeGlob` Tag

**Purpose:**  
Includes and processes multiple files that match a specified glob pattern.

**Arguments:**  
- **Scalar Node or Sequence Node:** A YAML scalar containing the glob pattern or a sequence of glob patterns.

**Structure:**
```yaml
globIncludedContent: !IncludeGlob "configs/*.yml"
```

**Example:**
```yaml
globIncludedContent: !IncludeGlob,Var globPattern
```

**Behavior:**
- Expands the glob pattern to match multiple files.
- Reads and processes each matched YAML file.
- Aggregates the content into a YAML sequence.

**Result:**
```yaml
globIncludedContent:
  - { message: "Content from glob_test_1.yml" }
  - { text: "Content from glob_test_2.yml" }
```

---

### `!IncludeText` Tag

**Purpose:**  
Includes the content of a text file as a string within the YAML document.

**Arguments:**  
- **Scalar Node:** A YAML scalar containing the file path to include.

**Structure:**
```yaml
textContent: !IncludeText "path/to/text_test.txt"
```

**Example:**
```yaml
textContent: !IncludeText,Var textFile
```

**Behavior:**
- Reads the specified text file.
- Inserts its content as a string into the YAML document.

**Result:**
```yaml
textContent: "Sample Text Content"
```

---

## `!Index` Tag

**Purpose:**  
Creates a dictionary from a list by indexing items based on a specified key or expression.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `over` (required): The collection (sequence) to index.
  - `by` (required): An expression that determines the key for each item.
  - `template` (optional): An expression to transform each item before indexing.
  - `as` (optional): The variable name for the current item during evaluation.
  - `duplicates` (optional): Determines behavior when duplicate keys are encountered (`error`, `warn`, `ignore`).
  - `result_as` (optional): The variable name for the result of the `template` expression.

**Structure:**
```yaml
indexedEmployees: !Index
  over: !Var employees
  by: !Lookup item.name
  template: !Lookup item.department
  as: employee
  duplicates: ignore
```

**Example:**
```yaml
!Defaults
employees:
  - name: Alice
    department: Engineering
  - name: Bob
    department: Marketing
  - name: Alice
    department: HR
---
indexedEmployees: !Index
  over: !Var employees
  by: !Lookup employee.name
  template: !Lookup employee.department
  as: employee
  duplicates: ignore
```

**Behavior:**
- Iterates over each item in the `over` collection.
- Evaluates the `by` expression to determine the key for each item.
- Optionally applies the `template` expression to transform each item before indexing.
- Handles duplicate keys based on the `duplicates` argument:
  - `error`: Throws an error on encountering duplicate keys.
  - `warn`: Logs a warning and skips duplicate entries.
  - `ignore`: Silently skips duplicate entries.

**Result:**
```yaml
indexedEmployees:
  "Alice": "Engineering"  # Second "Alice" entry is ignored
  "Bob": "Marketing"
```

---

## `!IsBoolean` Tag

**Purpose:**  
Checks if a given value is of boolean type.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to check.

**Structure:**
```yaml
isBoolean: !IsBoolean true
```

**Example:**
```yaml
isBooleanExample: !IsBoolean,Var someFlag
```

**Behavior:**
- Returns `true` if the value is a boolean (`true` or `false`); otherwise, returns `false`.

**Result:**
```yaml
isBoolean: true
```

---

## `!IsDict` Tag

**Purpose:**  
Checks if a given value is a dictionary (mapping) type.

**Arguments:**  
- **Any Node:** Can be a scalar, sequence, or mapping.

**Structure:**
```yaml
isDict: !IsDict {key: "value"}
```

**Example:**
```yaml
isDictExample: !IsDict,Var configSettings
```

**Behavior:**
- Returns `true` if the value is a mapping (`{}`); otherwise, returns `false`.

**Result:**
```yaml
isDict: true
```

---

## `!IsInteger` Tag

**Purpose:**  
Checks if a given value is an integer.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to check.

**Structure:**
```yaml
isInteger: !IsInteger 42
```

**Example:**
```yaml
isIntegerExample: !IsInteger,Var userAge
```

**Behavior:**
- Returns `true` if the value is an integer; otherwise, returns `false`.

**Result:**
```yaml
isInteger: true
```

---

## `!IsList` Tag

**Purpose:**  
Checks if a given value is a list (sequence) type.

**Arguments:**  
- **Any Node:** Can be a scalar, sequence, or mapping.

**Structure:**
```yaml
isList: !IsList [1, 2, 3]
```

**Example:**
```yaml
isListExample: !IsList,Var items
```

**Behavior:**
- Returns `true` if the value is a sequence (`[]`); otherwise, returns `false`.

**Result:**
```yaml
isList: true
```

---

## `!IsNone` Tag

**Purpose:**  
Checks if a given value is `null` (none).

**Arguments:**  
- **Any Node:** Can be a scalar, sequence, or mapping.

**Structure:**
```yaml
isNone: !IsNone null
```

**Example:**
```yaml
isNoneExample: !IsNone,Var optionalField
```

**Behavior:**
- Returns `true` if the value is `null` or explicitly set to `null`; otherwise, returns `false`.

**Result:**
```yaml
isNone: true
```

---

## `!IsNumber` Tag

**Purpose:**  
Checks if a given value is a number (integer or float).

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to check.

**Structure:**
```yaml
isNumber: !IsNumber 3.14
```

**Example:**
```yaml
isNumberExample: !IsNumber,Var responseTime
```

**Behavior:**
- Returns `true` if the value is a number; otherwise, returns `false`.

**Result:**
```yaml
isNumber: true
```

---

## `!IsString` Tag

**Purpose:**  
Checks if a given value is a string.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to check.

**Structure:**
```yaml
isString: !IsString "Hello"
```

**Example:**
```yaml
isStringExample: !IsString,Var userName
```

**Behavior:**
- Returns `true` if the value is a string; otherwise, returns `false`.

**Result:**
```yaml
isString: true
```

---

## `!Join` Tag

**Purpose:**  
Joins a list of items into a single string with an optional separator.

**Arguments:**  
- **Either:**
  - **Scalar Node:** The list of items to join.
  - **Mapping Node:** Must contain:
    - `items` (required): A YAML sequence (`[]`) of items to join.
    - `separator` (optional): A string to separate the joined items. Defaults to a space (`" "`).

**Structure:**
```yaml
joinedWords: !Join,Var words
```
Or using a mapping node:
```yaml
joinedCustomSeparatorWords: !Join
  items: ["hello", "world"]
  separator: ", "
```

**Example:**
```yaml
!Defaults
words: ["hello", "world"]
customSeparator: ", "
---
joinedWords: !Join,Var words
joinedWithCustomSeparator: !Join
  items: !Var words
  separator: !Var customSeparator
```

**Behavior:**
- Concatenates the list of items using the specified separator.
- If no separator is provided, defaults to a space.

**Result:**
```yaml
joinedWords: "hello world"
joinedWithCustomSeparator: "hello, world"
```

**Edge Cases:**
- **Empty List:** Returns an empty string.
- **Single Element:** Returns the element itself without a separator.
- **Non-String Elements:** Converts all items to strings before joining.

---

## `!Lookup` Tags

Emrichen provides two `!Lookup` tags to perform JSONPath queries on the current environment's variables, allowing for more complex variable access patterns.

### `!Lookup` Tag

**Purpose:**  
Performs a JSONPath lookup and returns the first matching result, supporting nested access and complex queries.

**Arguments:**  
- **String Node:** A YAML scalar representing the JSONPath expression to evaluate.

**Structure:**
```yaml
personName: !Lookup "people[0].name"
```

**Example:**
```yaml
!Defaults
people:
  - name: "Alice"
    age: 30
  - name: "Bob"
    age: 25
user:
  firstName: Jane
  lastName: Doe
---
firstPersonName: !Lookup "people[0].name"
userFirstName: !Lookup "user.firstName"
```

**Behavior:**
- Evaluates the JSONPath expression against the current environment.
- Returns the first match found.
- Supports nested access and complex queries.
- If no match is found, it results in an error unless handled otherwise.

**Result:**
```yaml
firstPersonName: "Alice"
userFirstName: "Jane"
```

---

### `!LookupAll` Tag

**Purpose:**  
Performs a JSONPath lookup and returns **all** matching results as a list, supporting nested access and complex queries.

**Arguments:**  
- **String Node:** A YAML scalar representing the JSONPath expression to evaluate.

**Structure:**
```yaml
personAges: !LookupAll "people[*].age"
```

**Example:**
```yaml
!Defaults
people:
  - name: "Alice"
    age: 30
  - name: "Bob"
    age: 25
---
allPersonAges: !LookupAll "people[*].age"
```

**Behavior:**
- Evaluates the JSONPath expression against the current environment.
- Returns all matches found as a YAML sequence.
- Supports nested access and complex queries.
- If no matches are found, returns an empty list.

**Result:**
```yaml
allPersonAges:
  - 30
  - 25
```

**Important Distinction:**
- Use `!Var` for direct, top-level variable access.
- Use `!Lookup` or `!LookupAll` for nested access, JSONPath queries, or more complex variable retrieval patterns.

---

## `!Loop` Tag

**Purpose:**  
Iterates over a collection (list or mapping) and applies a template to each item, generating a new collection based on the results.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `over` (required): The collection (sequence or mapping) to iterate over.
  - `template` (required): The template to apply to each item.
  - `as` (optional): The variable name for the current item. Defaults to `item`.
  - `index_as` (optional): The variable name for the current index or key.
  - `previous_as` (optional): The variable name for the previous item.
  - `index_start` (optional): An integer to specify the starting index. Defaults to `0`.
  - `as_documents` (optional): Treats each iteration's output as a separate YAML document. (Not supported yet)

**Structure:**
```yaml
loopedNumbers: !Loop
  over: !Var numbers
  as: number
  template: !Format "Number: {number}"
```

**Example:**
```yaml
!Defaults
numbers: [1, 2, 3, 4, 5]
---
loopedNumbers: !Loop
  over: !Var numbers
  as: num
  template: !Format "Number {num}"
```

**Behavior:**
- Iterates over each element in the `over` collection.
- Assigns each element to the variable specified by `as`.
- Optionally assigns the current index/key to `index_as` and the previous item to `previous_as`.
- Applies the `template` to each item.
- Aggregates the results into a new sequence or mapping.

**Result:**
```yaml
loopedNumbers:
  - "Number 1"
  - "Number 2"
  - "Number 3"
  - "Number 4"
  - "Number 5"
```

**Advanced Example with Index and Previous Value:**
```yaml
loopDetails: !Loop
  over: !Var numbers
  as: current
  previous_as: prev
  index_as: idx
  template: !Format "Index {idx}: Current={current}, Previous={prev}"
```

**Result:**
```yaml
loopDetails:
  - "Index 0: Current=1, Previous=null"
  - "Index 1: Current=2, Previous=1"
  - "Index 2: Current=3, Previous=2"
  - "Index 3: Current=4, Previous=3"
  - "Index 4: Current=5, Previous=4"
```

**Note:**  
The `as_documents` argument is reserved for future implementations and is currently not supported.

---

## `!MD5` Tag

**Purpose:**  
Generates an MD5 hash of the provided scalar value.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to hash.

**Structure:**
```yaml
hashedStringMD5: !MD5 "data to hash"
```

**Example:**
```yaml
hashedString: !MD5,Var sampleString
```

**Behavior:**
- Computes the MD5 hash of the input value.
- Returns the hash as a hexadecimal string.

**Result:**
```yaml
hashedStringMD5: "6cd3556deb0da54bca060b4c39479839"
```

---

## `!Merge` Tag

**Purpose:**  
Merges multiple YAML mapping nodes into a single mapping. Later mappings can override keys from earlier ones.

**Arguments:**  
- **Sequence of Mapping Nodes:** A YAML sequence (`[]`) where each item is a mapping to be merged.

**Structure:**
```yaml
mergedConfig: !Merge
  - {a: 1, b: 2}
  - {b: 3, c: 4}
  - {d: 5}
```

**Example:**
```yaml
!Defaults
dict1: {a: 1, b: 2}
dict2: {b: 3, c: 4}
dict3: {d: 5}
---
mergedConfig: !Merge
  - !Var dict1
  - !Var dict2
  - !Var dict3
```

**Behavior:**
- Iterates over each mapping in the sequence.
- Merges them into a single mapping.
- If duplicate keys are encountered, the value from the **last** mapping takes precedence.

**Result:**
```yaml
mergedConfig:
  a: 1
  b: 3
  c: 4
  d: 5
```

**Advanced Example with Nested Mappings:**
```yaml
!Defaults
nestedDict1: {x: {y: 1, z: 2}}
nestedDict2: {x: {z: 3, w: 4}}
---
mergedNested: !Merge
  - !Var nestedDict1
  - !Var nestedDict2
```

**Result:**
```yaml
mergedNested:
  x:
    y: 1
    z: 3
    w: 4
```

---

## `!Not` Tag

**Purpose:**  
Negates a boolean value.

**Arguments:**  
- **Any Node:** A YAML node that evaluates to a boolean.

**Structure:**
```yaml
negatedValue: !Not,Var isEnabled
```

**Example:**
```yaml
!Defaults
isEnabled: true
---
isDisabled: !Not,Var isEnabled
```

**Behavior:**
- Evaluates the input node.
- Returns the logical negation of the boolean value.

**Result:**
```yaml
isDisabled: false
```

---

## `!Op` Tag

**Purpose:**  
Performs a binary operation between two values based on the specified operator.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `a` (required): The first operand.
  - `op` (required): The operator as a string (e.g., `+`, `>`, `==`).
  - `b` (required): The second operand.

**Structure:**
```yaml
operationResult: !Op
  a: 5
  op: "+"
  b: 3
```

**Example:**
```yaml
!Defaults
number1: 10
number2: 4
---
addition: !Op
  a: !Var number1
  op: "+"
  b: !Var number2
subtraction: !Op
  a: !Var number1
  op: "-"
  b: !Var number2
```

**Behavior:**
- Evaluates the operands `a` and `b`.
- Performs the specified operation (`op`) between them.
- Supports arithmetic, logical, and string operations.

**Supported Operators:**
- **Arithmetic:** `+` (add), `-` (subtract), `*` (multiply), `/` (divide), `%` (modulo)
- **Comparison:** `==`, `!=`, `<`, `>`, `<=`, `>=`
- **String Operations:** `contains`, `startswith`, `endswith`, `matches`
- **Membership:** `in`, `not in`

**Result:**
```yaml
addition: 14
subtraction: 6
```

**Advanced Example with Logical Operations:**
```yaml
!Defaults
isAdmin: true
hasAccess: false
---
canAccess: !Op
  a: !Var isAdmin
  op: "&&"
  b: !Var hasAccess
```

**Result:**
```yaml
canAccess: false
```

---

## `!SHA1` Tag

**Purpose:**  
Generates a SHA1 hash of the provided scalar value.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to hash.

**Structure:**
```yaml
hashedStringSHA1: !SHA1 "data to hash"
```

**Example:**
```yaml
hashedString: !SHA1,Var sampleString
```

**Behavior:**
- Computes the SHA1 hash of the input value.
- Returns the hash as a hexadecimal string.

**Result:**
```yaml
hashedStringSHA1: "2ef7bde608ce5404e97d5f042f95f89f1c232871"
```

---

## `!SHA256` Tag

**Purpose:**  
Generates a SHA256 hash of the provided scalar value.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the value to hash.

**Structure:**
```yaml
hashedStringSHA256: !SHA256 "data to hash"
```

**Example:**
```yaml
hashedString: !SHA256,Var sampleString
```

**Behavior:**
- Computes the SHA256 hash of the input value.
- Returns the hash as a hexadecimal string.

**Result:**
```yaml
hashedStringSHA256: "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
```

---

## `!URLEncode` Tag

**Purpose:**  
Encodes a string for safe inclusion in URLs or combines a base URL with query parameters.

**Arguments:**  
- **Either:**
  - **Scalar Node:** A YAML scalar string to be URL-encoded.
  - **Mapping Node:** A YAML mapping containing:
    - `url` (required): The base URL.
    - `query` (optional): A mapping of query parameters.

**Structure:**
```yaml
encodedSimpleString: !URLEncode "hello world & special=characters"

fullURL: !URLEncode
  url: "https://example.com/"
  query:
    param1: "value1"
    param2: "value2"
```

**Example:**
```yaml
!Defaults
baseURL: "https://example.com/"
queryParams:
  key1: "value1"
  key2: "value2"
---
simpleEncode: !URLEncode "hello world & special=characters"
fullURL: !URLEncode
  url: !Var baseURL
  query: !Var queryParams
```

**Behavior:**
- **Scalar Input:** Encodes the string using URL encoding rules.
- **Mapping Input:** Appends query parameters to the base URL after encoding them.

**Result:**
```yaml
simpleEncode: "hello%20world%20%26%20special%3Dcharacters"
fullURL: "https://example.com/?key1=value1&key2=value2"
```

**Edge Cases:**
- **Long Strings:** Handles very long strings efficiently.
- **Invalid Inputs:** Returns an error if the input type is unsupported (e.g., sequences).

**Error Handling Example:**
```yaml
encodedInvalidInput: !URLEncode [1, 2, 3]  # Causes an error
```

---

## `!Var` Tag

**Purpose:**  
Substitutes the value of a directly referenced, predefined variable within the YAML document.

**Arguments:**  
- **Scalar Node:** A YAML scalar representing the exact variable name to substitute.

**Structure:**
```yaml
variableValue: !Var variableName
```

**Example:**
```yaml
!Defaults
app_name: myapp
image_tag: v1.0.0
user:
  firstName: Jane
  lastName: Doe
---
deploymentName: !Var app_name
image: !Var image_tag
```

**Behavior:**
- Retrieves the value of the specified variable directly from the current environment.
- Replaces the `!Var` tag with the variable's value.
- Does not support nested access or JSON paths (use `!Lookup` for those cases).

**Result:**
```yaml
deploymentName: "myapp"
image: "v1.0.0"
```

**Important Note:**
`!Var` only accesses top-level variables directly. For nested access or more complex queries, use `!Lookup`.

**Incorrect Usage (Will Not Work):**
```yaml
# This will not work as expected
fullName: !Var user.firstName
```

**Correct Usage for Nested Access:**
```yaml
# Use !Lookup for nested access
fullName: !Lookup user.firstName
```

---

## `!Void` Tag

**Purpose:**  
Removes or omits items from the output YAML document.

**Arguments:**  
- **None:** The `!Void` tag does not take any arguments.

**Structure:**
```yaml
optionalField: !Void
```

**Example:**
```yaml
steps:
  - name: Step 1
    action: Deploy
  - name: Step 2
    action: !Void
```

**Behavior:**
- When used within sequences or mappings, items tagged with `!Void` are omitted from the final output.
- Useful for conditionally removing parts of the YAML based on logic.

**Result:**
```yaml
steps:
  - name: Step 1
    action: Deploy
```

---

## `!With` Tag

**Purpose:**  
Defines a local scope with its own set of variables, allowing for scoped variable definitions and template applications.

**Arguments:**  
- **Mapping Node:** Must contain the following keys:
  - `vars` (required): A mapping of variable names to values, defining the local scope.
  - `template` (required): The YAML node to process within the defined scope.

**Structure:**
```yaml
withScope: !With
  vars:
    localVar: "Local Value"
  template: !Var localVar
```

**Example:**
```yaml
!Defaults
globalVar: "global value"
numbers: [1, 2, 3, 4, 5]
---
withScopeExample: !With
  vars:
    localVar: "Hello, Local!"
  template: !Var localVar
```

**Behavior:**
- Pushes a new variable frame with the specified `vars`.
- Processes the `template` within this local scope.
- Pops the variable frame after processing, reverting to the previous scope.

**Advanced Example with Nested `!With` Blocks:**
```yaml
nestedWith: !With
  vars:
    outerVar: "outer"
  template: !With
    vars:
      innerVar: "inner"
    template: !Join { items: [!Var outerVar, !Var innerVar], separator: " " }
```

**Result:**
```yaml
nestedWith: "outer inner"
```

**Usage with `!Loop`:**
```yaml
withLoopExample: !With
  vars:
    items: !Var numbers
  template: !Loop
    over: !Var items
    as: item
    template: !Var item
```

**Result:**
```yaml
withLoopExample:
  - 1
  - 2
  - 3
  - 4
  - 5
```

---

# Conclusion

This detailed specification serves as a comprehensive guide to the Emrichen language, outlining the functionality, structure, and usage of each custom tag. By understanding and utilizing these tags, users can create dynamic, flexible, and maintainable YAML configurations tailored to their specific needs.

For further information, refer to the [Emrichen Examples](test-data) and [Go Unit Tests](pkg/emrichen/) which provide exhaustive use cases and scenarios demonstrating the power and flexibility of Emrichen in various contexts.