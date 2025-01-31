---
Title: Emrichen
Slug: emrichen
Topics:
  - help
Commands:
  - help
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: GeneralTopic
---

# Emrichen: A YAML Templating Engine

Emrichen is a powerful templating engine specifically designed for generating and manipulating YAML configurations. Unlike traditional text-based templating systems, Emrichen understands YAML's structure, enabling robust and type-safe template processing while avoiding common pitfalls like indentation errors and type mismatches.

## Getting Help

Emrichen provides comprehensive built-in help features that you can access through the command line:

- `emrichen help --all`: Display all available documentation topics
- `emrichen help --examples`: Show practical examples for common use cases
- `emrichen help tag-<name>`: Get detailed information about a specific tag (e.g., `emrichen help tag-op`)

## Core Features

### Variable Management
- **Variable Definition**: Use `!Defaults` to define default values
- **Variable Substitution**: Access variables using `!Var`
- **Scoped Variables**: Create local variable scopes with `!With`

Variables in Emrichen serve as the foundation for template reusability and configuration. The `!Defaults` tag allows you to define a set of default values at the beginning of your template, which can be overridden when needed. These variables can be simple scalar values, lists, or complex nested structures.

The `!With` tag creates a new scope for variables, allowing you to work with temporary values without affecting the global scope. This is particularly useful when working with loops or when you need to transform data temporarily.

### Control Flow
- **Conditional Logic**: Use `!If` for conditional processing
- **Loops**: Iterate over collections with `!Loop`
- **Error Handling**: Control error behavior with `!Error` and debug with `!Debug`

Control flow in Emrichen is designed to be both powerful and intuitive. The `!If` tag supports complex conditions through the `!Op` tag, allowing you to create sophisticated conditional logic. For example, you can combine multiple conditions using logical operators or perform numerical comparisons.

The `!Loop` tag is a versatile tool for iterating over collections. It provides features like index tracking, access to the previous item, and the ability to filter items during iteration. This makes it ideal for generating repetitive configurations while maintaining flexibility.

### Data Manipulation
- **List Operations**: Concatenate lists with `!Concat`, join elements with `!Join`
- **Dictionary Operations**: Merge dictionaries with `!Merge`
- **Filtering**: Filter collections using `!Filter`
- **Grouping**: Group data with `!Group` and `!Index`
- **String Formatting**: Format strings using `!Format` with Go template syntax

Data manipulation in Emrichen goes beyond simple variable substitution. The `!Merge` tag allows you to combine multiple dictionaries, with later values taking precedence. This is particularly useful when working with configuration overlays or environment-specific settings.

The `!Format` tag supports Go's powerful template syntax, enabling complex string formatting operations. You can access nested data structures, apply formatting functions, and even use conditional logic within your format strings.

### Advanced Features
- **File Inclusion**: Include other YAML files with `!Include`
- **URL Handling**: Encode URLs and query parameters with `!URLEncode`
- **Type Checking**: Validate types with `!IsString`, `!IsNumber`, etc.
- **Lookup Operations**: Access nested data using `!Lookup`

File inclusion through the `!Include` tag enables modular template design. You can break down complex configurations into smaller, reusable components and include them as needed. This supports both YAML and binary files, with options for base64 encoding when required.

The lookup system, implemented through the `!Lookup` tag, provides a powerful way to access nested data structures using JSONPath-like syntax. This makes it easy to work with complex data structures while maintaining readable templates.

## Basic Usage Examples

### Variable Substitution
```yaml
!Defaults
name: John
greeting: Hello
---
message: !Format "{{.greeting}}, {{.name}}!"
```

In this example, we define two variables and use them in a format string. The `!Format` tag uses Go's template syntax, allowing for more complex formatting operations when needed.

### Conditional Processing
```yaml
!Defaults
is_production: true
---
environment: !If
  test: !Var is_production
  then: production
  else: development
```

The `!If` tag demonstrates Emrichen's conditional processing capabilities. You can use complex conditions by combining the `!Op` tag with various operators.

### Looping and Data Transformation
```yaml
!Defaults
ports: [80, 443, 8080]
---
container_ports: !Loop
  over: !Var ports
  template:
    port: !Var item
    protocol: TCP
```

This example shows how to transform a simple list into a more complex structure using the `!Loop` tag. The `item` variable automatically contains the current item being processed.

### Working with Complex Data
```yaml
!Defaults
services:
  - name: web
    port: 80
  - name: api
    port: 8080
---
deployments: !Loop
  over: !Var services
  template:
    name: !Format "{{.item.name}}-service"
    port: !Var item.port
```

Here we demonstrate working with structured data, showing how to access nested properties and transform them into a new format.

## Best Practices

1. **Variable Organization**
   - Keep default values organized in a single `!Defaults` section
   - Use meaningful variable names
   - Consider using `!With` for local variable scoping

Organizing variables effectively is crucial for maintaining readable and maintainable templates. The `!Defaults` section should be placed at the beginning of your template, clearly documenting all available configuration options.

2. **Error Handling**
   - Use `!Debug` for troubleshooting
   - Implement proper error handling with `!Error`
   - Validate data types when necessary

Error handling in Emrichen is proactive. The `!Debug` tag helps you inspect variables during template processing, while `!Error` allows you to fail fast when invalid conditions are detected. Type validation tags help ensure data consistency.

3. **Code Structure**
   - Break down complex templates into smaller, reusable components
   - Use `!Include` to maintain modular templates
   - Comment your YAML files for clarity

Modular template design is essential for managing complex configurations. Use the `!Include` tag to split your templates into logical components, and maintain a clear directory structure for included files.

4. **Performance Considerations**
   - Use `!Filter` before `!Loop` to reduce iterations
   - Avoid deeply nested template structures
   - Consider caching frequently used values in variables

Performance optimization in Emrichen often involves reducing the number of operations performed. The `!Filter` tag can significantly improve performance by reducing the dataset before processing.

## Common Use Cases

1. **Kubernetes Configurations**
   - Generate deployment manifests
   - Manage service configurations
   - Handle environment-specific settings

Kubernetes configurations benefit greatly from Emrichen's templating capabilities. You can maintain a single template that adapts to different environments while ensuring consistency across your infrastructure.

2. **Application Configuration**
   - Generate config files for different environments
   - Manage feature flags
   - Handle service dependencies

Application configuration management becomes more manageable with Emrichen's variable system and conditional processing. You can maintain different configurations for development, staging, and production environments while sharing common settings.

3. **Infrastructure as Code**
   - Template cloud resource definitions
   - Manage infrastructure configurations
   - Generate deployment scripts

Infrastructure as Code benefits from Emrichen's ability to handle complex data structures and its support for modular templates. You can maintain reusable infrastructure components while allowing for customization.

## Advanced Topics

### Custom Tags
Emrichen supports extending its functionality through custom tags, allowing you to add domain-specific templating capabilities. Custom tags can be implemented to handle specific use cases or to integrate with external systems.

### Template Composition
Complex templates can be broken down into smaller, reusable components using `!Include` and `!With` tags for better maintainability. This modular approach allows you to build a library of reusable template components.

### Error Handling and Debugging
The `!Debug` and `!Error` tags provide powerful tools for troubleshooting and handling edge cases in your templates. Combined with type validation tags, they help ensure your templates remain robust and maintainable.

## Conclusion

Emrichen provides a robust and flexible solution for YAML templating needs. Its understanding of YAML structure, combined with powerful features like variable management, control flow, and data manipulation, makes it an excellent choice for managing complex configurations in modern software development workflows. The comprehensive help system and extensive documentation make it easy to learn and master Emrichen's capabilities.

