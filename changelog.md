# Changelog

## Refactor tag handlers to use function maps

Refactored the interpreter to use a map of tag handlers instead of a large switch statement. This makes the code more maintainable and easier to extend with new tags.

- Moved all tag handlers to a defaultHandlers map
- Updated TagFunc signature to include interpreter instance
- Modified Process method to use the handlers map
- Updated RegisterTag to wrap old-style functions to match new signature 

## Comprehensive Documentation Update

Added comprehensive documentation for Emrichen in the topics section, providing a complete overview of the templating engine's capabilities and best practices.

- Added detailed explanation of core features and functionality
- Included practical examples for common use cases
- Added best practices and performance considerations
- Documented advanced topics and template composition strategies
- Added detailed explanations for each major feature
- Added information about built-in help commands
- Expanded examples with detailed explanations 