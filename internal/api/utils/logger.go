package utils

import (
	"github.com/mbndr/logo"
	"os"
)

func Logger(path string) *logo.Logger {
	cli := logo.NewReceiver(os.Stderr, "")
	cli.Color = true
	cli.Level = logo.DEBUG

	file, _ := logo.Open(path)
	out := logo.NewReceiver(file, "")
	out.Format = "%s: %s"

	return logo.NewLogger(cli, out)
}
