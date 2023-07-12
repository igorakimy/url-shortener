package slogdiscard

import (
	"context"

	"golang.org/x/exp/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Просто игнорируем запись журнала
	return nil
}

func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Возвращаем тот же обработчик, т.к. нет атрибутов для сохранения
	return h
}

func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	// Возвращаем тот же обработчик, т.к. нет группы для сохранения
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Всегда возвращать false, т.к. запись журнала игнорируется
	return false
}
