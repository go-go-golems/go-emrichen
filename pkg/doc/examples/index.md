---
Title: "!Index Tag"
Slug: tag-index
Short: |
  ```
  !Index
  over: [List of items]
  by: Expression to determine the key
  template: (Optional) Template to apply for each item
  duplicates: (Optional) How to handle duplicates ('error', 'warn', 'ignore')
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
# `!Index` Tag 

The `!Index` tag in Emrichen is used to create a dictionary out of a list by specifying a key for each item. It supports
optional templating for each item and handling of duplicate keys.

```yaml
!Index
  over: [List of items]
  by: Expression to determine the key
  template: (Optional) Template to apply for each item
  duplicates: (Optional) How to handle duplicates ('error', 'warn', 'ignore')
```

## Examples

### Basic Indexing

Creating a dictionary indexed by names:

```yaml
indexedEmployees: !Index
  over: !Var employees
  by: !Lookup item.name
  template: !Lookup item.department
```

**Output:**

```yaml
indexedEmployees:
  Alice: Engineering
  Bob: Marketing
  Carol: Engineering
  Dave: Marketing
```

### Indexing with Custom Key

Indexing products by their ID:

```yaml
indexedProducts: !Index
  over: !Var products
  by: !Lookup item.id
  template: !Lookup item.name
```

**Output:**

```yaml
indexedProducts:
  101: "Laptop"
  102: "Smartphone"
  103: "Tablet"
```

### Handling Duplicates in Index

Indexing orders by customer name and handling duplicate keys by ignoring them:

```yaml
indexedOrders: !Index
  over: !Var orders
  by: !Lookup item.customer
  duplicates: ignore
  template: !Lookup item.amount
```

**Output:**

```yaml
indexedOrders:
  John: 500
  Emma: 450
```

## Notes

- The `!Index` tag requires that its argument `over` is a list.
- The `by` expression is used to determine the unique key for each item in the list.
- The optional `template` can be used to specify how each item should be represented in the resulting dictionary.
- Duplicate keys can be handled by specifying `duplicates` as 'error', 'warn', or 'ignore'.