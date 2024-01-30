package emrichen

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
	"text/template"
)

func (ei *Interpreter) handleFormat(node *yaml.Node) (*yaml.Node, error) {
	formatString := node.Value

	ret, err := ei.renderFormatString(formatString)
	if err != nil {
		return nil, err
	}

	return ValueToNode(ret)
}

func (ei *Interpreter) renderFormatString(formatString string) (string, error) {
	// Transform the template to the Go template format.
	formatString, err := transformTemplate(formatString)
	if err != nil {
		return "", fmt.Errorf("error transforming template: %v", err)
	}

	tmpl := template.New("format")
	for _, funcMap := range ei.funcmaps {
		tmpl = tmpl.Funcs(funcMap)
	}

	tmpl = tmpl.Funcs(
		map[string]interface{}{
			"lookup": func(path string) interface{} {
				v, err := ei.LookupFirst(path)
				if err != nil {
					return nil
				}
				v_, _ := NodeToInterface(v)
				return v_
			},
			"lookupAll": func(path string) []interface{} {
				v, err := ei.LookupAll(path)
				if err != nil {
					return nil
				}
				v_, _ := NodeToSlice(v)
				return v_
			},
			"exists": func(path string) bool {
				_, err := ei.LookupFirst(path)
				return err == nil
			},
		},
	)
	tmpl, err = tmpl.Parse(formatString)
	if err != nil {
		return "", fmt.Errorf("error parsing format string: %v", err)
	}

	var formatted bytes.Buffer
	frame := ei.env.GetCurrentFrame()
	vars := map[string]interface{}{}
	if frame.Variables != nil {
		vars = frame.Variables
	}

	if err := tmpl.Execute(&formatted, vars); err != nil {
		return "", fmt.Errorf("error executing format template: %v", err)
	}

	return formatted.String(), nil
}

// transformTemplate converts templates from the old Emrichen format to Go template format.
//
// This function is designed to support the transition from the old Emrichen template format to the Go template format.
// It performs the following transformations:
//   - Converts simple variable expressions from `{variable}` to `{{.variable}}`.
//     This is for basic variable replacements.
//   - Transforms complex expressions, specifically JSONPath lookup expressions,
//     from `{expression}` to `{{lookup "expression"}}`.
//     This allows for more complex data retrieval and manipulation within the template.
//   - Preserves existing Go template expressions that are already in the `{{...}}` format.
//     This ensures that any Go template syntax present in the input is left unchanged.
//
// Note: This function assumes that the input string is a valid template in the old Emrichen format.
// It does not validate the correctness of the Emrichen or Go template syntax.
func transformTemplate(input string) (string, error) {
	varPattern := regexp.MustCompile(`(?m)(\{\{.*?\}\}|\{([^\{\}]+)\})`)

	transformFunc := func(match string) string {
		// Check if the match is already in the Go template format.
		if strings.HasPrefix(match, "{{") && strings.HasSuffix(match, "}}") {
			return match // Return the match as is, since it's already in the correct format.
		}

		// Remove the curly braces to get the variable name or expression.
		expression := match[1 : len(match)-1]

		// Check if the expression is simple (variable name) or complex.
		if strings.ContainsAny(expression, " .,;+-*/&|<>=()[]{}") {
			// For complex expressions, use the lookup function.
			return fmt.Sprintf(`{{lookup "%s"}}`, expression)
		}
		// For simple variable names, use the dot notation.
		return fmt.Sprintf("{{.%s}}", expression)
	}

	result := varPattern.ReplaceAllStringFunc(input, transformFunc)

	return result, nil
}
