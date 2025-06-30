package emrichen

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"io"
	"testing"
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

// Test types for interface testing
type testStringer struct {
	value string
}

func (ts testStringer) String() string {
	return ts.value
}

type testTextMarshaler struct {
	value string
}

func (tm testTextMarshaler) MarshalText() ([]byte, error) {
	return []byte("marshaled:" + tm.value), nil
}

type testTextMarshalerError struct{}

func (tm testTextMarshalerError) MarshalText() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}

func TestValueToNode_PointerSupport(t *testing.T) {
	// Test nil pointer
	var nilPtr *string
	node, err := ValueToNode(nilPtr)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!null", node.Tag)

	// Test pointer to string
	str := "hello"
	node, err = ValueToNode(&str)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "hello", node.Value)

	// Test pointer to int
	num := 42
	node, err = ValueToNode(&num)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!int", node.Tag)
	assert.Equal(t, "42", node.Value)

	// Test pointer to time.Time
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	node, err = ValueToNode(&testTime)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "2024-01-01T12:00:00Z", node.Value)
}

func TestValueToNode_DurationSupport(t *testing.T) {
	duration := 5*time.Minute + 30*time.Second

	node, err := ValueToNode(duration)
	require.NoError(t, err)

	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "5m30s", node.Value)
}

func TestValueToNode_NetIPSupport(t *testing.T) {
	// Test IPv4
	ipv4 := net.ParseIP("192.168.1.1")
	node, err := ValueToNode(ipv4)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "192.168.1.1", node.Value)

	// Test IPv6
	ipv6 := net.ParseIP("2001:db8::1")
	node, err = ValueToNode(ipv6)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "2001:db8::1", node.Value)
}

func TestValueToNode_URLSupport(t *testing.T) {
	// Test url.URL value
	u, err := url.Parse("https://example.com/path?query=value")
	require.NoError(t, err)

	node, err := ValueToNode(*u)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "https://example.com/path?query=value", node.Value)

	// Test *url.URL pointer
	node, err = ValueToNode(u)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "https://example.com/path?query=value", node.Value)

	// Test nil *url.URL
	var nilURL *url.URL
	node, err = ValueToNode(nilURL)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!null", node.Tag)
}

func TestValueToNode_UUIDSupport(t *testing.T) {
	testUUID := uuid.New()

	node, err := ValueToNode(testUUID)
	require.NoError(t, err)

	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, testUUID.String(), node.Value)
}

func TestValueToNode_TextMarshalerSupport(t *testing.T) {
	// Test successful TextMarshaler
	tm := testTextMarshaler{value: "test"}
	node, err := ValueToNode(tm)
	require.NoError(t, err)
	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "marshaled:test", node.Value)

	// Test TextMarshaler error
	tmErr := testTextMarshalerError{}
	_, err = ValueToNode(tmErr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal text")
}

func TestValueToNode_StringerSupport(t *testing.T) {
	stringer := testStringer{value: "custom string"}

	node, err := ValueToNode(stringer)
	require.NoError(t, err)

	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	assert.Equal(t, "custom string", node.Value)
}

// Type that implements both interfaces to test priority
type bothInterfaces struct{}

func (bi bothInterfaces) MarshalText() ([]byte, error) {
	return []byte("from-text-marshaler"), nil
}

func (bi bothInterfaces) String() string {
	return "from-stringer"
}

func TestValueToNode_InterfacePriority(t *testing.T) {
	// Create a type that implements both TextMarshaler and Stringer
	// TextMarshaler should take priority
	bi := bothInterfaces{}
	node, err := ValueToNode(bi)
	require.NoError(t, err)

	assert.Equal(t, yaml.ScalarNode, node.Kind)
	assert.Equal(t, "!!str", node.Tag)
	// Should use TextMarshaler, not Stringer
	assert.Equal(t, "from-text-marshaler", node.Value)
}

func TestValueToNode_InVarContext(t *testing.T) {
	// Test various types in !Var context
	duration := 2 * time.Hour
	ipAddr := net.ParseIP("10.0.0.1")
	testURL, _ := url.Parse("http://localhost:8080")
	testUUID := uuid.New()
	stringer := testStringer{value: "var-test"}

	tests := []testCase{
		{
			name: "duration in var",
			inputYAML: `
timeout: !Var duration
`,
			initVars: map[string]interface{}{
				"duration": duration,
			},
			expected: `
timeout: "2h0m0s"
`,
		},
		{
			name: "ip address in var",
			inputYAML: `
server_ip: !Var ip
`,
			initVars: map[string]interface{}{
				"ip": ipAddr,
			},
			expected: `
server_ip: "10.0.0.1"
`,
		},
		{
			name: "url in var",
			inputYAML: `
endpoint: !Var url
`,
			initVars: map[string]interface{}{
				"url": testURL,
			},
			expected: `
endpoint: "http://localhost:8080"
`,
		},
		{
			name: "uuid in var",
			inputYAML: `
id: !Var uuid
`,
			initVars: map[string]interface{}{
				"uuid": testUUID,
			},
			expected: fmt.Sprintf(`
id: "%s"
`, testUUID.String()),
		},
		{
			name: "stringer in var",
			inputYAML: `
message: !Var stringer
`,
			initVars: map[string]interface{}{
				"stringer": stringer,
			},
			expected: `
message: "var-test"
`,
		},
		{
			name: "pointer in var",
			inputYAML: `
value: !Var ptr
`,
			initVars: map[string]interface{}{
				"ptr": &duration,
			},
			expected: `
value: "2h0m0s"
`,
		},
	}

	runTests(t, tests)
}
