---
Title: Emrichen Tag Reference
Slug: emrichen-reference
Short: Detailed reference for all Emrichen tags with signatures and examples
Topics:
  - emrichen
  - templating
  - yaml
  - reference
Commands:
  - none
Flags:
  - none
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This document provides a comprehensive reference for all tags available in the Go implementation of Emrichen. Each tag is detailed with its purpose, signature (arguments), and usage examples.

## Tag Categories

Emrichen tags can be broadly categorized as follows:

- **Data Access & Manipulation**: Tags for accessing variables and transforming data structures (`!Var`, `!Lookup`, `!Format`, `!Loop`, `!With`).
- **Logic & Control Flow**: Tags for conditional logic and iteration control (`!If`, `!All`, `!Any`, `!Not`, `!Op`, `!Filter`).
- **Data Structure Operations**: Tags for working with lists and dictionaries (`!Concat`, `!Merge`, `!Group`, `!Index`, `!Join`).
- **Type Operations & Encoding**: Tags for type checking and data encoding (`!Is*`, `!Base64`, `!URLEncode`, `!MD5`, `!SHA1`, `!SHA256`).
- **File Operations**: Tags for including content from external files (`!Include*`).
- **Debugging & Error Handling**: Tags for development and error management (`!Debug`, `!Error`, `!Exists`, `!Void`).

---

## `!All`

**Purpose**: Checks if all items in a sequence are truthy. (`!And` is an alias for this tag).

**Signature**:

```yaml
!All sequence
```

- `sequence`: A YAML sequence (`[]`) of values or expressions to evaluate.

**Behavior**: Returns `true` if all items evaluate to true (non-empty string, non-zero number, true boolean, non-empty list/dict). Returns `false` otherwise. Short-circuits on the first falsy value.

**Examples**:

```yaml
!Defaults
  items_all_true: [true, 1, "hello"]
  items_one_false: [true, 0, "hello"]
---
check1: !All !Var items_all_true  # Output: true
check2: !All !Var items_one_false # Output: false
check3: !All [!Exists foo, true]  # Output: false (assuming foo is not defined)
```

---

## `!Any`

**Purpose**: Checks if at least one item in a sequence is truthy. (`!Or` is an alias for this tag).

**Signature**:

```yaml
!Any sequence
```

- `sequence`: A YAML sequence (`[]`) of values or expressions to evaluate.

**Behavior**: Returns `true` if at least one item evaluates to true. Returns `false` if all items are falsy. Short-circuits on the first truthy value.

**Examples**:

```yaml
!Defaults
  items_one_true: [false, 0, "hello"]
  items_all_false: [false, 0, ""]
---
check1: !Any !Var items_one_true  # Output: true
check2: !Any !Var items_all_false # Output: false
check3: !Any [!Exists foo, true]   # Output: true
```

---

## `!Base64`

**Purpose**: Encodes a scalar value into a Base64 string.

**Signature**:

```yaml
!Base64 scalar
```

- `scalar`: The string, number, or boolean value to encode.

**Examples**:

```yaml
encoded: !Base64 "Hello Emrichen!" # Output: SGVsbG8gRW1yaWNoZW4h
secret: !Base64 !Var api_key      # Encodes the value of the 'api_key' variable
```

---

## `!Concat`

**Purpose**: Concatenates multiple sequences into a single sequence.

**Signature**:

```yaml
!Concat sequence
```

- `sequence`: A sequence where each item is itself a sequence to be concatenated.

**Examples**:

```yaml
!Defaults
list1: [1, 2]
list2: [3, 4]
---
combined: !Concat
  - !Var list1
  - [5, 6]
  - !Var list2
# Output: [1, 2, 5, 6, 3, 4]
```

---

## `!Debug`

**Purpose**: Outputs the processed value of a node to stderr for debugging, without altering the final YAML output.

**Signature**:

```yaml
!Debug any
```

- `any`: The node (scalar, sequence, mapping) whose processed value should be printed to stderr.

**Examples**:

```yaml
!Defaults { user: { name: "Debug User", id: 123 } }
---
config:
  # The value of user.name will be printed to stderr during processing
  username: !Debug !Lookup user.name
  enabled: true
# Output YAML:
# config:
#   username: Debug User
#   enabled: true
# Stderr Output:
# DEBUG: Debug User
```

---

## `!Defaults`

**Purpose**: Defines default variables accessible throughout the document via `!Var`.

**Signature**:

```yaml
!Defaults mapping
```

- `mapping`: A YAML mapping (`{}`) where keys are variable names and values are their defaults.

**Behavior**: Variables are defined in a scope. Multiple `!Defaults` tags merge, with later definitions overriding earlier ones within the same document or included file.

**Examples**:

```yaml
!Defaults
app_name: my-app
replicas: 2
---
!Defaults # Override replicas
replicas: 3
image: my-image:latest
---
deployment:
  name: !Var app_name # Output: my-app
  replicas: !Var replicas # Output: 3
  image: !Var image # Output: my-image:latest
```

---

## `!Error`

**Purpose**: Halts processing and outputs a custom error message to stderr.

**Signature**:

```yaml
!Error scalar
```

- `scalar`: The error message (string) to display.

**Examples**:

```yaml
!Defaults { required_var: null }
---
# This will halt processing if required_var is null or not defined
config: !If
  test: !IsNone !Var required_var
  then: !Error "required_var must be set!"
  else: !Var required_var
```

---

## `!Exists`

**Purpose**: Checks if a variable or JSONPath exists in the current context.

**Signature**:

```yaml
!Exists scalar
```

- `scalar`: The variable name or JSONPath string to check.

**Examples**:

```yaml
!Defaults { user: { name: "Test" } }
---
name_exists: !Exists user.name # Output: true
email_exists: !Exists user.email # Output: false
global_var_exists: !Exists app_name # Output: false (if not defined in !Defaults)
```

---

## `!Filter`

**Purpose**: Filters elements of a sequence based on a predicate.

**Signature**:

```yaml
!Filter mapping
```

- `mapping`:
  - `over`: (Required) The sequence to filter.
  - `test`: (Required) An expression evaluated for each item. If truthy, the item is kept. The item is available as `item` within the `test` expression.
  - `as`: (Optional, default: `item`) The variable name for the current item within the `test` expression.

**Examples**:

```yaml
!Defaults
  numbers: [1, 2, 3, 4, 5, 6]
  users:
    - { name: "Alice", active: true }
    - { name: "Bob", active: false }
    - { name: "Charlie", active: true }
---
even_numbers: !Filter
  over: !Var numbers
  test: !Op { a: !Lookup item, op: '%', b: 2 }, op: '==', b: 0 }
# Output: [2, 4, 6]

active_users: !Filter
  over: !Var users
  as: user
  test: !Lookup user.active
# Output:
# - name: Alice
#   active: true
# - name: Charlie
#   active: true
```

---

## `!Format`

**Purpose**: Formats a string using Go `text/template`, with access to variables and special functions.

**Signature**:

```yaml
!Format scalar
```

- `scalar`: The format string.

**Behavior**:

- Simple placeholders `{var}` are transformed to `{{.var}}`.
- Complex placeholders `{path.to[0].value}` are transformed to `{{lookup "path.to[0].value"}}`.
- Existing Go template syntax `{{...}}` is preserved.
- Provides `lookup`, `lookupAll`, and `exists` functions within the template.
- Supports standard Go `text/template` functions and custom functions added to the interpreter.

**Examples**:

```yaml
!Defaults { name: "Emrichen", version: "1.0" }
---
message: !Format "Welcome to {name} v{version}!"
# Output: Welcome to Emrichen v1.0!

advanced: !Format "User: {{.user.name}}, Exists: {{exists "user.email"}}"
# Assuming user.name is defined, user.email is not
# Output: User: SomeUser, Exists: false

go_template: !Format "{{if .isAdmin}}Admin Access{{else}}User Access{{end}}"
# Output depends on the value of the 'isAdmin' variable
```

---

## `!Group`

**Purpose**: Groups items from a sequence into a dictionary based on a key expression.

**Signature**:

```yaml
!Group mapping
```

- `mapping`:
  - `over`: (Required) The sequence to group.
  - `by`: (Required) An expression evaluated for each item to determine its group key. The item is available as `item`.
  - `as`: (Optional, default: `item`) The variable name for the current item within the `by` expression.
  - `template`: (Optional) A template applied to each item before adding it to a group. If omitted, the original item is used.

**Examples**:

```yaml
!Defaults
items:
  - { category: "A", value: 1 }
  - { category: "B", value: 2 }
  - { category: "A", value: 3 }
---
grouped_items: !Group
  over: !Var items
  by: !Lookup item.category
# Output:
# A:
#   - { category: "A", value: 1 }
#   - { category: "A", value: 3 }
# B:
#   - { category: "B", value: 2 }

grouped_values: !Group
  over: !Var items
  by: !Lookup item.category
  template: !Lookup item.value
# Output:
# A: [1, 3]
# B: [2]
```

---

## `!If`

**Purpose**: Conditional logic: returns one value if a test is true, another if false.

**Signature**:

```yaml
!If mapping
```

- `mapping`:
  - `test`: (Required) The condition to evaluate.
  - `then`: (Required) The value to return if `test` is truthy.
  - `else`: (Optional) The value to return if `test` is falsy. If omitted and `test` is falsy, `!If` returns nothing (effectively removing the key/value or list item).

**Examples**:

```yaml
!Defaults { enabled: true, mode: "prod" }
---
status: !If { test: !Var enabled, then: "Active", else: "Inactive" } # Output: Active
config: !If
  test: !Op { a: !Var mode, op: "==", b: "prod" }
  then: { db: "prod_db" }
  else: { db: "dev_db" }
# Output: { db: "prod_db" }

optional_feature: !If
  test: !Exists feature_flag
  then: "Enabled"
# Output: (nothing if feature_flag doesn't exist)
```

---

## `!Include`

**Purpose**: Includes and processes another YAML file. If the file contains multiple YAML documents, they are returned as a sequence.

**Signature**:

```yaml
!Include scalar
```

- `scalar`: The path to the YAML file to include. Paths are relative to the current file or absolute.

**Example**:
`config.yml`:

```yaml
!Include partials/common.yml
---
specific: value
```

`partials/common.yml`:

```yaml
!Defaults { common_setting: "default" }
---
key: !Var common_setting
```

**Output of processing `config.yml`**:

```yaml
key: default
---
specific: value
```

---

## `!IncludeBase64`

**Purpose**: Includes a binary file's content as a Base64 encoded string.

**Signature**:

```yaml
!IncludeBase64 scalar
```

- `scalar`: Path to the file.

**Example**:

```yaml
# Assuming icon.png exists
icon_data: !IncludeBase64 assets/icon.png
# Output: (Base64 string of icon.png content)
```

---

## `!IncludeBinary`

**Purpose**: Includes the raw content of a binary file as a string. _Use with caution, as YAML is not designed for arbitrary binary data._

**Signature**:

```yaml
!IncludeBinary scalar
```

- `scalar`: Path to the file.

**Example**:

```yaml
# Assuming data.bin exists
binary_content: !IncludeBinary data/data.bin
# Output: (Raw bytes of data.bin interpreted as a string)
```

---

## `!IncludeGlob`

**Purpose**: Includes and processes multiple YAML files matching glob patterns. Results are combined into a single sequence.

**Signature**:

```yaml
!IncludeGlob scalar | sequence
```

- `scalar`: A single glob pattern string.
- `sequence`: A sequence of glob pattern strings.

**Example**:

```yaml
# Include all YAML files in the services directory
services_config: !IncludeGlob services/*.yml
```

---

## `!IncludeText`

**Purpose**: Includes the content of a text file as a single string.

**Signature**:

```yaml
!IncludeText scalar
```

- `scalar`: Path to the text file.

**Example**:

```yaml
script_content: !IncludeText scripts/init.sh
# Output: (Content of init.sh as a string)
```

---

## `!Index`

**Purpose**: Creates a dictionary (mapping) from a sequence, using a specified key expression.

**Signature**:

```yaml
!Index mapping
```

- `mapping`:
  - `over`: (Required) The sequence to index.
  - `by`: (Required) An expression evaluated for each item to determine its key in the output dictionary. The item is available as `item`.
  - `as`: (Optional, default: `item`) The variable name for the current item within the `by` expression.
  - `template`: (Optional) A template applied to each item to determine its value in the output dictionary. If omitted, the original item is used.
  - `duplicates`: (Optional, default: `error`) How to handle duplicate keys: `error` (halt), `first` (keep first), `last` (keep last), `merge` (deep merge values), `list` (collect values in a list).

**Examples**:

```yaml
!Defaults
users:
  - { id: 1, name: "Alice" }
  - { id: 2, name: "Bob" }
  - { id: 1, name: "Alicia" } # Duplicate ID
---
users_by_id: !Index
  over: !Var users
  by: !Lookup item.id
  duplicates: last # Keep the last item for duplicate ID 1
# Output:
# 1: { id: 1, name: "Alicia" }
# 2: { id: 2, name: "Bob" }

user_names_by_id: !Index
  over: !Var users
  by: !Lookup item.id
  template: !Lookup item.name
  duplicates: list # Collect names for duplicate ID 1
# Output:
# 1: ["Alice", "Alicia"]
# 2: ["Bob"]
```

---

## `!IsBoolean`, `!IsDict`, `!IsInteger`, `!IsList`, `!IsNone`, `!IsNumber`, `!IsString`

**Purpose**: Check the type of a processed value.

**Signature**:

```yaml
!Is[Type] any
```

- `any`: The value or expression whose type needs checking after processing.

**Behavior**: Returns `true` if the processed value matches the specified type, `false` otherwise.

- `!IsDict`: Checks for YAML mappings.
- `!IsList`: Checks for YAML sequences.
- `!IsNone`: Checks for `null`.
- `!IsNumber`: Checks for integers or floats.

**Example**:

```yaml
!Defaults { count: 5, name: "Test", config: { key: "value" }, items: [1, 2] }
---
is_count_int: !IsInteger !Var count # Output: true
is_name_string: !IsString !Var name   # Output: true
is_config_dict: !IsDict !Var config  # Output: true
is_items_list: !IsList !Var items   # Output: true
is_missing_none: !IsNone !Var missing # Output: true (assuming 'missing' is undefined)
is_count_number: !IsNumber !Var count # Output: true
```

---

## `!Join`

**Purpose**: Joins elements of a sequence into a single string using a separator.

**Signature**:

```yaml
!Join mapping | sequence
```

- `mapping`:
  - `items`: (Required) The sequence whose elements are to be joined. Elements are converted to strings.
  - `separator`: (Optional, default: `""`) The string to insert between elements.
- `sequence`: A shorthand where the first element is the separator string and the second is the sequence to join.

**Examples**:

```yaml
!Defaults { parts: ["a", "b", "c"] }
---
joined_map: !Join { items: !Var parts, separator: "-" } # Output: "a-b-c"
joined_seq: !Join [", ", !Var parts] # Output: "a, b, c"
```

---

## `!Lookup`

**Purpose**: Performs a JSONPath lookup to retrieve a single value.

**Signature**:

```yaml
!Lookup scalar
```

- `scalar`: The JSONPath expression string.

**Behavior**: Returns the first value matching the path. If the path does not exist or resolves to multiple values, an error occurs (unless combined with tags like `!Exists` or `!Loop`).

**Example**:

```yaml
!Defaults { config: { users: [{ name: "Alice" }, { name: "Bob" }] } }
---
first_user_name: !Lookup config.users[0].name # Output: Alice
```

---

## `!LookupAll`

**Purpose**: Performs a JSONPath lookup, returning all matching values as a sequence.

**Signature**:

```yaml
!LookupAll scalar
```

- `scalar`: The JSONPath expression string.

**Example**:

```yaml
!Defaults { config: { users: [{ name: "Alice" }, { name: "Bob" }] } }
---
all_user_names: !LookupAll config.users[*].name # Output: ["Alice", "Bob"]
```

---

## `!Loop`

**Purpose**: Iterates over a sequence or mapping, applying a template to each element.

**Signature**:

```yaml
!Loop mapping
```

- `mapping`:
  - `over`: (Required) The sequence or mapping to iterate over.
  - `template`: (Required) The template to process for each item.
  - `as`: (Optional, default: `item`) Variable name for the current element's value.
  - `index_as`: (Optional) Variable name for the current element's index (for sequences) or key (for mappings).
  - `previous_as`: (Optional) Variable name for the previous element's value (null for the first iteration).
  - `index_start`: (Optional, default: `0`) Starting index for sequence iteration (items before this are skipped).
  - `as_documents`: (Optional, default: `false`) **Not yet supported.** If true, treat each iteration's output as a separate YAML document.

**Behavior**:

- **For Sequences**: Iterates through the list. `template` is evaluated with `as` (item value), `index_as` (numeric index), and `previous_as` (previous item value) in scope. Returns a sequence of the results.
- **For Mappings**: Iterates through key-value pairs. `template` is evaluated with `as` (item value), `index_as` (key string), and `previous_as` (previous item value) in scope. Returns a mapping with original keys and processed values.

**Examples**:

```yaml
!Defaults
  items: ["apple", "banana"]
  config: { port: 80, host: "localhost" }
---
# Sequence iteration
item_list: !Loop
  over: !Var items
  index_as: idx
  template: !Format "{idx}: {item}"
# Output: ["0: apple", "1: banana"]

# Mapping iteration
config_list: !Loop
  over: !Var config
  as: value
  index_as: key
  template: !Format "{key}={value}"
# Output: { port: "port=80", host: "host=localhost" }

# Previous item example
diffs: !Loop
  over: [10, 12, 15]
  previous_as: prev
  template: !If
    test: !IsNone !Var prev
    then: 0 # First item has no previous
    else: !Op { a: !Var item, op: '-', b: !Var prev }
# Output: [0, 2, 3]
```

---

## `!MD5`, `!SHA1`, `!SHA256`

**Purpose**: Computes the hash of a scalar value.

**Signature**:

```yaml
![HashAlgorithm] scalar
```

- `scalar`: The string value to hash.

**Example**:

```yaml
data: "Emrichen Hashing"
---
md5_hash: !MD5 !Var data
sha1_hash: !SHA1 !Var data
sha256_hash: !SHA256 !Var data
# Output: (Hexadecimal hash strings)
```

---

## `!Merge`

**Purpose**: Deep merges multiple mappings (dictionaries).

**Signature**:

```yaml
!Merge sequence
```

- `sequence`: A sequence of mappings to merge. Later mappings override keys from earlier ones. Sequences within mappings are typically replaced, not merged element-wise (standard deep merge behavior).

**Examples**:

```yaml
!Defaults
base: { a: 1, b: { x: 10 } }
override: { b: { y: 20 }, c: 3 }
---
merged: !Merge
  - !Var base
  - !Var override
  - { a: 99 } # Further override
# Output:
# a: 99
# b: { x: 10, y: 20 }
# c: 3
```

---

## `!Not`

**Purpose**: Negates a boolean value.

**Signature**:

```yaml
!Not any
```

- `any`: The value whose truthiness should be negated.

**Behavior**: Evaluates the input's truthiness and returns the opposite boolean value (`true` if input is falsy, `false` if input is truthy).

**Examples**:

```yaml
!Defaults { enabled: false }
---
is_disabled: !Not !Var enabled # Output: true
check: !Not !Exists non_existent_var # Output: true
```

---

## `!Op`

**Purpose**: Performs binary operations between two values.

**Signature**:

```yaml
!Op mapping
```

- `mapping`:
  - `op`: (Required) The operator string.
  - `a`: (Required) The left operand.
  - `b`: (Required) The right operand.

**Supported Operators**:

- **Comparison**: `=`, `==`, `===`, `≠`, `!=`, `!==`, `ne`, `<`, `lt`, `>`, `gt`, `<=`, `le`, `lte`, `>=`, `ge`, `gte`
- **Arithmetic**: `+`, `plus`, `add`, `-`, `minus`, `sub`, `subtract`, `*`, `×`, `mul`, `times`, `/`, `÷`, `div`, `divide`, `truediv`, `//`, `floordiv`, `%`, `mod`, `modulo`
- **String**: `contains`, `startswith`, `endswith`, `matches` (regex)
- **Membership**: `in`, `∈`, `not in`, `∉`

**Behavior**:

- Operands `a` and `b` are processed first.
- Arithmetic and numeric comparisons attempt to convert operands to numbers (float64). Integer operations are used if both operands are integers where appropriate (e.g., `+`, `-`, `*`, `//`, `%`).
- String operations work on string representations.
- `in`/`not in`: Checks if `a` is present in sequence `b`.

**Examples**:

```yaml
!Defaults { x: 10, y: 5, text: "hello world", pattern: "^hello" }
---
sum: !Op { a: !Var x, op: "+", b: !Var y } # Output: 15
is_greater: !Op { a: !Var x, op: ">", b: !Var y } # Output: true
floor_div: !Op { a: !Var x, op: "//", b: 3 } # Output: 3
contains_hello: !Op { a: !Var text, op: "contains", b: "hello" } # Output: true
matches_pattern: !Op { a: !Var text, op: "matches", b: !Var pattern } # Output: true
is_in_list: !Op { a: 5, op: "in", b: [1, 5, 10] } # Output: true
```

---

## `!URLEncode`

**Purpose**: Encodes a string for URL query parameters or builds a URL with query parameters.

**Signature**:

```yaml
!URLEncode scalar | mapping
```

- `scalar`: A string to be URL-encoded (query parameter encoding).
- `mapping`:
  - `url`: (Required) The base URL string.
  - `query`: (Required) A mapping representing the query parameters (key-value pairs). Values will be URL-encoded.

**Examples**:

```yaml
encoded_param: !URLEncode "a value with spaces & symbols"
# Output: a+value+with+spaces+%26+symbols

full_url: !URLEncode
  url: "http://example.com/search"
  query:
    q: "emrichen rocks"
    page: 2
# Output: http://example.com/search?page=2&q=emrichen+rocks
```

---

## `!Var`

**Purpose**: Substitutes the value of a variable defined by `!Defaults` or environment variables.

**Signature**:

```yaml
!Var scalar
```

- `scalar`: The name of the variable (or JSONPath to access nested values within a variable).

**Behavior**: Looks up the variable in the current scope (innermost first). If not found in `!Defaults`, it may fall back to environment variables depending on interpreter configuration.

**Examples**:

```yaml
!Defaults { user: { name: "Default", id: 0 }, backup_host: "backup.local" }
---
username: !Var user.name # Output: Default
host: !Var backup_host # Output: backup.local
# Assuming PWD env var is /home/user
current_dir: !Var PWD # Output: /home/user (if env fallback enabled)
```

---

## `!Void`

**Purpose**: Represents nothing. Useful for conditionally removing items or keys.

**Signature**:

```yaml
!Void
```

**Behavior**: When a tag like `!If` returns `!Void`, the corresponding key-value pair in a mapping or the item in a sequence is completely removed from the output.

**Examples**:

```yaml
!Defaults { enable_feature_x: false }
---
config:
  feature_x:
    !If { test: !Var enable_feature_x, then: { enabled: true }, else: !Void  }
  always_present: true
# Output:
# config:
#   always_present: true
# (feature_x key is removed because !If returned !Void)
```

---

## `!With`

**Purpose**: Defines a local scope with temporary variables for processing a template.

**Signature**:

```yaml
!With mapping
```

- `mapping`:
  - `vars`: (Required) A mapping defining local variables.
  - `template`: (Required) The template to process within the local scope.

**Behavior**: The variables defined in `vars` are available only when processing `template`. They temporarily override any variables with the same name from outer scopes.

**Examples**:

```yaml
!Defaults { global_msg: "Hello" }
---
result: !With
  vars:
    name: "Local Scope"
    global_msg: "Overridden" # Overrides the default
  template:
    message: !Format "{global_msg}, {name}!"
    original: !Var global_msg # Accesses the original default outside the 'template'
# Output:
# result:
#   message: "Overridden, Local Scope!"
#   original: "Hello"
```
