# Go-Emrichen: Documentation and Tutorial

## Introduction

Go-Emrichen is a powerful templating engine designed for generating YAML configurations with ease and precision. It's a Go implementation of the original Python Emrichen, bringing the same flexibility and robustness to Go developers. Go-Emrichen allows you to dynamically generate configuration files for a wide range of applications, including Kubernetes deployments, configuration management, and more.

## Installation

To use go-emrichen in your Go project, you first need to install it. Run the following command:

```
go get github.com/go-go-golems/go-emrichen
```

## Basic Usage

### Importing the Library

To use go-emrichen in your Go program, import it as follows:

```go
import (
    "github.com/go-go-golems/go-emrichen/pkg/emrichen"
    "gopkg.in/yaml.v3"
)
```

(make sure you have the yaml v3 import if you are using go modules.)

### Creating an Interpreter

The core of go-emrichen is the `Interpreter`. To create a new interpreter:

```go
interpreter, err := emrichen.NewInterpreter()
if err != nil {
    // Handle error
}
```

### Processing YAML

To process a YAML file with go-emrichen:

```go
func processFile(interpreter *emrichen.Interpreter, filePath string, w io.Writer) error {
    f, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    decoder := yaml.NewDecoder(f)

    for {
        var document interface{}
        err = decoder.Decode(interpreter.CreateDecoder(&document))
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        // Skip empty documents
        if document == nil {
            continue
        }

        processedYAML, err := yaml.Marshal(&document)
        if err != nil {
            return err
        }

        _, err = w.Write(processedYAML)
        if err != nil {
            return err
        }
    }

    return nil
}
```

This function processes a YAML file, applying the Emrichen transformations, and writes the result to the provided writer.

## Advanced Usage

### Adding Variables

You can add variables to the Interpreter to use in your YAML templates:

```go
vars := map[string]interface{}{
    "environment": "production",
    "replicas": 3,
}

interpreter, err := emrichen.NewInterpreter(emrichen.WithVars(vars))
if err != nil {
    // Handle error
}
```

### Custom Functions

Go-Emrichen allows you to add custom functions to extend its capabilities:

```go
import "text/template"

customFuncs := template.FuncMap{
    "uppercase": strings.ToUpper,
    "lowercase": strings.ToLower,
}

interpreter, err := emrichen.NewInterpreter(emrichen.WithFuncMap(customFuncs))
if err != nil {
    // Handle error
}
```

### Additional Tags

You can add custom tags to enhance the functionality of go-emrichen:

```go
customTags := map[string]func(node *yaml.Node) (*yaml.Node, error){
    "!CustomTag": func(node *yaml.Node) (*yaml.Node, error) {
        // Implement custom tag logic
        return node, nil
    },
}

interpreter, err := emrichen.NewInterpreter(emrichen.WithAdditionalTags(customTags))
if err != nil {
    // Handle error
}
```

## Tutorial: Using Go-Emrichen in a Kubernetes Deployment

Let's create a simple program that uses go-emrichen to generate a Kubernetes deployment configuration.

1. Create a new Go file named `main.go`:

```go
package main

import (
    "fmt"
    "os"

    "github.com/go-go-golems/go-emrichen/pkg/emrichen"
    "gopkg.in/yaml.v3"
)

func main() {
    // Create an Interpreter with variables
    vars := map[string]interface{}{
        "APP_NAME": "myapp",
        "REPLICAS": 3,
        "IMAGE": "myapp:latest",
    }

    interpreter, err := emrichen.NewInterpreter(emrichen.WithVars(vars))
    if err != nil {
        fmt.Println("Error creating interpreter:", err)
        os.Exit(1)
    }

    // YAML template
    yamlTemplate := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .APP_NAME }}
spec:
  replicas: {{ .REPLICAS }}
  selector:
    matchLabels:
      app: {{ .APP_NAME }}
  template:
    metadata:
      labels:
        app: {{ .APP_NAME }}
    spec:
      containers:
      - name: {{ .APP_NAME }}
        image: {{ .IMAGE }}
`

    // Process the template
    var result interface{}
    err = yaml.Unmarshal([]byte(yamlTemplate), interpreter.CreateDecoder(&result))
    if err != nil {
        fmt.Println("Error processing template:", err)
        os.Exit(1)
    }

    // Marshal the result back to YAML
    output, err := yaml.Marshal(result)
    if err != nil {
        fmt.Println("Error marshaling result:", err)
        os.Exit(1)
    }

    // Print the result
    fmt.Println(string(output))
}
```

2. Run the program:

```
go run main.go
```

This will output a Kubernetes deployment YAML with the variables interpolated.

## Conclusion

Go-Emrichen provides a powerful way to template and generate YAML configurations in Go programs. By leveraging its features like variable interpolation, custom functions, and additional tags, you can create flexible and dynamic configuration generation systems.

For more advanced usage and a complete list of available tags, refer to the official documentation and examples in the go-emrichen repository:

You can find detailed documentation for each tag in the [doc section](pkg/doc/examples)
as well as an exhaustive list of examples in [the examples yamls](test-data)
and [in the go unit tests](pkg/emrichen/).

