---
Title: "!URLEncode Tag"
Slug: tag-urlencode
Short: |
  ```
  !URLEncode
  - A string to encode
    OR
  - url: The URL to combine query parameters into
  - query: An object of query string parameters to add OR a string of query string parameters
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
# `!URLEncode` Tag

The `!URLEncode` tag in Emrichen is designed for URL encoding of strings and for combining base URLs with query
parameters. It accepts either a single string to encode or a combination of a base URL and query parameters.

- **A string to encode**: Encodes the given string for safe inclusion in a URL.
- **url**: Specifies the base URL to which query parameters will be added.
- **query**: Defines the query parameters to add to the base URL. This can be an object of query string parameters or a string of query string parameters.

## Examples

### Basic URL Encoding

Encode a simple string to be safely included in a URL.

```yaml
encodedString: !URLEncode "hello world & special=characters"
```

**Output:**

```yaml
encodedString: "hello+world+%26+special%3Dcharacters"
```

### URL with Query Parameters

Combine a base URL with query parameters.

```yaml
fullURL: !URLEncode
  url: "https://example.com/"
  query:
    param1: "value1"
    param2: "value2"
```

**Output:**

```yaml
fullURL: "https://example.com/?param1=value1&param2=value2"
```

### Complex String Encoding

Encode a complex string containing an email address and parameters.

```yaml
encodedComplexString: !URLEncode "email=example@example.com&param=value"
```

**Output:**

```yaml
encodedComplexString: "email%3Dexample%40example.com%26param%3Dvalue"
```
