package logs

// LogContext represents a hierarchical logging context. Contexts can be
// nested and will automatically generate a prefix based on their
// parent/child relationship.
type LogContext struct {
	// Name is the name of the context segment.
	Name string

	// Parent references the parent context. Nil means this context is the
	// root one.
	Parent *LogContext
}

// Prefix returns the context prefix in the form "parent:child". An empty
// context results in an empty string.
func (c *LogContext) Prefix() string {
	if c == nil {
		return ""
	}

	if c.Parent == nil {
		return c.Name
	}

	parent := c.Parent.Prefix()
	if parent == "" {
		return c.Name
	}

	if c.Name == "" {
		return parent
	}

	return parent + ":" + c.Name
}

// NewLogContext creates a new LogContext optionally linked to a parent.
func NewLogContext(name string, parent *LogContext) *LogContext {
	return &LogContext{Name: name, Parent: parent}
}
