package util

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultErrorExitCode = 1
)

type ErrHandler func(string, int)

var (
	ErrExit         = fmt.Errorf("exit")
	FatalErrHandler = fatal
)

func CheckErr(err error) {
	checkErr(err, FatalErrHandler)
}

func fatal(msg string, code int) {
	if msg != "" && !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprint(os.Stderr, msg)
	os.Exit(code)
}

func checkErr(err error, handleErr ErrHandler) {
	if err == nil {
		return
	}

	if err == ErrExit {
		handleErr("", DefaultErrorExitCode)
	} else {
		msg := err.Error()
		if !strings.HasPrefix(msg, "error: ") {
			msg = fmt.Sprintf("error: %s", msg)
		}
		handleErr(msg, DefaultErrorExitCode)
	}
}

func DefaultSubCommandRun(out io.Writer) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cmd.SetOutput(out)
		cmd.Help()
		CheckErr(ErrExit)
	}
}
