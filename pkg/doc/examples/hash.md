---
Title: "String Processing Tags: Base64, SHA1, SHA256, MD5"
Slug: string-processing-tags
Short: |
  ```
  !Base64,Var value
  !SHA1,Var value
  !SHA256,Var value
  !MD5,Var value
  ```
Command:
  - emrichen
Topics:
  - tags
IsTemplate: false
IsTopLevel: true
ShowPerDefault: true
SectionType: Example
---
# String Processing Tags: Base64, SHA1, SHA256, MD5

The string processing tags `!Base64`, `!SHA1`, `!SHA256`, and `!MD5` in Emrichen are used to encode or hash strings.
These tags accept a single parameter: the string to be processed. The `!Base64` tag encodes the string into Base64
format, while `!SHA1`, `!SHA256`, and `!MD5` hash the string using their respective algorithms. When multiple tags are
directly in sequence, they must be composed using commas (e.g., `!Base64,Var value`) to maintain valid YAML tags.

## Examples

### Base64 Encoding

Encode a string using Base64.

```yaml
!Defaults
sampleString: "Hello, World!"
---
encodedString: !Base64,Var sampleString
```

**Output:**

```yaml
encodedString: "SGVsbG8sIFdvcmxkIQ=="
```

### SHA256 Hashing

Hash a string using SHA256.

```yaml
!Defaults
sampleString: "Hello, World!"
---
hashedStringSHA256: !SHA256,Var sampleString
```

**Output:**

```yaml
hashedStringSHA256: "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
```

### Combined Hashing and Encoding

Hash a string using SHA256 and then encode it with Base64.

```yaml
!Defaults
sampleString: "Hello, World!"
---
combinedHashAndEncode: !Base64,SHA256,Var sampleString
```

**Output:**

```yaml
combinedHashAndEncode: "pZGm1Av0IEBKAQczz7exkNYsZb8LzaMrtXsn2a2fFG4="
```