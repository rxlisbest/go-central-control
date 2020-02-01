package utils

import (
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}%{color:reset}`,
)

func initLogging() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(-1, "")
	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}
