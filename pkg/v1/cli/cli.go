package cli

import (
	"fmt"
	"time"

	builder "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"
	"github.com/mirkobrombin/go-cli-builder/v2/pkg/help"
	"github.com/mirkobrombin/go-cli-builder/v2/pkg/parser"
	"github.com/vanilla-os/sdk/pkg/v1/roff"
)

// Base is an alias for builder.Base to be used by consumers
type Base = builder.Base

// Command represents a CLI command.
type Command struct {
	Use   string
	Short string
	Long  string

	root any
	app  *builder.App
}

// Name returns the name of the command
func (c *Command) Name() string {
	return c.Use
}

// Execute runs the command
func (c *Command) Execute() error {
	if c.app == nil {
		return fmt.Errorf("no application initialized. Use NewCommandFromStruct")
	}
	return c.app.Run()
}

// AddCommand adds a dynamic command to the application.
func (c *Command) AddCommand(name string, cmd *parser.CommandNode) {
	c.app.AddCommand(name, cmd)
}

// SetTranslator sets the translator for the application.
func (c *Command) SetTranslator(tr help.Translator) {
	c.app.SetTranslator(tr)
}

// SetName sets the name of the root command.
func (c *Command) SetName(name string) {
	c.app.SetName(name)
}

// Reload re-parses the root struct to pick up dynamic changes.
func (c *Command) Reload() error {
	return c.app.Reload()
}

// NewCommandFromStruct returns a new Command created from a struct.
//
// Example:
//
//	type RootCmd struct {
//		cli.Base
//		Poll PollCmd `cmd:"poll" help:"Ask the user preferred hero"`
//		Man  ManCmd  `cmd:"man" help:"Generate man page"`
//	}
//
//	cmd, err := cli.NewCommandFromStruct(&RootCmd{})
//	if err != nil {
//		return nil, err
//	}
//
//	err := cmd.Execute()
func NewCommandFromStruct(s any) (*Command, error) {
	app, err := builder.New(s)
	if err != nil {
		return nil, err
	}

	node := app.RootNode

	c := &Command{
		Use:   node.Name,
		Short: node.Description,
		root:  s,
		app:   app,
	}
	return c, nil
}

// GenerateManPage generates a man page for the declarative struct
//
// Example:
//
//	type RootCmd struct {
//		cli.Base
//		Poll PollCmd `cmd:"poll" help:"Ask the user preferred hero"`
//		Man  ManCmd  `cmd:"man" help:"Generate man page"`
//	}
//
//	man, err := cli.GenerateManPage(&RootCmd{})
//	if err != nil {
//		return "", err
//	}
func GenerateManPage(root any) (string, error) {
	node, err := parser.Parse("root", root)
	if err != nil {
		return "", err
	}

	d := roff.NewDocument()
	d.Heading(1, node.Name, node.Description, time.Now())

	docNode(d, node)

	return d.String(), nil
}

// docNode recursively documents a command node and its children.
func docNode(d *roff.Document, node *parser.CommandNode) {
	d.Section("subcommand " + node.Name)
	d.Indent(4)
	d.Text(node.Description)
	d.IndentEnd()
	d.EndSection()

	// Options
	if len(node.Flags) > 0 {
		d.SubSection("Options")
		for name, meta := range node.Flags {
			short := ""
			if meta.Short != "" {
				short = fmt.Sprintf("-%s, ", meta.Short)
			}
			d.Text(fmt.Sprintf("  %s--%s  %s\n", short, name, meta.Description))
		}
		d.EndSection()
	}

	// Commands
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			docNode(d, child)
		}
	}
}

// GetRoot returns the underlying root struct of the command.
func (c *Command) GetRoot() any {
	return c.root
}
