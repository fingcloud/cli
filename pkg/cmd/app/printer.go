package app

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/util"
)

const (
	defaultFormat = "{{.Name}}\t{{.Platform}}\t{{.Status}}\t{{.Image}}\t{{.CreatedAt}}"
)

var tableHeaders = map[string]string{
	"ID":        "ID",
	"Name":      "NAME",
	"Label":     "LABEL",
	"Platform":  "PLATFORM",
	"Status":    "STATUS",
	"Image":     "IMAGE",
	"CreatedAt": "CREATED",
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
		if err := tmpl.Execute(buf, printContext{app}); err != nil {
			return err
		}
		buf.WriteString("\n")
	}
	buf.WriteTo(t)

	return t.Flush()
}

type printContext struct {
	app *api.App
}

func (c printContext) ID() string {
	return fmt.Sprintf("%d", c.app.ID)
}

func (c printContext) Name() string {
	return c.app.Name
}

func (c printContext) Platform() string {
	return c.app.Platform
}

func (c printContext) Image() string {
	return c.app.Image
}

func (c printContext) Status() string {
	if c.app.Status == "" {
		return "-"
	}
	return string(c.app.Status)
}

func (c printContext) CreatedAt() string {
	if c.app.CreatedAt == nil {
		return "-"
	}
	return util.FuzzyAgo(time.Now().Sub(*c.app.CreatedAt))
}
