package emrichen

import (
	"testing"
)

func TestTransformTemplate(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple variable",
			input:    "Hello, {name}!",
			expected: "Hello, {{.name}}!",
		},
		{
			name:     "Complex expression (JSONPath lookup)",
			input:    "User: {$.users[?(@.name=='John')]}",
			expected: `User: {{lookup "$.users[?(@.name=='John')]"}}`,
		},
		{
			name:     "Preserve Go template expression",
			input:    "Date: {{.Date}}",
			expected: "Date: {{.Date}}",
		},
		{
			name:     "Mixed content",
			input:    "{greeting}, {name}! Your balance is {$.accounts[?(@.id==123)].balance}.",
			expected: "{{.greeting}}, {{.name}}! Your balance is {{lookup \"$.accounts[?(@.id==123)].balance\"}}.",
		},
		{
			name:     "No transformation needed",
			input:    "Just a regular string.",
			expected: "Just a regular string.",
		},
		{
			name: "Multiline with mixed content",
			input: `Dear {name},
Your account balance is {$.accounts[?(@.id==123)].balance}.
Best regards, {{.CompanyName}}`,
			expected: `Dear {{.name}},
Your account balance is {{lookup "$.accounts[?(@.id==123)].balance"}}.
Best regards, {{.CompanyName}}`,
		},
		{
			name: "Mixed preserve and update",
			input: `Hello, {name}!
Your task "{{.TaskName}}" is due on {{.DueDate}}.
Details: {$.details[?(@.important==true)]}`,
			expected: `Hello, {{.name}}!
Your task "{{.TaskName}}" is due on {{.DueDate}}.
Details: {{lookup "$.details[?(@.important==true)]"}}`,
		},
		{
			name: "Multiline with complex expressions",
			input: `{user}
- Email: {$.user.email}
- Phone: {$.user.phone}
{{.Footer}}`,
			expected: `{{.user}}
- Email: {{lookup "$.user.email"}}
- Phone: {{lookup "$.user.phone"}}
{{.Footer}}`,
		},
	}

	// Iterate over each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := transformTemplate(tt.input)
			if err != nil {
				t.Errorf("transformTemplate(%q) returned an error: %v", tt.input, err)
			}
			if actual != tt.expected {
				t.Errorf("transformTemplate(%q) = %q, want %q", tt.input, actual, tt.expected)
			}
		})
	}
}
