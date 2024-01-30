---
Title: "Error and Debug Helpers"
Slug: error-debug-helpers
Short: |
  ```
  !Error "Error message"
  !Debug,Var variableName
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
# Error and Debug Helpers

The `!Error` and `!Debug` tags are tools in Emrichen for error handling and debugging templates. The `!Error`
tag interrupts template processing and outputs a custom error message, making it useful for validating template
conditions. The `!Debug` tag prints the value of its argument to stderr, aiding in debugging complex templates.

- `!Error`: Accepts a single argument, a string that is the error message to be displayed.
- `!Debug`: Can be combined with other tags (e.g., `!Var`) to output their processed values for debugging purposes.

## Examples

### Example 1: Using `!Error` for Conditional Validation

This example demonstrates using `!Error` to enforce that a variable must be set.

```yaml
!Defaults
requiredVar: null
---
conditionalError: !If
  test: !IsNone,Var requiredVar
  then: !Error "requiredVar must not be null"
  else: !Var requiredVar
```

### Example 2: Debugging Variable Values

Showcases how to use `!Debug` to print the value of a variable.

```yaml
!Defaults
debugVar: "Debugging with Emrichen"
---
debuggedValue: !Debug,Var debugVar
```

### Example 3: Combining `!Error` and `!Debug`

Combines `!Debug` and `!Error` to first debug a variable's value, then conditionally throw an error if another variable is not set.

```yaml
!Defaults
debugExampleVar: "Debug Example"
errorExampleVar: null
---
steps:
  - debugStep: !Debug,Var debugExampleVar
  - errorStep: !If
      test: !IsNone,Var errorExampleVar
      then: !Error "errorExampleVar must not be null"
      else: !Var errorExampleVar
```
