package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

func InitLogger() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(i.(string))
		switch level {
		case "ERROR":
			return colorRed + "ERROR" + colorReset
		case "WARN":
			return colorYellow + "WARN " + colorReset
		case "INFO":
			return colorCyan + "INFO " + colorReset
		case "DEBUG":
			return colorBlue + "DEBUG" + colorReset
		default:
			return level
		}
	}

	output.FormatMessage = func(i interface{}) string {
		return colorGreen + i.(string) + colorReset
	}

	log.Logger = log.Output(output).With().Timestamp().Logger()
}
