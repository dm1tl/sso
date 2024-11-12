package slogerr

import (
	"log/slog"
)

// func for adding an error for current log
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
