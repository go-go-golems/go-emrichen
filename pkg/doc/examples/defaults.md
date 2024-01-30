---
Title: "!Defaults Tag"
Slug: tag-defaults
Short: |
  ```
  !Defaults
  var1: value1
  var2: value2
  ...
  ```
Command:
  - emrichen
Topics:
  - tags
IsTemplate: false
IsTopLevel: true
ShowPerDefault: false
SectionType: Example
---
# `!Defaults` Tag

The `!Defaults` tag in Emrichen is used to define default values for variables.
These default values can be overridden by other variable sources or explicitly in the template.

```yaml
!Defaults
var1: value1
var2: value2
...
```

## Examples

### Basic Defaults

Define default values for variables:

```yaml
!Defaults
image_tag: latest
replica_count: 3
```

### Overriding Defaults

Override default values later in the template:

```yaml
!Defaults
image_tag: latest
replica_count: 3
---
image_tag: v1.2.3
replicas: !Var replica_count
```

**Output:**

```yaml
image_tag: v1.2.3
replicas: 3
```

### Defaults with Multiple Documents

Use `!Defaults` in a multi-document YAML file:

```yaml
!Defaults
image_tag: latest
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: nginx
          image: !Format "nginx:{image_tag}"
```

**Output:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: nginx
          image: nginx:latest
```

## Notes

- The `!Defaults` tag must be in a separate document (preceded by `---`) in the YAML file.
- Variables defined in `!Defaults` can be overridden by other variable sources or explicitly in the template.
- If a variable is not defined elsewhere and no default is provided, the behavior depends on the template's error handling configuration.