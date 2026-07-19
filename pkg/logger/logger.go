package logger

import (
	"io"
	"log"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	log *zerolog.Logger
}

func NewLogger(level string, isDev bool, serviceName string) *Logger {

	zlevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Printf("Ошибка уровня: %v", err)
	}

	var output io.Writer = os.Stdout

	if isDev {
		output = zerolog.ConsoleWriter{ Out: os.Stdout }
	}

	internalLog := zerolog.New(output).
		Level(zlevel).
		With().
		Str("service: ", serviceName).
		Timestamp().
		Logger()

	return &Logger{
		log: &internalLog,
	}
}

func (l *Logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	l.log.Error().Err(err).Msg(msg)
}