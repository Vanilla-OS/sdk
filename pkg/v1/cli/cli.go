package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"fmt"
	"reflect"
	"strings"
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

	root   any
	app    *builder.App
	manCmd *ManCmd
}

// ManCmd is the command to generate the man page
type ManCmd struct {
	Base
	root       any
	translator help.Translator
}

// Run runs the man command
//
// Example:
//
//	manCmd := &cli.ManCmd{root: s}
//	err := parser.Run(manCmd)
func (c *ManCmd) Run() error {
	man, err := GenerateManPage(c.root, c.translator)
	if err != nil {
		return err
	}

	fmt.Println(man)
	return nil
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
	if c.manCmd != nil {
		c.manCmd.translator = tr
	}
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

	// We inject the man command
	manCmd := &ManCmd{root: s}
	manNode, err := parser.Parse("man", manCmd)
	if err == nil {
		manNode.Description = "Generate man page"
		app.AddCommand("man", manNode)
	}

	node := app.RootNode

	c := &Command{
		Use:    node.Name,
		Short:  node.Description,
		root:   s,
		app:    app,
		manCmd: manCmd,
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
//	man, err := cli.GenerateManPage(&RootCmd{}, nil)
//	if err != nil {
//		return "", err
//	}
//
// GenerateManPage automatically uses a zero-value instance of the root struct
// to exclude any dynamic commands.
func GenerateManPage(root any, tr help.Translator) (string, error) {
	t := reflect.TypeOf(root)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	cleanRoot := reflect.New(t).Interface()

	node, err := parser.Parse("root", cleanRoot)
	if err != nil {
		return "", err
	}

	description := node.Description
	if tr != nil {
		description = tr(cleanKey(description))
	}

	d := roff.NewDocument()
	d.Heading(1, node.Name, description, time.Now())

	docNode(d, node, tr)

	return d.String(), nil
}

// docNode recursively documents a command node and its children.
func docNode(d *roff.Document, node *parser.CommandNode, tr help.Translator) {
	description := node.Description
	if tr != nil {
		description = tr(cleanKey(description))
	}

	d.Section("subcommand " + node.Name)
	d.Indent(4)
	d.Text(description)
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
			desc := meta.Description
			if tr != nil {
				desc = tr(cleanKey(desc))
			}
			d.Text(fmt.Sprintf("  %s--%s  %s\n", short, name, desc))
		}
		d.EndSection()
	}

	// Commands
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			docNode(d, child, tr)
		}
	}
}

func cleanKey(key string) string {
	if strings.HasPrefix(key, "pr:") {
		return strings.TrimPrefix(key, "pr:")
	}
	return key
}

// GetRoot returns the underlying root struct of the command.
func (c *Command) GetRoot() any {
	return c.root
}
