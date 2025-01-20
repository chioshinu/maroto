// Package page implements creation of pages.
package page

import (
	"fmt"
	"strings"

	"github.com/johnfercher/go-tree/node"

	"github.com/chioshinu/maroto/v2/pkg/core"
	"github.com/chioshinu/maroto/v2/pkg/core/entity"
	"github.com/chioshinu/maroto/v2/pkg/props"
)

type Page struct {
	number int
	total  int
	rows   []core.Row
	footer []core.Row
	config *entity.Config
	prop   props.PageNumber
}

// New is responsible to create a core.Page.
func New(ps ...props.PageNumber) core.Page {
	prop := props.PageNumber{}
	if len(ps) > 0 {
		prop = ps[0]
	}

	return &Page{
		prop: prop,
	}
}

// Render renders a Page into a PDF context.
func (p *Page) Render(provider core.Provider, cell entity.Cell) {
	innerCell := cell.Copy()

	prop := &props.Rect{}
	prop.MakeValid()

	if p.config.BackgroundImage != nil {
		provider.AddBackgroundImageFromBytes(p.config.BackgroundImage.Bytes, &innerCell, prop, p.config.BackgroundImage.Extension)
	}

	for _, row := range p.rows {
		row.Render(provider, innerCell)
		innerCell.Y += row.GetHeight(provider, &innerCell)
	}

	for _, row := range p.footer {
		if p.prop.Pattern != "" {
			for _, col := range row.GetColumns() {
				col.SetTransform(func(v string) string {
					return strings.ReplaceAll(v, p.prop.Pattern, fmt.Sprintf("%d", p.number))
				})
			}
		}
		row.Render(provider, innerCell)
		innerCell.Y += row.GetHeight(provider, &innerCell)
	}
}

// SetConfig sets the Page configuration.
func (p *Page) SetConfig(config *entity.Config) {
	p.config = config
	for _, row := range p.rows {
		row.SetConfig(config)
	}
}

// SetNumber sets the Page number and total.
func (p *Page) SetNumber(number int, total int) {
	p.number = number
	p.total = total
}

// GetNumber returns the Page number.
func (p *Page) GetNumber() int {
	return p.number
}

// Add adds one or more rows to the Page.
func (p *Page) Add(rows ...core.Row) core.Page {
	p.rows = append(p.rows, rows...)
	return p
}

// AddFooter adds one or more rows to the Page footer.
func (p *Page) AddFooter(rows ...core.Row) {
	p.footer = append(p.footer, rows...)
}

// GetRows returns the rows of the Page.
func (p *Page) GetRows() []core.Row {
	return p.rows
}

// GetStructure returns the Structure of a Page.
func (p *Page) GetStructure() *node.Node[core.Structure] {
	str := core.Structure{
		Type: "page",
	}

	n := node.New(str)
	for _, r := range p.rows {
		inner := r.GetStructure()
		n.AddNext(inner)
	}

	return n
}
