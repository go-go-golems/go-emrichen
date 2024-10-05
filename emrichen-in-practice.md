# How to Use Emrichen in Practice

Emrichen empowers users to create dynamic and flexible YAML configurations through its robust templating capabilities. This section provides a step-by-step guide on installing Emrichen, understanding its basic to advanced features, and implementing practical examples that range from simple to complex use cases.

## Table of Contents

1. [Installation](#installation)
2. [Basic Usage](#basic-usage)
   - [Defining Defaults](#defining-defaults)
   - [Variable Substitution](#variable-substitution)
3. [Intermediate Usage](#intermediate-usage)
   - [Conditional Logic with `!If`](#conditional-logic-with-if)
   - [Looping with `!Loop`](#looping-with-loop)
   - [Filtering Data with `!Filter`](#filtering-data-with-filter)
4. [Advanced Usage](#advanced-usage)
   - [Grouping Data with `!Group`](#grouping-data-with-group)
   - [Merging Dictionaries with `!Merge`](#merging-dictionaries-with-merge)
   - [Composing Tags for Complex Operations](#composing-tags-for-complex-operations)
5. [Practical Examples](#practical-examples)
   - [Simple Example: Greeting Message](#simple-example-greeting-message)
   - [Complex Example: Dynamic Deployment Configuration](#complex-example-dynamic-deployment-configuration)

---

## Installation

Emrichen is implemented in Go and can be integrated as a library in your Go projects or used as a standalone tool for processing YAML templates. Follow the steps below to install Emrichen.

### Prerequisites

- **Go Language:** Ensure that Go is installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Installing as a Go Library

To include Emrichen in your Go project, use the `go get` command:

```bash
go get github.com/go-go-golems/go-emrichen
```

### Building the Standalone Tool

If Emrichen provides a command-line interface (CLI), you can build it using the following commands. (Assuming the CLI is part of the repository.)

```bash
git clone https://github.com/go-go-golems/go-emrichen.git
cd go-emrichen
go build -o emrichen cmd/emrichen/main.go
```

After building, you can run the `emrichen` executable to process your YAML templates.

### Verifying the Installation

To verify that Emrichen is installed correctly, you can display its version or help information:

```bash
emrichen --version
```

```bash
emrichen --help
```

*Note:* The exact commands may vary based on the implementation details of Emrichen's CLI.

---

## Basic Usage

Begin by understanding the foundational features of Emrichen: defining default variables and substituting them within your YAML templates.

### Defining Defaults

Emrichen allows you to set default variables that can be referenced throughout your YAML document. This promotes reusability and easier configuration management.

**Example:**

```yaml
# defaults.yml
!Defaults
app_name: myapp
image_tag: v1.0.0
replica_count: 3
```

### Variable Substitution

Once defaults are defined, you can substitute these variables in other parts of your YAML using the `!Var` tag for top-level variables, and `!Lookup` for nested access.

**Example:**

```yaml
# deployment.yml
!Include "defaults.yml"
---
deployment:
  name: !Var app_name
  image: !Format "myapp:{image_tag}"
  replicas: !Var replica_count
  config:
    database_url: !Lookup "database.url"
```

**Resulting YAML After Processing:**

```yaml
deployment:
  name: "myapp"
  image: "myapp:v1.0.0"
  replicas: 3
  config:
    database_url: "postgres://localhost:5432/mydb"
```

**Explanation:**

1. **Defining Defaults:** The `defaults.yml` file sets default values for `app_name`, `image_tag`, `replica_count`, and nested values like `database.url`.
2. **Including Defaults:** The `!Include` tag incorporates `defaults.yml` into `deployment.yml`.
3. **Substituting Variables:** 
   - The `!Var` tag replaces placeholders with their corresponding top-level default values.
   - The `!Lookup` tag is used for accessing nested values like `database.url`.

---

## Intermediate Usage

Once comfortable with basic variable substitution, Emrichen's intermediate features like conditional logic, looping, and data filtering can enhance your YAML configurations' dynamism.

### Conditional Logic with `!If`

The `!If` tag allows you to include different values based on a condition, enabling dynamic configurations.

**Structure:**

```yaml
result: !If
  test: <condition>
  then: <value_if_true>
  else: <value_if_false>
```

**Example:**

```yaml
# conditional.yml
!Defaults
is_production: true
---
environment:
  type: !If
    test: !Var is_production
    then: "production"
    else: "development"
```

**Resulting YAML After Processing:**

```yaml
environment:
  type: "production"
```

**Explanation:**

- **Test Condition:** Evaluates the `is_production` variable.
- **Then Clause:** If `is_production` is `true`, sets `environment.type` to `"production"`.
- **Else Clause:** If `is_production` is `false`, sets `environment.type` to `"development"`.

### Looping with `!Loop`

The `!Loop` tag iterates over a collection, applying a template to each item. This is useful for generating repetitive configurations like multiple Kubernetes resources.

**Structure:**

```yaml
result: !Loop
  over: <collection>
  as: <current_item_variable>
  template: <template_to_apply>
```

**Example:**

```yaml
# loop.yml
!Defaults
services:
  - name: frontend
    port: 80
  - name: backend
    port: 8080
---
serviceDefinitions: !Loop
  over: !Var services
  as: service
  template: !Var service
```

**Resulting YAML After Processing:**

```yaml
serviceDefinitions:
  - name: frontend
    port: 80
  - name: backend
    port: 8080
```

**Explanation:**

- **Over Clause:** Iterates over the `services` list.
- **As Clause:** Assigns each service to the `service` variable.
- **Template Clause:** Inserts the service definition into `serviceDefinitions`.

### Filtering Data with `!Filter`

The `!Filter` tag allows you to filter elements from a collection based on a predicate, enabling conditional inclusion of data.

**Structure:**

```yaml
result: !Filter
  test: <predicate_expression>
  over: <collection>
  as: <current_item_variable>  # Optional
```

**Example:**

```yaml
# filter.yml
!Defaults
numbers: [1, 2, 3, 4, 5, 6]
---
evenNumbers: !Filter
  test: !Op
    a: !Lookup item
    op: "%"
    b: 2
  over: !Var numbers
  as: item
```

**Resulting YAML After Processing:**

```yaml
evenNumbers:
  - 2
  - 4
  - 6
```

**Explanation:**

- **Test Clause:** Uses the `!Op` tag to check if a number is even (`number % 2 == 0`).
- **Over Clause:** Applies the filter to the `numbers` list.
- **As Clause:** Assigns each number to the `item` variable for evaluation.

---

## Advanced Usage

Emrichen's advanced features enable complex transformations and data manipulations, providing granular control over your YAML configurations.

### Grouping Data with `!Group`

The `!Group` tag aggregates items from a collection based on a specified key or expression, allowing for organized and hierarchical configurations.

**Structure:**

```yaml
result: !Group
  over: <collection>
  by: <grouping_expression>
  template: <optional_transformation>
  as: <current_item_variable>  # Optional
```

**Example:**

```yaml
# group.yml
!Defaults
employees:
  - name: Alice
    department: Engineering
  - name: Bob
    department: Marketing
  - name: Carol
    department: Engineering
---
groupedByDepartment: !Group
  over: !Var employees
  by: !Lookup employee.department
  as: employee
  template: !Lookup employee.name
```

**Resulting YAML After Processing:**

```yaml
groupedByDepartment:
  Engineering:
    - "Alice"
    - "Carol"
  Marketing:
    - "Bob"
```

**Explanation:**

- **Over Clause:** Iterates over the `employees` list.
- **By Clause:** Groups employees by their `department`.
- **As Clause:** Assigns each employee to the `employee` variable.
- **Template Clause:** Extracts the `name` of each employee for the grouped list.

### Merging Dictionaries with `!Merge`

The `!Merge` tag combines multiple dictionaries into a single mapping, allowing for the consolidation of configurations.

**Structure:**

```yaml
result: !Merge
  - <first_mapping>
  - <second_mapping>
  - ...
```

**Example:**

```yaml
# merge.yml
!Defaults
defaultSettings:
  timeout: 30
  retries: 3
overrideSettings:
  retries: 5
  verbose: true
---
mergedSettings: !Merge
  - !Var defaultSettings
  - !Var overrideSettings
```

**Resulting YAML After Processing:**

```yaml
mergedSettings:
  timeout: 30
  retries: 5
  verbose: true
```

**Explanation:**

- **First Mapping:** Sets default `timeout` and `retries`.
- **Second Mapping:** Overrides `retries` and adds `verbose`.
- **Merged Result:** Combines both mappings, with the second mapping's `retries` overriding the first.

### Composing Tags for Complex Operations

Emrichen allows for the combination of multiple tags to perform sophisticated data manipulations, enabling intricate configurations tailored to specific requirements.

**Example: Dynamic Deployment Configuration**

```yaml
# dynamic_deployment.yml
!Defaults
env: production
image_tags:
  frontend: v1.2.3
  backend: v2.3.4
replicas:
  frontend: 3
  backend: 2
---
deploymentConfig: !With
  vars:
    services:
      - name: frontend
        image: !Lookup image_tags.frontend
        replicas: !Lookup replicas.frontend
      - name: backend
        image: !Lookup image_tags.backend
        replicas: !Lookup replicas.backend
  template: !Loop
    over: !Var services
    as: service
    template: !Group
      over: !Var service
      by: !Format "{service.name}_deployment"
      template: !Merge
        - { image: !Lookup service.image }
        - { replicas: !Lookup service.replicas }
```

**Resulting YAML After Processing:**

```yaml
deploymentConfig:
  frontend_deployment:
    image: "v1.2.3"
    replicas: 3
  backend_deployment:
    image: "v2.3.4"
    replicas: 2
```

**Explanation:**

1. **Defaults:**
   - Sets the environment to `production`.
   - Defines `image_tags` for `frontend` and `backend`.
   - Specifies `replicas` for each service.

2. **With Scope:**
   - Defines a local scope with `services`, each containing `name`, `image`, and `replicas`.
   
3. **Loop:**
   - Iterates over each service in `services`.
   
4. **Group within Loop:**
   - Groups each service by a dynamically formatted deployment name (`{service.name}_deployment`).
   
5. **Merge within Group:**
   - Merges the `image` and `replicas` into a single mapping for each deployment.

**Benefits:**

- **Scalability:** Easily add more services without modifying the core deployment logic.
- **Maintainability:** Centralizes configurations like image tags and replica counts.
- **Flexibility:** Combines multiple tags (`!With`, `!Loop`, `!Group`, `!Format`, `!Merge`) for dynamic and context-aware configurations.

---

## Practical Examples

To solidify your understanding of Emrichen's capabilities, let's explore two practical examples: a simple greeting message and a more complex dynamic deployment configuration.

### Simple Example: Greeting Message

**Objective:**  
Generate personalized greeting messages based on user data.

**Files:**

1. **defaults.yml**

    ```yaml
    !Defaults
    user:
      firstName: Alice
      lastName: Smith
      isAdmin: false
    ```

2. **greeting.yml**

    ```yaml
    !Include "defaults.yml"
    ---
    greetingMessage: !If
      test: !Lookup user.isAdmin
      then: !Format "Welcome back, Admin {user.firstName} {user.lastName}!"
      else: !Format "Hello, {user.firstName} {user.lastName}!"
    ```

**Processing Steps:**

1. **Include Defaults:** Incorporate `defaults.yml` to set user variables.
2. **Conditional Greeting:** Use `!If` to determine if the user is an admin.
3. **Format Message:** Generate a greeting message based on the user's admin status.

**Resulting YAML:**

```yaml
greetingMessage: "Hello, Alice Smith!"
```

**Changing User Role:**

To generate an admin greeting, update `isAdmin` to `true` in `defaults.yml`:

```yaml
!Defaults
user:
  firstName: Alice
  lastName: Smith
  isAdmin: true
```

**Resulting YAML After Processing:**

```yaml
greetingMessage: "Welcome back, Admin Alice Smith!"
```

**Explanation:**

- **Flexibility:** Easily switch between admin and regular user greetings by toggling the `isAdmin` flag.
- **Reusability:** The same template can generate different messages based on user roles.

### Complex Example: Dynamic Deployment Configuration

**Objective:**  
Generate Kubernetes deployment configurations for multiple services with varying settings.

**Files:**

1. **defaults.yml**

    ```yaml
    !Defaults
    environment: production
    services:
      - name: frontend
        image: frontend:latest
        replicas: 3
        ports: [80, 443]
      - name: backend
        image: backend:latest
        replicas: 2
        ports: [8080]
      - name: cache
        image: redis:alpine
        replicas: 1
        ports: [6379]
    ```

2. **deployment_template.yml**

    ```yaml
    !Include "defaults.yml"
    ---
    deployments: !Loop
      over: !Var services
      as: service
      template: !Group
        over: 
          - name: !Var service.name
            image: !Var service.image
            replicas: !Var service.replicas
            ports: !Var service.ports
            environment: !Var environment
        by: !Format "{name}_deployment"
        template: !Merge
          - { image: !Var image }
          - { replicas: !Var replicas }
          - { ports: !Var ports }
          - { environment: !Var environment }
    ```

**Processing Steps:**

1. **Include Defaults:** Load `defaults.yml` to define environment and service configurations.
2. **Loop Over Services:** Iterate over each service defined in `services`.
3. **Group and Merge:** For each service, group the deployment by a formatted name and merge relevant configurations.
4. **Generate Deployments:** Create a `deployments` mapping containing configurations for each service.

**Resulting YAML:**

```yaml
deployments:
  frontend_deployment:
    image: "frontend:latest"
    replicas: 3
    ports:
      - 80
      - 443
    environment: "production"
  backend_deployment:
    image: "backend:latest"
    replicas: 2
    ports:
      - 8080
    environment: "production"
  cache_deployment:
    image: "redis:alpine"
    replicas: 1
    ports:
      - 6379
    environment: "production"
```

**Explanation:**

- **Dynamic Iteration:** The `!Loop` tag processes each service, ensuring scalability as more services are added.
- **Structured Grouping:** The `!Group` tag organizes deployments by service names, creating clear and maintainable configurations.
- **Comprehensive Merging:** The `!Merge` tag consolidates various configuration aspects (image, replicas, ports, environment) into each deployment entry.
- **Consistency:** All deployments inherit the `environment` variable, ensuring uniform settings across services.

**Benefits:**

- **Scalability:** Easily add or remove services in `defaults.yml` without altering the deployment template.
- **Maintainability:** Centralizes configuration logic, reducing redundancy and potential errors.
- **Flexibility:** Supports complex configurations with minimal adjustments, catering to diverse deployment scenarios.

---

# Conclusion

Emrichen offers a robust and flexible approach to YAML templating, enabling dynamic and maintainable configurations through its extensive suite of custom tags. By leveraging features ranging from basic variable substitution to advanced data grouping and merging, Emrichen streamlines the process of generating complex YAML structures tailored to your specific needs.

Whether you're managing simple configurations or orchestrating intricate deployment setups, Emrichen provides the tools necessary to enhance your workflow, improve scalability, and ensure consistency across your projects.

For further exploration, refer to the [Detailed Specification of Emrichen Language](#detailed-specification-of-emrichen-language) and experiment with the provided practical examples to harness the full potential of Emrichen in your YAML configurations.