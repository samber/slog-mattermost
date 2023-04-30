package slogmattermost

import (
	"encoding"
	"fmt"
	"strconv"

	"github.com/nafisfaysal/matterhook"
	"golang.org/x/exp/slog"
)

type Converter func(loggerAttr []slog.Attr, record *slog.Record) *matterhook.Message

func DefaultConverter(loggerAttr []slog.Attr, record *slog.Record) *matterhook.Message {
	message := &matterhook.Message{}
	message.Text = record.Message
	message.Attachments = []matterhook.Attachment{
		{
			Color:  colorMap[record.Level],
			Fields: []matterhook.Field{},
		},
	}

	attrToMattermostMessage("", loggerAttr, message)
	record.Attrs(func(attr slog.Attr) bool {
		attrToMattermostMessage("", []slog.Attr{attr}, message)
		return true
	})

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
			field.Value = attrToValue(v)
			message.Attachments[0].Fields = append(message.Attachments[0].Fields, field)
		}

	}
}

func attrToValue(v slog.Value) string {
	kind := v.Kind()

	switch kind {
	case slog.KindAny:
		return anyValueToString(v)
	case slog.KindLogValuer:
		return anyValueToString(v)
	case slog.KindGroup:
		// not expected to reach this line
		return anyValueToString(v)
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().UTC().String()
	default:
		return anyValueToString(v)
	}
}

func anyValueToString(v slog.Value) string {
	if tm, ok := v.Any().(encoding.TextMarshaler); ok {
		data, err := tm.MarshalText()
		if err != nil {
			return ""
		}

		return string(data)
	}

	return fmt.Sprintf("%+v", v.Any())
}
