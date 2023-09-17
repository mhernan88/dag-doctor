package shared

import (
	"log"
	"log/slog"
	"os"
)

func GetLogger() (*slog.Logger, *os.File) {
	logFilename, err := GetLogFilename()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logFilename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return slog.New(slog.NewJSONHandler(f , nil)), f
}
