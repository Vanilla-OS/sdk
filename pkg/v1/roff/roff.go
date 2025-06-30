package roff

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/vanilla-os/sdk/pkg/v1/roff/types"
)

// Document wraps types.Document adding helper methods.
type Document struct {
	types.Document
}

// NewDocument creates a new roff document.
func NewDocument() *Document {
	return &Document{}
}

func (d *Document) writef(format string, args ...any) {
	if bytes.HasSuffix(d.Buffer.Bytes(), []byte("\n")) && strings.HasPrefix(format, "\n") {
		format = strings.TrimPrefix(format, "\n")
	}
	fmt.Fprintf(&d.Buffer, format, args...)
}

func (d *Document) writelnf(format string, args ...any) {
	d.writef(format+"\n", args...)
}

// Heading writes the document heading.
func (d *Document) Heading(section uint, title, desc string, ts time.Time) {
	d.writef(types.TitleHeading, strings.ToUpper(title), section, title, ts.Format("2006-01-02"), desc)
}

// Paragraph starts a new paragraph.
func (d *Document) Paragraph() { d.writelnf(types.Paragraph) }

// Indent increases the indentation level.
func (d *Document) Indent(n int) {
	if n >= 0 {
		d.writelnf(types.Indent+" %d", n)
	} else {
		d.writelnf(types.Indent)
	}
}

// IndentEnd decreases the indentation level.
func (d *Document) IndentEnd() { d.writelnf(types.IndentEnd) }

// TaggedParagraph starts a tagged paragraph.
func (d *Document) TaggedParagraph(indentation int) {
	if indentation >= 0 {
		d.writelnf(types.TaggedParagraph+" %d", indentation)
	} else {
		d.writelnf(types.TaggedParagraph)
	}
}

// List writes a list item.
func (d *Document) List(text string) {
	d.writelnf(types.IndentedParagraph+" \\(bu 3\n%s", escapeText(strings.TrimSpace(text)))
}

// Section writes a section heading.
func (d *Document) Section(text string) {
	d.writelnf(types.SectionHeading, strings.ToUpper(text))
}

// SubSection writes a subsection heading.
func (d *Document) SubSection(text string) {
	d.writelnf(types.SubSectionHeading, strings.ToUpper(text))
}

// EndSection ends the current section.
func (d *Document) EndSection() { d.writelnf("") }

// EndSubSection ends the current subsection.
func (d *Document) EndSubSection() { d.writelnf("") }

// Text writes text handling basic lists.
func (d *Document) Text(text string) {
	inList := false
	for i, line := range strings.Split(text, "\n") {
		if i > 0 && !inList {
			d.Paragraph()
		}
		if strings.HasPrefix(line, "*") {
			if !inList {
				d.Indent(-1)
				inList = true
			}
			d.List(line[1:])
		} else {
			if inList {
				d.IndentEnd()
				inList = false
			}
			d.writef(escapeText(line))
		}
	}
}

// TextBold writes bold text.
func (d *Document) TextBold(text string) {
	d.writef(types.Bold)
	d.Text(text)
	d.writef(types.PreviousFont)
}

// TextItalic writes italic text.
func (d *Document) TextItalic(text string) {
	d.writef(types.Italic)
	d.Text(text)
	d.writef(types.PreviousFont)
}

// String returns the document as a string.
func (d *Document) String() string { return d.Buffer.String() }

func escapeText(s string) string {
	s = strings.ReplaceAll(s, `\`, `\e`)
	s = strings.ReplaceAll(s, ".", "\\&.")
	return s
}
