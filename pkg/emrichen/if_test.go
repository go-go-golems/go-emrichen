package emrichen

import (
	"testing"
)

func TestEmrichenIfAndFilterTags(t *testing.T) {
	tests := []testCase{
		{
			name: "Basic boolean condition - true",
			inputYAML: `!If
  test: true
  then: 'Yes'
  else: 'No'`,
			expected: "'Yes'",
		},
		{
			name: "Basic boolean condition - false",
			inputYAML: `!If
  test: false
  then: 'Yes'
  else: 'No'`,
			expected: "'No'",
		},
		{
			name: "Nested !If statements",
			inputYAML: `!If
  test: true
  then: !If {test: false, then: "Inner Yes", else: "Inner No"}
  else: 'Outer No'`,
			expected: "'Inner No'",
		},
		{
			name: "Variable substitution in conditions",
			inputYAML: `!If
  test: !Var condition
  then: 'Variable Yes'
  else: 'Variable No'`,
			expected: "'Variable Yes'",
			initVars: map[string]interface{}{
				"condition": true,
			},
		},

		{
			name: "If tag true condition",
			inputYAML: `!If
  test: !Op
    a: 10
    op: ">"
    b: 5
  then: "True Condition"
  else: "False Condition"`,
			expected: "\"True Condition\"",
		},
		{
			name: "If tag false condition",
			inputYAML: `!If
  test: !Op
    a: 3
    op: ">"
    b: 5
  then: "True Condition"
  else: "False Condition"`,
			expected: "\"False Condition\"",
		},
		{
			name: "Type-checking with string in condition",
			inputYAML: `!If
  test: 'true'  # String 'true', not boolean
  then: 'String Yes'
  else: 'String No'`,
			expected: "'String Yes'",
		},
		{
			name: "Type-checking with number in condition",
			inputYAML: `!If
  test: 1  # Number 1
  then: 'Number Yes'
  else: 'Number No'`,
			expected: "'Number Yes'",
		},
		// Test 9: Conditional Evaluation with `null` and Empty Values
		{
			name: "Evaluation with null condition",
			inputYAML: `!If
  test: null
  then: 'Null Yes'
  else: 'Null No'`,
			expected: "'Null No'",
		},
		{
			name: "Evaluation with empty string condition",
			inputYAML: `!If
  test: ''
  then: 'Empty Yes'
  else: 'Empty No'`,
			expected: "'Empty No'",
		},
		{
			name: "Evaluation with empty list condition",
			inputYAML: `!If
  test: []
  then: 'Empty List Yes'
  else: 'Empty List No'`,
			expected: "'Empty List No'",
		},
		{
			name: "Omitting 'then' branch - true condition",
			inputYAML: `!If
  test: true
  else: 'No'`,
			expected: "null", // Expecting null or equivalent when 'then' is missing and condition is true
		},
		{
			name: "Omitting 'else' branch - false condition",
			inputYAML: `!If
  test: false
  then: 'Yes'`,
			expected: "null", // Expecting null or equivalent when 'else' is missing and condition is false
		},
		{
			name: "Omitting both 'then' and 'else' - true condition",
			inputYAML: `!If
  test: true`,
			expected: "null", // Expecting null or equivalent when both branches are missing
		},
		{
			name: "Omitting both 'then' and 'else' - false condition",
			inputYAML: `!If
  test: false`,
			expected: "null", // Expecting null or equivalent when both branches are missing
		},
		{
			name: "Variable substitution in condition with omitted 'then'",
			inputYAML: `!If
  test: !Var condition
  else: 'Variable No'`,
			expected: "null",
			initVars: map[string]interface{}{
				"condition": true,
			},
		},
		{
			name: "Variable substitution in condition with omitted 'else'",
			inputYAML: `!If
  test: !Var condition
  then: 'Variable Yes'`,
			expected: "null",
			initVars: map[string]interface{}{
				"condition": false,
			},
		},
		{
			name: "Nested Op test",
			inputYAML: `!If
  test: !Op
    a: 10
    op: ">"
    b: 5
  then: "True Condition"
  else: "False Condition"`,
			expected: "\"True Condition\"",
		},
		{
			name: "Nested All in test",
			inputYAML: `!If
  test: !All
    - !Op
      a: 10
      op: ">"
      b: 5
    - !Op
      a: [ ]
      op: '!='
      b: [ ]
  then: "True Condition"
  else: "False Condition"`,
			expected: "\"False Condition\"",
		},
	}

	runTests(t, tests)
}
