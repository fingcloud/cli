package formatter

import (
	"bytes"
	"io"
	"text/tabwriter"
	"text/template"

	"github.com/Masterminds/sprig"
)

type formatter struct {
	output io.Writer
	Format string
	Header interface{}
	buffer *bytes.Buffer
}

func New(output io.Writer, header interface{}, format string) *formatter {
	return &formatter{
		output: output,
		Format: format,
		Header: header,
	}
}

type Executer func(w io.Writer, tmpl *template.Template) error

func (f *formatter) Write(executer Executer) error {
	f.buffer = new(bytes.Buffer)

	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(f.Format)
	if err != nil {
		return err
	}

	if err := executer(f.buffer, tmpl); err != nil {
		return err
	}

	t := tabwriter.NewWriter(f.output, 10, 1, 3, ' ', 0)
	buffer := new(bytes.Buffer)

	if err := tmpl.Execute(buffer, f.Header); err != nil {
		return err
	}

	buffer.WriteTo(t)
	t.Write([]byte{'\n'})

	f.buffer.WriteTo(t)
	t.Flush()

	return nil
}
