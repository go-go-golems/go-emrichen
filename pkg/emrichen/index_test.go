package emrichen

import (
	"testing"
)

func TestIndexTag(t *testing.T) {
	tests := []testCase{
		{
			name: "Basic Indexing",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.9
    - name: John
      score: 9.8
  as: foo
  by: !Lookup foo.name
  duplicates: ignore
  template: !Lookup foo.score`,
			expected: `{"manifold": 7.8, "John": 9.8}`,
		},
		{
			name: "Basic Indexing (default as)",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.8
  by: !Lookup item.name
  template: !Lookup item.score`,
			expected: `{"manifold": 7.8, "John": 9.8}`,
		},
		{
			name: "Basic Indexing (duplicate error)",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.9
    - name: John
      score: 9.8
  by: !Lookup item.name
  template: !Lookup item.score`,
			expectError:        true,
			expectErrorMessage: "Duplicate key encountered: John",
		},
		{
			name: "Indexing Without Template",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.9
    - name: John
      score: 9.8
  as: foo
  by: !Lookup foo.name
  duplicates: ignore`,
			expected: `{"manifold": {"name": "manifold", "score": 7.8}, "John": {"name": "John", "score": 9.8}}`,
		},
		{
			name: "Index with Result As",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.9
    - name: John
      score: 9.8
  as: foo
  template:
    NAME: !Lookup foo.name
    SCORE: !Lookup foo.score
  result_as: result
  by: !Lookup result.NAME
  duplicates: ignore`,
			expected: `{"manifold": {"NAME": "manifold", "SCORE": 7.8}, "John": {"NAME": "John", "SCORE": 9.8}}`,
		},
		{
			name: "Duplicate Keys Error",
			inputYAML: `!Index
  over:
    - name: manifold
      score: 7.8
    - name: John
      score: 9.9
    - name: John
      score: 9.8
  as: item
  by: !Lookup item.name
  duplicates: error
  template: !Lookup item.score`,
			expectError:        true,
			expectErrorMessage: "Duplicate key encountered: John",
		},
		{
			name: "Invalid Key Path",
			inputYAML: `!Index
  over:
    - name: Alpha
      value: 1
    - name: Beta
      value: 2
  as: item
  by: !Lookup item.nonexistent
  template: !Lookup item.value`,
			expectError:        true,
			expectErrorMessage: "nonexistent is not found",
		},
	}

	runTests(t, tests)
}
