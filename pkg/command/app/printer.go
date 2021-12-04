package app

import (
	"bytes"
	"io"
	"text/tabwriter"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/fingcloud/cli/pkg/api"
)

const (
	defaultFormat = "{{.Name}}\t{{.Platform}}\t{{.Status}}"
)

var tableHeaders = map[string]string{
	"ID":       "ID",
	"Name":     "NAME",
	"Label":    "LABEL",
	"Platform": "PLATFORM",
	"Status":   "STATUS",
}

func PrintFormat(output io.Writer, format string, apps []*api.App) error {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(format)
	if err != nil {
		return err
	}

	t := tabwriter.NewWriter(output, 10, 1, 3, ' ', 0)

	if err := tmpl.Execute(buf, tableHeaders); err != nil {
		return err
	}
	buf.WriteString("\n")

	for _, app := range apps {
		if err := tmpl.Execute(buf, app); err != nil {
			return err
		}
		buf.WriteString("\n")
	}
	buf.WriteTo(t)

	return t.Flush()
}
