package emrichen

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"io"
	"testing"
	"time"
)

type testCase struct {
	name               string
	inputYAML          string
	expected           string
	initVars           map[string]interface{} // Adding a new field for initial variable bindings
	expectError        bool
	expectErrorMessage string
	expectPanic        bool
}

func runTests(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ei, err := NewInterpreter(WithVars(tc.initVars))
			require.NoError(t, err)

			decoder := yaml.NewDecoder(bytes.NewReader([]byte(tc.inputYAML)))

			hadError := false
			var resultNode *yaml.Node
			if tc.expectPanic {
				_, err = expectPanic(t, func() (*yaml.Node, error) {
					for {
						inputNode := yaml.Node{}
						// Parse input YAML
						err2 := decoder.Decode(ei.CreateDecoder(&inputNode))
						require.NoError(t, err2)
					}
				})
				require.Error(t, err)
				require.Equal(t, "paniced", err.Error())
				return
			}

			for {
				inputNode := yaml.Node{}
				// Parse input YAML
				err2 := decoder.Decode(ei.CreateDecoder(&inputNode))
				if err2 == io.EOF {
					break
				}
				err = err2
				if err != nil {
					hadError = true
					break
				}

				// empty document node after defaults
				if inputNode.Kind == 0 {
					continue
				}
				resultNode, err = ei.Process(&inputNode)
				if err != nil {
					hadError = true
					break
				}
			}

			if hadError {
				if tc.expectError {
					require.Error(t, err, "Expected an error but got none")
					if tc.expectErrorMessage != "" {
						assert.Equal(t, tc.expectErrorMessage, err.Error())
					}
				} else {
					require.NoError(t, err, "Unexpected error encountered", err)
				}
				return
			} else {
				require.NoError(t, err, "Unexpected error encountered", err)
			}

			expectedNode := yaml.Node{}
			err = yaml.Unmarshal([]byte(tc.expected), &expectedNode)
			require.NoError(t, err)

			expected_ := convertNodeToInterface(&expectedNode)
			actual_ := convertNodeToInterface(resultNode)

			assert.Equal(t, expected_, actual_)
		})
	}
}

func expectPanic(t *testing.T, f func() (*yaml.Node, error)) (*yaml.Node, error) {
	didPanic := false

	_, _ = func() (*yaml.Node, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("paniced")
				didPanic = true
			}
		}()
		return f()
	}()

	if !didPanic {
		t.Errorf("Expected a panic to occur, but none did")
	}

	return nil, errors.New("paniced")
}

func TestValueToNode_TimeSupport(t *testing.T) {
	// Test time.Time conversion
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	
	node, err := ValueToNode(testTime)
	require.NoError(t, err)
	
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "2023-12-25T15:30:45Z", node.Value)
}

func TestValueToNode_TimeInVar(t *testing.T) {
	// Test time.Time in !Var context
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	
	tests := []testCase{
		{
			name: "time in var",
			inputYAML: `
timestamp: !Var myTime
`,
			initVars: map[string]interface{}{
				"myTime": testTime,
			},
			expected: `
timestamp: "2024-01-01T12:00:00Z"
`,
		},
	}
	
	runTests(t, tests)
}
