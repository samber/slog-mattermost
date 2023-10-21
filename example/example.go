package main

import (
	"fmt"
	"time"

	slogmattermost "github.com/samber/slog-mattermost/v2"

	"log/slog"
)

func main() {
	url := "https://your-mattermost-server.com/hooks/xxx-generatedkey-xxx"
	channel := "alerts"

	logger := slog.New(slogmattermost.Option{Level: slog.LevelDebug, WebhookURL: url, Channel: channel}.NewMattermostHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("A message")
}
