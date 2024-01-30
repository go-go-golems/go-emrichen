package emrichen

import (
	"testing"
)

func TestEmrichenMergeTag(t *testing.T) {
	tests := []testCase{
		{
			name: "Basic Merge",
			inputYAML: `!Merge
  - a: 1
  - b: 2`,
			expected: "{a: 1, b: 2}",
		},
		{
			name: "Overlapping Keys",
			inputYAML: `!Merge
  - a: 1
  - a: 2`,
			expected: "{a: 2}",
		},
		{
			name: "Nested Merge",
			inputYAML: `!Merge
  - a: {x: 1, y: 2}
  - a: {y: 3, z: 4}`,
			expected: "{a: {y: 3, z: 4}}",
		},
		{
			name: "Empty Dictionaries",
			inputYAML: `!Merge
  - {}
  - {a: 1}`,
			expected: "{a: 1}",
		},
		{
			name: "Null Values",
			inputYAML: `!Merge
  - a: 1
  - a: null`,
			expected: "{a: null}",
		},
		{
			name: "List Values",
			inputYAML: `!Merge
  - a: [1, 2]
  - a: [3, 4]`,
			expected: "{a: [3, 4]}",
		},
		{
			name: "Complex Scenario",
			inputYAML: `!Merge
  - a: {x: 1, y: [1, 2]}
  - b: 2
  - a: {y: [3], z: 3}`,
			expected: "{a: {y: [3], z: 3}, b: 2}",
		},
		{
			name: "Merge With Variable Substitution",
			inputYAML: `!Merge
  - a: !Var valueA
  - b: !Var valueB`,
			initVars: map[string]interface{}{
				"valueA": 1,
				"valueB": 2,
			},
			expected: "{a: 1, b: 2}",
		},
		{
			name: "Merge With Variables Referring to Dictionaries",
			inputYAML: `!Merge
  - !Var dict1
  - !Var dict2`,
			initVars: map[string]interface{}{
				"dict1": map[string]interface{}{"a": 1, "b": 2},
				"dict2": map[string]interface{}{"b": 3, "c": 4},
			},
			expected: "{a: 1, b: 3, c: 4}",
		},
	}

	runTests(t, tests)
}
