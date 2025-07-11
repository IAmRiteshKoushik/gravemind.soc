package pkg

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log *LoggerService

type LoggerService struct {
	log zerolog.Logger
	env string
}

func NewLoggerService(env string, file *os.File) *LoggerService {
	var output io.Writer

	if env == "development" {
		// Logging to both file and std.out during development
		fileOut := zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}
		consoleOut := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		output = zerolog.MultiLevelWriter(consoleOut, fileOut)

	} else if env == "production" {
		// Logging only to file during production
		output = zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339}

	} else {
		panic(errors.New("could not identify environment"))
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	return &LoggerService{
		log: logger,
		env: env,
	}
}

func (l *LoggerService) Info(msg string) {
	l.log.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) Warn(msg string) {
	l.log.WithLevel(zerolog.WarnLevel).Msg(msg)
}

func (l *LoggerService) Error(msg string, err error) {
	l.log.WithLevel(zerolog.ErrorLevel).Err(err).Msg(msg)
}

func (l *LoggerService) Fatal(msg string) {
	l.log.WithLevel(zerolog.FatalLevel).Msg(msg)
}
