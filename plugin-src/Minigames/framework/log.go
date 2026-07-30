package framework

import "log/slog"


type LogWrapper struct {
	*slog.Logger
}


func NewLogWrapper(logger *slog.Logger) *LogWrapper {
	return &LogWrapper{Logger: logger}
}
