---
Title: "!Include Tags"
Slug: tag-include-family
Short: |
  ```
  !Include
  !IncludeBase64
  !IncludeBinary
  !IncludeGlob
  !IncludeText
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
# `!Include` Tags Family

The `!Include` tags in Emrichen are a set of constructs designed for incorporating external content into the
current YAML or JSON document. 

## `!Include`

Includes the content of a YAML or JSON file.

- **Usage**: `!Include path/to/file.yml`
- **Purpose**: To include and render a YAML or JSON file within the current document.

**Example**:

```yaml
config: !Include config.yml
```

## `!IncludeBase64`

Includes a file's content encoded in Base64.

- **Usage**: `!IncludeBase64 path/to/binary.file`
- **Purpose**: To include binary files (e.g., images, documents) as a Base64 encoded string.

**Example**:

```yaml
imageData: !IncludeBase64 path/to/image.png
```

## `!IncludeBinary`

Includes a binary file's content directly.

- **Usage**: `!Base64,IncludeBinary path/to/binary.file`
- **Purpose**: Mainly used for hashing or when binary content needs to be processed in some way.

**Example**:

```yaml
binaryContent: !IncludeBinary path/to/executable.bin
```

## `!IncludeGlob`

Includes multiple files matching a glob pattern.

- **Usage**: `!IncludeGlob "path/to/files/*.yml"`
- **Purpose**: To include and render multiple YAML or JSON files that match a specified pattern, 
  useful for batch processing or when dealing with multiple configuration files. The included context 
  also gets expanded with emrichen.

**Example**:

```yaml
documents: !IncludeGlob "configs/*.yml"
```

## `!IncludeText`

Includes the content of a text file as a string.

- **Usage**: `!IncludeText path/to/text.file`
- **Purpose**: To include plain text files, allowing their content to be used directly within the document.

**Example**:

```yaml
welcomeMessage: !IncludeText welcome.txt
```