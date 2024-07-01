# Development Guidelines

This section covers general guidelines for developing with and contributing to
the Vanilla OS SDK.

## Inclusion of Features in the SDK

Deciding whether a feature should be included in the SDK is a crucial step in
maintaining consistency and avoiding unnecessary bloat. Consider the following
criteria:

- **Reusability:**\
   Is the feature likely to be reused by multiple applications or components?
- **Dependency Management:**\
   Does the feature introduce new dependencies? Minimize dependencies to keep
  the SDK lightweight. If the feature is considered relevant to the SDK but
  adds extra depedencies, please open an issue with your suggestion so the
  maintainers can decide whether it should be implemented. This will avoid
  unnecessary work that may not be merged.
- **Command-Line Tools**:
   Use of command-line tools under the hood should always be considered a
  last resort. Whenever command-line tools are used, it is useful to map a
  wrapper to allow more practical usage. In such cases, the wrapper should be
  placed in the `pkg/v1/legacy` package.

## Package Structure

When developing a new feature, it is important to choose the right package
within the SDK. Follow these guidelines:

- **Logical Grouping:**\
   Group related functionalities together to provide a clear structure, for
  example, the `system` package contains functions for interacting with the
  system in many ways, such as managing processes, users etc., we are not
  placing all the functions in the same file, but we are grouping them in
  multiple files, for example, the `system` package contains the following
  files: `proc.go`, `user.go` and so on, this provides a clear structure
  and makes it easier to find the functions you are looking for by looking
  at the package structure.
- **Avoid Overcrowding:**\
   If a package becomes too large, consider breaking it into sub-packages
  for better organization, but keep the structure as flat as possible.
- **Naming Conventions:**\
   Use descriptive names for packages and sub-packages, for example, if
  a package contains a set of functions for interacting with the file system,
  name it `fs` or `filesystem`.
- **Avoid Redundancy:**\
   Avoid creating packages that duplicate existing functionalities.

## Naming Conventions

Follow these naming conventions to ensure consistency and readability:

- **Package Names:**\
   Use lowercase, snake_case names for packages, for example, `fs` or
  `filesystem` or `file_system` not `FS` or `FileSystem`.
- **Function Names:**\
   Use camel case names for functions, for example, `ListFiles` not
  `list_files`.
- **Type Names, Variable Names and Constant Names:**\
   Use camel case self-descriptive names for types, for example, `ServiceStatus`
  not `status`.
- **Acronyms:**\
   Use camel case names for acronyms, for example, `HTTPServer` not
  `HttpServer`.
- **Abbreviations:**\
   Avoid abbreviations, for example, use `FileSystem` not `FileSys`. If you
  want to use abbreviations for long names, you are probably using the wrong
  name and it's worth considering stepping back and rethinking the logic.

## Code Formatting

A clear and consistent code formatting style is essential for readability and
maintainability. Follow these guidelines:

- **Formatting:**\
   Use the in-built `gofmt` tool to format the code, for example:

  ```bash
  gofmt -w .
  ```

## Documentation Standards

Consistent and comprehensive documentation is essential for developers using
the SDK. Follow these documentation standards:

- **Clear Descriptions:**\
   Provide clear and concise descriptions for functions, types, and packages,
  remember to keep the descriptions up-to-date when making changes.
- **Examples:**\
   Include usage examples to demonstrate how to use different features.

Please refer to the [Go Doc Comments](https://go.dev/doc/comment) for more
information on how to write documentation comments.

## Testing Practices

Ensure the reliability and stability of the SDK by following effective
testing practices:

- **Unit Testing:**\
   Write unit tests for individual functions and features and locate them in
  the tests directory of the package.
- **Integration Testing:**\
   Test interactions between different components within the SDK.

## Compatibility Considerations

When introducing new features or making changes, consider the impact on
existing applications:

- **Backward Compatibility:**\
   Strive to maintain backward compatibility when introducing new versions.
  It's fine to change how features works, but provide clear documentation
  and migration guides for developers.
- **Deprecation Warnings:**\
   If a feature is deprecated, provide clear deprecation warnings and guidance
  on alternatives.

## Error Handling

Error handling is an important part of developing reliable applications.
Follow these guidelines:

- **Be hungry for errors:**\
   Always check for errors and handle them appropriately, even if they are
  unlikely to occur.
- **Return errors:**\
   Return errors to the caller instead of logging them, so that the caller
  can decide how to handle them. In case a function is changing its behavior
  due to an error (e.g. using a built-in file if the host file is not found),
  it should return a simple message with the comment `// @sdk:hint` before
  the message, so that it is easy to track. See an example in the
  `pkg/v1/hardware/hardware.go (LoadPCIDeviceMap)` function.
- **Provide context:**\
   Provide context when returning errors, for example, if a function fails
  to open a file, return an error with a message like "failed to open file
  'foo.txt'", this will help the caller to understand the cause of the error
  without having to dig into the implementation details.
- **Avoid panic:**\
   Avoid panic during development, but also avoid using panic to handle
  errors, as it can cause the application to crash.

## Usage of Custom Types

To enhance the readability and maintainability of the SDK, use custom types
to represent complex data structures. Follow these guidelines:

- **Avoid using primitive types:**\
   Avoid using primitive types such as `int` and `string` to represent
  complex data structures, instead, use custom types, for example, if a
  function named `ListCPUs` returns a list of CPUs, define a custom type
  named `CPU` and use `[]CPU` to represent the list.
- **Avoid using `interface{}`:**\
   Avoid using `interface{}` to represent complex data structures, as it
  makes the code harder to read and maintain.
- **Enumerate types:**\
   Enumerate types to represent a set of values, for example, if a function
  returns a service status, define a custom type named `ServiceStatus` and use
  `ServiceStatus` to represent the status, for example:

  ```go
    type ServiceStatus int

    const (
        ServiceStatusError ServiceStatus = 0
        ServiceStatusOK ServiceStatus = 1
    )
  ```

## SDK Versioning

The SDK follows the [Semantic Versioning](https://semver.org/) standard.
The version number is defined in the `VERSION` file in the root directory
of the SDK version, for example, the version number of the SDK version
`v1` is defined in `v1/VERSION`.

The version number is defined as follows:

```text
<major>.<minor>.<patch>
```

- **Major:**\
   The major version number is incremented when introducing breaking
  changes.
- **Minor:**\
   The minor version number is incremented when introducing new features
  or functionalities.
- **Patch:**\
  The patch version number is incremented when making bug fixes or
  performance improvements.

### Defining SDK Versions

When defining a new SDK version, consider the following:

- **Breaking Changes:**\
   If the new version introduces breaking changes, increment the major
  version number. This is crucial for maintaining backward compatibility and
  avoiding third-party applications from assuming that the new version is
  backward compatible.
- **New Features:**\
   If the new version introduces new features or functionalities, increment
  the minor version number.
- **Bug Fixes:**\
   If the new version only contains bug fixes or performance improvements,
  increment the patch version number.

By adhering to these guidelines, developers can contribute to a consistent,
well-documented, and reliable SDK.
