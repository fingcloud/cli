package printer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
)

type tablePrinter struct {
	buffer     *bytes.Buffer
	format     string
	headers    interface{}
	RowHandler TableRowHandler
}

// TableRowHandler handles table row
type TableRowHandler func(interface{}) error

func NewTablePrinter(format string, headers interface{}) *tablePrinter {
	return &tablePrinter{
		format:  format,
		buffer:  new(bytes.Buffer),
		headers: headers,
	}
}

func (p *tablePrinter) HandleRow(fn TableRowHandler) {
	p.RowHandler = fn
}

func (p *tablePrinter) printRow(tmpl *template.Template, data interface{}) error {
	err := tmpl.Execute(p.buffer, data)
	if err != nil {
		return fmt.Errorf("can't print row: %v", err)
	}

	p.buffer.WriteString("\n")
	return nil
}

func (p *tablePrinter) Print() error {
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(p.format)
	if err != nil {
		return err
	}

	p.RowHandler()
	return nil
}
