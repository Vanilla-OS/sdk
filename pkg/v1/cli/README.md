# CLI Package

The `cli` package provides a declarative framework for building CLI applications within the Vanilla OS SDK. It handles command hierarchy, flag parsing, argument validation, and help text generation.

## Usage

CLI applications are defined by mapping structs to commands, using struct fields and their associated tags.

### Command Definition

A struct representing a command must include `cli.Base` (an alias for `go-cli-builder`'s Base struct).

- `cmd`: Defines the name of a subcommand
- `help`: Defines the help description for the command or flag

```go
type RootCmd struct {
    cli.Base
    Config ConfigCmd `cmd:"config" help:"pr:vso.cmd.config"`
}
```

### Flags

Flags are defined as struct fields with the `flag` tag.

- `short`: Defines a single-letter alias
- `long`: Defines the full name of the flag
- `name`: Defines the display name in help text (can be a translation key)

```go
type RootCmd struct {
    cli.Base
    Version bool `flag:"short:v, long:version" help:"pr:vso.msg.version"`
}
```

### Arguments

Positional arguments are defined as struct fields with the `arg` tag.

- `optional`: Marks an argument as not required
- `name`: Defines the argument placeholder in help text

```go
type ConfigSetCmd struct {
    cli.Base
    Key   string `arg:"" name:"key" help:"pr:vso.arg.key"`
    Value string `arg:"" name:"value" help:"pr:vso.arg.value"`
}
```

## i18n Integration

The package integrates with the SDK's `i18n` component through the `pr:` (that stands for "provider")prefix in tags.

Strings starting with `pr:` are treated as translation keys, when an application is initialized via `app.App`, the SDK automatically registers a translator that resolves these keys using the application's locale files. Command titles, flag descriptions, and argument names are all translatable.

## Dynamic Commands

While the primary structure is static, commands can be added or modified at runtime using the `AddCommand` method on a `Command` instance.

```go
cmd, _ := cli.NewCommandFromStruct(&RootCmd{})
cmd.AddCommand("dynamic-command", myDynamicNode)
```

## UI Components

The package provides several pre-styled UI components built with Bubble Tea for common CLI interactions:

| Component | Purpose |
| :--- | :--- |
| `PromptText` | Collects string input from the user. |
| `SelectOption` | Displays a list for selection. Supports `%d` formatting for option counts. |
| `ConfirmAction` | Handles Yes/No confirmations. |
| `ProgressBar` | Visualizes the progress of a task. |
| `Spinner` | Indicates background activity for indeterminate tasks. |
| `Table` | Renders structured data in tabular format. |

## Execution

To initialize and run a CLI application:

```go
func main() {
    myApp := app.New("my-app", "1.0.0")
    // ... setup app ...
    
    cmd, err := cli.NewCommandFromStruct(&RootCmd{})
    if err != nil {
        panic(err)
    }
    
    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
    }
}
```
