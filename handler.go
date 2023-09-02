package slogmattermost

import (
	"context"

	"log/slog"

	"github.com/nafisfaysal/matterhook"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Mattermost webhook url
	WebhookURL string
	// Mattermost channel (default: webhook channel)
	Channel string
	// bot username (default: webhook username)
	Username string
	// bot emoji (default: webhook emoji)
	IconEmoji string

	// optional: customize Mattermost event builder
	Converter Converter
}

func (o Option) NewMattermostHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.WebhookURL == "" {
		panic("missing Mattermost webhook url")
	}

	return &MattermostHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*MattermostHandler)(nil)

type MattermostHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *MattermostHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *MattermostHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	message := converter(h.attrs, &record)

	if h.option.Channel != "" {
		message.Channel = h.option.Channel
	}

	if h.option.Username != "" {
		message.Username = h.option.Username
	}

	if h.option.IconEmoji != "" {
		message.IconEmoji = h.option.IconEmoji
	}

	return matterhook.Send(h.option.WebhookURL, *message)
}

func (h *MattermostHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MattermostHandler{
		option: h.option,
		attrs:  appendAttrsToGroup(h.groups, h.attrs, attrs),
		groups: h.groups,
	}
}

func (h *MattermostHandler) WithGroup(name string) slog.Handler {
	return &MattermostHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
