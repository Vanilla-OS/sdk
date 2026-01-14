package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import "bytes"

// Document represents a roff document buffer.
type Document struct {
	Buffer bytes.Buffer
}

const (
	// TitleHeading formats the document title header.
	TitleHeading = `.TH %[1]s %[2]d "%[4]s" "%[3]s" "%[5]s"`
	// Paragraph macro.
	Paragraph = "\n.PP"
	// Indent begins a relative indent block.
	Indent = "\n.RS"
	// IndentEnd ends an indent block.
	IndentEnd = "\n.RE"
	// IndentedParagraph writes an indented paragraph.
	IndentedParagraph = "\n.IP"
	// SectionHeading macro.
	SectionHeading = "\n.SH %s"
	// SubSectionHeading macro.
	SubSectionHeading = "\n.SS %s"
	// TaggedParagraph macro.
	TaggedParagraph = "\n.TP"

	// Bold escape sequence.
	Bold = `\fB`
	// Italic escape sequence.
	Italic = `\fI`
	// PreviousFont resets the font.
	PreviousFont = `\fP`
)
