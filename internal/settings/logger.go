package settings

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	slogmulti "github.com/samber/slog-multi"
	slogwebhook "github.com/samber/slog-webhook/v2"
)

type DiscordLogger struct {
	Endpoint  string `yaml:"endpoint"`
	AvatarURL string `yaml:"avatarURL"`
	Embed     bool   `yaml:"embed"`
}

type WebhookLogger struct {
	Endpoint string `yaml:"endpoint"`
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
	Webhook *WebhookLogger `yaml:"webhook,omitempty"`
	File    *FileLogger    `yaml:"file,omitempty"`
	Stdout  *StdoutLogger  `yaml:"std,omitempty"`
}

func NewLogger(settingsPath string, loggers []LoggerConfig, with ...any) (*slog.Logger, error) {
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

			h := option.NewWebhookHandler()
			h = h.WithAttrs([]slog.Attr{{
				Key:   "avatarURL",
				Value: slog.StringValue(logger.Discord.AvatarURL),
			}})

			handlers = append(handlers, h)
		case logger.Webhook != nil:
			option := slogwebhook.Option{
				Level:    logger.Level,
				Endpoint: logger.Discord.Endpoint,
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

	logger := slog.New(slogmulti.Fanout(handlers...)).With(with...)

	return logger, nil
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordFooter struct {
	Text string `json:"text"`
}

type discordEmbed struct {
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Timestamp   string         `json:"timestamp"`
	Color       int            `json:"color"`
	Footer      discordFooter  `json:"footer"`
	Fields      []discordField `json:"fields"`
}

// type discordWebhookPayload struct {
// 	Content   string         `json:"content"`
// 	Username  string         `json:"username"`
// 	AvatarURL string         `json:"avatar_url"`
// 	Embeds    []discordEmbed `json:"embeds"`
// }

func DiscordEmbedConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	embed := discordEmbed{
		Title:       record.Level.String(),
		Type:        "rich",
		Description: record.Message,
		Color:       levelToDiscorColor(record.Level),
		Timestamp:   record.Time.Format(time.RFC3339),
	}

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "op" {
			embed.Footer = discordFooter{
				Text: fmt.Sprint(attr.Value),
			}
			return true
		}

		embed.Fields = append(embed.Fields, discordField{
			Name:  attr.Key,
			Value: fmt.Sprint(attr.Value),
		})
		return true
	})

	avatarURL := ""
	for _, attr := range loggerAttr {
		if attr.Key == "avatarURL" {
			avatarURL = fmt.Sprint(attr.Value)
			continue
		}

		embed.Fields = append(embed.Fields, discordField{
			Name:  attr.Key,
			Value: fmt.Sprint(attr.Value),
		})
	}

	return map[string]any{
		"content":    recordToText(record),
		"embeds":     []discordEmbed{embed},
		"avatar_url": avatarURL,
	}
}

func DiscordTextConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	text := recordToText(record)

	avatarURL := ""
	for _, attr := range loggerAttr {
		if attr.Key == "avatarURL" {
			avatarURL = fmt.Sprint(attr.Value)
			continue
		}

		text += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
	}

	return map[string]any{
		"content":    text,
		"avatar_url": avatarURL,
	}
}

func recordToText(record *slog.Record) string {
	text := fmt.Sprintf("time=%s level=%s msg=%s", record.Time.Format(time.RFC3339), record.Level, record.Message)
	record.Attrs(func(attr slog.Attr) bool {
		text += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
		return true
	})
	return text
}

func levelToDiscorColor(level slog.Level) int {
	switch level {
	case slog.LevelDebug:
		return 0x3498db
	case slog.LevelInfo:
		return 0x2ecc71
	case slog.LevelWarn:
		return 0xf1c40f
	case slog.LevelError:
		return 0xe74c3c
	default:
		return 0
	}
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
