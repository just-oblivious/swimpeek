package config

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func GetLogger(prefix string) *log.Logger {
	styles := log.DefaultStyles()
	styles.Levels[log.FatalLevel] = styles.Levels[log.ErrorLevel]

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
		ReportCaller:    false,
		TimeFormat:      time.TimeOnly,
		Level:           log.InfoLevel,
		Prefix:          prefix,
	})

	logger.SetStyles(styles)
	return logger
}
