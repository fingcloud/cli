package auth

import (
	"bytes"
	"io"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/util"
)

const (
	defaultFormat = "{{.Email}}\t{{.Status}}\t{{.LoginAt}}\t{{.LastUsedAt}}"
)

var tableHeaders = map[string]string{
	"Email":      "EMAIL",
	"Status":     "STATUS",
	"LoginAt":    "LOGIN AT",
	"LastUsedAt": "LAST USED AT",
}

func PrintFormat(output io.Writer, format string, sessions []session.Session) error {
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

	for _, sess := range sessions {
		if err := tmpl.Execute(buf, sessionPrint{sess}); err != nil {
			return err
		}
		buf.WriteString("\n")
	}
	buf.WriteTo(t)

	return t.Flush()
}

type sessionPrint struct {
	sess session.Session
}

func (s sessionPrint) Email() string {
	return s.sess.Email
}

func (s sessionPrint) Status() string {
	if !s.sess.Default {
		return "-"
	}
	return "ACTIVE"
}

func (s sessionPrint) LoginAt() string {
	if s.sess.LoginAt == nil {
		return "NOT SET"
	}
	return util.FuzzyAgo(time.Now().Sub(*s.sess.LoginAt))
}

func (s sessionPrint) LastUsedAt() string {
	if s.sess.LastUsedAt == nil {
		return "NOT SET"
	}
	return util.FuzzyAgo(time.Now().Sub(*s.sess.LastUsedAt))
}
