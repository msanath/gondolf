package ctxslog

import (
	"fmt"
	"log/slog"
	"strings"
)

// NewCustomHandler returns a custom handler which adds source,
// replaces source file path with just the file name and its penultimate directory.
func NewCustomHandler(logLevel slog.Level) *slog.HandlerOptions {
	replaceAttrFunc := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.SourceKey:
			if source, ok := a.Value.Any().(*slog.Source); ok {
				// The source file contains the entire path. We replace it here with
				// just the file name and its penultimate directory.
				idx := strings.LastIndexByte(source.File, '/')
				if idx == -1 {
					break
				}
				idx = strings.LastIndexByte(source.File[:idx], '/')
				if idx == -1 {
					break
				}
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File[idx+1:], source.Line))
			}
		}
		return a
	}

	return &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: replaceAttrFunc,
	}
}
