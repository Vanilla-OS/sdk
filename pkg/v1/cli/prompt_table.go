package cli

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Table prompt implementation using lipgloss for standardized output.
*/

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Table renders a styled table with the given headers and data using lipgloss.
//
// Example:
//
//	err := myApp.CLI.Table(
//		[]string{"Name", "Age"},
//		[][]string{
//			{"Batman", "35"},
//			{"Robin", "25"},
//		},
//	)
func (c *Command) Table(headers []string, data [][]string) error {
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(0, 1)
	styledHeaders := make([]string, len(headers))
	for i, h := range headers {
		styledHeaders[i] = headerStyle.Render(h)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(styledHeaders...).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			default:
				if row%2 == 0 {
					return lipgloss.NewStyle().
						Foreground(lipgloss.Color("252")).
						Padding(0, 1)
				}
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("246")).
					Padding(0, 1)
			}
		})

	fmt.Println(t)
	return nil
}
