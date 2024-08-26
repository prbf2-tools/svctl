package settings

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	slogmulti "github.com/samber/slog-multi"
	slogwebhook "github.com/samber/slog-webhook/v2"
)

type DiscordLogger struct {
	Endpoint string `yaml:"endpoint"`
	Embed    bool   `yaml:"embed"`
}

type logType string

const (
	textLogger logType = "text"
	jsonLogger logType = "json"
)

type FileLogger struct {
	Path string  `yaml:"path"`
	Type logType `yaml:"type"`
}

type StdoutLogger struct {
	Type logType `yaml:"type"`
}

type LoggerConfig struct {
	Level   slog.Level     `yaml:"level"`
	Discord *DiscordLogger `yaml:"discord,omitempty"`
	File    *FileLogger    `yaml:"file,omitempty"`
	Stdout  *StdoutLogger  `yaml:"std,omitempty"`
}

func NewLogger(settingsPath string, loggers []LoggerConfig) (*slog.Logger, error) {
	var handlers []slog.Handler

	for _, logger := range loggers {
		switch {
		case logger.Discord != nil:
			option := slogwebhook.Option{
				Level:    logger.Level,
				Endpoint: logger.Discord.Endpoint,
			}

			if logger.Discord.Embed {
				option.Converter = DiscordEmbedConverter
			} else {
				option.Converter = DiscordTextConverter
			}

			handlers = append(handlers, option.NewWebhookHandler())
		case logger.File != nil:
			path := logger.File.Path
			if !filepath.IsAbs(path) {
				path = filepath.Join(settingsPath, path)
			}

			file, err := openOrCreateFile(path)
			if err != nil {
				return nil, err
			}

			options := &slog.HandlerOptions{
				Level: logger.Level,
			}

			switch logger.File.Type {
			case textLogger:
				handlers = append(handlers, slog.NewTextHandler(file, options))
			case jsonLogger:
				handlers = append(handlers, slog.NewJSONHandler(file, options))
			default:
				return nil, errors.New("invalid logger type")
			}
		case logger.Stdout != nil:
			options := &slog.HandlerOptions{
				Level: logger.Level,
			}

			switch logger.Stdout.Type {
			case textLogger:
				handlers = append(handlers, slog.NewTextHandler(os.Stdout, options))
			case jsonLogger:
				handlers = append(handlers, slog.NewJSONHandler(os.Stdout, options))
			default:
				return nil, errors.New("invalid logger type")
			}
		}
	}

	return slog.New(slogmulti.Fanout(handlers...)), nil
}

func DiscordEmbedConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	return nil
}

func DiscordTextConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	return nil
}

func openOrCreateFile(path string) (*os.File, error) {
	_, err := os.Open(path)

	if os.IsNotExist(err) {
		return os.Create(path)
	}

	if err != nil {
		return nil, err
	}

	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
}
