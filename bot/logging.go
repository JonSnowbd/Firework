package bot

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("BotStatus")
var logfile io.Writer

func init() {
	// Initiate the logging folder.
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir("./logs", 0777)
	}

	logfile, err := os.Create("./logs/" + fmt.Sprint("", int32(time.Now().Unix())) + ".log")
	if err != nil {
		log.Error(err)
	}

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	stdBackend := logging.NewLogBackend(os.Stderr, "", 0)
	fileBackend := logging.NewLogBackend(logfile, "", 0)

	fstdBackend := logging.NewBackendFormatter(stdBackend, format)
	ffileBackend := logging.NewBackendFormatter(fileBackend, format)
	logging.SetBackend(fstdBackend, ffileBackend)
}
