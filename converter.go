package slogmattermost

import (
	"log/slog"

	"github.com/nafisfaysal/matterhook"
	slogcommon "github.com/samber/slog-common"
)

var SourceKey = "source"

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) *matterhook.Message

func DefaultConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) *matterhook.Message {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	if addSource {
		attrs = append(attrs, slogcommon.Source(SourceKey, record))
	}
	attrs = slogcommon.ReplaceAttrs(replaceAttr, []string{}, attrs...)

	// handler formatter
	message := &matterhook.Message{}
	message.Text = record.Message
	message.Attachments = []matterhook.Attachment{
		{
			Color:  colorMap[record.Level],
			Fields: []matterhook.Field{},
		},
	}

	attrToMattermostMessage("", attrs, message)
	return message
}

func attrToMattermostMessage(base string, attrs []slog.Attr, message *matterhook.Message) {
	for i := range attrs {
		attr := attrs[i]
		k := attr.Key
		v := attr.Value
		kind := attr.Value.Kind()

		if kind == slog.KindGroup {
			attrToMattermostMessage(base+k+".", v.Group(), message)
		} else {
			field := matterhook.Field{}
			field.Title = base + k
			field.Value = slogcommon.ValueToString(v)
			message.Attachments[0].Fields = append(message.Attachments[0].Fields, field)
		}

	}
}
